package protocol

import "golang.org/x/net/websocket"

var clientHandlers = []recvHandler[ClientVisitor]{
	createRecvHandler(packetTypeOutputVideoPacket, func(visitor ClientVisitor, obj OutputVideoPacket, payload []byte) error {
		return visitor.OutputVideoPacket(obj, payload)
	}),
	createRecvHandler(packetTypeExpiringSoon, func(visitor ClientVisitor, obj ExpiringSoon, _payload []byte) error {
		return visitor.ExpringSoon(obj)
	}),
}

func RecvClient(ws *websocket.Conn, visitor ClientVisitor) error {
	packet, err := recv(ws)
	if err != nil {
		return err
	}

	return handle(packet, visitor, clientHandlers)
}
