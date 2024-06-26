package assets

import (
	"bytes"
	"math/rand"
	"strings"
)

var dudeNames []string

func init() {
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
}

func GetRandomName() string {
	return dudeNames[rand.Intn(len(dudeNames))]
}
