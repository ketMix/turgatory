package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type UIOptions struct {
	Scale  float64
	Width  int
	Height int
}

func (o UIOptions) CoordsToScreen(x, y float64) (float64, float64) {
	return x * o.Scale, y * o.Scale
}

func (o UIOptions) ScreenToCoords(x, y float64) (float64, float64) {
	return x / o.Scale, y / o.Scale
}

type UI struct {
	dudePanel DudePanel
	roomPanel RoomPanel
	options   *UIOptions
}

func NewUI() *UI {
	ui := &UI{}

	{
		panelSprite := Must(render.NewSprite("ui/panels"))
		ui.dudePanel = DudePanel{
			top:      Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
			topright: Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
			mid:      Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
			midright: Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
			bot:      Must(render.NewSubSprite(panelSprite, 0, 32, 16, 16)),
			botright: Must(render.NewSubSprite(panelSprite, 16, 32, 16, 16)),
		}
	}
	{
		panelSprite := Must(render.NewSprite("ui/botPanel"))
		ui.roomPanel = RoomPanel{
			topleft:  Must(render.NewSubSprite(panelSprite, 0, 0, 16, 32)),
			left:     Must(render.NewSubSprite(panelSprite, 0, 16, 16, 32)),
			topmid:   Must(render.NewSubSprite(panelSprite, 16, 0, 16, 32)),
			mid:      Must(render.NewSubSprite(panelSprite, 16, 16, 16, 32)),
			topright: Must(render.NewSubSprite(panelSprite, 32, 0, 16, 32)),
			right:    Must(render.NewSubSprite(panelSprite, 32, 16, 16, 32)),
		}
	}
	return ui
}

func (ui *UI) Layout(o *UIOptions) {
	ui.options = o
	ui.dudePanel.Layout(o)
	ui.roomPanel.Layout(o)
}

func (ui *UI) Update(o *UIOptions) {
	ui.dudePanel.Update(o)
	ui.roomPanel.Update(o)
}

func (ui *UI) Draw(o *render.Options) {
	ui.dudePanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.roomPanel.Draw(o)
}

type DudePanel struct {
	render.Originable
	render.Positionable
	drawered     bool
	width        float64
	height       float64
	top          *render.Sprite
	topright     *render.Sprite
	mid          *render.Sprite
	midright     *render.Sprite
	bot          *render.Sprite
	botright     *render.Sprite
	drawerInterp InterpNumber
	dudeProfiles []*DudeProfile
}

type DudeProfile struct {
	render.Positionable
	stack      *render.Stack
	dude       *Dude
	hovered    bool
	height     float64
	width      float64
	stackScale float64
}

func (dp *DudeProfile) Draw(o *render.Options) {
	x, y := dp.Position()
	// Save these top options for drawing dude profiles
	profileOptions := render.Options{
		Screen: o.Screen,
		Pitch:  2,
	}
	profileOptions.DrawImageOptions.GeoM.Concat(o.DrawImageOptions.GeoM)
	profileOptions.DrawImageOptions.GeoM.Scale(dp.stackScale, dp.stackScale)
	profileOptions.DrawImageOptions.GeoM.Translate(x, y)
	// Also shove 'em to the right a little.
	profileOptions.DrawImageOptions.GeoM.Translate(dp.width/2, 0)
	dp.stack.Draw(&profileOptions)

	if dp.hovered {
		op := &render.TextOptions{
			Screen: o.Screen,
			Font:   assets.DisplayFont,
			Color:  color.White,
		}
		op.GeoM.Translate(x+dp.width*2.5, y)
		render.DrawText(op, dp.dude.Name())
		op.Font = assets.BodyFont
		op.GeoM.Reset()
		op.GeoM.Translate(x+dp.width*2.5, y+assets.DisplayFont.LineHeight-assets.BodyFont.LineHeight/2)
		render.DrawText(op, fmt.Sprintf("Level %d %s", dp.dude.Level(), dp.dude.Profession()))
	}
}

func InBounds(x, y, width, height, mx, my float64) bool {
	if mx > x && mx < x+width && my > y && my < y+height {
		return true
	}
	return false
}

func (dp *DudePanel) Layout(o *UIOptions) {
	// eww
	dp.bot.Scale = o.Scale
	dp.botright.Scale = o.Scale
	dp.mid.Scale = o.Scale
	dp.midright.Scale = o.Scale
	dp.top.Scale = o.Scale
	dp.topright.Scale = o.Scale

	partWidth, partHeight := dp.top.Size()
	dp.width = partWidth * 2
	dp.height = float64(o.Height - o.Height/3)

	// Position at vertical center.
	dp.SetPosition(0, float64(o.Height/2)-dp.height/2)

	// Position dude faces
	dpx, dpy := dp.Position()
	dpy += partHeight / 2 // Pad the top a bit
	y := 0.0
	for _, p := range dp.dudeProfiles {
		p.SetPosition(dpx, dpy+y)
		p.stackScale = o.Scale + 1
		p.width = float64(p.stack.Width()) * p.stackScale
		p.height = float64(p.stack.Height()) * p.stackScale * 1.5

		y += p.height + 4
	}
}

func (dp *DudePanel) Update(o *UIOptions) {
	dp.drawerInterp.Update()

	dpx, dpy := dp.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())

	maxX := dpx + dp.width
	maxY := dpy + dp.height

	if mx > dpx && mx < maxX && my > dpy && my < maxY {
		if dp.drawered {
			dp.drawered = false
			dp.drawerInterp.Set(0, 3)
		}
	} else {
		if !dp.drawered {
			dp.drawered = true
			dp.drawerInterp.Set(-(dp.width - dp.width/4), 3)
		}
	}

	//if !dp.drawered {
	for _, p := range dp.dudeProfiles {
		px, py := p.Position()
		if InBounds(px, py, dp.width, p.height, mx, my) {
			fmt.Println("!!! hovered over my guy: ", p.dude.name)
			p.hovered = true
		} else {
			p.hovered = false
		}
	}
	//}
}

func (dp *DudePanel) Draw(o *render.Options) {
	pw, ph := dp.top.Size()
	o.DrawImageOptions.GeoM.Translate(dp.drawerInterp.Current, 0)
	o.DrawImageOptions.GeoM.Translate(dp.Position())
	// top
	dp.top.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	dp.topright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(0, ph)
	o.DrawImageOptions.GeoM.Translate(-pw, 0)

	// mid
	parts := int(math.Floor(dp.height/ph)) - 2
	for y := 0; y < parts; y++ {
		dp.mid.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
		dp.midright.Draw(o)
		o.DrawImageOptions.GeoM.Translate(-pw, 0)
		o.DrawImageOptions.GeoM.Translate(0, ph)
	}
	// bottom
	dp.bot.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	dp.botright.Draw(o)

	// Draw dudes, but offset also by the drawerInterp
	o.DrawImageOptions.GeoM.Reset()
	o.DrawImageOptions.GeoM.Translate(dp.drawerInterp.Current, 0)
	for _, p := range dp.dudeProfiles {
		p.Draw(o)
	}
}

func (dp *DudePanel) SyncDudes(dudes []*Dude) {
	for _, dude := range dudes {
		stack := render.CopyStack(dude.stack)
		stack.SetPosition(0, 0)
		stack.SetOriginToCenter()
		stack.SetRotation(math.Pi/2 - math.Pi/4)

		dp.dudeProfiles = append(dp.dudeProfiles, &DudeProfile{
			dude:   dude,
			stack:  stack,
			width:  float64(stack.Width()),
			height: float64(stack.Height()) * 2, // x2 for slice pitch of 1
		})
	}
}

type RoomPanel struct {
	render.Originable
	render.Positionable
	drawered     bool
	drawerInterp InterpNumber
	width        float64
	height       float64
	left         *render.Sprite
	topleft      *render.Sprite
	mid          *render.Sprite
	topmid       *render.Sprite
	right        *render.Sprite
	topright     *render.Sprite
}

func (rp *RoomPanel) Layout(o *UIOptions) {
	rp.left.Scale = o.Scale
	rp.topleft.Scale = o.Scale
	rp.mid.Scale = o.Scale
	rp.topmid.Scale = o.Scale
	rp.right.Scale = o.Scale
	rp.topright.Scale = o.Scale

	_, ph := rp.topleft.Size()

	rp.width = float64(o.Width - o.Width/3)
	rp.height = ph * 2
	rp.SetPosition(float64(o.Width/2)-float64(rp.width)/2, float64(o.Height)-96)
}

func (rp *RoomPanel) Update(o *UIOptions) {
	rp.drawerInterp.Update()

	rpx, rpy := rp.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())

	maxX := rpx + float64(rp.width)
	maxY := rpy + rp.height

	_, ph := rp.topleft.Size()

	if mx > rpx && mx < maxX && my > rpy && my < maxY {
		if rp.drawered {
			rp.drawered = false
			rp.drawerInterp.Set(0, 3.5)
		}
	} else {
		if !rp.drawered {
			rp.drawered = true
			rp.drawerInterp.Set(ph, 3.5)
		}
	}
}

func (rp *RoomPanel) Draw(o *render.Options) {
	pw, ph := rp.topleft.Size()
	o.DrawImageOptions.GeoM.Translate(0, rp.drawerInterp.Current)
	o.DrawImageOptions.GeoM.Translate(rp.Position())
	// topleft
	rp.topleft.Draw(o)
	o.DrawImageOptions.GeoM.Translate(0, ph)
	// left
	rp.left.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, -ph)
	// mid
	parts := int(math.Floor(float64(rp.width)/float64(pw))) - 2
	for i := 0; i < parts; i++ {
		rp.topmid.Draw(o)
		o.DrawImageOptions.GeoM.Translate(0, ph)
		rp.mid.Draw(o)
		o.DrawImageOptions.GeoM.Translate(0, -ph)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	// topright
	rp.topright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(0, ph)
	// right
	rp.right.Draw(o)

	o.DrawImageOptions.GeoM.Reset()
	o.DrawImageOptions.GeoM.Translate(0, rp.drawerInterp.Current)
	o.DrawImageOptions.GeoM.Translate(rp.Position())
	o.DrawImageOptions.GeoM.Translate(pw/2, 8)
	// Quick hacky test.
	rd := GetRoomDef(HealingShrine, Small)
	o.Screen.DrawImage(rd.image, &o.DrawImageOptions)
	o.DrawImageOptions.GeoM.Translate(float64(rd.image.Bounds().Dx())+8, 0)
	rd = GetRoomDef(Library, Medium)
	o.Screen.DrawImage(rd.image, &o.DrawImageOptions)
	o.DrawImageOptions.GeoM.Translate(float64(rd.image.Bounds().Dx())+8, 0)
	rd = GetRoomDef(Armory, Medium)
	o.Screen.DrawImage(rd.image, &o.DrawImageOptions)
	o.DrawImageOptions.GeoM.Translate(float64(rd.image.Bounds().Dx())+8, 0)
	rd = GetRoomDef(Treasure, Small)
	o.Screen.DrawImage(rd.image, &o.DrawImageOptions)
	o.DrawImageOptions.GeoM.Translate(float64(rd.image.Bounds().Dx())+8, 0)
	rd = GetRoomDef(Combat, Small)
	o.Screen.DrawImage(rd.image, &o.DrawImageOptions)
}
