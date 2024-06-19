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

func (s *Stack) Draw(o *Options) {
	if s.currentFrame == nil {
		return
	}

	opts := ebiten.DrawImageOptions{}

	rotation := s.rotation

	// calculate the new position based on the rotation
	dx := math.Cos(rotation) * s.rotationDistance
	dy := math.Sin(rotation) * s.rotationDistance

	// apply the rotation and rotation offset for the stack
	// rotation offset potentially unique to stack, hardcoded to pie for now
	// account for rotation distance
	opts.GeoM.Rotate(rotation)

	// translate the stack to the new position
	screen := o.Screen
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	centerX, centerY := float64(screenWidth/2), float64(screenHeight/2)
	opts.GeoM.Translate(centerX+dx, centerY+dy)

	// Uh... this might come before? FIXME later
	opts.GeoM.Concat(o.DrawImageOptions.GeoM)
	// Draw our slices from!
	for _, slice := range s.currentFrame.Slices {
		o.Screen.DrawImage(slice.Image, &opts)
		opts.GeoM.Translate(0, -o.Camera.Pitch)
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
