package protocol

type Preset struct {
	Width      int `json:"w"`
	Height     int `json:"h"`
	Bitrate    int `json:"bitrate"`
	ProfileIdc int `json:"profileIdc"`
	LevelIdc   int `json:"levelIdc"`
	Framerate  int `json:"fps"`
}

// from client to server

const (
	packetTypeInit        packetType = "init"
	packetTypeVideoHeader packetType = "videoHeader"
	packetTypeVideoPacket packetType = "videoPacket"
	packetTypeEof         packetType = "eof"
)

type Init struct {
	Presets                 []Preset `json:"presets"`
	TicksPerSecond          int64    `json:"ticksPerSecond"`
	NextSegmentIndex        int      `json:"nextSegmentIndex"`
	TargetSegmentDurationMs int      `json:"targetSegmentDurationMs"`
}

type VideoPacket struct {
	PacketPts      int64 `json:"pts"`
	PacketDts      int64 `json:"dts"`
	PacketDuration int64 `json:"duration"`
	KeyFrame       bool  `json:"keyFrame"`
}

type ServerVisitor interface {
	Init(obj Init) error
	VideoHeader(payload []byte) error
	VideoPacket(obj VideoPacket, payload []byte) error
	Eof() error

	UnknownPacket(packetType string) error
}

// from server to client

const (
	packetTypeOutputVideoHeader packetType = "outputVideoHeader"
	packetTypeOutputVideoPacket packetType = "outputVideoPacket"
	packetTypeExpiringSoon      packetType = "expiringSoon"
)

type OutputVideoHeader struct {
	PresetIndex int `json:"presetIndex"`
}

type OutputVideoPacket struct {
	PresetIndex  int  `json:"presetIndex"`
	SegmentIndex int  `json:"segmentIndex"`
	SegmentEnd   bool `json:"segmentEnd"`
	DurationMs   int  `json:"durationMs"`
}

type ExpiringSoon struct {
	FinalSegmentIndex int `json:"finalSegmentIndex"`
}

type ClientVisitor interface {
	OutputVideoHeader(obj OutputVideoHeader, payload []byte) error
	OutputVideoPacket(obj OutputVideoPacket, payload []byte) error
	ExpringSoon(obj ExpiringSoon) error

	UnknownPacket(packetType string) error
}

// sending packet

type packetType string // should not be used externally, ServerVisitor/ClientVisitor should be used for reading packetType, Make* to write packetType

type Packet struct {
	Type   packetType
	JSON   interface{}
	Binary []byte
}
