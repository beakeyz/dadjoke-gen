package crypto

import (
	"crypto/sha1"
	"fmt"
)

// central function for hashing shit so we are sure to use the same thing every time
func HashString (str string) string {
  hasher := sha1.New()
  hasher.Write([]byte(str))
  funnie := fmt.Sprintf("%x", hasher.Sum(nil))

  return funnie
}
