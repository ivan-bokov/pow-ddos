package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ivan-bokov/pow-ddos/config"
	"github.com/ivan-bokov/pow-ddos/internal/service/ddos"
	"github.com/ivan-bokov/pow-ddos/internal/service/pow/hashcash"
	"github.com/ivan-bokov/pow-ddos/internal/storage/file"
)

type PoW interface {
	Challenge() ([]byte, error)
	Verify(challenge []byte, nonce uint64, difficulty uint8) bool
}

type Storage interface {
	GetQuote() string
}

type Guard interface {
	Take()
	Rate() uint64
	Reset()
	Difficulty() uint8
}

type Server struct {
	writeTimeout  time.Duration
	readTimeout   time.Duration
	idleTimeout   time.Duration
	keepAlive     time.Duration
	pow           PoW
	ddos          Guard
	maxDifficulty uint8
	storage       Storage
}

func New(cfg config.ServerConfig) (*Server, error) {
	s := &Server{
		ddos:          ddos.New(cfg.DDOS.Window, uint8(cfg.Pow.DifficultyStart), cfg.DDOS.Rate),
		pow:           hashcash.New(cfg.Pow.ChallengeLength),
		writeTimeout:  cfg.WriteTimeout,
		readTimeout:   cfg.ReadTimeout,
		idleTimeout:   cfg.IdleTimeout,
		keepAlive:     cfg.KeepAlive,
		maxDifficulty: cfg.Pow.MaxDifficulty,
	}
	storage, err := file.New(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}
	s.storage = storage
	return s, nil
}

func (s *Server) difficulty() uint8 {
	return min(s.ddos.Difficulty(), s.maxDifficulty)
}

func (s *Server) Run(ctx context.Context, addr string) error {
	lc := net.ListenConfig{
		KeepAlive: s.keepAlive,
	}
	l, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()
	slog.Info("Start listening", "addr", addr)
	for {
		if err = ctx.Err(); err != nil {
			return fmt.Errorf("context error: %w", err)
		}
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept: %w", err)
		}
		go s.handle(ctx, conn)
	}
}

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			slog.Info("[PANIC]", "error", err)
		}
	}()
	s.ddos.Take()
	if err := s.setTimeouts(conn); err != nil {
		slog.Info("failed to set timeouts", "error", err)
		return
	}
	ctrl := NewControl(conn, s.pow, s.storage, s.difficulty())
	defer ctrl.Close()
	ctrl.Run(ctx)
	return
}

func (s *Server) setTimeouts(conn net.Conn) error {
	if s.readTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
			return fmt.Errorf("failed to set read deadline: %w", err)
		}
	}
	if s.idleTimeout > 0 {
		if err := conn.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			return fmt.Errorf("failed to set deadline: %w", err)
		}
	}
	if s.writeTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
			return fmt.Errorf("failed to set write deadline: %w", err)
		}
	}

	return nil
}
