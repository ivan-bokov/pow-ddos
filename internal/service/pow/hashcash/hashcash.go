package hashcash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type HashCash struct {
	r             io.Reader
	sizeChallenge int
}

func New(size int) *HashCash {
	return &HashCash{
		r:             rand.Reader,
		sizeChallenge: size,
	}
}
func NewSolver() *HashCash {
	return &HashCash{}
}

func (_ *HashCash) Verify(challenge []byte, nonce uint64, difficulty uint8) bool {
	target := make([]byte, difficulty)
	return checkChallenge(challenge, nonce, target)
}

func (_ *HashCash) Calculate(challenge []byte, difficulty uint8) (uint64, error) {
	if difficulty < 1 {
		return 0, fmt.Errorf("difficulty must be greater than 0")
	}
	if difficulty > 20 {
		return 0, fmt.Errorf("difficulty must be less than 20")
	}
	target := make([]byte, difficulty)
	for i := uint64(0); i < math.MaxUint64; i++ {
		if checkChallenge(challenge, i, target) {
			return i, nil
		}
	}
	return 0, fmt.Errorf("failed to find nonce")
}

func (h *HashCash) Challenge() ([]byte, error) {
	c := make([]byte, h.sizeChallenge)
	if _, err := io.ReadFull(h.r, c[:]); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return c[:], nil
}

func checkChallenge(challenge []byte, nonce uint64, target []byte) bool {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	challenge = append(challenge, nonceBytes...)
	h := sha1.New()
	h.Write(challenge)
	hash := h.Sum(nil)
	return bytes.Compare(hash[:len(target)], target) == 0
}
