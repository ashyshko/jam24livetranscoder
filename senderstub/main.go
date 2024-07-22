package main

import (
	"flag"
	"fmt"
	"log"

	"jam24livetranscoder/protocol"

	"golang.org/x/net/websocket"
)

func main() {
	fmt.Println("Hello, World!")
	videoFile := flag.String("i", "", "input video file")
	ffprobeOutput := flag.String("c", "", "ffprobe output with flags: '-select_streams 0 -print_format compact= -show_frames'")
	transcoderUrl := flag.String("o", "ws://localhost:8898/ws", "websocket url of transcoder")
	flag.Parse()

	if *videoFile == "" {
		flag.Usage()
		log.Fatalf("no video file provided")
	}

	if *ffprobeOutput == "" {
		flag.Usage()
		log.Fatalf("no ffprobe output provided")
	}

	origin := "http://localhost/"
	ws, err := websocket.Dial(*transcoderUrl, "", origin)
	if err != nil {
		log.Fatalf("websocket dial failed: %s", err)
	}

	protocol.Send(
		ws,
		protocol.MakeInitPacket(protocol.Init{
			Presets:               []int{540, 360, 240},
			OutputTimestampOffset: 0,
		}),
	)

}
