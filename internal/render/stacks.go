package render

type Stacks []*Stack

func (s *Stacks) Draw(o *Options) {
	for _, stack := range *s {
		stack.Draw(o)
	}
}

func (s *Stacks) Update() {
	for _, stack := range *s {
		stack.Update()
	}
}

func (s *Stacks) Add(stack *Stack) {
	*s = append(*s, stack)
}

func (s *Stacks) Remove(stack *Stack) {
	for i, v := range *s {
		if v == stack {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

func (s *Stacks) SetRotations(r float64) {
	for _, stack := range *s {
		stack.SetRotation(r)
	}
}

func (s *Stacks) SetPositions(x, y float64) {
	for _, stack := range *s {
		stack.SetPosition(x, y)
	}
}

func (s *Stacks) SetOrigins(x, y float64) {
	for _, stack := range *s {
		stack.SetOrigin(x, y)
	}
}
