package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/internal/render"
)

type UIOptions struct {
	Scale  float64
	Width  int
	Height int
}

type UI struct {
	dudePanel DudePanel
	options   *UIOptions
}

func NewUI() *UI {
	ui := &UI{}

	panelSprite := Must(render.NewSprite("ui/panels"))

	ui.dudePanel = DudePanel{
		top:      Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
		topright: Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
		mid:      Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
		midright: Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
		bot:      Must(render.NewSubSprite(panelSprite, 0, 32, 16, 16)),
		botright: Must(render.NewSubSprite(panelSprite, 16, 32, 16, 16)),
	}
	return ui
}

func (ui *UI) Layout(o *UIOptions) {
	ui.options = o
	ui.dudePanel.Layout(o)
}

func (ui *UI) Update(o *UIOptions) {
	ui.dudePanel.Update(o)
}

func (ui *UI) Draw(o *render.Options) {
	o.DrawImageOptions.GeoM.Scale(ui.options.Scale, ui.options.Scale)
	ui.dudePanel.Draw(o)
}

type DudePanel struct {
	render.Originable
	render.Positionable
	drawered bool
	height   int
	top      *render.Sprite
	topright *render.Sprite
	mid      *render.Sprite
	midright *render.Sprite
	bot      *render.Sprite
	botright *render.Sprite
}

func (dp *DudePanel) Layout(o *UIOptions) {
	dp.height = o.Height - o.Height/3
	// Position at vertical center.
	dp.SetPosition(0, float64(o.Height/2)-float64(dp.height)/2)
}

func (dp *DudePanel) Update(o *UIOptions) {
	dpx, dpy := dp.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())

	maxX := (dpx + 32) * o.Scale
	maxY := (dpy + float64(dp.height)) * o.Scale

	if mx > dpx && mx < maxX && my > dpy && my < maxY {
		dp.drawered = false
	} else {
		dp.drawered = true
	}

	if !dp.drawered {
		// TODO: CRAZY DUDE LIST!!!!
	}
}

func (dp *DudePanel) Draw(o *render.Options) {
	if dp.drawered {
		o.DrawImageOptions.GeoM.Translate(-48, 0)
	}
	y := 0
	o.DrawImageOptions.GeoM.Translate(dp.Position())
	// top
	dp.top.Draw(o)
	o.DrawImageOptions.GeoM.Translate(16, 0)
	dp.topright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(0, 16)
	y += 16
	o.DrawImageOptions.GeoM.Translate(-16, 0)
	// mid
	for ; y < dp.height-16; y += 16 {
		dp.mid.Draw(o)
		o.DrawImageOptions.GeoM.Translate(16, 0)
		dp.midright.Draw(o)
		o.DrawImageOptions.GeoM.Translate(-16, 0)
		o.DrawImageOptions.GeoM.Translate(0, 16)
	}
	// bottom
	dp.bot.Draw(o)
	o.DrawImageOptions.GeoM.Translate(16, 0)
	dp.botright.Draw(o)
}
