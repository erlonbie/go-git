package plumbing

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type PackObject struct {
	Type    string
	Size    int
	Content []byte
}

func ParsePackfile(data []byte) ([]PackObject, error) {
	reader := bytes.NewReader(data)

	header := make([]byte, 12)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("error reading packfile header: %w", err)
	}

	if string(header[0:4]) != "PACK" {
		return nil, fmt.Errorf("invalid packfile signature: %s", string(header[0:4]))
	}

	numObjects := binary.BigEndian.Uint32(header[8:12])
	fmt.Printf("Packfile contains %d objects\n", numObjects)

	var objects []PackObject

	for i := uint32(0); i < numObjects; i++ {
		obj, err := parseObject(reader)
		if err != nil {
			return nil, fmt.Errorf("error parsing object %d: %w", i, err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

func parseObject(reader *bytes.Reader) (PackObject, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return PackObject{}, err
	}

	objTypeInt := (b >> 4) & 7
	
	size := uint64(b & 15)
	shift := uint(4)

	for b&0x80 != 0 {
		b, err = reader.ReadByte()
		if err != nil {
			return PackObject{}, err
		}
		size |= uint64(b&0x7f) << shift
		shift += 7
	}

	var objType string
	switch objTypeInt {
	case 1:
		objType = "commit"
	case 2:
		objType = "tree"
	case 3:
		objType = "blob"
	case 4:
		objType = "tag"
	case 6:
		objType = "ofs_delta"
	case 7:
		objType = "ref_delta"
	default:
		return PackObject{}, fmt.Errorf("unknown object type: %d", objTypeInt)
	}

	if objTypeInt >= 1 && objTypeInt <= 4 {
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return PackObject{}, fmt.Errorf("error creating zlib reader: %w", err)
		}
		defer zlibReader.Close()

		content, err := io.ReadAll(zlibReader)
		if err != nil {
			return PackObject{}, fmt.Errorf("error reading zlib stream: %w", err)
		}

		return PackObject{
			Type:    objType,
			Size:    int(size),
			Content: content,
		}, nil
	}

	return PackObject{}, fmt.Errorf("delta objects (%s) parsing is not fully implemented yet", objType)
}
