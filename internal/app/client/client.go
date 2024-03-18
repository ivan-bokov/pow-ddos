package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/ivan-bokov/pow-ddos/internal/protocol"
	"github.com/ivan-bokov/pow-ddos/internal/util"
	"github.com/ivan-bokov/pow-ddos/internal/version"
)

type Solver interface {
	Calculate(challenge []byte, difficulty uint8) (uint64, error)
}

type Client struct {
	addr   string
	solver Solver
}

func New(serverAddr string, solver Solver) *Client {
	return &Client{
		addr:   serverAddr,
		solver: solver,
	}
}

func (c *Client) GetQuote(ctx context.Context) (string, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}()
	id, err := util.GenID(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate id: %w", err)
	}
	if err = protocol.WriteMsg(conn, &protocol.Auth{id, version.Version}); err != nil {
		return "", fmt.Errorf("failed to write message: %w", err)
	}
	var msg protocol.Message
	for {
		if err = ctx.Err(); err != nil {
			return "", fmt.Errorf("context error: %w", err)
		}
		if msg, err = protocol.ReadMsg(conn); err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}
		switch m := msg.(type) {
		case *protocol.Quote:
			return m.Quote, nil
		case *protocol.Quit:
			return "", fmt.Errorf("server error: %s", m.Reason)
		case *protocol.ChallengeRequest:
			nonce, err := c.solver.Calculate(m.Challenge, m.Difficulty)
			if err != nil {
				return "", fmt.Errorf("failed to calculate: %w", err)
			}
			if err = protocol.WriteMsg(conn, &protocol.ChallengeResponse{Nonce: nonce}); err != nil {
				return "", fmt.Errorf("failed to write message: %w", err)
			}
		}
	}
}
