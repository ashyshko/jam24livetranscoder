package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"

	"jam24livetranscoder/protocol"

	"golang.org/x/net/websocket"
)

func main() {
	videoFile := flag.String("i", "", "input video file")
	ffprobeOutput := flag.String("c", "", "ffprobe output with flags: '-select_streams 0 -print_format compact= -show_frames'")
	transcoderUrl := flag.String("o", "ws://localhost:8898/ws", "websocket url of transcoder")
	savePath := flag.String("d", "", "path prefix for saving output video")
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

	init := protocol.Init{
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
	}

	err = protocol.Send(
		ws,
		protocol.MakeInit(init),
	)

	if err != nil {
		log.Fatalf("send init failed: %s", err)
	}

	ffprobeOutputFile, err := os.OpenFile(*ffprobeOutput, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open ffprobe output error: %v", err)
		return
	}
	defer ffprobeOutputFile.Close()

	inputFile, err := os.OpenFile(*videoFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open input file error: %v", err)
		return
	}
	defer inputFile.Close()

	recv := newReceiver(*savePath, init.Presets)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			err := protocol.RecvClient(ws, &recv)
			if err != nil {
				log.Printf("recv failed: %s", err)
				break
			}
		}
	}()

	//

	ffprobeOutputScanner := bufio.NewScanner(ffprobeOutputFile)
	for ffprobeOutputScanner.Scan() {
		line := ffprobeOutputScanner.Text()
		mediaType, err := parseFFprobeParam("media_type", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		if mediaType != "video" {
			log.Printf("Ignore media type %s", mediaType)
			continue
		}

		keyFrame, err := parseFFprobeParamInt("key_frame", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		ptsTime, err := parseFFprobeParamFloat("pts_time", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		dtsTime, err := parseFFprobeParamFloat("pkt_dts_time", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		durationTime, err := parseFFprobeParamFloat("pkt_duration_time", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		pktPos, err := parseFFprobeParamInt("pkt_pos", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		pktSize, err := parseFFprobeParamInt("pkt_size", line)
		if err != nil {
			log.Printf("Ignore line: %s (%s)", line, err)
			continue
		}

		_, err = inputFile.Seek(pktPos, 0)
		if err != nil {
			log.Fatalf("seek error: %s", err)
		}

		buf := make([]byte, pktSize)

		_, err = inputFile.Read(buf)
		if err != nil {
			log.Fatalf("read error: %s", err)
		}

		err = protocol.Send(ws, protocol.MakeVideoPacket(
			protocol.VideoPacket{
				PacketPts:      int64(ptsTime * 1000.0),
				PacketDts:      int64(dtsTime * 1000.0),
				PacketDuration: int64(durationTime * 1000.0),
				KeyFrame:       keyFrame > 0,
			},
			buf,
		))
	}

	if err := ffprobeOutputScanner.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

	log.Printf("Sending eof")
	err = protocol.Send(ws, protocol.MakeEof())
	if err != nil {
		log.Fatalf("send eof failed: %s", err)
	}

	log.Printf("waiting read complete")
	wg.Wait()

	log.Printf("Completed!")
}

func parseFFprobeParam(param string, line string) (string, error) {
	regExp := regexp.MustCompile(fmt.Sprintf(`\|%s=([^\|]+)\|`, param))
	res := regExp.FindStringSubmatch(line)
	if res == nil || len(res) < 2 {
		return "", fmt.Errorf("no string submatch found for param %s: %s (%+v)", param, line, res)
	}

	if res[1] == "" {
		return "", fmt.Errorf("empty submatch for param %s: %s (%+v)", param, line, res)
	}

	return res[1], nil
}

func parseFFprobeParamInt(param string, line string) (int64, error) {
	str, err := parseFFprobeParam(param, line)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(str, 10, 64)
}

func parseFFprobeParamFloat(param string, line string) (float64, error) {
	str, err := parseFFprobeParam(param, line)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(str, 64)
}
