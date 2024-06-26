package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type UIOptions struct {
	Scale    float64
	Width    int
	Height   int
	Messages []string
}

func (o UIOptions) CoordsToScreen(x, y float64) (float64, float64) {
	return x * o.Scale, y * o.Scale
}

func (o UIOptions) ScreenToCoords(x, y float64) (float64, float64) {
	return x / o.Scale, y / o.Scale
}

type UI struct {
	dudePanel    DudePanel
	roomPanel    RoomPanel
	speedPanel   SpeedPanel
	messagePanel MessagePanel
	options      *UIOptions
}

func NewUI() *UI {
	ui := &UI{}

	{
		panelSprite := Must(render.NewSprite("ui/panels"))
		ui.dudePanel = DudePanel{
			top:      Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
			topright: Must(render.NewSubSprite(panelSprite, 32, 0, 16, 16)),
			mid:      Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
			midright: Must(render.NewSubSprite(panelSprite, 32, 16, 16, 16)),
			bot:      Must(render.NewSubSprite(panelSprite, 16, 32, 16, 16)),
			botright: Must(render.NewSubSprite(panelSprite, 32, 32, 16, 16)),
		}
		ui.dudePanel.dudeDetails = &DudeDetails{
			top:      ui.dudePanel.top,
			topright: ui.dudePanel.topright,
			topleft:  Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
			mid:      ui.dudePanel.mid,
			midright: ui.dudePanel.midright,
			midleft:  Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
			bot:      ui.dudePanel.bot,
			botright: ui.dudePanel.botright,
			botleft:  Must(render.NewSubSprite(panelSprite, 0, 32, 16, 16)),
		}
	}
	{
		panelSprite := Must(render.NewSprite("ui/panels"))
		ui.messagePanel = MessagePanel{
			maxLines: 50,
			top:      Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
			topleft:  Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
			topright: Must(render.NewSubSprite(panelSprite, 32, 0, 16, 16)),
			mid:      Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
			midleft:  Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
			midright: Must(render.NewSubSprite(panelSprite, 32, 16, 16, 16)),
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
	{
		ui.speedPanel = SpeedPanel{
			musicButton:  NewButton("music", "music on"),
			soundButton:  NewButton("sound", "sound on"),
			pauseButton:  NewButton("play", "playing"),
			speedButton:  NewButton("fast", "fast"),
			cameraButton: NewButton("story", "camera: story"),
		}
	}
	return ui
}

func (ui *UI) Layout(o *UIOptions) {
	ui.options = o
	ui.dudePanel.Layout(o)
	ui.roomPanel.Layout(o)
	ui.speedPanel.Layout(o)
	ui.messagePanel.Layout(o)
}

func (ui *UI) Update(o *UIOptions) {
	ui.dudePanel.Update(o)
	ui.roomPanel.Update(o)
	ui.speedPanel.Update(o)
	ui.messagePanel.Update(o)
}

func (ui *UI) Draw(o *render.Options) {
	ui.dudePanel.Draw(o)
	// o.DrawImageOptions.GeoM.Reset()
	// ui.roomPanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.speedPanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.messagePanel.Draw(o)
}

type DudeDetails struct {
	render.Positionable
	dude             *Dude
	width            float64
	height           float64
	top              *render.Sprite
	topleft          *render.Sprite
	topright         *render.Sprite
	mid              *render.Sprite
	midleft          *render.Sprite
	midright         *render.Sprite
	bot              *render.Sprite
	botleft          *render.Sprite
	botright         *render.Sprite
	equipmentDetails []*EquipmentDetails
}

type EquipmentDetails struct {
	render.Positionable
	equipmentType EquipmentType
	equipment     *Equipment
	height        float64
	width         float64
	hovered       bool
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
	drawerInterp render.InterpNumber
	dudeProfiles []*DudeProfile
	dudeDetails  *DudeDetails
	onDudeClick  func(*Dude)
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
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		op.Color = color.RGBA{200, 50, 50, 255}
		render.DrawText(op, fmt.Sprintf("HP: %d/%d", dp.dude.stats.currentHp, dp.dude.stats.totalHp))
		op.Color = color.RGBA{200, 200, 200, 255}
		op.GeoM.Translate(0, assets.BodyFont.LineHeight*2)
		render.DrawText(op, fmt.Sprintf("%s strength", PaddedIntString(dp.dude.stats.strength, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s agility", PaddedIntString(dp.dude.stats.agility, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s defense", PaddedIntString(dp.dude.stats.defense, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s wisdom", PaddedIntString(dp.dude.stats.wisdom, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s cowardice", PaddedIntString(dp.dude.stats.cowardice, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s luck", PaddedIntString(dp.dude.stats.luck, 4)))
	}
}

func (dd *DudeDetails) SetDude(dude *Dude) {
	dd.dude = dude
	if dude == nil {
		return
	}

	// Set equipment details
	dd.equipmentDetails = []*EquipmentDetails{}
	// Armor
	dd.equipmentDetails = append(dd.equipmentDetails, &EquipmentDetails{
		equipmentType: EquipmentTypeArmor,
		equipment:     dude.equipped[EquipmentTypeArmor],
	})

	// Weapon
	dd.equipmentDetails = append(dd.equipmentDetails, &EquipmentDetails{
		equipmentType: EquipmentTypeWeapon,
		equipment:     dude.equipped[EquipmentTypeWeapon],
	})

	// Accessory
	dd.equipmentDetails = append(dd.equipmentDetails, &EquipmentDetails{
		equipmentType: EquipmentTypeAccessory,
		equipment:     dude.equipped[EquipmentTypeAccessory],
	})
}

func (dd *DudeDetails) Layout(o *UIOptions, dp *DudePanel) {
	// eww
	scale := o.Scale * 1.1
	dd.bot.Scale = scale
	dd.botleft.Scale = scale
	dd.botright.Scale = scale
	dd.mid.Scale = scale
	dd.midleft.Scale = scale
	dd.midright.Scale = scale
	dd.top.Scale = scale
	dd.topleft.Scale = scale
	dd.topright.Scale = scale

	partWidth, _ := dd.top.Size()
	dd.width = partWidth * 9
	dd.height = dd.width * 0.5

	// Position at vertical center + width of dude panel
	x, y := dp.width+dd.width/3, float64(o.Height/2)-dd.height/2
	dd.SetPosition(x, y)

	// Layout equipment details
	for i, ed := range dd.equipmentDetails {
		ed.Layout(dd, i)
	}
}

func (dd *DudeDetails) Draw(o *render.Options) {
	if dd.dude == nil {
		return
	}
	x, y := dd.Position()
	pw, ph := dd.top.Size()
	o.DrawImageOptions.GeoM.Translate(x, y)

	// Calculate parts
	verticalParts := int(math.Floor(dd.height/ph)) - 2
	horizontalParts := int(math.Floor(dd.width/pw)) - 2

	// top
	dd.topleft.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	for x := 0; x < horizontalParts; x++ {
		dd.top.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	dd.topright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(-dd.width+pw, ph)

	// mid
	for y := 0; y < verticalParts; y++ {
		dd.midleft.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
		for x := 0; x < horizontalParts; x++ {
			dd.mid.Draw(o)
			o.DrawImageOptions.GeoM.Translate(pw, 0)
		}
		dd.midright.Draw(o)
		o.DrawImageOptions.GeoM.Translate(-dd.width+pw, ph)
	}

	// bottom
	dd.botleft.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	for x := 0; x < horizontalParts; x++ {
		dd.bot.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	dd.botright.Draw(o)

	// Details
	op := &render.TextOptions{
		Screen: o.Screen,
		Font:   assets.DisplayFont,
		Color:  color.White,
	}
	op.GeoM.Reset()
	op.GeoM.Translate(x+15, y+10)
	render.DrawText(op, dd.dude.Name())
	op.Font = assets.BodyFont
	op.GeoM.Reset()
	op.GeoM.Translate(x+15, y+10+assets.DisplayFont.LineHeight-assets.BodyFont.LineHeight/2)
	render.DrawText(op, fmt.Sprintf("Level %d %s", dd.dude.Level(), dd.dude.Profession()))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	op.Color = color.RGBA{200, 50, 50, 255}
	render.DrawText(op, fmt.Sprintf("HP: %d/%d", dd.dude.stats.currentHp, dd.dude.stats.totalHp))
	op.Color = color.RGBA{200, 200, 200, 255}
	op.GeoM.Translate(0, assets.BodyFont.LineHeight*2)
	render.DrawText(op, fmt.Sprintf("%s strength", PaddedIntString(dd.dude.stats.strength, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s agility", PaddedIntString(dd.dude.stats.agility, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s defense", PaddedIntString(dd.dude.stats.defense, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s wisdom", PaddedIntString(dd.dude.stats.wisdom, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s cowardice", PaddedIntString(dd.dude.stats.cowardice, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s luck", PaddedIntString(dd.dude.stats.luck, 4)))

	// Equipment
	op.GeoM.Reset()
	op.GeoM.Translate(x+dd.width*0.40, y+10)
	op.Font = assets.DisplayFont
	op.Color = color.White
	render.DrawText(op, "Equipment")
	op.GeoM.Translate(0, assets.DisplayFont.LineHeight)

	hoveredEquipment := (*Equipment)(nil)
	for _, ed := range dd.equipmentDetails {
		ed.Draw(o)
		if ed.hovered {
			hoveredEquipment = ed.equipment
		}
	}

	if hoveredEquipment != nil {
		drawEquipmentDescription(o, dd, hoveredEquipment)
	}
}

func drawEquipmentDescription(o *render.Options, dd *DudeDetails, hoveredEquipment *Equipment) {
	x, y := dd.Position()

	// hasPerk := hoveredEquipment.perk != nil

	// Draw equipment description below
	o.DrawImageOptions.GeoM.Reset()
	o.DrawImageOptions.GeoM.Translate(x, y+dd.height-20)

	// top
	dd.topleft.Draw(o)
	pw, ph := dd.top.Size()
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	horizontalParts := int(math.Floor(dd.width/pw)) - 2
	for x := 0; x < horizontalParts; x++ {
		dd.top.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	dd.topright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(-dd.width+pw, ph)

	// mid
	dd.midleft.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	for x := 0; x < horizontalParts; x++ {
		dd.mid.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	dd.midright.Draw(o)
	o.DrawImageOptions.GeoM.Translate(-dd.width+pw, ph)

	// bottom
	dd.botleft.Draw(o)
	o.DrawImageOptions.GeoM.Translate(pw, 0)
	for x := 0; x < horizontalParts; x++ {
		dd.bot.Draw(o)
		o.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	dd.botright.Draw(o)

	// Details
	op := &render.TextOptions{
		Screen: o.Screen,
		Font:   assets.DisplayFont,
		Color:  color.White,
	}

	// LEFT SIDE
	op.GeoM.Reset()
	op.GeoM.Translate(x+15, y+dd.height-10)
	op.Font = assets.DisplayFont
	op.Color = hoveredEquipment.quality.TextColor()
	render.DrawText(op, hoveredEquipment.Name())
	op.GeoM.Translate(0, assets.DisplayFont.LineHeight+1)

	op.Font = assets.BodyFont
	if hoveredEquipment.perk != nil {
		op.GeoM.Translate(10, 0)
		op.Color = hoveredEquipment.perk.Quality().TextColor()
		render.DrawText(op, "of "+hoveredEquipment.perk.Name())
		op.GeoM.Translate(0, assets.DisplayFont.LineHeight+1)
		op.GeoM.Translate(-10, 0)
	} else {
		op.GeoM.Translate(0, assets.DisplayFont.LineHeight+1)
	}

	// Draw equipment description
	op.Color = color.RGBA{200, 200, 200, 255}
	render.DrawText(op, hoveredEquipment.Description())
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)

	// Draw equipment perks
	if hoveredEquipment.perk != nil {
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		op.Color = hoveredEquipment.perk.Quality().TextColor()
		render.DrawText(op, hoveredEquipment.perk.Description())
	}

	// RIGHT SIDE
	// Draw equipment type
	op.GeoM.Reset()
	op.GeoM.Translate(x+dd.width-50, y+dd.height-10)
	op.Color = color.RGBA{200, 200, 200, 255}
	render.DrawText(op, hoveredEquipment.Type().String())
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	op.GeoM.Translate(-50, 0)

	// Draw Stats
	op.Color = color.RGBA{200, 200, 200, 255}
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s strength", PaddedIntString(hoveredEquipment.stats.strength, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s agility", PaddedIntString(hoveredEquipment.stats.agility, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s defense", PaddedIntString(hoveredEquipment.stats.defense, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s wisdom", PaddedIntString(hoveredEquipment.stats.wisdom, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s cowardice", PaddedIntString(hoveredEquipment.stats.cowardice, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s luck", PaddedIntString(hoveredEquipment.stats.luck, 4)))

}

func (ed *EquipmentDetails) Layout(dd *DudeDetails, i int) {
	ed.height = (assets.BodyFont.LineHeight + 2) * 3
	ed.width = dd.width * 0.45

	ddX, ddY := dd.Position()
	yOffset := ed.height * float64(i+1)
	edX, edY := ddX+ed.width, ddY+yOffset
	ed.SetPosition(edX, edY)
}

func (ed *EquipmentDetails) Update(o *UIOptions) {
	ed.hovered = false
	x, y := ed.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())
	if InBounds(x, y, ed.width, ed.height, mx, my) {
		ed.hovered = true
	} else {
		ed.hovered = false
	}
}

func (ed *EquipmentDetails) Draw(o *render.Options) {
	// Details
	op := &render.TextOptions{
		Screen: o.Screen,
		Font:   assets.DisplayFont,
		Color:  color.White,
	}

	op.GeoM.Translate(ed.Position())
	op.Font = assets.BodyFont
	op.Color = color.RGBA{200, 200, 200, 255}

	render.DrawText(op, ed.equipmentType.String())
	op.GeoM.Translate(10, assets.BodyFont.LineHeight+1)
	equipment := ed.equipment

	if equipment != nil {
		op.Color = equipment.quality.TextColor()
		render.DrawText(op, equipment.Name())
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		if equipment.perk != nil {
			op.GeoM.Translate(10, 0)
			op.Color = equipment.perk.Quality().TextColor()
			render.DrawText(op, equipment.perk.Name())
			op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
			op.GeoM.Translate(-10, 0)
		} else {
			op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		}
	} else {
		op.GeoM.Translate(0, (assets.BodyFont.LineHeight+1)*2)
	}
	op.GeoM.Translate(-10, 0)
}

func PaddedIntString(i int, pad int) string {
	str := fmt.Sprintf("%d", i)
	for len(str) < pad {
		str = " " + str
	}
	return str
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

	dp.dudeDetails.Layout(o, dp)
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

	selectedDude := false
	for _, p := range dp.dudeProfiles {
		px, py := p.Position()
		if InBounds(px, py, dp.width, p.height, mx, my) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				selectedDude = true
				dp.dudeDetails.SetDude(p.dude)
				if dp.onDudeClick != nil {
					dp.onDudeClick(p.dude)
				}
			}
			p.hovered = true
		} else {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if !selectedDude {
					dp.dudeDetails.SetDude(nil)
				}
			}
			p.hovered = false
		}
	}

	for _, ed := range dp.dudeDetails.equipmentDetails {
		ed.Update(o)
	}
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

	o.DrawImageOptions.GeoM.Reset()

	// Draw dude details
	dp.dudeDetails.Draw(o)
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
	drawerInterp render.InterpNumber
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

// ===============================================
type Button struct {
	baseSprite  *render.Sprite
	sprite      *render.Sprite
	onClick     func()
	wobbler     float64
	tooltip     string
	showTooltip bool
}

func NewButton(name string, tooltip string) *Button {
	return &Button{
		baseSprite: Must(render.NewSpriteFromStaxie("ui/button", "base")),
		sprite:     Must(render.NewSpriteFromStaxie("ui/button", name)),
		tooltip:    tooltip,
	}
}

func (b *Button) Layout(o *UIOptions) {
	b.baseSprite.Scale = o.Scale
	b.sprite.Scale = o.Scale
}

func (b *Button) Update() {
	x, y := b.Position()
	w, h := b.sprite.Size()
	mx, my := IntToFloat2(ebiten.CursorPosition())
	if InBounds(x, y, w, h, mx, my) {
		b.showTooltip = true
		b.wobbler += 0.1
	} else {
		b.showTooltip = false
		b.wobbler = 0
	}
	b.sprite.SetRotation(math.Sin(b.wobbler) * 0.05)
	b.baseSprite.SetRotation(math.Sin(b.wobbler) * 0.05)
}

func (b *Button) Check(mx, my float64) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := b.Position()
		w, h := b.sprite.Size()
		if mx > x && mx < x+w && my > y && my < y+h {
			if b.onClick != nil {
				b.onClick()
			}
		}
	}
}

func (b *Button) SetPosition(x, y float64) {
	b.baseSprite.SetPosition(x, y)
	b.sprite.SetPosition(x, y)
}

func (b *Button) Position() (float64, float64) {
	return b.baseSprite.Position()
}

func (b *Button) SetImage(name string) {
	b.sprite.SetStaxie("ui/button", name)
}

func (b *Button) Draw(o *render.Options) {
	b.baseSprite.Draw(o)
	b.sprite.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	if b.tooltip != "" && b.showTooltip {
		op := &render.TextOptions{
			Screen: o.Screen,
			Font:   assets.DisplayFont,
			Color:  color.NRGBA{184, 152, 93, 200},
		}
		width, _ := text.Measure(b.tooltip, assets.DisplayFont.Face, assets.BodyFont.LineHeight)
		x, y := b.Position()
		w, h := b.sprite.Size()
		x += w
		op.GeoM.Translate(x-width, y+h)
		render.DrawText(op, b.tooltip)
	}
}

type SpeedPanel struct {
	render.Positionable
	width        float64
	height       float64
	cameraButton *Button
	pauseButton  *Button
	speedButton  *Button
	musicButton  *Button
	soundButton  *Button
}

func (sp *SpeedPanel) Layout(o *UIOptions) {
	sp.cameraButton.Layout(o)
	sp.pauseButton.Layout(o)
	sp.speedButton.Layout(o)
	sp.musicButton.Layout(o)
	sp.soundButton.Layout(o)

	bw, bh := sp.pauseButton.sprite.Size()

	sp.width = bw*5 + 4*5
	sp.height = bh + bh/4

	x := float64(o.Width) - sp.width
	y := 4.0

	sp.SetPosition(x, y)

	sp.musicButton.SetPosition(x, y)
	x += bw + 4
	sp.soundButton.SetPosition(x, y)
	x += bw + 4
	sp.cameraButton.SetPosition(x, y)
	x += bw + 4
	sp.pauseButton.SetPosition(x, y)
	x += bw + 4
	sp.speedButton.SetPosition(x, y)
}

func (sp *SpeedPanel) Update(o *UIOptions) {
	mx, my := IntToFloat2(ebiten.CursorPosition())
	x, y := sp.Position()
	if InBounds(x, y, sp.width, sp.height, mx, my) {
		sp.musicButton.Update()
		sp.soundButton.Update()
		sp.cameraButton.Update()
		sp.pauseButton.Update()
		sp.speedButton.Update()
		sp.musicButton.Check(mx, my)
		sp.soundButton.Check(mx, my)
		sp.cameraButton.Check(mx, my)
		sp.pauseButton.Check(mx, my)
		sp.speedButton.Check(mx, my)
	}
}

func (sp *SpeedPanel) Draw(o *render.Options) {
	sp.musicButton.Draw(o)
	sp.soundButton.Draw(o)
	sp.cameraButton.Draw(o)
	sp.pauseButton.Draw(o)
	sp.speedButton.Draw(o)
}

type MessagePanel struct {
	render.Positionable
	width        float64
	height       float64
	drawered     bool
	pinned       bool
	maxLines     int
	drawerInterp render.InterpNumber
	top          *render.Sprite
	topleft      *render.Sprite
	topright     *render.Sprite
	mid          *render.Sprite
	midleft      *render.Sprite
	midright     *render.Sprite
}

func (mp *MessagePanel) Layout(o *UIOptions) {
	// eww
	mp.mid.Scale = o.Scale
	mp.midleft.Scale = o.Scale
	mp.midright.Scale = o.Scale
	mp.top.Scale = o.Scale
	mp.topleft.Scale = o.Scale
	mp.topright.Scale = o.Scale

	mp.width = float64(o.Width) * 0.75
	mp.height = assets.BodyFont.LineHeight*float64(mp.maxLines) + 15 // buffer
	mp.SetPosition((float64(o.Width))/2-(mp.width/2), float64(o.Height)-mp.height+50)
}

func (mp *MessagePanel) Update(o *UIOptions) {
	mp.drawerInterp.Update()

	rpx, rpy := mp.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())

	maxX := rpx + float64(mp.width)
	maxY := rpy + mp.height

	_, ph := mp.topleft.Size()

	if mx > rpx && mx < maxX && my > rpy && my < maxY {
		if mp.drawered {
			mp.drawered = false
			mp.drawerInterp.Set(0, 4)
		}
	} else {
		if !mp.drawered {
			mp.drawered = true
			mp.drawerInterp.Set(mp.height-ph*2, 4)
		}
	}
}

func (mp *MessagePanel) Draw(o *render.Options) {
	// Draw the panel
	x, y := mp.Position()
	pw, ph := mp.top.Size()

	op := &render.Options{
		Screen: o.Screen,
	}
	op.DrawImageOptions.GeoM.Translate(0, mp.drawerInterp.Current)
	op.DrawImageOptions.GeoM.Translate(x, y)

	// top
	mp.topleft.Draw(op)
	op.DrawImageOptions.GeoM.Translate(pw, 0)
	for x := 0; x < int(mp.width/pw)-2; x++ {
		mp.top.Draw(op)
		op.DrawImageOptions.GeoM.Translate(pw, 0)
	}
	mp.topright.Draw(op)
	op.DrawImageOptions.GeoM.Translate(-mp.width+pw, ph)

	// mid
	for y := 0; y < int(mp.height/ph)-2; y++ {
		mp.midleft.Draw(op)
		op.DrawImageOptions.GeoM.Translate(pw, 0)
		for x := 0; x < int(mp.width/pw)-2; x++ {
			mp.mid.Draw(op)
			op.DrawImageOptions.GeoM.Translate(pw, 0)
		}
		mp.midright.Draw(op)
		op.DrawImageOptions.GeoM.Translate(-mp.width+pw, ph)
	}

	if mp.drawered && !mp.pinned {
		return
	}

	messages := GetMessages()

	// Set initial position to bottom right of message panel
	baseX := x + mp.width - 10
	baseY := y + mp.height - 10 // Bottom edge minus padding

	// Calculate the number of messages to display
	maxLines := min(mp.maxLines-1, len(messages))

	// Render messages from bottom to top
	for i := 0; i < maxLines; i++ {
		messageIndex := len(messages) - 1 - i
		if messageIndex < 0 {
			break
		}
		message := messages[messageIndex]

		tOp := &render.TextOptions{
			Screen: o.Screen,
			Font:   assets.BodyFont,
			Color:  message.kind.Color(),
		}

		w, h := text.Measure(message.text, assets.BodyFont.Face, assets.BodyFont.LineHeight)
		posX := baseX - w
		posY := baseY - float64(h*float64(i+1))

		// Ensure the text doesn't go above the panel
		if posY < y {
			break
		}

		tOp.GeoM.Translate(posX, posY)
		render.DrawText(tOp, message.text)
	}
}
