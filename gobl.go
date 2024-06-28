package main

import (
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	runArgs := append([]interface{}{}, "./game"+exe)
	var wasmSrc string

	Task("build").
		Exec("go", "build", "./cmd/game")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("cmd/game/*", "internal/game/*", "internal/render/*", "assets/*", "assets/walls/*", "assets/floors/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Task("build-web").
		Env("GOOS=js", "GOARCH=wasm").
		Exec("go", "build", "-o", "web/game.wasm", "./cmd/game").
		Exec("go", "env", "GOROOT").
		Result(func(i interface{}) {
			goRoot := strings.TrimSpace(i.(string))
			wasmSrc = filepath.Join(goRoot, "misc/wasm/wasm_exec.js")
		}).
		Exec("cp", &wasmSrc, "web/")
	Go()
}
