package file

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Storage struct {
	quotes []string
}

func New(path string) (*Storage, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	quotes := make([]string, 0)
	for _, line := range strings.Split(string(f), "\n") {
		if line == "" {
			continue
		}
		quotes = append(quotes, line)
	}
	return &Storage{
		quotes: quotes,
	}, nil
}

func (s *Storage) GetQuote() string {
	return s.quotes[rand.Intn(len(s.quotes))]
}
