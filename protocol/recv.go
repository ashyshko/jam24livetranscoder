package protocol

import (
	"encoding/binary"
	"io"
)

type recvPacket struct {
	Type   packetType
	JSON   []byte
	Binary []byte
}

func recv(reader io.Reader) (packet recvPacket, err error) {
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

	packet.Type = packetType(packetTypeBuffer)

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

	return
}
