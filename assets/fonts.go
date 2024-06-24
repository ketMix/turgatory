package assets

import "github.com/hajimehoshi/ebiten/v2/text/v2"

type Font struct {
	Face       *text.GoTextFace
	LineHeight float64
}

var DisplayFont Font
var BodyFont Font

func init() {
	{
		f, err := FS.Open("fonts/antiquity-print.ttf")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		s, err := text.NewGoTextFaceSource(f)
		if err != nil {
			panic(err)
		}

		DisplayFont.Face = &text.GoTextFace{
			Source: s,
			Size:   13,
		}
		DisplayFont.LineHeight = 13 * 1.5
	}

	{
		f, err := FS.Open("fonts/BitPotionExt.ttf")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		s, err := text.NewGoTextFaceSource(f)
		if err != nil {
			panic(err)
		}

		BodyFont.Face = &text.GoTextFace{
			Source: s,
			Size:   16,
		}
		BodyFont.LineHeight = 9
	}
}
