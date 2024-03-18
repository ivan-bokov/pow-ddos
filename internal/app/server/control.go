package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/ivan-bokov/pow-ddos/internal/protocol"
	"github.com/ivan-bokov/pow-ddos/internal/version"
)

type Control struct {
	conn       net.Conn
	pow        PoW
	difficulty uint8
	challenge  []byte
	in         chan protocol.Message
	storage    Storage
}

func NewControl(conn net.Conn, pow PoW, storage Storage, difficulty uint8) *Control {
	return &Control{conn: conn, difficulty: difficulty, in: make(chan protocol.Message), pow: pow, storage: storage}
}

func (c *Control) Close() {
	if c.conn == nil {
		return
	}
	if err := c.conn.Close(); err != nil {
		slog.Info("failed to close connection", "error", err)
		return
	}
	c.conn = nil
	return
}

func (c *Control) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		c.reader(ctx)
	}()
	go func() {
		defer wg.Done()
		c.manager(ctx)
	}()
	wg.Wait()
	return
}
func (c *Control) reader(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			slog.Info("[PANIC]", "error", err)
		}
	}()
	for {
		if err := ctx.Err(); err != nil {
			slog.Info("context error", "error", err)
			return
		}
		if c.conn == nil {
			slog.Info("connection is nil")
			return
		}
		msg, err := protocol.ReadMsg(c.conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("connection closed")
				return
			}
			slog.Info("failed to read message", "error", err)
			return
		}
		slog.Debug("Received message", "message", msg)
		c.in <- msg
	}
}

func (c *Control) manager(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			slog.Info("[PANIC]", "error", err)
		}
	}()
	for {
		if err := ctx.Err(); err != nil {
			return
		}
		if c.conn == nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case msg := <-c.in:
			switch m := msg.(type) {
			case *protocol.Auth:
				c.handleAuth(m)
			case *protocol.ChallengeResponse:
				c.handleChallenge(m)
			default:
				slog.Info("unknown message type", "message", msg)
				c.Close()
			}
		}
	}
}

func (c *Control) handleChallenge(msg *protocol.ChallengeResponse) {
	if c.challenge == nil {
		if err := protocol.WriteMsg(c.conn, &protocol.Quit{
			Reason: "Invalid challenge",
		}); err != nil {
			slog.Info("failed to write message", "error", err)
		}
		c.Close()
		return
	}
	if !c.pow.Verify(c.challenge, msg.Nonce, c.difficulty) {
		if err := protocol.WriteMsg(c.conn, &protocol.Quit{
			Reason: "Invalid challenge",
		}); err != nil {
			slog.Info("failed to write message", "error", err)
		}
		return
	}
	if err := protocol.WriteMsg(c.conn, &protocol.Quote{
		Quote: c.storage.GetQuote(),
	}); err != nil {
		slog.Info("failed to write message", "error", err)
		return
	}
}

func (c *Control) handleAuth(msg *protocol.Auth) {
	if msg.Version != version.Version {
		if err := protocol.WriteMsg(c.conn, &protocol.Quit{
			Reason: fmt.Sprintf("Unsupported version %s, current version %s", msg.Version, version.Version),
		}); err != nil {
			slog.Info("failed to write message", "error", err)
			c.Close()
			return
		}
	}
	var err error
	if c.challenge, err = c.pow.Challenge(); err != nil {
		slog.Info("failed to generate challenge", "error", err)
		return
	}
	if err = protocol.WriteMsg(c.conn, &protocol.ChallengeRequest{
		Challenge:  c.challenge,
		Difficulty: c.difficulty}); err != nil {
		slog.Info("failed to write message", "error", err)
		return
	}
}
