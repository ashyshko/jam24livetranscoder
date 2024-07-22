package protocol

type Init struct {
	Presets               []int `json:"presets"`
	OutputTimestampOffset int64 `json:"output_timestamp_offset"`
}

type ClientVideo struct {
	PacketPts int64 `json:"packet_pts"`
	PacketDts int64 `json:"packet_dts"`
}

type ServerVideo struct {
	PresetIndex  int   `json:"preset_index"`
	SegmentIndex int64 `json:"segment_index"`

	PacketPts int64 `json:"packet_pts"`
	PacketDts int64 `json:"packet_dts"`
}

type packetType string // should not be used externally, ServerVisitor/ClientVisitor should be used for reading packetType, Make*Packet to write packetType

type Packet struct {
	Type   packetType
	JSON   interface{}
	Binary []byte
}

type ServerVisitor interface {
	Init(obj Init) error
	Video(obj ClientVideo, payload []byte) error

	UnknownPacket(packetType string, packetData map[string]interface{}, payload []byte) error
}

type ClientVisitor interface {
	Video(obj ServerVideo)

	UnknownPacket(packetType string, packetData map[string]interface{}, payload []byte) error
}
