package protocol

import (
	"encoding/json"
	"fmt"
	"io"
)

func RecvServer(reader io.Reader, visitor ServerVisitor) error {
	packet, err := recv(reader)
	if err != nil {
		return err
	}

	switch packet.Type {
	case packetTypeInit:
		init := Init{}
		err = json.Unmarshal(packet.JSON, &init)
		if err != nil {
			return fmt.Errorf("init unmarshal failed: %s", err)
		}

		return visitor.Init(init)

	case packetTypeClientVideo:
		clientVideo := ClientVideo{}
		err = json.Unmarshal(packet.JSON, &clientVideo)
		if err != nil {
			return fmt.Errorf("clientVideo unmarshal failed: %s", err)
		}

		return visitor.Video(clientVideo, packet.Binary)

	default:
		packetData := map[string]interface{}{}
		err = json.Unmarshal(packet.JSON, &packetData)
		if err != nil {
			return fmt.Errorf("unknown packet unmarshal failed: %s", err)
		}

		return visitor.UnknownPacket(string(packet.Type), packetData, packet.Binary)
	}

}
