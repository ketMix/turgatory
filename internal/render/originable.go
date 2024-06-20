package render

type Originable struct {
	ox, oy float64
}

func (p *Originable) SetOrigin(x, y float64) {
	p.ox = x
	p.oy = y
}

func (p *Originable) Origin() (float64, float64) {
	return p.ox, p.oy
}
