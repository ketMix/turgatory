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
	dudeInfoPanel  DudeInfoPanel
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
		panelSprite := Must(render.NewSprite("ui/altPanels"))
		ui.messagePanel = MessagePanel{
			maxLines: 8,
			top:      Must(render.NewSubSprite(panelSprite, 16, 0, 16, 16)),
			topleft:  Must(render.NewSubSprite(panelSprite, 0, 0, 16, 16)),
			topright: Must(render.NewSubSprite(panelSprite, 32, 0, 16, 16)),
			mid:      Must(render.NewSubSprite(panelSprite, 16, 16, 16, 16)),
			midleft:  Must(render.NewSubSprite(panelSprite, 0, 16, 16, 16)),
			midright: Must(render.NewSubSprite(panelSprite, 32, 16, 16, 16)),
			pinned:   false,
		}
	}
	ui.gameInfoPanel = MakeGameInfoPanel()
	ui.speedPanel = MakeSpeedPanel()
	ui.dudePanel = MakeDudePanel()
	ui.dudeInfoPanel = MakeDudeInfoPanel()
	ui.equipmentPanel = MakeEquipmentPanel()
	ui.roomPanel = MakeRoomPanel()
	ui.roomInfoPanel = MakeRoomInfoPanel()
	ui.feedback = MakeFeedbackPopup()
	ui.buttonPanel = MakeButtonPanel(assets.DisplayFont, PanelStyleButton)
	ui.buttonPanel.Disable()

	ui.bossPanel = MakeBossPanel()

	return ui
}

func (ui *UI) Layout(o *UIOptions) {
	ui.options = o
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
	h := float64(o.Height)/2 - float64(o.Height)/12
	ui.dudePanel.panel.SetSize(
		96*o.Scale,
		h-8*o.Scale,
	)
	ui.dudePanel.panel.SetPosition(
		8,
		float64(o.Height)/2-h-8*o.Scale,
	)

	ui.equipmentPanel.panel.SetSize(
		96*o.Scale,
		h-8*o.Scale,
	)
	ui.equipmentPanel.panel.SetPosition(
		8,
		float64(o.Height)/2+8*o.Scale,
	)

	ui.dudePanel.Layout(o)

	ui.dudeInfoPanel.Layout(o)

	h = 64
	h = 64

	ui.dudeInfoPanel.panel.SetSize(
		64*o.Scale,
		h*o.Scale,
	)
	ts := ui.dudeInfoPanel.title.Width() + 8*o.Scale
	if ts > ui.dudeInfoPanel.panel.Width() {
		ts = math.Ceil(ts/ui.dudeInfoPanel.panel.center.Width()) * ui.dudeInfoPanel.panel.center.Width()
		ui.dudeInfoPanel.panel.SetSize(ts, ui.dudeInfoPanel.panel.Height())
	}
	ui.dudeInfoPanel.panel.SetPosition(
		ui.dudePanel.panel.X()+ui.dudePanel.panel.Width()+4*o.Scale,
		ui.dudePanel.panel.Y(),
	)

	ui.equipmentPanel.Layout(o)

	// Manually position roomPanel
	ui.roomPanel.panel.SetSize(
		96*o.Scale,
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
		ui.roomPanel.panel.X()-ui.roomInfoPanel.panel.Width()-4*o.Scale,
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
	ui.buttonPanel.Layout(nil, o)

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
	ui.dudeInfoPanel.Update(o)
	ui.equipmentPanel.Update(o)
	ui.roomPanel.Update(o)
	ui.speedPanel.Update(o)
	ui.messagePanel.Update(o)
	ui.feedback.Update(o)
	ui.buttonPanel.Update(o)
}

func (ui *UI) Check(mx, my float64, kind UICheckKind) bool {
	if ui.dudePanel.Check(mx, my, kind) {
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
	ui.dudePanel.Draw(o)
	ui.dudeInfoPanel.Draw(o)

	ui.gameInfoPanel.Draw(o)

	ui.bossPanel.Draw(o)

	ui.feedback.Draw(o)
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
	width    float64
	height   float64
	drawered bool
	pinned   bool
	maxLines int
	//drawerInterp render.InterpNumber
	top      *render.Sprite
	topleft  *render.Sprite
	topright *render.Sprite
	mid      *render.Sprite
	midleft  *render.Sprite
	midright *render.Sprite
}

func (mp *MessagePanel) Layout(o *UIOptions) {
	// eww
	mp.mid.Scale = o.Scale
	mp.midleft.Scale = o.Scale
	mp.midright.Scale = o.Scale
	mp.top.Scale = o.Scale
	mp.topleft.Scale = o.Scale
	mp.topright.Scale = o.Scale

	//mp.width = float64(o.Width) * 0.75
	mp.width = 208 * o.Scale
	lines := mp.maxLines
	if mp.pinned {
		lines *= 3
	}

	mp.height = assets.BodyFont.LineHeight*float64(lines) + 15 // buffer
	mp.SetPosition((float64(o.Width))/2-(mp.width/2), float64(o.Height)-mp.height)
}

func (mp *MessagePanel) Update(o *UIOptions) {
	//mp.drawerInterp.Update()

	rpx, rpy := mp.Position()
	mx, my := IntToFloat2(ebiten.CursorPosition())

	maxX := rpx + float64(mp.width)
	maxY := rpy + mp.height

	//_, ph := mp.topleft.Size()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mx > rpx && mx < maxX && my > rpy && my < maxY {
			mp.pinned = !mp.pinned
			/*if mp.drawered {
				mp.drawered = false
				mp.drawerInterp.Set(0, 4)
			}*/
		} else {
			/*if !mp.drawered {
				mp.drawered = true
				mp.drawerInterp.Set(mp.height-ph*2, 4)
			}*/
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
	//op.DrawImageOptions.GeoM.Translate(0, mp.drawerInterp.Current)
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

	messages := GetMessages()

	// Set initial position to bottom right of message panel
	//baseX := x + mp.width - 10
	baseX := x + 10
	baseY := y + mp.height - 17 // Bottom edge minus padding

	lines := mp.maxLines
	if mp.pinned {
		lines *= 3
	}

	// Calculate the number of messages to display
	maxLines := min(lines-1, len(messages))

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

		_, h := text.Measure(message.text, assets.BodyFont.Face, assets.BodyFont.LineHeight)
		posX := baseX
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
	buyButton   *ButtonPanel
	onBuyClick  func()
}

func MakeRoomPanel() RoomPanel {
	rp := RoomPanel{
		panel: NewUIPanel(PanelStyleInteractive),
		title: NewUIText("Rooms", assets.DisplayFont, assets.ColorHeading),
		count: NewUIText("0", assets.BodyFont, assets.ColorHeading),
		list:  NewUIItemList(DirectionVertical),
	}
	btn := MakeButtonPanel(assets.BodyFont, PanelStyleButtonAttached)
	rp.buyButton = &btn
	rp.buyButton.text.center = true
	rp.buyButton.text.SetText("Reroll Rooms\ngp")
	rp.list.spaceBetween = -2

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

	rp.buyButton.SetSize(rp.panel.Width(), 48)
	rp.buyButton.Layout(nil, o)
	rp.buyButton.text.SetPosition(rp.buyButton.text.X(), rp.buyButton.text.Y()+4*o.Scale)
	rp.buyButton.SetPosition(rp.panel.X()+rp.panel.Width()/2-rp.buyButton.Width()/2, rp.panel.Y()+rp.panel.Height()-10*o.Scale)
}

func (rp *RoomPanel) Update(o *UIOptions) {
	rp.panel.Update(o)
}

func (rp *RoomPanel) Check(mx, my float64, kind UICheckKind) bool {
	if rp.panel.Check(mx, my, kind) {
		return true
	}
	if rp.buyButton.Check(mx, my, kind) {
		if kind == UICheckClick && rp.onBuyClick != nil {
			rp.onBuyClick()
		}
		return true
	}
	return false
}

func (rp *RoomPanel) Draw(o *render.Options) {
	rp.buyButton.Draw(o)
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
		panel:       NewUIPanel(PanelStyleTransparent),
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

type DudePanel struct {
	panel       *UIPanel
	list        *UIItemList
	count       *UIText
	dudeSprites []*UIImage
	dudes       []*Dude
	title       *UIText
	buyButton   *ButtonPanel
	onItemClick func(index int)
	onItemHover func(index int)
	onBuyClick  func()
}

func MakeDudePanel() DudePanel {
	dp := DudePanel{
		panel: NewUIPanel(PanelStyleInteractive),
		title: NewUIText("Dudes", assets.DisplayFont, assets.ColorHeading),
		list:  NewUIItemList(DirectionVertical),
		count: NewUIText("0", assets.BodyFont, assets.ColorHeading),
	}
	btn := MakeButtonPanel(assets.BodyFont, PanelStyleButtonAttached)
	dp.buyButton = &btn
	dp.buyButton.text.center = true
	dp.buyButton.text.SetText("Buy\nRandom Dude")

	dp.panel.AddChild(dp.title)
	dp.panel.AddChild(dp.list)
	dp.panel.sizeChildren = true
	dp.panel.centerChildren = true
	dp.list.centerItems = true
	dp.list.centerList = true

	return dp
}

func (dp *DudePanel) SetDudes(dudes []*Dude) {
	dp.list.Clear()
	dp.dudeSprites = nil
	dp.dudes = nil
	for index, dude := range dudes {
		img := dp.DudeToImage(dude)
		img.scale = 1

		img.onCheck = func(kind UICheckKind) {
			if kind == UICheckClick && dp.onItemClick != nil {
				dp.onItemClick(index)
				dp.list.selected = index
			}
			if kind == UICheckHover && dp.onItemHover != nil {
				dp.onItemHover(index)
			}
		}

		dp.list.AddItem(img)
		dp.dudeSprites = append(dp.dudeSprites, img)
		dp.dudes = dudes
	}
	dp.count.text = fmt.Sprintf("%d", len(dudes))
}

func (dp *DudePanel) DudeToImage(dude *Dude) *UIImage {
	stack := render.CopyStack(dude.stack)

	img := ebiten.NewImage(int(float64(stack.Width())*1.25), int(float64(stack.Height())*2))
	img.Clear()

	profileOptions := render.Options{
		Screen: img,
		Pitch:  1,
	}
	stack.SetPosition(0, float64(stack.Height())*1.25)
	stack.SetRotation(math.Pi / 3)

	// Absoltue criminality.
	profileOptions.DrawImageOptions.ColorScale.Scale(0, 0, 0, 1)
	for x := -1; x < 2; x += 2 {
		for y := -1; y < 2; y += 2 {
			profileOptions.DrawImageOptions.GeoM.Scale(1.25, 1)
			profileOptions.DrawImageOptions.GeoM.Translate(float64(x), float64(y))
			stack.Draw(&profileOptions)
			profileOptions.DrawImageOptions.GeoM.Reset()
		}
	}
	profileOptions.DrawImageOptions.ColorScale.Reset()

	profileOptions.DrawImageOptions.GeoM.Scale(1.25, 1)
	stack.Draw(&profileOptions)

	for _, eq := range dude.equipped {
		estack := render.CopyStack(eq.stack)
		estack.SetOrigin(stack.Origin())
		estack.SetPosition(stack.Position())
		estack.SetRotation(stack.Rotation())
		estack.Draw(&profileOptions)
	}

	return NewUIImage(profileOptions.Screen)
}

func (dp *DudePanel) Layout(o *UIOptions) {
	dp.panel.padding = 6 * o.Scale
	dp.list.SetSize(dp.panel.Width(), dp.panel.Height()-dp.panel.padding*2-dp.title.Height())
	for _, ds := range dp.dudeSprites {
		ds.Layout(dp.list, o)
	}
	dp.list.Layout(nil, o)

	dp.panel.Layout(nil, o)

	dp.count.SetPosition(dp.list.X(), dp.list.Y()-dp.count.Height()/4)
	dp.count.Layout(nil, o)

	dp.buyButton.SetSize(dp.panel.Width(), 48)
	dp.buyButton.Layout(nil, o)
	dp.buyButton.text.SetPosition(dp.buyButton.text.X(), dp.buyButton.text.Y()+4*o.Scale)
	dp.buyButton.SetPosition(dp.panel.X()+dp.panel.Width()/2-dp.buyButton.Width()/2, dp.panel.Y()+dp.panel.Height()-10*o.Scale)
}

func (dp *DudePanel) Update(o *UIOptions) {
	for i, d := range dp.dudes {
		if d.dirtyEquipment {
			dp.dudeSprites[i] = dp.DudeToImage(d)
			d.dirtyEquipment = false
		}
	}
	dp.panel.Update(o)
}

func (dp *DudePanel) Check(mx, my float64, kind UICheckKind) bool {
	if dp.panel.Check(mx, my, kind) {
		return true
	}
	if dp.buyButton.Check(mx, my, kind) {
		if kind == UICheckClick && dp.onBuyClick != nil {
			dp.onBuyClick()
		}
		return true
	}
	return false
}

func (dp *DudePanel) Draw(o *render.Options) {
	dp.buyButton.Draw(o)
	dp.panel.Draw(o)
	dp.count.Draw(o)
}

type DudeInfoPanel struct {
	panel       *UIPanel
	title       *UIText
	description *UIText
	cost        *UIText

	dude *Dude

	level     *UIText
	xp        *UIText
	hp        *UIText
	strength  *UIText
	agility   *UIText
	defense   *UIText
	wisdom    *UIText
	cowardice *UIText
	luck      *UIText

	hidden bool
}

func MakeDudeInfoPanel() DudeInfoPanel {
	dip := DudeInfoPanel{
		panel:     NewUIPanel(PanelStyleTransparent),
		title:     NewUIText("Mah Dude", assets.DisplayFont, assets.ColorDudeTitle),
		level:     NewUIText("Level 0 sucker", assets.BodyFont, assets.ColorDudeLevel),
		xp:        NewUIText("0/0 xp", assets.BodyFont, assets.ColorDudeXP),
		hp:        NewUIText("0/0 hp", assets.BodyFont, assets.ColorDudeHP),
		strength:  NewUIText("0 strength", assets.BodyFont, assets.ColorDudeStrength),
		agility:   NewUIText("0 agility", assets.BodyFont, assets.ColorDudeAgility),
		defense:   NewUIText("0 defense", assets.BodyFont, assets.ColorDudeDefense),
		wisdom:    NewUIText("0 wisdom", assets.BodyFont, assets.ColorDudeWisdom),
		cowardice: NewUIText("0 cowardice", assets.BodyFont, assets.ColorDudeCowardice),
		luck:      NewUIText("0 luck", assets.BodyFont, assets.ColorDudeLuck),
		hidden:    false,
	}
	dip.panel.AddChild(dip.title)
	dip.title.ignoreScale = true
	//dip.panel.AddChild(dip.description)

	dip.panel.AddChild(dip.level)
	dip.level.ignoreScale = true
	dip.panel.AddChild(dip.xp)
	dip.xp.ignoreScale = true
	dip.panel.AddChild(dip.hp)
	dip.hp.ignoreScale = true
	dip.panel.AddChild(dip.strength)
	dip.strength.ignoreScale = true
	dip.panel.AddChild(dip.agility)
	dip.agility.ignoreScale = true
	dip.panel.AddChild(dip.defense)
	dip.defense.ignoreScale = true
	dip.panel.AddChild(dip.wisdom)
	dip.wisdom.ignoreScale = true
	dip.panel.AddChild(dip.cowardice)
	dip.cowardice.ignoreScale = true
	dip.panel.AddChild(dip.luck)
	dip.luck.ignoreScale = true

	dip.panel.sizeChildren = true
	//dip.panel.centerChildren = true
	return dip
}

func (dip *DudeInfoPanel) SetDude(dude *Dude) {
	dip.dude = dude
	if dude == nil {
		dip.hidden = true
		return
	}
	dip.hidden = false
	dip.SyncDude()
}

func (dip *DudeInfoPanel) SyncDude() {
	if dip.dude == nil {
		return
	}

	dip.title.SetText(dip.dude.Name())

	dip.level.SetText(fmt.Sprintf("Level %d %s", dip.dude.Level(), dip.dude.Profession()))
	dip.xp.SetText(fmt.Sprintf("%d/%d xp", dip.dude.XP(), dip.dude.NextLevelXP()))
	stats := dip.dude.GetCalculatedStats()
	dip.hp.SetText(fmt.Sprintf("%d/%d hp", stats.currentHp, stats.totalHp))
	dip.strength.SetText(fmt.Sprintf("%s strength", PaddedIntString(stats.strength, 4)))
	dip.agility.SetText(fmt.Sprintf("%s agility", PaddedIntString(stats.agility, 4)))
	dip.defense.SetText(fmt.Sprintf("%s defense", PaddedIntString(stats.defense, 4)))
	dip.wisdom.SetText(fmt.Sprintf("%s wisdom", PaddedIntString(stats.wisdom, 4)))
	dip.cowardice.SetText(fmt.Sprintf("%s cowardice", PaddedIntString(stats.cowardice, 4)))
	dip.luck.SetText(fmt.Sprintf("%s luck", PaddedIntString(stats.luck, 4)))
}

func (dip *DudeInfoPanel) Layout(o *UIOptions) {
	dip.panel.padding = 6 * o.Scale
	dip.panel.spaceBetween = -2 * o.Scale
	dip.panel.Layout(nil, o)
}

func (dip *DudeInfoPanel) Update(o *UIOptions) {
	dip.panel.Update(o)
	if dip.dude != nil {
		if dip.dude.dirtyStats {
			dip.SyncDude()
			dip.dude.dirtyStats = false
		}
	}
}

func (dip *DudeInfoPanel) Check(mx, my float64, kind UICheckKind) bool {
	if dip.hidden {
		return false
	}
	return dip.panel.Check(mx, my, kind)
}

func (dip *DudeInfoPanel) Draw(o *render.Options) {
	if dip.hidden {
		return
	}
	dip.panel.Draw(o)
}

type EquipmentPanel struct {
	panel       *UIPanel
	list        *UIItemList
	equipment   []*Equipment
	title       *UIText
	buyButton   *ButtonPanel
	showDetails bool

	details *EquipmentDetailsPanel

	onBuyClick  func()
	onItemClick func(index int)
	onItemHover func(index int)
}

func MakeEquipmentPanel() EquipmentPanel {
	ep := EquipmentPanel{
		panel:   NewUIPanel(PanelStyleInteractive),
		title:   NewUIText("Loot", assets.DisplayFont, assets.ColorHeading),
		list:    NewUIItemList(DirectionVertical),
		details: NewEquipmentDetailsPanel(),
	}
	btn := MakeButtonPanel(assets.BodyFont, PanelStyleButtonAttached)
	ep.buyButton = &btn
	ep.buyButton.text.center = true
	ep.buyButton.text.SetText("Buy\nRandom Loot")
	ep.list.spaceBetween = -2
	ep.panel.AddChild(ep.title)
	ep.panel.AddChild(ep.list)
	ep.panel.sizeChildren = true
	ep.panel.centerChildren = true

	return ep
}

func (ep *EquipmentPanel) SetEquipment(equipment []*Equipment) {
	ep.equipment = equipment
	ep.list.Clear()
	for index, eq := range equipment {
		txt := NewUIText(eq.Name(), assets.BodyFont, eq.quality.TextColor())
		txt.ignoreScale = true
		ep.list.AddItem(txt)

		txt.onCheck = func(kind UICheckKind) {
			if kind == UICheckClick && ep.onItemClick != nil {
				ep.onItemClick(index)
				ep.list.selected = index
			}
			if kind == UICheckHover && ep.onItemHover != nil {
				ep.onItemHover(index)
			}
		}
	}
}

func (ep *EquipmentPanel) Layout(o *UIOptions) {
	ep.panel.padding = 6 * o.Scale
	ep.list.SetSize(ep.panel.Width(), ep.panel.Height()-ep.panel.padding*2-ep.title.Height())
	ep.buyButton.SetSize(ep.panel.Width(), 48)
	ep.buyButton.Layout(nil, o)
	ep.buyButton.text.SetPosition(ep.buyButton.text.X(), ep.buyButton.text.Y()+4*o.Scale)
	ep.buyButton.SetPosition(ep.panel.X()+ep.panel.Width()/2-ep.buyButton.Width()/2, ep.panel.Y()+ep.panel.Height()-10*o.Scale)

	ep.panel.Layout(nil, o)
	ep.details.Layout(o)
	ep.details.panel.SetSize(128*o.Scale, 96*o.Scale)
	ep.details.panel.SetPosition(ep.panel.X()+ep.panel.Width()+6*o.Scale, ep.panel.Y())

	// Dynamically size our details panel.
	newHeight := (ep.details.luck.Y() + ep.details.luck.Height()) - ep.details.panel.Y()
	newHeight = math.Ceil(newHeight/ep.details.panel.center.Height()) * ep.details.panel.center.Height()
	ep.details.panel.SetSize(128*o.Scale, newHeight)
}

func (ep *EquipmentPanel) Update(o *UIOptions) {
	ep.panel.Update(o)
	if ep.showDetails {
		ep.details.Update(o)
	}
}

func (ep *EquipmentPanel) Check(mx, my float64, kind UICheckKind) bool {
	if ep.panel.Check(mx, my, kind) {
		return true
	}
	if ep.buyButton.Check(mx, my, kind) {
		if kind == UICheckClick {
			if ep.onBuyClick != nil {
				ep.onBuyClick()
			}
		}
		return true
	}
	if ep.showDetails && ep.details.Check(mx, my, kind) {
		return true
	}
	return false
}

func (ep *EquipmentPanel) Draw(o *render.Options) {
	ep.buyButton.Draw(o)
	ep.panel.Draw(o)
	if ep.showDetails {
		ep.details.Draw(o)
	}
}

type EquipmentDetailsPanel struct {
	panel           *UIPanel
	title           *UIText
	description     *UIText
	level           *UIText
	equipment       *Equipment
	perk            *UIText
	perkDescription *UIText
	uses            *UIText

	agility   *UIText
	strength  *UIText
	defense   *UIText
	wisdom    *UIText
	cowardice *UIText
	luck      *UIText

	swapButton  *ButtonPanel
	sellButton  *ButtonPanel
	onSwapClick func()
	onSellClick func()

	hidden bool
}

func NewEquipmentDetailsPanel() *EquipmentDetailsPanel {
	edp := &EquipmentDetailsPanel{
		panel:           NewUIPanel(PanelStyleTransparent),
		title:           NewUIText("My Steeze", assets.DisplayFont, assets.ColorHeading),
		description:     NewUIText("", assets.BodyFont, assets.ColorItemDescription),
		level:           NewUIText("", assets.BodyFont, assets.ColorItemLevel),
		perk:            NewUIText("", assets.BodyFont, assets.ColorItemPerk),
		perkDescription: NewUIText("", assets.BodyFont, assets.ColorItemPerkDescription),
		uses:            NewUIText("", assets.BodyFont, assets.ColorItemUses),
		agility:         NewUIText("", assets.BodyFont, assets.ColorDudeAgility),
		strength:        NewUIText("", assets.BodyFont, assets.ColorDudeStrength),
		defense:         NewUIText("", assets.BodyFont, assets.ColorDudeDefense),
		wisdom:          NewUIText("", assets.BodyFont, assets.ColorDudeWisdom),
		cowardice:       NewUIText("", assets.BodyFont, assets.ColorDudeCowardice),
		luck:            NewUIText("", assets.BodyFont, assets.ColorDudeLuck),
	}
	{
		btn := MakeButtonPanel(assets.BodyFont, PanelStyleButtonAttached)
		edp.sellButton = &btn
		edp.sellButton.text.center = true
		edp.sellButton.text.SetText("Sell for\ngp")
		edp.sellButton.onClick = func() {
			if edp.onSellClick != nil {
				edp.onSellClick()
			}
		}
	}
	{
		btn := MakeButtonPanel(assets.BodyFont, PanelStyleButtonAttached)
		edp.swapButton = &btn
		edp.swapButton.text.center = true
		edp.swapButton.text.SetText("Swap to\nDude")
		edp.swapButton.onClick = func() {
			if edp.onSwapClick != nil {
				edp.onSwapClick()
			}
		}
	}

	edp.title.ignoreScale = true
	edp.description.ignoreScale = true
	edp.level.ignoreScale = true
	edp.perk.ignoreScale = true
	edp.perkDescription.ignoreScale = true
	edp.uses.ignoreScale = true
	edp.agility.ignoreScale = true
	edp.strength.ignoreScale = true
	edp.defense.ignoreScale = true
	edp.wisdom.ignoreScale = true
	edp.cowardice.ignoreScale = true
	edp.luck.ignoreScale = true
	edp.panel.AddChild(edp.title)
	edp.panel.AddChild(edp.level)
	edp.panel.AddChild(edp.perk)
	edp.panel.AddChild(edp.uses)
	edp.panel.AddChild(edp.description)
	edp.panel.AddChild(edp.perkDescription)
	edp.panel.AddChild(edp.agility)
	edp.panel.AddChild(edp.strength)
	edp.panel.AddChild(edp.defense)
	edp.panel.AddChild(edp.wisdom)
	edp.panel.AddChild(edp.cowardice)
	edp.panel.AddChild(edp.luck)
	edp.panel.sizeChildren = true
	return edp
}

func (edp *EquipmentDetailsPanel) SetEquipment(equipment *Equipment) {
	edp.equipment = equipment
	if equipment != nil {
		edp.title.SetText(equipment.Name())
		edp.title.textOptions.Color = equipment.quality.TextColor()
		edp.description.SetText(equipment.Description())
		edp.level.SetText(fmt.Sprintf("Level %d %s", equipment.stats.level, equipment.Type()))

		edp.uses.SetText(fmt.Sprintf("%d/%d uses", equipment.uses, equipment.totalUses))

		if equipment.perk != nil {
			edp.perk.SetText(equipment.perk.Name())
			edp.perk.textOptions.Color = equipment.perk.Quality().TextColor()
			edp.perkDescription.SetText(equipment.perk.Description())
		} else {
			edp.perk.SetText("")
			edp.perkDescription.SetText("")
		}

		edp.agility.SetText(fmt.Sprintf("%s agility", PaddedIntString(equipment.stats.agility, 4)))
		edp.strength.SetText(fmt.Sprintf("%s strength", PaddedIntString(equipment.stats.strength, 4)))
		edp.defense.SetText(fmt.Sprintf("%s defense", PaddedIntString(equipment.stats.defense, 4)))
		edp.wisdom.SetText(fmt.Sprintf("%s wisdom", PaddedIntString(equipment.stats.wisdom, 4)))
		edp.cowardice.SetText(fmt.Sprintf("%s cowardice", PaddedIntString(equipment.stats.cowardice, 4)))
		edp.luck.SetText(fmt.Sprintf("%s luck", PaddedIntString(equipment.stats.luck, 4)))

		edp.sellButton.text.SetText(fmt.Sprintf("Sell for\n%.0fgp", equipment.GoldValue()))

		edp.hidden = false
	} else {
		edp.hidden = true
	}
}

func (edp *EquipmentDetailsPanel) Layout(o *UIOptions) {
	if edp.hidden {
		return
	}

	edp.panel.padding = 6 * o.Scale
	edp.panel.spaceBetween = -1 * o.Scale
	edp.panel.Layout(nil, o)

	edp.sellButton.SetSize(edp.panel.Width(), 48)
	edp.sellButton.Layout(nil, o)
	edp.sellButton.text.SetPosition(edp.sellButton.text.X(), edp.sellButton.text.Y()+4*o.Scale)
	edp.sellButton.SetPosition(edp.panel.X(), edp.panel.Y()+edp.panel.Height()-10*o.Scale)

	edp.swapButton.SetSize(edp.panel.Width(), 48)
	edp.swapButton.Layout(nil, o)
	edp.swapButton.text.SetPosition(edp.swapButton.text.X(), edp.swapButton.text.Y()+4*o.Scale)
	edp.swapButton.SetPosition(edp.panel.X()+edp.panel.Width()-edp.swapButton.Width(), edp.panel.Y()+edp.panel.Height()-10*o.Scale)
}

func (edp *EquipmentDetailsPanel) Update(o *UIOptions) {
	if edp.hidden {
		return
	}

	edp.panel.Update(o)
}

func (edp *EquipmentDetailsPanel) Check(mx, my float64, kind UICheckKind) bool {
	if edp.hidden {
		return false
	}
	if edp.panel.Check(mx, my, kind) {
		return true
	}
	if edp.sellButton.Check(mx, my, kind) {
		return true
	}
	if edp.swapButton.Check(mx, my, kind) {
		return true
	}
	return false
}

func (edp *EquipmentDetailsPanel) Draw(o *render.Options) {
	if edp.hidden {
		return
	}
	edp.sellButton.Draw(o)
	edp.swapButton.Draw(o)
	edp.panel.Draw(o)
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

func MakeButtonPanel(font assets.Font, style PanelStyle) ButtonPanel {
	bp := ButtonPanel{
		panel: NewUIPanel(style),
		text:  NewUIText("arghh", font, assets.ColorHeading),
	}
	bp.panel.AddChild(bp.text)
	bp.panel.sizeChildren = false
	bp.panel.centerChildren = true
	return bp
}

func (bp *ButtonPanel) Layout(parent UIElement, o *UIOptions) {
	bp.panel.padding = 1 * o.Scale
	bp.text.Layout(nil, o)
	bp.panel.Layout(nil, o)
	bp.doSize(o)
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

func (bp *ButtonPanel) X() float64 {
	return bp.panel.X()
}

func (bp *ButtonPanel) Y() float64 {
	return bp.panel.Y()
}

func (bp *ButtonPanel) Position() (float64, float64) {
	return bp.panel.Position()
}

func (bp *ButtonPanel) SetPosition(x, y float64) {
	bp.panel.SetPosition(x, y)
}

func (bp *ButtonPanel) Size() (float64, float64) {
	return bp.panel.Size()
}

func (bp *ButtonPanel) SetSize(w, h float64) {
	bp.panel.SetSize(w, h)
}

func (bp *ButtonPanel) Width() float64 {
	return bp.panel.Width()
}

func (bp *ButtonPanel) Height() float64 {
	return bp.panel.Height()
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
