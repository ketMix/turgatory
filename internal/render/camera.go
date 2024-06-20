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
		Zoom:         3.0,
	}
}

func (c *Camera) Transform(options *Options) {
	cx, cy := c.Position()

	// camera origin do be center of the screen
	ox, oy := float64(options.Screen.Bounds().Dx())/2, float64(options.Screen.Bounds().Dy())/2

	cx += ox
	cy += oy

	options.DrawImageOptions.GeoM.Translate(-cx, -cy)

	options.DrawImageOptions.GeoM.Rotate(c.Rotation())

	options.DrawImageOptions.GeoM.Scale(c.Zoom, c.Zoom)

	options.DrawImageOptions.GeoM.Translate(ox, oy)

	options.Pitch = c.Pitch * c.Zoom
}
