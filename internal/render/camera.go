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
	//Rotateable

	x InterpNumber
	y InterpNumber

	Pitch      float64
	zoom       InterpNumber
	Mode       CameraMode
	textOffset InterpNumber
	rotate     InterpNumber
	lastLevel  int
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		x:     InterpNumber{Current: x, Target: x, Speed: 0.1},
		y:     InterpNumber{Current: y, Target: y, Speed: 0.1},
		Pitch: 1.0,
		zoom:  InterpNumber{Current: 3.0, Target: 3.0, Speed: 0.1},
	}
}

func (c *Camera) Story() int {
	return c.lastLevel
}

func (c *Camera) SetStory(level int) {
	if level < 0 {
		level = 0
	}
	switch c.Mode {
	case CameraModeTower:
		c.SetPosition(c.x.Target, -65+float64(level)*c.GetMultiplier()*c.zoom.Target)
	case CameraModeStack:
		c.SetPosition(c.x.Target, float64(level)*c.GetMultiplier()*c.zoom.Target)
	case CameraModeSuperZoom:
		c.SetPosition(c.x.Target, 61+float64(level)*c.GetMultiplier()*c.zoom.Target)
	}
	c.lastLevel = level
}

func (c *Camera) GetMultiplier() float64 {
	switch c.Mode {
	case CameraModeTower:
		return 18.5
	case CameraModeStack:
		return 9.32
	case CameraModeSuperZoom:
		return 3.5
	}
	return 0
}

func (c *Camera) SetMode(mode CameraMode) {
	c.Mode = mode
	switch mode {
	case CameraModeTower:
		c.SetZoom(1.5)
		c.SetTextOffset(75)
	case CameraModeStack:
		c.SetZoom(3)
		c.SetTextOffset(0)
	case CameraModeSuperZoom:
		c.SetZoom(8.0)
		c.SetTextOffset(-52)
	}
	c.SetStory(c.lastLevel)
}

func (c *Camera) Update() {
	c.x.Update()
	c.y.Update()
	c.zoom.Update()
	c.textOffset.Update()
	c.rotate.Update()
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
	dx := math.Abs(c.x.Current - x)
	dy := math.Abs(c.y.Current - y)
	targetTicks := 20.0

	// Calculate the number of ticks it will take to get to the target.
	dxTicks := dx / targetTicks
	dyTicks := dy / targetTicks

	c.x.Set(x, dxTicks)
	c.y.Set(y, dyTicks)
}

func (c *Camera) Position() (float64, float64) {
	return c.x.Current, c.y.Current
}

func (c *Camera) SetZoom(zoom float64) {
	dz := math.Abs(c.zoom.Current - zoom)
	targetTicks := 20.0

	dzTicks := dz / targetTicks

	c.zoom.Set(zoom, dzTicks)
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
	do := math.Abs(c.textOffset.Current - offset)
	targetTicks := 20.0

	doTicks := do / targetTicks

	c.textOffset.Set(offset, doTicks)
}

func (c *Camera) TextOffset() float64 {
	return c.textOffset.Current
}

func (c *Camera) Rotation() float64 {
	return c.rotate.Current
}

func (c *Camera) SetRotation(rotation float64) {
	dr := math.Abs(c.rotate.Current - rotation)
	targetTicks := 20.0

	drTicks := dr / targetTicks

	c.rotate.Set(rotation, drTicks)
}

func (c *Camera) SetRotationAt(rotation float64, speed float64) {
	dr := math.Abs(c.rotate.Current - rotation)
	targetTicks := speed

	drTicks := dr / targetTicks

	c.rotate.Set(rotation, drTicks)
}
