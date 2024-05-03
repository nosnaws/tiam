package gmo

import (
	"math/rand"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

func genRandomName() string {
	rand.Seed(time.Now().UnixNano())

	return petname.Generate(3, "-")
}
