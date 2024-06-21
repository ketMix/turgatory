package render

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// VGroup manages a slice of images intended to be rendered at offsets. This is basically an offscreen framebuffer for stack slice rendering.
type VGroup struct {
	Positionable
	Images []*ebiten.Image
	Width  int
	Height int
	Debug  bool
}

// NewVGroup creates a new VGroup. Destroy() _MUST_ be called once the VGroup is no longer needed.
func NewVGroup(w, h, n int) *VGroup {
	vg := &VGroup{
		Width:  w,
		Height: h,
	}

	for i := 0; i < n; i++ {
		img := ebiten.NewImageWithOptions(image.Rect(0, 0, w, h), &ebiten.NewImageOptions{
			Unmanaged: true,
		})
		img.Clear() // iirc, this is needed to prevent garbage contents on certain platforms/gpus
		vg.Images = append(vg.Images, img)
	}

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

	// TODO: We could actually do some matrix math here to "tilt" the rendered layers "away" from the camera, which would enhance the 3D look. Shame I'm bad at math.

	for _, img := range vg.Images {
		//img.Fill(color.NRGBA{255, 0, 0, 255})
		if vg.Debug {
			ebitenutil.DrawLine(img, 1, 1, float64(vg.Width), 1, color.NRGBA{255, 0, 255, 64})
			ebitenutil.DrawLine(img, 1, 1, 1, float64(vg.Height), color.NRGBA{255, 0, 255, 64})
			ebitenutil.DrawLine(img, float64(vg.Width), 1, float64(vg.Width), float64(vg.Height), color.NRGBA{255, 0, 255, 64})
			ebitenutil.DrawLine(img, 1, float64(vg.Height), float64(vg.Width), float64(vg.Height), color.NRGBA{255, 0, 255, 64})
		}
		// lol, this might be okay...
		w, h := img.Bounds().Dx(), img.Bounds().Dy()
		for i := 0; i < h; i++ {
			opts2 := ebiten.DrawImageOptions{}

			opts2.GeoM.Concat(opts.GeoM)
			opts2.GeoM.Translate(0, float64(h)) // It seems okay to shunt it down like this..
			opts2.GeoM.Translate(0, float64(i))

			o.Screen.DrawImage(img.SubImage(image.Rect(0, i, w, i+1)).(*ebiten.Image), &opts2)
		}
		//o.Screen.DrawImage(img, &opts)
		opts.GeoM.Translate(0, -o.Pitch)
	}
}
