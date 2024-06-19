package render

type Positionable struct {
	x, y float64
}

func (p *Positionable) SetPosition(x, y float64) {
	p.x = x
	p.y = y
}

func (p *Positionable) Position() (float64, float64) {
	return p.x, p.y
}
