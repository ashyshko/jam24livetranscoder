package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

func Send(ws *websocket.Conn, packet Packet) error {
	commandPayload, err := json.Marshal(packet.JSON)
	if err != nil {
		return fmt.Errorf("Can't marshal command: %s", err)
	}

	packetType := []byte(packet.Type)

	output := make([]byte, 0, 4+len(packetType)+4+len(commandPayload)+4+len(packet.Binary))
	output = binary.LittleEndian.AppendUint32(output, uint32(len(packetType)))
	output = append(output, packetType...)
	output = binary.LittleEndian.AppendUint32(output, uint32(len(commandPayload)))
	output = append(output, commandPayload...)
	output = binary.LittleEndian.AppendUint32(output, uint32(len(packet.Binary)))
	if packet.Binary != nil {
		output = append(output, packet.Binary...)
	}

	err = websocket.Message.Send(ws, output)
	return err
}
