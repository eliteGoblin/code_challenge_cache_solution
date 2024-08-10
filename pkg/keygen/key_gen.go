package keygen

import (
	"fmt"
	"math/rand"
)

func RandomKey(min, max int) string {
	n := rand.Intn(max-min) + min
	return fmt.Sprintf("key%d", n)
}
