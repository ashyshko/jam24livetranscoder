package protocol

import (
	"encoding/json"
	"fmt"
	"io"
)

type visitorBase interface {
	UnknownPacket(packetType string) error
}

type recvHandler[Visitor visitorBase] struct {
	packetType packetType
	action     func(packet recvPacket, visitor Visitor) error
}

func createRecvHandler[K any, Visitor visitorBase](packetType packetType, handlerFn func(visitor Visitor, obj K, payload []byte) error) recvHandler[Visitor] {
	return recvHandler[Visitor]{
		packetType: packetType,
		action: func(packet recvPacket, visitor Visitor) error {
			var obj K
			err := json.Unmarshal(packet.JSON, &obj)
			if err != nil {
				return fmt.Errorf("%s unmarshal failed: %s", string(packetType), err)
			}

			return handlerFn(visitor, obj, packet.Binary)
		},
	}
}

var serverHandlers = []recvHandler[ServerVisitor]{
	createRecvHandler(packetTypeInit, func(visitor ServerVisitor, obj Init, _payload []byte) error {
		return visitor.Init(obj)
	}),
	createRecvHandler(packetTypeVideoPacket, func(visitor ServerVisitor, obj VideoPacket, payload []byte) error {
		return visitor.VideoPacket(obj, payload)
	}),
	createRecvHandler(packetTypeEof, func(visitor ServerVisitor, obj interface{}, payload []byte) error {
		return visitor.Eof()
	}),
}

func handle[Visitor visitorBase](packet recvPacket, visitor Visitor, handlers []recvHandler[Visitor]) error {
	for _, handler := range handlers {
		if packet.Type != string(handler.packetType) {
			continue
		}

		return handler.action(packet, visitor)
	}

	return visitor.UnknownPacket(packet.Type)
}

func RecvServer(reader io.Reader, visitor ServerVisitor) error {
	packet, err := recv(reader)
	if err != nil {
		return err
	}

	return handle(packet, visitor, serverHandlers)
}
