package render

type InterpNumber struct {
	Current, Target, Speed float64
}

func (i *InterpNumber) Set(target, speed float64) {
	i.Target = target
	i.Speed = speed
}

func (i *InterpNumber) Update() bool {
	if i.Current < i.Target {
		i.Current += i.Speed
		if i.Current > i.Target {
			i.Current = i.Target
			return true
		}
	} else if i.Current > i.Target {
		i.Current -= i.Speed
		if i.Current < i.Target {
			i.Current = i.Target
			return true
		}
	}
	return false
}
