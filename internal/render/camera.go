package render

import (
	"math"
)

type CameraMode int

const (
	CameraModeTower CameraMode = iota
	CameraModeStack
	CameraModeSuperZoom
)

type Camera struct {
	Originable
	Rotateable

	x InterpNumber
	y InterpNumber

	Pitch      float64
	zoom       InterpNumber
	Mode       CameraMode
	textOffset InterpNumber
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		x:     InterpNumber{Current: x, Target: x, Speed: 0.1},
		y:     InterpNumber{Current: y, Target: y, Speed: 0.1},
		Pitch: 1.0,
		zoom:  InterpNumber{Current: 3.0, Target: 3.0, Speed: 0.1},
	}
}

func (c *Camera) SetMode(mode CameraMode) {
	c.Mode = mode
	switch mode {
	case CameraModeTower:
		c.SetPosition(0, -100)
		c.SetZoom(1.5)
		c.SetTextOffset(125)
	case CameraModeStack:
		c.SetPosition(0, 0)
		c.SetZoom(3)
		c.SetTextOffset(0)
	case CameraModeSuperZoom:
		c.SetPosition(0, 90)
		c.SetZoom(8.0)
		c.SetTextOffset(-80)
	}
}

func (c *Camera) Update() {
	c.x.Update()
	c.y.Update()
	c.zoom.Update()
	c.textOffset.Update()
}

func (c *Camera) Transform(options *Options) {
	cx, cy := c.Position()
	ox, oy := c.Origin()

	zoom := c.Zoom()

	options.DrawImageOptions.GeoM.Translate(cx, cy)
	//options.DrawImageOptions.GeoM.Rotate(c.Rotation()) // c'ya camera rotation
	options.DrawImageOptions.GeoM.Scale(zoom, zoom)
	options.DrawImageOptions.GeoM.Translate(ox, oy)

	options.TowerRotation = c.Rotation()

	options.Pitch = c.Pitch * zoom
}

func (c *Camera) ScreenToWorld(x, y float64) (float64, float64) {
	cx, cy := c.Position()
	ox, oy := c.Origin()
	zoom := c.Zoom()
	rads := c.Rotation()

	x = (x - cx - ox) / zoom
	y = (y - cy - oy) / zoom

	x, y = x*math.Cos(rads)+y*math.Sin(rads), y*math.Cos(rads)-x*math.Sin(rads)

	return x, y
}

func (c *Camera) WorldToScreen(x, y float64) (float64, float64) {
	cx, cy := c.Position()
	ox, oy := c.Origin()
	zoom := c.Zoom()
	rads := c.Rotation()

	cx = cx * zoom
	cy = cy * zoom

	x, y = x*math.Cos(-rads)+y*math.Sin(-rads), y*math.Cos(-rads)-x*math.Sin(-rads)

	x = x*zoom + cx + ox
	y = y*zoom + cy + oy

	return x, y
}

func (c *Camera) SetPosition(x, y float64) {
	c.x.Set(x, c.zoom.Current)
	c.y.Set(y, c.zoom.Current)
}

func (c *Camera) Position() (float64, float64) {
	return c.x.Current, c.y.Current
}

func (c *Camera) SetZoom(zoom float64) {
	c.zoom.Set(zoom, 0.1)
}

func (c *Camera) Zoom() float64 {
	return c.zoom.Current
}

func (c *Camera) ZoomIn() {
	c.zoom.Set(c.zoom.Target+c.zoom.Target*0.01, 1)
}

func (c *Camera) ZoomOut() {
	c.zoom.Set(c.zoom.Target-c.zoom.Target*0.01, 1)
}

func (c *Camera) SetTextOffset(offset float64) {
	c.textOffset.Set(offset, c.zoom.Target)
}

func (c *Camera) TextOffset() float64 {
	return c.textOffset.Current
}
