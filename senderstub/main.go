package main

import (
	"flag"
	"log"

	"jam24livetranscoder/protocol"

	"golang.org/x/net/websocket"
)

func main() {
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

	err = protocol.Send(
		ws,
		protocol.MakeInit(protocol.Init{
			Presets: []protocol.Preset{
				{
					Width:   640,
					Height:  360,
					Bitrate: 1000,
				},
				{
					Width:   426,
					Height:  240,
					Bitrate: 500,
				},
			},
			TicksPerSecond:          1000,
			NextSegmentIndex:        0,
			TargetSegmentDurationMs: 3000,
		}),
	)

	if err != nil {
		log.Fatalf("send init failed: %s", err)
	}

	// send frames

	err = protocol.Send(ws, protocol.MakeEof())
	if err != nil {
		log.Fatalf("send eof failed: %s", err)
	}

}
