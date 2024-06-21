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
