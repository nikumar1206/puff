package utils

import (
	"fmt"
	"math/rand/v2"
)

func RandomNanoID() string {
	id := ""
	for range 4 {
		r := rand.IntN(25) + 1
		id += fmt.Sprintf("%c", ('A' - 1 + r))
	}
	id += "-"
	for range 4 {
		r := rand.IntN(9)
		id += fmt.Sprint(r)
	}
	return id
}
