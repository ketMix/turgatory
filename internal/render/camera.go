package render

import "math"

type Camera struct {
	Originable
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
	ox, oy := c.Origin()

	options.DrawImageOptions.GeoM.Translate(cx, cy)
	//options.DrawImageOptions.GeoM.Rotate(c.Rotation()) // c'ya camera rotation
	options.DrawImageOptions.GeoM.Scale(c.Zoom, c.Zoom)
	options.DrawImageOptions.GeoM.Translate(ox, oy)

	options.TowerRotation = c.Rotation()

	options.Pitch = c.Pitch * c.Zoom
}

func (c *Camera) ScreenToWorld(x, y float64) (float64, float64) {
	cx, cy := c.Position()
	ox, oy := c.Origin()
	rads := c.Rotation()

	x = (x - cx - ox) / c.Zoom
	y = (y - cy - oy) / c.Zoom

	x, y = x*math.Cos(rads)+y*math.Sin(rads), y*math.Cos(rads)-x*math.Sin(rads)

	return x, y
}

func (c *Camera) WorldToScreen(x, y float64) (float64, float64) {
	cx, cy := c.Position()
	ox, oy := c.Origin()
	rads := c.Rotation()

	cx = cx * c.Zoom
	cy = cy * c.Zoom

	x, y = x*math.Cos(-rads)+y*math.Sin(-rads), y*math.Cos(-rads)-x*math.Sin(-rads)

	x = x*c.Zoom + cx + ox
	y = y*c.Zoom + cy + oy

	return x, y
}
