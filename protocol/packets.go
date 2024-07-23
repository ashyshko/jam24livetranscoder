package protocol

func MakeInit(obj Init) Packet {
	return Packet{
		Type: packetTypeInit,
		JSON: obj,
	}
}

func MakeVideoHeader(payload []byte) Packet {
	return Packet{
		Type:   packetTypeOutputVideoHeader,
		JSON:   map[string]interface{}{},
		Binary: payload,
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

func MakeOutputVideoHeader(obj OutputVideoHeader, payload []byte) Packet {
	return Packet{
		Type:   packetTypeOutputVideoHeader,
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
