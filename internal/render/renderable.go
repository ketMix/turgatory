package render

type Renderable interface {
	Draw(*Options)
	Update()
	Position() (float64, float64)
	SetPosition(float64, float64)
	Rotation() float64
	SetRotation(float64)
	SetRotationDistance(float64)
}
