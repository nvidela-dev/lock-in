package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func RandomID(prefix string) string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes[:]))
	}
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func SystemClock() time.Time {
	return time.Now().UTC()
}
