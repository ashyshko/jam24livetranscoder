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

func (this *session) Init(obj protocol.Init) error {
	if this.init != nil {
		return fmt.Errorf("init received twice")
	}

	this.init = &obj
	log.Printf("init %+v", obj)
	return nil
}

func (this *session) VideoPacket(obj protocol.VideoPacket, payload []byte) error {
	if this.init == nil {
		return fmt.Errorf("video received before init")
	}

	log.Printf("video %+v %+v", obj, payload)

	return nil
}

func (this *session) Eof() error {
	if this.init == nil {
		return fmt.Errorf("eof received before init")
	}

	log.Printf("eof")

	return nil
}

func (this *session) UnknownPacket(packetType string) error {
	log.Printf("unknown packet %s", packetType)
	return nil
}
