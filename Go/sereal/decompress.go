package sereal

import "fmt"

func DecmmpressDocument(b, dst []byte) (r []byte, err error) {
	header, err := readHeader(b)

	if err != nil {
		return nil, err
	}

	bodyStart := headerSize + header.suffixSize

	if bodyStart > len(b) || bodyStart < 0 {
		return nil, ErrCorrupt{errBadOffset}
	}

	switch header.version {
	case 1:
		break
	case 2:
		break
	case 3:
		break
	case 4:
		break
	default:
		return nil, fmt.Errorf("document version '%d' not yet supported", header.version)
	}

	var decomp decompressor

	switch header.doctype {
	case serealRaw:
		// nothing

	case serealSnappy:
		if header.version != 1 {
			return nil, ErrBadSnappy
		}
		decomp = SnappyCompressor{Incremental: false}

	case serealSnappyIncremental:
		decomp = SnappyCompressor{Incremental: true}

	case serealZlib:
		if header.version < 3 {
			return nil, ErrBadZlibV3
		}
		decomp = ZlibCompressor{}

	case serealZstd:
		if header.version < 4 {
			return nil, ErrBadZstdV4
		}
		decomp = ZstdCompressor{}

	default:
		return nil, fmt.Errorf("document type '%d' not yet supported", header.doctype)
	}

	if decomp != nil {
		decompBody, err := decomp.decompress(b[bodyStart:])
		if err != nil {
			return nil, err
		}

		dst = append(dst[:0], b[:headerSize+header.suffixSize]...)
		dst = append(dst, decompBody...)
	} else {
		dst = append(dst[:0], b...)
	}

	// set type to 0
	dst[4] = dst[4] & 0x0f

	return dst, nil
}
