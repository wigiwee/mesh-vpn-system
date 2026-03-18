package helper

import (
	"math/rand/v2"
	"strconv"
)

func GenerateRandomIPAddr() string {

	ipAddr := "100."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255)) + "."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255)) + "."
	ipAddr = ipAddr + strconv.Itoa(rand.IntN(255))
	return ipAddr
}
