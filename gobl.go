package main

import (
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	runArgs := append([]interface{}{}, "./game"+exe)

	Task("build").
		Exec("go", "build", "./cmd/game")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("cmd/game/*", "internal/game/*", "internal/render/*", "assets/*", "assets/walls/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Go()
}
