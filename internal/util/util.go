package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GenID(size int) (string, error) {
	c := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, c[:]); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	return hex.EncodeToString(c[:]), nil
}
