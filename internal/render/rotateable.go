package render

import "math"

type Rotateable struct {
	rotation         float64
	rotationDistance float64
}

func (r *Rotateable) SetRotation(rotation float64) {
	r.rotation = math.Mod(rotation, 2*math.Pi)
}

func (r *Rotateable) Rotation() float64 {
	return r.rotation
}

func (r *Rotateable) SetRotationDistance(rotationDistance float64) {
	r.rotationDistance = rotationDistance
}

func (r *Rotateable) RotationDistance() float64 {
	return r.rotationDistance
}
