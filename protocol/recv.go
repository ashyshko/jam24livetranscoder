package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

type recvPacket struct {
	Type   string
	JSON   []byte
	Binary []byte
}

func recv(ws *websocket.Conn) (packet recvPacket, err error) {
	var data []byte
	err = websocket.Message.Receive(ws, &data)
	if err != nil {
		return
	}

	reader := bytes.NewReader(data)

	sizeBuffer := make([]byte, 4)
	_, err = reader.Read(sizeBuffer)
	if err != nil {
		return
	}

	packetTypeBuffer := make([]byte, binary.LittleEndian.Uint32(sizeBuffer))
	_, err = reader.Read(packetTypeBuffer)
	if err != nil {
		return
	}

	packet.Type = string(packetTypeBuffer)

	_, err = reader.Read(sizeBuffer)
	if err != nil {
		return
	}

	packet.JSON = make([]byte, binary.LittleEndian.Uint32(sizeBuffer))
	_, err = reader.Read(packet.JSON)
	if err != nil {
		return
	}

	_, err = reader.Read(sizeBuffer)
	if err != nil {
		return
	}

	payloadSize := binary.LittleEndian.Uint32(sizeBuffer)
	packet.Binary = make([]byte, int(payloadSize))

	if payloadSize > 0 {
		_, err = reader.Read(packet.Binary)
		if err != nil {
			return
		}
	}

	if reader.Len() != 0 {
		err = fmt.Errorf("left unread data bytes %d", reader.Len())
		return
	}

	return
}

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

func handle[Visitor visitorBase](packet recvPacket, visitor Visitor, handlers []recvHandler[Visitor]) error {
	for _, handler := range handlers {
		if packet.Type != string(handler.packetType) {
			continue
		}

		return handler.action(packet, visitor)
	}

	return visitor.UnknownPacket(packet.Type)
}
