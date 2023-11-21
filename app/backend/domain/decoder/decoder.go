package decoder

import (
	"bytes"
	"encoding/binary"
	"fractalmri/gofract/pipes"
)

func DecodeMFS(data []byte) (pipes.Image64, error) {
	reader := bytes.NewReader(data)
	var width, height int32
	if err := binary.Read(reader, binary.NativeEndian, &width); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.NativeEndian, &height); err != nil {
		return nil, err
	}
	img := pipes.NewImage64(int(width), int(height))
	var pixel float64
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			if err := binary.Read(reader, binary.NativeEndian, &pixel); err != nil {
				return nil, err
			}
			img.Set(i, j, pixel)
		}
	}
	return img, nil
}
