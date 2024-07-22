package protocol

const (
	packetTypeInit        packetType = "init"
	packetTypeClientVideo packetType = "client_video"
	packetTypeServerVideo packetType = "server_video"
)

func MakeInitPacket(init Init) Packet {
	return Packet{
		Type: packetTypeInit,
		JSON: init,
	}
}

func MakeClientVideoPacket(clientVideo ClientVideo, payload []byte) Packet {
	return Packet{
		Type:   packetTypeClientVideo,
		JSON:   clientVideo,
		Binary: payload,
	}
}

func MakeServerVideoPacket(serverVideo ServerVideo, payload []byte) Packet {
	return Packet{
		Type:   packetTypeServerVideo,
		JSON:   serverVideo,
		Binary: payload,
	}
}
