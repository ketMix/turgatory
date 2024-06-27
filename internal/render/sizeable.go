package render

type Sizeable struct {
	width, height float64
}

func (s *Sizeable) SetSize(width, height float64) {
	s.width = width
	s.height = height
}

func (s *Sizeable) Size() (float64, float64) {
	return s.width, s.height
}

func (s *Sizeable) Width() float64 {
	return s.width
}

func (s *Sizeable) Height() float64 {
	return s.height
}
