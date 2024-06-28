package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	gameInfoPanel  GameInfoPanel
	dudePanel      DudePanel
	dudePanel2     DudePanel2
	equipmentPanel EquipmentPanel
	speedPanel     SpeedPanel
	messagePanel   MessagePanel
	roomPanel      RoomPanel
	roomInfoPanel  RoomInfoPanel
	options        *UIOptions
	feedback       FeedbackPopup
	buttonPanel    ButtonPanel
	bossPanel      BossPanel
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
			maxLines: 15,
			top:      Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
			topleft:  Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
			topright: Must(render.NewSubSprite(panelSprite, 32, 0, 16, 16)),
			mid:      Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
			midleft:  Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
			midright: Must(render.NewSubSprite(panelSprite, 32, 16, 16, 16)),
		}
	}
	ui.gameInfoPanel = MakeGameInfoPanel()
	ui.speedPanel = MakeSpeedPanel()
	ui.dudePanel2 = MakeDudePanel2()
	ui.equipmentPanel = MakeEquipmentPanel()
	ui.roomPanel = MakeRoomPanel()
	ui.roomInfoPanel = MakeRoomInfoPanel()
	ui.feedback = MakeFeedbackPopup()
	ui.buttonPanel = MakeButtonPanel()
	ui.buttonPanel.Disable()

	ui.bossPanel = MakeBossPanel()

	return ui
}

func (ui *UI) Layout(o *UIOptions) {
	ui.options = o
	ui.dudePanel.Layout(o)
	ui.speedPanel.Layout(nil, o)
	ui.messagePanel.Layout(o)

	// Position info.
	ui.gameInfoPanel.dudePanel.SetSize(
		96*o.Scale,
		8*o.Scale,
	)
	ui.gameInfoPanel.storyPanel.SetSize(
		96*o.Scale,
		8*o.Scale,
	)
	ui.gameInfoPanel.goldPanel.SetSize(
		96*o.Scale,
		8*o.Scale,
	)

	ui.gameInfoPanel.panel.SetSize(
		ui.gameInfoPanel.dudePanel.Width()+ui.gameInfoPanel.storyPanel.Width()+ui.gameInfoPanel.goldPanel.Width(),
		32*o.Scale,
	)
	ui.gameInfoPanel.panel.SetPosition(
		float64(o.Width/2)-ui.gameInfoPanel.panel.Width()/2,
		-ui.gameInfoPanel.panel.Height()/3,
	)

	ui.gameInfoPanel.Layout(o)

	// Manually position dude panel and equipment panel
	h := float64(o.Height)/2 - float64(o.Height)/10
	ui.dudePanel2.panel.SetSize(
		64*o.Scale,
		h,
	)
	ui.dudePanel2.panel.SetPosition(
		8,
		float64(o.Height)/2-h-8,
	)

	ui.equipmentPanel.panel.SetSize(
		64*o.Scale,
		h,
	)
	ui.equipmentPanel.panel.SetPosition(
		8,
		float64(o.Height)/2+8,
	)

	ui.dudePanel2.Layout(o)
	ui.equipmentPanel.Layout(o)

	// Manually position roomPanel
	ui.roomPanel.panel.SetSize(
		64*o.Scale,
		float64(o.Height)-float64(o.Height)/3,
	)
	ui.roomPanel.panel.SetPosition(
		float64(o.Width)-ui.roomPanel.panel.Width()-8,
		float64(o.Height)/2-ui.roomPanel.panel.Height()/2,
	)
	ui.roomPanel.Layout(o)

	ui.roomInfoPanel.panel.SetSize(
		224*o.Scale,
		64*o.Scale,
	)
	ui.roomInfoPanel.panel.SetPosition(
		ui.roomPanel.panel.X()-ui.roomInfoPanel.panel.Width()-8,
		ui.roomPanel.panel.Y()+(ui.roomPanel.panel.Height()-ui.roomInfoPanel.panel.Height()),
	)
	ui.roomInfoPanel.Layout(o)

	ui.feedback.panel.SetSize(
		320*o.Scale,
		16*o.Scale,
	)
	ui.feedback.panel.SetPosition(
		float64(o.Width)/2-ui.feedback.panel.Width()/2,
		float64(o.Height)/2-ui.feedback.panel.Height()/2,
	)
	ui.feedback.Layout(o)

	ui.buttonPanel.panel.SetPosition(
		float64(o.Width)/2-ui.buttonPanel.panel.Width()/2,
		float64(o.Height)/2-ui.buttonPanel.panel.Height()/2+float64(o.Height)/4,
	)
	ui.buttonPanel.Layout(o)

	ui.bossPanel.panel.SetSize(
		float64(o.Width)/3,
		24*o.Scale,
	)
	ui.bossPanel.panel.SetPosition(
		float64(o.Width)/2-ui.bossPanel.panel.Width()/2,
		float64(o.Height)/8-ui.bossPanel.panel.Height()/2,
	)
	ui.bossPanel.panel.padding = 2 * o.Scale
	ui.bossPanel.Layout(o)
}

func (ui *UI) Update(o *UIOptions) {
	ui.dudePanel.Update(o)
	ui.dudePanel2.Update(o)
	ui.equipmentPanel.Update(o)
	ui.roomPanel.Update(o)
	ui.speedPanel.Update(o)
	ui.messagePanel.Update(o)
	ui.feedback.Update(o)
	ui.buttonPanel.Update(o)
}

func (ui *UI) Check(mx, my float64, kind UICheckKind) bool {
	if ui.dudePanel2.Check(mx, my, kind) {
		return true
	}
	if ui.equipmentPanel.Check(mx, my, kind) {
		return true
	}
	if ui.roomPanel.Check(mx, my, kind) {
		return true
	}

	if ui.speedPanel.Check(mx, my, kind) {
		return true
	}

	if ui.buttonPanel.Check(mx, my, kind) {
		return true
	}
	return false
}

func (ui *UI) Draw(o *render.Options) {
	ui.buttonPanel.Draw(o)

	ui.dudePanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.equipmentPanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.speedPanel.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	ui.messagePanel.Draw(o)

	o.DrawImageOptions.GeoM.Reset()
	ui.roomPanel.Draw(o)

	ui.roomInfoPanel.Draw(o)

	o.DrawImageOptions.GeoM.Reset()
	ui.dudePanel2.Draw(o)

	ui.gameInfoPanel.Draw(o)

	ui.bossPanel.Draw(o)

	ui.feedback.Draw(o)
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
	dude       *Dude
	stack      *render.Stack
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
		stats := dp.dude.GetCalculatedStats()
		render.DrawText(op, fmt.Sprintf("HP: %d/%d", stats.currentHp, stats.totalHp))
		op.Color = color.RGBA{200, 200, 200, 255}
		op.GeoM.Translate(0, assets.BodyFont.LineHeight*2)
		render.DrawText(op, fmt.Sprintf("%s strength", PaddedIntString(stats.strength, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s agility", PaddedIntString(stats.agility, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s defense", PaddedIntString(stats.defense, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s wisdom", PaddedIntString(stats.wisdom, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s cowardice", PaddedIntString(stats.cowardice, 4)))
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("%s luck", PaddedIntString(stats.luck, 4)))
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

	// Use calculated stats for display
	stats := dd.dude.GetCalculatedStats()
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
	render.DrawText(op, fmt.Sprintf("HP: %d/%d", stats.currentHp, stats.totalHp))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	op.Color = color.RGBA{255, 215, 0, 255}
	render.DrawText(op, fmt.Sprintf("%.2f gold", dd.dude.gold))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)

	// XP
	op.Color = color.RGBA{200, 200, 255, 255}
	render.DrawText(op, fmt.Sprintf("%d/%d XP", dd.dude.xp, dd.dude.NextLevelXP()))

	op.Color = color.RGBA{200, 200, 200, 255}
	op.GeoM.Translate(0, assets.BodyFont.LineHeight*2)
	render.DrawText(op, fmt.Sprintf("%s strength", PaddedIntString(stats.strength, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s agility", PaddedIntString(stats.agility, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s defense", PaddedIntString(stats.defense, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s wisdom", PaddedIntString(stats.wisdom, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s cowardice", PaddedIntString(stats.cowardice, 4)))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
	render.DrawText(op, fmt.Sprintf("%s luck", PaddedIntString(stats.luck, 4)))

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
		op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)
		render.DrawText(op, fmt.Sprintf("Uses: %d/%d", hoveredEquipment.uses, hoveredEquipment.totalUses))
	}

	// RIGHT SIDE
	// Draw equipment type
	op.GeoM.Reset()
	op.GeoM.Translate(x+dd.width-100, y+dd.height-10)
	op.Color = color.RGBA{200, 200, 200, 255}
	render.DrawText(op, fmt.Sprintf("Level %d %s", hoveredEquipment.stats.level, hoveredEquipment.Type().String()))
	op.GeoM.Translate(0, assets.BodyFont.LineHeight+1)

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
	dp.dudeProfiles = nil // Reset them dude profiles
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

type SpeedPanel struct {
	render.Positionable
	render.Sizeable
	cameraButton *UIButton
	pauseButton  *UIButton
	speedButton  *UIButton
	musicButton  *UIButton
	soundButton  *UIButton
	buttons      []UIElement
}

func MakeSpeedPanel() SpeedPanel {
	sp := SpeedPanel{}
	sp.musicButton = NewUIButton("music", "music on")
	sp.soundButton = NewUIButton("sound", "sound on")
	sp.pauseButton = NewUIButton("play", "playing")
	sp.speedButton = NewUIButton("fast", "fast")
	sp.cameraButton = NewUIButton("story", "camera: story")
	sp.buttons = append(sp.buttons, sp.musicButton)
	sp.buttons = append(sp.buttons, sp.soundButton)
	sp.buttons = append(sp.buttons, sp.cameraButton)
	sp.buttons = append(sp.buttons, sp.pauseButton)
	sp.buttons = append(sp.buttons, sp.speedButton)
	return sp
}

func (sp *SpeedPanel) Layout(parent UIElement, o *UIOptions) {
	for _, b := range sp.buttons {
		b.Layout(sp, o)
	}

	bw, bh := sp.pauseButton.sprite.Size()

	sp.SetSize(
		bw*float64(len(sp.buttons))+4*float64(len(sp.buttons)),
		bh+bh/4,
	)

	x := float64(o.Width) - sp.Width()
	y := 4.0

	sp.SetPosition(x, y)

	for _, b := range sp.buttons {
		b.SetPosition(x, y)
		x += bw + 4
	}
}

func (sp *SpeedPanel) Update(o *UIOptions) {
	for _, b := range sp.buttons {
		b.Update(o)
	}
}

func (sp *SpeedPanel) Check(mx, my float64, kind UICheckKind) bool {
	inBounds := InBounds(sp.X(), sp.Y(), sp.Width(), sp.Height(), mx, my)
	if kind == UICheckHover && !inBounds {
		return false
	}
	for _, b := range sp.buttons {
		if b.Check(mx, my, kind) {
			return true
		}
	}
	return false
}

func (sp *SpeedPanel) Draw(o *render.Options) {
	for _, b := range sp.buttons {
		b.Draw(o)
	}
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
	mp.SetPosition((float64(o.Width))/2-(mp.width/2), float64(o.Height)-mp.height)
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
	for y := 0; y < int(mp.height/ph); y++ {
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
	baseY := y + mp.height - 17 // Bottom edge minus padding

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
		posY := baseY - float64(h*float64(i))

		// Ensure the text doesn't go above the panel
		if posY < y {
			break
		}

		tOp.GeoM.Translate(posX, posY)
		render.DrawText(tOp, message.text)
	}
}

type RoomPanel struct {
	panel       *UIPanel
	list        *UIItemList
	title       *UIText
	count       *UIText
	roomDefs    []*RoomDef
	onItemClick func(index int)
}

func MakeRoomPanel() RoomPanel {
	rp := RoomPanel{
		panel: NewUIPanel(PanelStyleInteractive),
		title: NewUIText("Rooms", assets.DisplayFont, assets.ColorHeading),
		count: NewUIText("0", assets.BodyFont, assets.ColorHeading),
		list:  NewUIItemList(DirectionVertical),
	}
	rp.panel.AddChild(rp.title)
	rp.panel.AddChild(rp.list)
	rp.panel.sizeChildren = true
	rp.panel.centerChildren = true
	rp.list.centerItems = true
	rp.list.centerList = true

	return rp
}

func (rp *RoomPanel) SetRoomDefs(roomDefs []*RoomDef) {
	rp.roomDefs = roomDefs
	rp.list.Clear()
	for index, rd := range roomDefs {
		img := NewUIImage(rd.image)
		img.ignoreScale = true
		img.onCheck = func(kind UICheckKind) {
			if kind == UICheckClick && rp.onItemClick != nil {
				rp.onItemClick(index)
				rp.list.selected = index
			}
		}
		rp.list.AddItem(img)
	}
	rp.count.text = fmt.Sprintf("%d", len(rp.roomDefs))
}

func (rp *RoomPanel) Layout(o *UIOptions) {
	rp.panel.padding = 6 * o.Scale
	rp.list.SetSize(rp.panel.Width(), rp.panel.Height()-rp.panel.padding*2-rp.title.Height())

	rp.panel.Layout(nil, o)

	rp.count.Layout(nil, o)
	rp.count.SetPosition(rp.list.X(), rp.list.Y()-rp.count.Height()/4)
}

func (rp *RoomPanel) Update(o *UIOptions) {
	rp.panel.Update(o)
}

func (rp *RoomPanel) Check(mx, my float64, kind UICheckKind) bool {
	return rp.panel.Check(mx, my, kind)
}

func (rp *RoomPanel) Draw(o *render.Options) {
	rp.panel.Draw(o)
	rp.count.Draw(o)
}

type RoomInfoPanel struct {
	panel       *UIPanel
	title       *UIText
	description *UIText
	cost        *UIText

	hidden bool
}

func MakeRoomInfoPanel() RoomInfoPanel {
	rip := RoomInfoPanel{
		panel:       NewUIPanel(PanelStyleNormal),
		title:       NewUIText("Room Info", assets.DisplayFont, assets.ColorHeading),
		description: NewUIText("Description", assets.BodyFont, assets.ColorRoomDescription),
		cost:        NewUIText("Cost: 0", assets.BodyFont, assets.ColorRoomCost),
		hidden:      true,
	}
	rip.panel.AddChild(rip.title)
	rip.panel.AddChild(rip.description)
	rip.panel.AddChild(rip.cost)
	rip.panel.sizeChildren = true
	//rip.panel.centerChildren = true
	return rip
}

func (rip *RoomInfoPanel) Layout(o *UIOptions) {
	rip.panel.padding = 6 * o.Scale
	rip.panel.Layout(nil, o)
}

func (rip *RoomInfoPanel) Update(o *UIOptions) {
	rip.panel.Update(o)
}

func (rip *RoomInfoPanel) Check(mx, my float64, kind UICheckKind) bool {
	if rip.hidden {
		return false
	}
	return rip.panel.Check(mx, my, kind)
}

func (rip *RoomInfoPanel) Draw(o *render.Options) {
	if rip.hidden {
		return
	}
	rip.panel.Draw(o)
}

type GameInfoPanel struct {
	panel      *UIPanel
	storyPanel *UIPanel
	storyText  *UIText
	dudePanel  *UIPanel
	dudeText   *UIText
	goldPanel  *UIPanel
	goldText   *UIText
}

func MakeGameInfoPanel() GameInfoPanel {
	gip := GameInfoPanel{
		panel:      NewUIPanel(PanelStyleNormal),
		storyPanel: NewUIPanel(PanelStyleNormal),
		dudePanel:  NewUIPanel(PanelStyleNormal),
		goldPanel:  NewUIPanel(PanelStyleNormal),
	}

	gip.storyText = NewUIText("Story 0/0", assets.BodyFont, assets.ColorStory)
	gip.storyPanel.AddChild(gip.storyText)
	gip.storyPanel.flowDirection = DirectionHorizontal
	gip.storyPanel.hideBackground = true

	gip.dudeText = NewUIText("Dudes 0", assets.BodyFont, assets.ColorDude)
	gip.dudePanel.AddChild(gip.dudeText)
	gip.dudePanel.flowDirection = DirectionHorizontal
	gip.dudePanel.hideBackground = true

	gip.goldText = NewUIText("Gold 0", assets.BodyFont, assets.ColorGold)
	gip.goldPanel.AddChild(gip.goldText)
	gip.goldPanel.flowDirection = DirectionHorizontal
	gip.goldPanel.hideBackground = true

	gip.panel.AddChild(gip.storyPanel)
	gip.panel.AddChild(gip.dudePanel)
	gip.panel.AddChild(gip.goldPanel)
	gip.panel.flowDirection = DirectionHorizontal
	gip.panel.sizeChildren = true
	gip.panel.centerChildren = true

	return gip
}

func (gip *GameInfoPanel) Layout(o *UIOptions) {
	gip.panel.padding = 6 * o.Scale
	gip.storyPanel.padding = 6 * o.Scale
	gip.dudePanel.padding = 6 * o.Scale
	gip.goldPanel.padding = 6 * o.Scale
	gip.panel.Layout(nil, o)
}

func (gip *GameInfoPanel) Draw(o *render.Options) {
	gip.panel.Draw(o)
}

type DudePanel2 struct {
	panel *UIPanel
	list  *UIItemList
	title *UIText
	//dudeDefs []*DudeDef
}

func MakeDudePanel2() DudePanel2 {
	dp := DudePanel2{
		panel: NewUIPanel(PanelStyleInteractive),
		title: NewUIText("Dudes", assets.DisplayFont, assets.ColorHeading),
		list:  NewUIItemList(DirectionVertical),
	}
	dp.panel.AddChild(dp.title)
	dp.panel.AddChild(dp.list)
	dp.panel.sizeChildren = true
	dp.panel.centerChildren = true

	return dp
}

/*func (dp *DudePanel2) SetDudeDefs(dudeDefs []*DudeDef) {
	dp.list.Clear()
	for _, dd := range dudeDefs {
		img := NewUIImage(dd.image)
		img.ignoreScale = true
		dp.list.AddItem(img)
	}
}*/

func (dp *DudePanel2) Layout(o *UIOptions) {
	dp.panel.padding = 6 * o.Scale
	dp.list.SetSize(dp.panel.Width(), dp.panel.Height()-dp.panel.padding*2-dp.title.Height())

	dp.panel.Layout(nil, o)
}

func (dp *DudePanel2) Update(o *UIOptions) {
	dp.panel.Update(o)
}

func (dp *DudePanel2) Check(mx, my float64, kind UICheckKind) bool {
	return dp.panel.Check(mx, my, kind)
}

func (dp *DudePanel2) Draw(o *render.Options) {
	dp.panel.Draw(o)
}

type EquipmentPanel struct {
	panel *UIPanel
	list  *UIItemList
	title *UIText
}

func MakeEquipmentPanel() EquipmentPanel {
	ep := EquipmentPanel{
		panel: NewUIPanel(PanelStyleInteractive),
		title: NewUIText("Loot", assets.DisplayFont, assets.ColorHeading),
		list:  NewUIItemList(DirectionVertical),
	}
	ep.panel.AddChild(ep.title)
	ep.panel.AddChild(ep.list)
	ep.panel.sizeChildren = true
	ep.panel.centerChildren = true

	return ep
}

func (ep *EquipmentPanel) Layout(o *UIOptions) {
	ep.panel.padding = 6 * o.Scale
	ep.list.SetSize(ep.panel.Width(), ep.panel.Height()-ep.panel.padding*2-ep.title.Height())

	ep.panel.Layout(nil, o)
}

func (ep *EquipmentPanel) Update(o *UIOptions) {
	ep.panel.Update(o)
}

func (ep *EquipmentPanel) Check(mx, my float64, kind UICheckKind) bool {
	return ep.panel.Check(mx, my, kind)
}

func (ep *EquipmentPanel) Draw(o *render.Options) {
	ep.panel.Draw(o)
}

type FeedbackPopup struct {
	panel *UIPanel
	text  *UIText
	ticks int
	kind  FeedbackKind
}

func MakeFeedbackPopup() FeedbackPopup {
	fp := FeedbackPopup{
		panel: NewUIPanel(PanelStyleInteractive),
		text:  NewUIText("", assets.BodyFont, assets.ColorHeading),
	}
	fp.text.center = true
	fp.panel.AddChild(fp.text)
	fp.panel.hideBackground = true
	fp.panel.centerChildren = true
	return fp
}

func (fp *FeedbackPopup) Layout(o *UIOptions) {
	fp.panel.Layout(nil, o)
}

func (fp *FeedbackPopup) Update(o *UIOptions) {
	fp.ticks--
	// Fade out the last 10 ticks
	if fp.ticks < 10 {
		fp.text.textOptions.Color = color.NRGBA{
			R: fp.kind.R,
			G: fp.kind.G,
			B: fp.kind.B,
			A: uint8(float64(fp.ticks) / 10 * 255),
		}
	}

	fp.panel.Update(o)
}

func (fp *FeedbackPopup) Draw(o *render.Options) {
	if fp.ticks > 0 {
		fp.panel.Draw(o)
	}
}

func (fp *FeedbackPopup) Msg(kind FeedbackKind, text string) {
	fp.kind = kind
	fp.text.textOptions.Color = color.NRGBA(kind)
	fp.text.SetText(text)
	fp.ticks = len(text) * 5
}

type FeedbackKind color.NRGBA

var (
	FeedbackGeneric = FeedbackKind{255, 255, 255, 255}
	FeedbackGood    = FeedbackKind{0, 255, 0, 255}
	FeedbackBad     = FeedbackKind{255, 0, 0, 255}
	FeedbackWarning = FeedbackKind{255, 255, 0, 255}
)

type ButtonPanel struct {
	panel    *UIPanel
	text     *UIText
	hidden   bool
	disabled bool
	onClick  func()
}

func MakeButtonPanel() ButtonPanel {
	bp := ButtonPanel{
		panel: NewUIPanel(PanelStyleButton),
		text:  NewUIText("arghh", assets.DisplayFont, assets.ColorHeading),
	}
	bp.panel.AddChild(bp.text)
	bp.panel.sizeChildren = false
	bp.panel.centerChildren = true
	return bp
}

func (bp *ButtonPanel) Layout(o *UIOptions) {
	bp.panel.padding = 1 * o.Scale
	bp.text.Layout(nil, o)
	bp.panel.Layout(nil, o)
}

func (bp *ButtonPanel) doSize(o *UIOptions) {
	// lol size according to highest divisible
	width := math.Ceil(bp.text.Width()/bp.panel.left.Width())*bp.panel.left.Width() + bp.panel.left.Width() + bp.panel.right.Width()
	height := math.Ceil(bp.text.Height()/bp.panel.top.Height()) * bp.panel.top.Height()
	bp.panel.SetSize(width, height)

	bp.panel.Layout(nil, o)
}

func (bp *ButtonPanel) Update(o *UIOptions) {
	if bp.hidden {
		return
	}
	// eh...
	bp.doSize(o)
	bp.panel.Update(o)
}

func (bp *ButtonPanel) Check(mx, my float64, kind UICheckKind) bool {
	if bp.hidden || bp.disabled {
		return false
	}
	if InBounds(bp.panel.X(), bp.panel.Y(), bp.panel.Width(), bp.panel.Height(), mx, my) {
		if kind == UICheckClick {
			if bp.onClick != nil {
				bp.onClick()
				return true
			}
		}
		if kind == UICheckHover {
			return true
		}
	}
	return false
}

func (bp *ButtonPanel) Draw(o *render.Options) {
	if bp.hidden || bp.disabled {
		return
	}
	bp.panel.Draw(o)
}

func (bp *ButtonPanel) Disable() {
	bp.panel.SetStyle(PanelStyleButtonDisabled)
	bp.disabled = true
	bp.text.textOptions.ColorScale.ScaleAlpha(0.5)
}

func (bp *ButtonPanel) Enable() {
	bp.panel.SetStyle(PanelStyleButton)
	bp.disabled = false
	bp.text.textOptions.ColorScale.ScaleAlpha(1.0)
}

type BossPanel struct {
	panel   *UIPanel
	text    *UIText
	current float64
	hidden  bool
}

func MakeBossPanel() BossPanel {
	bp := BossPanel{
		panel:   NewUIPanel(PanelStyleBar),
		text:    NewUIText("baus", assets.DisplayFont, assets.ColorBoss),
		current: 1.0,
		hidden:  true,
	}
	bp.panel.AddChild(bp.text)
	return bp
}

func (bp *BossPanel) Layout(o *UIOptions) {
	if bp.hidden {
		return
	}

	bp.text.Layout(nil, o)
	bp.panel.Layout(nil, o)
	// Force text's position
	bp.text.SetPosition(bp.panel.X()+bp.panel.Width()/2-bp.text.Width()/2, bp.panel.Y()+bp.panel.Height()/2-bp.text.Height()/2-3)
}

func (bp *BossPanel) Update(o *UIOptions) {
	if bp.hidden {
		return
	}

	bp.panel.Update(o)
}

func (bp *BossPanel) Draw(o *render.Options) {
	if bp.hidden {
		return
	}

	// Draw that bar.
	x, y := bp.panel.Position()
	w, h := bp.panel.Size()

	x += bp.panel.padding
	y += bp.panel.padding
	w -= bp.panel.padding * 2
	h -= bp.panel.padding * 2

	w *= bp.current

	vector.DrawFilledRect(o.Screen, float32(x), float32(y), float32(w), float32(h), color.NRGBA{200, 20, 20, 200}, true)

	bp.panel.Draw(o)
}
