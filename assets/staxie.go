package assets

// Staxie is the structure extracted from a Staxie PNG file.
type Staxie struct {
	Stacks      map[string]StaxieStack
	FrameWidth  int
	FrameHeight int
}

// FromBytes reads the given PNG bytes into a staxie structure, providing it has a stAx section.
func (s *Staxie) FromBytes(data []byte) error {
	s.Stacks = make(map[string]StaxieStack)
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

			for i := 0; i < int(stackCount); i++ {
				stack := StaxieStack{
					Animations: make(map[string]StaxieAnimation),
				}
				name := readString()
				sliceCount := readUint16()
				animationCount := readUint16()

				stack.SliceCount = int(sliceCount)

				for j := 0; j < int(animationCount); j++ {
					animation := StaxieAnimation{}
					animationName := readString()
					frameTime := readUint32()
					frameCount := readUint16()

					animation.Frametime = frameTime

					for k := 0; k < int(frameCount); k++ {
						frame := StaxieFrame{}
						for l := 0; l < int(sliceCount); l++ {
							slice := StaxieSlice{}
							if offset+1 > len(data) {
								panic("Out of bounds")
							}
							slice.Shading = data[offset]
							offset++
							frame.Slices = append(frame.Slices, slice)
						}
						animation.Frames = append(animation.Frames, frame)
					}
					stack.Animations[animationName] = animation
				}
				s.Stacks[name] = stack
			}
		default: // Skip non-stAx sections
			offset += int(chunkSize)
		}
		offset += 4 // Skip CRC32
	}

	return nil
}

type StaxieStack struct {
	SliceCount int
	Animations map[string]StaxieAnimation
}

type StaxieAnimation struct {
	Frametime uint32
	Frames    []StaxieFrame
}

type StaxieFrame struct {
	Slices []StaxieSlice
}

type StaxieSlice struct {
	Shading uint8
	X       int
	Y       int
}
