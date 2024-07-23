package main

import (
	"fmt"
	"jam24livetranscoder/protocol"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

type session struct {
	ws            *websocket.Conn
	transcoderObj WdpInstance
	init          *protocol.Init
	closed        *atomic.Bool
	readCompleted *sync.WaitGroup
}

func newSession(ws *websocket.Conn, wdpWrap WdpWrap) session {
	transcoderObj := wdpWrap.Alloc()

	readCompletedWg := new(sync.WaitGroup)
	closed := new(atomic.Bool)

	readCompletedWg.Add(1)

	go func() {
		defer readCompletedWg.Done()
		for !closed.Load() {
			packet, err := transcoderObj.TakePacket()
			if err != nil {
				log.Printf("take frame failed: %s", err)
				break
			}

			if packet == nil {
				time.Sleep(5 * time.Millisecond)
				continue
			}

			// log.Printf("recv packet: stream=%d pts=%d dts=%d, size=%d", packet.PresetIndex, packet.Pts, packet.Dts, len(packet.Payload))
			err = protocol.Send(ws, protocol.MakeOutputVideoPacket(protocol.OutputVideoPacket{
				PresetIndex:  packet.PresetIndex,
				SegmentIndex: 0,
				SegmentEnd:   false,
				DurationMs:   3000,
				PacketPts:    packet.Pts,
				PacketDts:    packet.Dts,
				KeyFrame:     packet.Keyframe,
			}, packet.Payload))
			if err != nil {
				log.Printf("protocol.Send failed: %s", err)
				break
			}
		}
	}()

	return session{
		ws:            ws,
		transcoderObj: transcoderObj,
		closed:        closed,
		readCompleted: readCompletedWg,
	}
}

func (this *session) Close() {
	this.closed.Store(true)
	this.readCompleted.Wait()
	this.transcoderObj.Close()
}

func (this *session) Init(obj protocol.Init) error {
	if this.init != nil {
		return fmt.Errorf("init received twice")
	}

	this.init = &obj

	for _, preset := range obj.Presets {
		err := this.transcoderObj.AddPreset(preset.Width, preset.Height, preset.Bitrate, preset.Framerate)
		if err != nil {
			return fmt.Errorf("can't add preset: %s", err)
		}
	}

	err := this.transcoderObj.Init(int(obj.TicksPerSecond), obj.NextSegmentIndex, obj.TargetSegmentDurationMs)
	if err != nil {
		return fmt.Errorf("can't init: %s", err)
	}

	log.Printf("init completed %d: %+v", this.transcoderObj, obj)

	return nil
}

func (this *session) VideoHeader(payload []byte) error {
	log.Printf("video header %d", len(payload))

	err := this.transcoderObj.OnHeader(payload)
	if err != nil {
		return fmt.Errorf("onHeader failed: %s", err)
	}

	return nil
}

func (this *session) VideoPacket(obj protocol.VideoPacket, payload []byte) error {
	if this.init == nil {
		return fmt.Errorf("video received before init")
	}

	keyFrame := int8(0)
	if obj.KeyFrame {
		keyFrame = 1
	}

	err := this.transcoderObj.OnVideo(obj.PacketPts, obj.PacketDts, keyFrame, payload)
	if err != nil {
		return fmt.Errorf("onVideo failed: %s", err)
	}

	return nil
}

func (this *session) Eof() error {
	if this.init == nil {
		return fmt.Errorf("eof received before init")
	}

	log.Printf("eof")

	err := this.transcoderObj.OnEof()
	if err != nil {
		return fmt.Errorf("onEof failed: %s", err)
	}

	return nil
}

func (this *session) UnknownPacket(packetType string) error {
	log.Printf("unknown packet %s", packetType)
	return nil
}
