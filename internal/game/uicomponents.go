package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ebijam24/assets"
	"github.com/kettek/ebijam24/internal/render"
)

type UIElement interface {
	Position() (float64, float64)
	SetPosition(float64, float64)
	Size() (float64, float64)
	SetSize(float64, float64)
	X() float64
	Y() float64
	Width() float64
	Height() float64
	Layout(parent UIElement, o *UIOptions)
	Draw(o *render.Options)
	Update(o *UIOptions)
	Check(float64, float64) bool
}

// ======== BUTTON ========
type UIButton struct {
	render.Positionable
	render.Sizeable
	noBackdrop  bool
	baseSprite  *render.Sprite
	sprite      *render.Sprite
	onClick     func()
	wobbler     float64
	tooltip     string
	showTooltip bool
}

func NewUIButton(name string, tooltip string) *UIButton {
	return &UIButton{
		baseSprite: Must(render.NewSpriteFromStaxie("ui/button", "base")),
		sprite:     Must(render.NewSpriteFromStaxie("ui/button", name)),
		tooltip:    tooltip,
	}
}

func (b *UIButton) Layout(parent UIElement, o *UIOptions) {
	b.baseSprite.Scale = o.Scale
	b.sprite.Scale = o.Scale
}

func (b *UIButton) Update(o *UIOptions) {
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

func (b *UIButton) Check(mx, my float64) bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := b.Position()
		w, h := b.sprite.Size()
		if mx > x && mx < x+w && my > y && my < y+h {
			if b.onClick != nil {
				b.onClick()
				return true
			}
		}
	}
	return false
}

func (b *UIButton) SetPosition(x, y float64) {
	b.Positionable.SetPosition(x, y)
	b.baseSprite.SetPosition(x, y)
	b.sprite.SetPosition(x, y)
}

func (b *UIButton) Size() (float64, float64) {
	return b.baseSprite.Size()
}

func (b *UIButton) Width() float64 {
	w, _ := b.baseSprite.Size()
	return w
}

func (b *UIButton) Height() float64 {
	_, h := b.baseSprite.Size()
	return h
}

func (b *UIButton) SetImage(name string) {
	b.sprite.SetStaxie("ui/button", name)
}

func (b *UIButton) Draw(o *render.Options) {
	if !b.noBackdrop {
		b.baseSprite.Draw(o)
	}
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

// ======== UIItemList ========
const (
	DirectionVertical = iota
	DirectionHorizontal
)

type UIItemList struct {
	render.Positionable
	render.Sizeable
	items             []UIElement
	selected          int
	itemOffset        int
	direction         int
	lastVisibleIndex  int
	itemsAllVisible   bool
	changed           bool
	centerItems       bool // Center items on the opposite axis to direction
	centerList        bool // Center the items to be visually centered
	decrementUIButton *UIButton
	incrementUIButton *UIButton
}

func NewUIItemList(direction int) *UIItemList {
	l := &UIItemList{
		direction:   direction,
		centerItems: true,
	}

	l.decrementUIButton = NewUIButton("arrow", "")
	l.decrementUIButton.noBackdrop = true
	l.incrementUIButton = NewUIButton("arrow", "")
	l.incrementUIButton.noBackdrop = true

	if direction == DirectionVertical {
		l.decrementUIButton.sprite.SetStaxieAnimation("ui/button", "arrow", "up")
		l.incrementUIButton.sprite.SetStaxieAnimation("ui/button", "arrow", "down")
	} else {
		l.decrementUIButton.sprite.SetStaxieAnimation("ui/button", "arrow", "left")
		l.incrementUIButton.sprite.SetStaxieAnimation("ui/button", "arrow", "right")
	}

	l.decrementUIButton.onClick = func() {
		if l.itemOffset > 0 {
			l.itemOffset--
		}
		l.changed = true
	}
	l.incrementUIButton.onClick = func() {
		if l.itemsAllVisible {
			return
		}
		if l.itemOffset < len(l.items)-1 {
			l.itemOffset++
		}
		l.changed = true
	}

	return l
}
func (l *UIItemList) adjustButtons() {
	if l.itemOffset == 0 {
		l.decrementUIButton.sprite.Transparency = 0.25
	} else {
		l.decrementUIButton.sprite.Transparency = 0
	}
	if l.itemsAllVisible {
		l.incrementUIButton.sprite.Transparency = 0.25
	} else {
		l.incrementUIButton.sprite.Transparency = 0
	}
}
func (l *UIItemList) Layout(parent UIElement, o *UIOptions) {
	l.decrementUIButton.Layout(l, o)
	l.incrementUIButton.Layout(l, o)

	l.changed = true

	w, h := l.Size()
	bw, bh := l.decrementUIButton.Size()
	if l.direction == DirectionVertical {
		l.decrementUIButton.SetPosition(l.X()+w/2-bw/2, l.Y()-bw/2+4)
		l.incrementUIButton.SetPosition(l.X()+w/2-bw/2, l.Y()+l.Height()-bh/2-4)
	} else {
		l.decrementUIButton.SetPosition(l.X()-bw/2+4, l.Y()+h/2-bh/2)
		l.incrementUIButton.SetPosition(l.X()+l.Width()-bw/2-4, l.Y()+h/2-bh/2)
	}
}
func (l *UIItemList) Update(o *UIOptions) {
	l.decrementUIButton.Update(o)
	l.incrementUIButton.Update(o)

	if l.changed {
		v := 0.0
		maxSize := 0.0
		itemsSize := 0.0
		if l.direction == DirectionVertical {
			maxSize = l.Height() - l.decrementUIButton.Height()/2
		} else {
			maxSize = l.Width() - l.decrementUIButton.Width()/2
		}
		l.lastVisibleIndex = len(l.items)
		for i := l.itemOffset; i < len(l.items); i++ {
			if l.direction == DirectionVertical {
				cs := 0.0
				if l.centerItems {
					cs = l.Width()/2 - l.items[i].Width()/2
				}
				l.items[i].SetPosition(l.X()+cs, l.Y()+float64(v))
			} else {
				cs := 0.0
				if l.centerItems {
					cs = l.Height()/2 - l.items[i].Height()/2
				}
				l.items[i].SetPosition(l.X()+float64(v), l.Y()+cs)
			}
			l.items[i].Layout(l, o)
			itemsSize = v

			if l.direction == DirectionVertical {
				v += l.items[i].Height()
			} else {
				v += l.items[i].Width()
			}

			if v >= maxSize {
				l.lastVisibleIndex = i
				break
			}
		}
		if l.lastVisibleIndex == len(l.items) {
			l.itemsAllVisible = true
		} else {
			l.itemsAllVisible = false
		}

		// Yeah, this isn't great, but I want items centered.
		for i := l.itemOffset; i < l.lastVisibleIndex; i++ {
			if l.centerList {
				if l.direction == DirectionVertical {
					offset := maxSize - itemsSize
					l.items[i].SetPosition(l.items[i].X(), l.items[i].Y()+offset/2)
				} else {
					offset := maxSize - itemsSize
					l.items[i].SetPosition(l.items[i].X()+offset/2, l.items[i].Y())
				}
			}
		}

		l.adjustButtons()

		l.changed = false
	}

	for i := l.itemOffset; i < l.lastVisibleIndex; i++ {
		l.items[i].Update(o)
	}
}
func (l *UIItemList) Check(mx, my float64) bool {
	if l.decrementUIButton.Check(mx, my) {
		return true
	}
	if l.incrementUIButton.Check(mx, my) {
		return true
	}
	for i := l.itemOffset; i < l.lastVisibleIndex; i++ {
		if l.items[i].Check(mx, my) {
			return true
		}
	}
	return false
}
func (l *UIItemList) Draw(o *render.Options) {
	l.decrementUIButton.Draw(o)
	l.incrementUIButton.Draw(o)
	o.DrawImageOptions.GeoM.Reset()
	for i := l.itemOffset; i < l.lastVisibleIndex; i++ {
		l.items[i].Draw(o)
	}
}
func (l *UIItemList) AddItem(item UIElement) {
	l.items = append(l.items, item)
	l.changed = true
}

// ======== UIPanel ========
type UIPanel struct {
	render.Positionable
	render.Sizeable

	children []UIElement

	padding       float64
	flowDirection int
	sizeChildren  bool

	top         *render.Sprite
	bottom      *render.Sprite
	left        *render.Sprite
	right       *render.Sprite
	topleft     *render.Sprite
	topright    *render.Sprite
	bottomleft  *render.Sprite
	bottomright *render.Sprite
	center      *render.Sprite
}

func NewUIPanel() *UIPanel {
	sp := Must(render.NewSprite("ui/panels"))
	return &UIPanel{
		topleft:     Must(render.NewSubSprite(sp, 0, 0, 16, 16)),
		top:         Must(render.NewSubSprite(sp, 16, 0, 16, 16)),
		topright:    Must(render.NewSubSprite(sp, 32, 0, 16, 16)),
		left:        Must(render.NewSubSprite(sp, 0, 16, 16, 16)),
		center:      Must(render.NewSubSprite(sp, 16, 16, 16, 16)),
		right:       Must(render.NewSubSprite(sp, 32, 16, 16, 16)),
		bottomleft:  Must(render.NewSubSprite(sp, 0, 32, 16, 16)),
		bottom:      Must(render.NewSubSprite(sp, 16, 32, 16, 16)),
		bottomright: Must(render.NewSubSprite(sp, 32, 32, 16, 16)),
	}
}

func (p *UIPanel) Layout(parent UIElement, o *UIOptions) {
	// Grosse
	p.topleft.Scale = o.Scale
	p.top.Scale = o.Scale
	p.topright.Scale = o.Scale
	p.left.Scale = o.Scale
	p.center.Scale = o.Scale
	p.right.Scale = o.Scale
	p.bottomleft.Scale = o.Scale
	p.bottom.Scale = o.Scale
	p.bottomright.Scale = o.Scale

	x := p.X() + p.padding
	y := p.Y() + p.padding
	for _, child := range p.children {
		child.SetPosition(x, y)

		if p.sizeChildren {
			if p.flowDirection == DirectionVertical {
				child.SetSize(p.Width()-p.padding*2, child.Height())
			} else {
				child.SetSize(child.Width(), p.Height()-p.padding*2)
			}
		}

		child.Layout(p, o)
		if p.flowDirection == DirectionVertical {
			y += child.Height()
		} else {
			x += child.Width()
		}
	}
}

func (p *UIPanel) Update(o *UIOptions) {
	for _, child := range p.children {
		child.Update(o)
	}
}

func (p *UIPanel) Check(mx, my float64) bool {
	if !InBounds(p.X(), p.Y(), p.Width(), p.Height(), mx, my) {
		return false
	}
	for _, child := range p.children {
		if child.Check(mx, my) {
			return true
		}
	}
	return false
}

func (p *UIPanel) Draw(o *render.Options) {
	x, y := p.Position()
	w, h := p.Size()

	op := &render.Options{
		Screen: o.Screen,
	}
	op.DrawImageOptions.GeoM.Concat(o.DrawImageOptions.GeoM)

	op.DrawImageOptions.GeoM.Translate(x, y)

	geom := ebiten.GeoM{}
	geom.Concat(op.DrawImageOptions.GeoM)

	// Draw corners
	p.topleft.Draw(op)
	op.DrawImageOptions.GeoM.Translate(w-p.topright.Width(), 0)
	p.topright.Draw(op)
	op.DrawImageOptions.GeoM.Translate(0, h-p.bottomright.Height())
	p.bottomright.Draw(op)
	op.DrawImageOptions.GeoM.Translate(-(w - p.bottomleft.Width()), 0)
	p.bottomleft.Draw(op)

	op.DrawImageOptions.GeoM.Reset()
	op.DrawImageOptions.GeoM.Concat(geom)
	// Draw sides
	op.DrawImageOptions.GeoM.Translate(p.topleft.Width(), 0)
	for i := 0; i < int(w-p.topleft.Width()-p.topright.Width()); i += int(p.top.Width()) {
		p.top.Draw(op)
		op.DrawImageOptions.GeoM.Translate(p.top.Width(), 0)
	}
	op.DrawImageOptions.GeoM.Translate(0, p.topright.Height())
	for i := 0; i < int(h-p.topright.Height()-p.bottomright.Height()); i += int(p.right.Height()) {
		p.right.Draw(op)
		op.DrawImageOptions.GeoM.Translate(0, p.right.Height())
	}
	op.DrawImageOptions.GeoM.Reset()
	op.DrawImageOptions.GeoM.Concat(geom)
	for i := 0; i < int(h-p.bottomleft.Height()-p.topleft.Height()); i += int(p.left.Height()) {
		op.DrawImageOptions.GeoM.Translate(0, p.left.Height())
		p.left.Draw(op)
	}
	op.DrawImageOptions.GeoM.Reset()
	op.DrawImageOptions.GeoM.Concat(geom)
	op.DrawImageOptions.GeoM.Translate(p.bottomleft.Width(), h-p.bottomleft.Height())
	for i := 0; i < int(w-p.bottomright.Width()-p.bottomleft.Width()); i += int(p.bottom.Width()) {
		p.bottom.Draw(op)
		op.DrawImageOptions.GeoM.Translate(p.bottom.Width(), 0)
	}

	// Draw center.
	op.DrawImageOptions.GeoM.Reset()
	op.DrawImageOptions.GeoM.Concat(geom)
	op.DrawImageOptions.GeoM.Translate(p.topleft.Width(), p.topleft.Height())
	maxWidth := w - p.topleft.Width() - p.topright.Width()
	maxHeight := h - p.topleft.Height() - p.bottomleft.Height()
	for y := 0; y < int(maxHeight); y += int(p.center.Height()) {
		for x := 0; x < int(maxWidth); x += int(p.center.Width()) {
			p.center.Draw(op)
			op.DrawImageOptions.GeoM.Translate(p.center.Width(), 0)
		}
		op.DrawImageOptions.GeoM.Translate(-maxWidth, p.center.Height())
	}

	op.DrawImageOptions.GeoM.Reset()
	// Draw children.
	for _, child := range p.children {
		child.Draw(op)
	}
}

// ======== UIText ========
type UIText struct {
	render.Positionable
	render.Sizeable

	text        string
	textWidth   float64
	textHeight  float64
	textScale   float64
	scale       float64
	textOptions render.TextOptions
}

func NewUIText(txt string, font assets.Font, color color.Color) *UIText {
	t := &UIText{
		text: txt,
		textOptions: render.TextOptions{
			Font:  font,
			Color: color,
		},
		textScale: 1,
	}

	w, h := text.Measure(txt, font.Face, font.LineHeight)
	t.textWidth = float64(w)
	t.textHeight = float64(h)

	return t
}

func (t *UIText) Layout(parent UIElement, o *UIOptions) {
	t.scale = o.Scale * t.textScale
	t.SetSize(t.textWidth*t.scale, t.textHeight*t.scale)
}

func (t *UIText) Update(o *UIOptions) {
}

func (t *UIText) Check(mx, my float64) bool {
	return false
}

func (t *UIText) Draw(o *render.Options) {
	t.textOptions.Screen = o.Screen
	t.textOptions.GeoM.Reset()
	t.textOptions.GeoM.Scale(t.scale, t.scale)
	t.textOptions.GeoM.Translate(t.X(), t.Y())
	render.DrawText(&t.textOptions, t.text)
}

func (t *UIText) SetScale(scale float64) {
	t.textScale = scale
}

type UIImage struct {
	render.Positionable
	render.Sizeable

	scale       float64
	finalScale  float64
	image       *ebiten.Image
	ignoreScale bool
	onClick     func()
}

func NewUIImage(img *ebiten.Image) *UIImage {
	return &UIImage{
		image: img,
		scale: 1,
	}
}

func (i *UIImage) Layout(parent UIElement, o *UIOptions) {
	if i.ignoreScale {
		i.finalScale = 1
	} else {
		i.finalScale = o.Scale * i.scale
	}
	i.SetSize(float64(i.image.Bounds().Dx())*i.finalScale, float64(i.image.Bounds().Dy())*i.finalScale)
}

func (i *UIImage) Update(o *UIOptions) {
}

func (i *UIImage) Check(mx, my float64) bool {
	if i.onClick != nil {
		if InBounds(i.X(), i.Y(), i.Width(), i.Height(), mx, my) {
			i.onClick()
			return true
		}
	}
	return false
}

func (i *UIImage) Draw(o *render.Options) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(i.finalScale, i.finalScale)
	op.GeoM.Translate(i.X(), i.Y())
	o.Screen.DrawImage(i.image, op)
}