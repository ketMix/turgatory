package assets

import (
	"bytes"
	"math/rand"
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
		dudeNames = append(dudeNames, string(name))
	}
}

func GetRandomName() string {
	return dudeNames[rand.Intn(len(dudeNames))]
}
