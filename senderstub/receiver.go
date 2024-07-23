package main

import (
	"encoding/binary"
	"fmt"
	"jam24livetranscoder/protocol"
	"log"
	"os"
	"path"
)

type receiver struct {
	writers []*os.File
}

func newReceiver(savePath string, presets []protocol.Preset) receiver {
	writers := []*os.File{}

	for _, preset := range presets {
		fileName := path.Join(savePath, fmt.Sprintf("output_%dp.h264", preset.Height))
		writer, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("can't open file %s: %s", fileName, err)
		}
		writers = append(writers, writer)
	}

	return receiver{
		writers: writers,
	}
}

func (this *receiver) OutputVideoHeader(obj protocol.OutputVideoHeader, payload []byte) error {
	log.Printf("VideoHeader %v", obj)
	return nil
}

func (this *receiver) OutputVideoPacket(obj protocol.OutputVideoPacket, payload []byte) error {
	err := replaceStartCodes(payload)
	if err != nil {
		return err
	}

	_, err = this.writers[obj.PresetIndex].Write(payload)
	return err
}

func (this *receiver) ExpringSoon(obj protocol.ExpiringSoon) error {
	log.Printf("ExpiringSoone %v", obj)
	return nil
}

func (this *receiver) UnknownPacket(packetType string) error {
	log.Printf("UnknownPacket %s", packetType)
	return nil
}

func replaceStartCodes(buffer []byte) error {
	pos := 0

	for pos+4 < len(buffer) {
		value := binary.BigEndian.Uint32(buffer[pos : pos+4])

		if value == 1 {
			log.Printf("start code is already 0001, skipping replacing")
			return nil
		}

		binary.BigEndian.PutUint32(buffer[pos:pos+4], 1)

		pos += int(value) + 4
	}

	if pos != len(buffer) {
		return fmt.Errorf("Unexpected buffer end")
	}

	return nil

}
