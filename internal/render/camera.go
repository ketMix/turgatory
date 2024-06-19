package render

type Camera struct {
	Positionable
	Rotateable
	Pitch float64
	Zoom  float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		Positionable: Positionable{x: x, y: y},
		Pitch:        1.0,
		Zoom:         1.0,
	}
}

func (c *Camera) Transform(options *Options) {
	cx, cy := c.Position()
	options.DrawImageOptions.GeoM.Translate(-cx, -cy)

	options.DrawImageOptions.GeoM.Rotate(c.Rotation())

	options.DrawImageOptions.GeoM.Scale(c.Zoom, c.Zoom)
}
