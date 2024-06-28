package render

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// VGroup manages a slice of images intended to be rendered at offsets. This is basically an offscreen framebuffer for stack slice rendering.
type VGroup struct {
	Positionable
	Images []*ebiten.Image
	Depth  int
	Width  int
	Height int
	Debug  bool
}

// NewVGroup creates a new VGroup. Destroy() _MUST_ be called once the VGroup is no longer needed.
func NewVGroup(w, h, n int) *VGroup {
	vg := &VGroup{
		Width:  w,
		Height: h,
		Depth:  n,
	}

	img := ebiten.NewImageWithOptions(image.Rect(0, 0, w, h*n), &ebiten.NewImageOptions{
		Unmanaged: true,
	})
	vg.Images = append(vg.Images, img)
	/*for i := 0; i < n; i++ {
		img := ebiten.NewImageWithOptions(image.Rect(0, 0, w, h), &ebiten.NewImageOptions{
			Unmanaged: true,
		})
		img.Clear() // iirc, this is needed to prevent garbage contents on certain platforms/gpus
		vg.Images = append(vg.Images, img)
	}*/

	return vg
}

// Destroy deallocates the underlying images. This _MUST_ be called.
func (vg *VGroup) Destroy() {
	for _, img := range vg.Images {
		img.Deallocate()
	}
	vg.Images = nil
}

// Clear clears the internal images.
func (vg *VGroup) Clear() {
	for _, img := range vg.Images {
		img.Clear()
	}
}

// Draw draws the internal images to the provided screen, applying geom and otherwise.
func (vg *VGroup) Draw(o *Options) {
	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(vg.Position())

	// Render from center?
	opts.GeoM.Translate(-float64(vg.Width)/2, -float64(vg.Height)/2)

	opts.GeoM.Concat(o.DrawImageOptions.GeoM)

	// Do not render if we're out of screen bounds.
	y := opts.GeoM.Element(1, 2)
	sy := opts.GeoM.Element(1, 1)
	h := float64(vg.Height) * sy * 1.25

	if y+h < 0 || y-h*.25 > float64(o.Screen.Bounds().Dy()) {
		return
	}

	//for _, img := range vg.Images {
	img := vg.Images[0]
	for index := 0; index < vg.Depth; index++ {
		// lol, this might be okay...
		w, h := vg.Width, vg.Height
		for i := 0; i < h; i++ {
			opts2 := ebiten.DrawImageOptions{}
			opts2.GeoM.Concat(opts.GeoM)
			opts2.GeoM.Translate(0, float64(h)) // It seems okay to shunt it down like this..
			opts2.GeoM.Translate(0, float64(i))

			o.Screen.DrawImage(img.SubImage(image.Rect(0, i+index*h, w, i+1+index*h)).(*ebiten.Image), &opts2)
		}
		//o.Screen.DrawImage(img, &opts)
		opts.GeoM.Translate(0, -o.Pitch)
	}
	{
		opts2 := ebiten.DrawImageOptions{}
		opts2.GeoM.Concat(opts.GeoM)
		opts2.GeoM.Translate(0, float64(vg.Height))
	}
}
