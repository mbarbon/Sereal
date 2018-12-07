package sereal

func LooksLikeSereal(b []byte) int {
	if len(b) < 7 {
		return -1
	}

	header, err := readHeader(b)
	if err != nil {
		return -1
	}

	return int(header.version)
}
