package guid

import "math/rand"

func GetRandInt32() int32 {
	return rand.Int31()
}

func getRandInt8() int8  {
	return int8(rand.Int31n(256))
}