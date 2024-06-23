package game

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IntToFloat2(a, b int) (float64, float64) {
	return float64(a), float64(b)
}

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
