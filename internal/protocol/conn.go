package protocol

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"reflect"
)

func readMsgShared(c net.Conn) ([]byte, error) {
	slog.Debug("Waiting to read message")

	var sz int64

	if err := binary.Read(c, binary.BigEndian, &sz); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}
	slog.Debug("Reading message with length", "length", sz)

	buffer := make([]byte, sz)
	n, err := c.Read(buffer)
	slog.Debug("Read message", "message", string(buffer))

	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	if int64(n) != sz {
		return nil, fmt.Errorf("expected to read %d bytes, but only read %d", sz, n)
	}

	return buffer, nil
}

func ReadMsg(c net.Conn) (Message, error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	msg, err := Unpack(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack message: %w", err)
	}

	return msg, nil
}

func WriteMsg(c net.Conn, msg interface{}) error {
	if reflect.ValueOf(msg).Type().Kind() != reflect.Ptr {
		return fmt.Errorf("msg must be a pointer")
	}
	buffer, err := Pack(msg)
	if err != nil {
		return fmt.Errorf("failed to pack message: %w", err)
	}
	slog.Debug("Writing message", "message", string(buffer))

	if err = binary.Write(c, binary.BigEndian, int64(len(buffer))); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	if _, err = c.Write(buffer); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
