package main

import (
	"fmt"
	"io"
	"jam24livetranscoder/protocol"
	"log"
)

type session struct {
	writer io.Writer
	init   *protocol.Init
}

func newSession(writer io.Writer) session {
	return session{
		writer: writer,
	}
}

func (this *session) Init(init protocol.Init) error {
	if this.init != nil {
		return fmt.Errorf("init received twice")
	}

	this.init = &init
	log.Printf("init %+v", init)
	return nil
}

func (this *session) Video(clientVideo protocol.ClientVideo, payload []byte) error {
	if this.init == nil {
		return fmt.Errorf("video received before init")
	}

	log.Printf("video %+v %+v", clientVideo, payload)

	return nil
}

func (this *session) UnknownPacket(packetType string, packetData map[string]interface{}, payload []byte) error {
	log.Printf("unknown packet %s %+v %+v", packetType, packetData, payload)
	return nil
}
