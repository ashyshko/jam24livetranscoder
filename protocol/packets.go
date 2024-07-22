package protocol

func MakeInit(obj Init) Packet {
	return Packet{
		Type: packetTypeInit,
		JSON: obj,
	}
}
func MakeVideoPacket(obj VideoPacket, payload []byte) Packet {
	return Packet{
		Type:   packetTypeVideoPacket,
		JSON:   obj,
		Binary: payload,
	}
}
func MakeEof() Packet {
	return Packet{
		Type: packetTypeEof,
		JSON: map[string]interface{}{},
	}
}
func MakeVideoHeader(obj VideoHeader, payload []byte) Packet {
	return Packet{
		Type:   packetTypeVideoHeader,
		JSON:   obj,
		Binary: payload,
	}
}
func MakeOutputVideoPacket(obj OutputVideoPacket, payload []byte) Packet {
	return Packet{
		Type:   packetTypeOutputVideoPacket,
		JSON:   obj,
		Binary: payload,
	}
}
func MakeExpiringSoon(obj ExpiringSoon) Packet {
	return Packet{
		Type: packetTypeExpiringSoon,
		JSON: obj,
	}
}
