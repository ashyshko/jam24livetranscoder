package protocol

import (
	"golang.org/x/net/websocket"
)

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

func RecvServer(ws *websocket.Conn, visitor ServerVisitor) error {
	packet, err := recv(ws)
	if err != nil {
		return err
	}

	return handle(packet, visitor, serverHandlers)
}
