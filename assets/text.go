package assets

import (
	"bytes"
	"math/rand"
	"strings"
)

var dudeNames []string
var hints []string

func init() {
	// Load names
	b, err := FS.ReadFile("dudes/names.txt")
	if err != nil {
		panic(err)
	}

	names := bytes.Split(b, []byte("\n"))
	for _, name := range names {
		if len(name) == 0 || name[0] == '#' {
			continue
		}
		n := strings.TrimSpace(string(name)) // This is necessary for line differences on Windows.

		dudeNames = append(dudeNames, n)
	}

	// Load hints
	b, err = FS.ReadFile("ui/hints.txt")
	if err != nil {
		panic(err)
	}

	hintText := bytes.Split(b, []byte("\n"))
	hints = []string{}
	for _, hint := range hintText {
		if len(hint) == 0 || hint[0] == '#' {
			continue
		}
		n := strings.TrimSpace(string(hint)) // This is necessary for line differences on Windows.

		hints = append(hints, n)
	}
}

func GetRandomName() string {
	return dudeNames[rand.Intn(len(dudeNames))]
}

func GetHints() []string {
	return hints
}
