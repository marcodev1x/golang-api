package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
)

func GenerateHash(id int) string {
	var randomInput int

	if id == 0 {
		randomInput = rand.Intn(10000000)
	} else {
		randomInput = id
	}

	hashBytes := sha256.Sum256([]byte(strconv.Itoa(randomInput)))
	hash := base64.URLEncoding.EncodeToString(hashBytes[:])

	return hash
}
