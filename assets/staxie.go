package assets

import (
	"bytes"
	"errors"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var stax = make(map[string]*Staxie)

func LoadStaxie(name string) (*Staxie, error) {
	if staxie, ok := stax[name]; ok {
		return staxie, nil
	}

	b, err := FS.ReadFile(name + ".png")
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	i, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	// Convert the image to an Ebiten image.
	eimg := ebiten.NewImageFromImage(i)

	// Read our staxie PNG data.
	staxie := &Staxie{}
	err = staxie.FromBytes(b)
	if err != nil {
		return nil, err
	}

	staxie.image = eimg

	for _, stack := range staxie.Stacks {
		staxie.StackNames = append(staxie.StackNames, stack.Name)
	}

	staxie.acquireSliceImages()

	stax[name] = staxie

	return staxie, nil
}

// Staxie is the structure extracted from a Staxie PNG file.
type Staxie struct {
	Stacks      map[string]*StaxieStack
	StackNames  []string
	FrameWidth  int
	FrameHeight int
	image       *ebiten.Image
}

// FromBytes reads the given PNG bytes into a staxie structure, providing it has a stAx section.
func (s *Staxie) FromBytes(data []byte) error {
	s.Stacks = make(map[string]*StaxieStack)
	s.FrameWidth = 0
	s.FrameHeight = 0

	offset := 0

	readUint32 := func() uint32 {
		if offset+4 > len(data) {
			panic("Out of bounds")
		}
		value := uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		return value
	}
	readUint16 := func() uint16 {
		if offset+2 > len(data) {
			panic("Out of bounds")
		}
		value := uint16(data[offset])<<8 | uint16(data[offset+1])
		offset += 2
		return value
	}
	readString := func() string {
		if offset+1 > len(data) {
			panic("Out of bounds")
		}
		length := data[offset]
		offset++
		if offset+int(length) > len(data) {
			panic("Out of bounds")
		}
		value := string(data[offset : offset+int(length)])
		offset += int(length)
		return value
	}
	readSection := func() string {
		if offset+4 > len(data) {
			panic("Out of bounds")
		}
		section := string(data[offset : offset+4])
		offset += 4
		return section
	}

	offset += 8 // Skip PNG header

	for offset < len(data) {
		chunkSize := readUint32()
		section := readSection()
		switch section {
		case "stAx":
			if offset+1 > len(data) {
				panic("Out of bounds")
			}
			version := data[offset]
			offset++
			if version != 0 {
				panic("Unsupported stAx version")
			}
			frameWidth := readUint16()
			frameHeight := readUint16()
			stackCount := readUint16()

			s.FrameWidth = int(frameWidth)
			s.FrameHeight = int(frameHeight)

			y := 0
			for i := 0; i < int(stackCount); i++ {
				stack := StaxieStack{
					Animations: make(map[string]StaxieAnimation),
				}
				name := readString()
				sliceCount := readUint16()
				animationCount := readUint16()

				stack.SliceCount = int(sliceCount)
				stack.Name = name

				for j := 0; j < int(animationCount); j++ {
					animation := StaxieAnimation{}
					animationName := readString()
					frameTime := readUint32()
					frameCount := readUint16()

					animation.Frametime = frameTime
					animation.Name = animationName

					for k := 0; k < int(frameCount); k++ {
						frame := StaxieFrame{}
						for l := 0; l < int(sliceCount); l++ {
							slice := StaxieSlice{
								X: l * int(frameWidth),
								Y: y,
							}
							if offset+1 > len(data) {
								panic("Out of bounds")
							}
							slice.Shading = data[offset]
							offset++
							frame.Slices = append(frame.Slices, slice)
						}
						if sliceCount > 0 {
							y += int(frameHeight)
						}
						frame.Index = k
						animation.Frames = append(animation.Frames, frame)
					}
					stack.Animations[animationName] = animation
				}
				s.Stacks[name] = &stack
			}
		default: // Skip non-stAx sections
			offset += int(chunkSize)
		}
		offset += 4 // Skip CRC32
	}

	return nil
}

// acquireSliceImages acquires the subimages for each slice from the Stack's Image.
func (s *Staxie) acquireSliceImages() {
	for _, stack := range s.Stacks {
		for _, animation := range stack.Animations {
			for _, frame := range animation.Frames {
				for i, slice := range frame.Slices {
					// FIXME: This would be more efficient, _however_ it renders 1px wide vertical lines at random positions! (Well, this isn't random, presumably it's because of some texture wrapping with a rounding error) -- TODO: bring this up as an issue in Ebitengine!
					//img := s.image.SubImage(image.Rect(slice.X, slice.Y, slice.X+s.FrameWidth, slice.Y+s.FrameHeight)).(*ebiten.Image)

					img := ebiten.NewImage(s.FrameWidth, s.FrameHeight)
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Translate(float64(-slice.X), float64(-slice.Y))
					img.DrawImage(s.image, opts)

					frame.Slices[i].Image = img
				}
			}
		}
	}
}

func (s *Staxie) GetStack(name string) (*StaxieStack, bool) {
	stack, ok := s.Stacks[name]
	return stack, ok
}

type StaxieStack struct {
	Name       string // For convenience
	SliceCount int
	Animations map[string]StaxieAnimation
}

func (s *StaxieStack) GetAnimation(name string) (StaxieAnimation, bool) {
	animation, ok := s.Animations[name]
	return animation, ok
}

type StaxieAnimation struct {
	Name      string // For convenience
	Frametime uint32
	Frames    []StaxieFrame
}

func (s *StaxieAnimation) GetFrame(index int) (*StaxieFrame, bool) {
	if index < 0 || index >= len(s.Frames) {
		return nil, false
	}
	return &s.Frames[index], true
}

type StaxieFrame struct {
	Index  int // ehh... why not
	Slices []StaxieSlice
}

func (s *StaxieFrame) GetSlice(index int) (*StaxieSlice, bool) {
	if index < 0 || index >= len(s.Slices) {
		return nil, false
	}
	return &s.Slices[index], true
}

type StaxieSlice struct {
	Shading uint8
	X       int
	Y       int
	// Might as well store the ebitengine subimage here for efficiency's sake.
	Image *ebiten.Image
}

var (
	ErrStackNotFound     = errors.New("stack not found")
	ErrAnimationNotFound = errors.New("animation not found")
	ErrFrameNotFound     = errors.New("frame not found")
	ErrSliceNotFound     = errors.New("slice not found")
)
