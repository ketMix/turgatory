package render

type Stacks []*Stack

func (s *Stacks) Draw(o *Options) {
	for _, stack := range *s {
		stack.Draw(o)
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
