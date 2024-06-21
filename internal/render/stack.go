package render

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ebijam24/assets"
)

type Stack struct {
	Positionable
	Rotateable
	Originable
	data             *assets.Staxie // Reference to the underlying stack data for subimages, etc.
	currentStack     *assets.StaxieStack
	currentAnimation *assets.StaxieAnimation
	currentFrame     *assets.StaxieFrame
	frameCounter     int
	MaxSliceIndex    int
	SliceOffset      int
	HeightOffset     float64
}

func NewStack(name string, stackName string, animationName string) (*Stack, error) {
	staxie, err := assets.LoadStaxie(name)
	if err != nil {
		return nil, err
	}
	if stackName == "" {
		for k := range staxie.Stacks {
			stackName = k
			break
		}
	}

	stack, ok := staxie.Stacks[stackName]
	if !ok {
		return nil, fmt.Errorf("stack %s does not exist in %s", stackName, name)
	}

	if animationName == "" {
		for k := range stack.Animations {
			animationName = k
			break
		}
	}
	animation, ok := stack.Animations[animationName]
	if !ok {
		return nil, fmt.Errorf("animation %s does not exist in %s", animationName, stackName)
	}

	frame, ok := animation.GetFrame(0)
	if !ok {
		return nil, fmt.Errorf("frame 0 does not exist in %s", animationName)
	}

	return &Stack{data: staxie, currentStack: &stack, currentAnimation: &animation, currentFrame: frame}, nil
}

func (s *Stack) Draw(o *Options) {
	if s.currentFrame == nil {
		return
	}

	opts := ebiten.DrawImageOptions{}

	// Rotate about origin.
	ox, oy := s.Origin()
	opts.GeoM.Translate(-ox, -oy)
	opts.GeoM.Rotate(s.Rotation())
	opts.GeoM.Translate(ox, oy)

	// Translate to position.
	opts.GeoM.Translate(s.Position())

	// Add additional transforms.
	opts.GeoM.Concat(o.DrawImageOptions.GeoM)

	opts.GeoM.Translate(0, s.HeightOffset)

	for index := 0; index < len(s.currentFrame.Slices); index++ {
		if index+s.SliceOffset >= len(s.currentFrame.Slices) {
			break
		}
		if s.MaxSliceIndex != 0 && index >= s.MaxSliceIndex {
			break
		}
		slice := s.currentFrame.Slices[index+s.SliceOffset]
		i := index

		// TODO: Make this configurable
		c := float64(index) / float64(len(s.currentFrame.Slices))
		c = math.Min(1.0, math.Max(0.5, c))
		color := float32(c)

		opts.ColorScale.Reset()
		opts.ColorScale.Scale(color, color, color, 1.0)

		if o.VGroup != nil {
			o.VGroup.Images[i].DrawImage(slice.Image, &opts)
		} else if o.Screen != nil {
			o.Screen.DrawImage(slice.Image, &opts)
			opts.GeoM.Translate(0, -o.Pitch)
		}
		//opts.GeoM.Skew(-0.002, 0.002) // Might be able to sine this with delta to create a wave effect...
	}
}

func (s *Stack) Update() {
	s.frameCounter++
	if s.frameCounter >= int(s.currentAnimation.Frametime) {
		s.frameCounter = 0
		nextFrame, ok := s.currentAnimation.GetFrame(s.currentFrame.Index + 1)
		if !ok {
			nextFrame, _ = s.currentAnimation.GetFrame(0)
		}
		s.currentFrame = nextFrame
	}
}

func (s *Stack) SliceCount() int {
	return len(s.currentFrame.Slices)
}

func (s *Stack) SetStack(name string) error {
	stack, ok := s.data.Stacks[name]
	if !ok {
		return fmt.Errorf("stack %s", name)
	}
	s.currentStack = &stack

	return s.SetAnimation(s.currentAnimation.Name)
}

func (s *Stack) SetAnimation(name string) error {
	animation, ok := s.currentStack.GetAnimation(name)
	if !ok {
		return fmt.Errorf("animation %s", name)
	}
	s.currentAnimation = &animation

	return s.SetFrame(0)
}

func (s *Stack) SetFrame(index int) error {
	frame, ok := s.currentAnimation.GetFrame(index)
	if !ok {
		return fmt.Errorf("frame %d", index)
	}
	s.currentFrame = frame
	return nil
}

func (s *Stack) SetOriginToCenter() {
	s.SetOrigin(float64(s.currentFrame.Slices[0].Image.Bounds().Dx())/2, float64(s.currentFrame.Slices[0].Image.Bounds().Dy())/2)
}
