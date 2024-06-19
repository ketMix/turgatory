package render

type Rotateable struct {
	rotation float64
}

func (r *Rotateable) SetRotation(rotation float64) {
	r.rotation = rotation
}

func (r *Rotateable) Rotation() float64 {
	return r.rotation
}
