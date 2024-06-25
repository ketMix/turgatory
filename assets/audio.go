package assets

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

func LoadSound(soundType string, name string) (*vorbis.Stream, error) {
	data, err := FS.ReadFile(fmt.Sprintf("audio/%s/%s.ogg", soundType, name))
	if err != nil {
		return nil, err
	}

	stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return stream, nil
}
