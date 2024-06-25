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
