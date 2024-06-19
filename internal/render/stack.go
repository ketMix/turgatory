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
	data             *assets.Staxie // Reference to the underlying stack data for subimages, etc.
	currentAnimation *assets.StaxieAnimation
	currentFrame     *assets.StaxieFrame
	frameCounter     int
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

	return &Stack{data: staxie, currentAnimation: &animation, currentFrame: frame}, nil
}

func (s *Stack) Draw(o Options) {
	if s.currentFrame == nil {
		return
	}

	opts := ebiten.DrawImageOptions{}
	// Rotate about center.
	hw := math.Round(float64(s.data.FrameWidth) / 2)
	hh := math.Round(float64(s.data.FrameHeight) / 2)
	opts.GeoM.Scale(2, 2)
	opts.GeoM.Translate(-hw, -hh)
	opts.GeoM.Rotate(s.rotation)
	opts.GeoM.Translate(hw, hh)
	// Position.
	opts.GeoM.Translate(float64(s.x), float64(s.y))
	// Uh... this might come before? FIXME later
	opts.GeoM.Concat(o.DrawImageOptions.GeoM)
	// Draw our slices from!
	for _, slice := range s.currentFrame.Slices {
		o.Screen.DrawImage(slice.Image, &opts)
		opts.GeoM.Translate(0, -1)
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
