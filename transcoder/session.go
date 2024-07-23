package main

import (
	"fmt"
	"jam24livetranscoder/protocol"
	"log"

	"golang.org/x/net/websocket"
)

type session struct {
	ws   *websocket.Conn
	init *protocol.Init
}

func newSession(ws *websocket.Conn) session {
	return session{
		ws: ws,
	}
}

func (this *session) Init(obj protocol.Init) error {
	if this.init != nil {
		return fmt.Errorf("init received twice")
	}

	this.init = &obj
	log.Printf("init %+v", obj)
	return nil
}

func (this *session) VideoHeader(payload []byte) error {
	log.Printf("video header %d", len(payload))
	return nil
}

func (this *session) VideoPacket(obj protocol.VideoPacket, payload []byte) error {
	if this.init == nil {
		return fmt.Errorf("video received before init")
	}

	segmentIndex := int(float64(obj.PacketPts) / (float64(this.init.TicksPerSecond) * 3))
	nextSegmentIndex := int((float64(obj.PacketPts) + 0.034*float64(this.init.TicksPerSecond)) / (float64(this.init.TicksPerSecond) * 3))
	segmentEnd := nextSegmentIndex > segmentIndex

	log.Printf("video %+v %+v", obj, len(payload))

	for presetIndex := range this.init.Presets {
		err := protocol.Send(
			this.ws,
			protocol.MakeOutputVideoPacket(
				protocol.OutputVideoPacket{
					PresetIndex:  presetIndex,
					SegmentIndex: segmentIndex,
					SegmentEnd:   segmentEnd,
					DurationMs:   3000,
					PacketPts:    obj.PacketPts,
					PacketDts:    obj.PacketDts,
				},
				payload),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (this *session) Eof() error {
	if this.init == nil {
		return fmt.Errorf("eof received before init")
	}

	log.Printf("eof")

	this.ws.Close()

	return nil
}

func (this *session) UnknownPacket(packetType string) error {
	log.Printf("unknown packet %s", packetType)
	return nil
}
