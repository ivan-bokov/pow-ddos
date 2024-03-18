package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func unpack(buffer []byte, msgIn Message) (Message, error) {
	var env Envelope
	if err := json.Unmarshal(buffer, &env); err != nil {
		return nil, fmt.Errorf("failed to unpack: %w", err)
	}

	var msg Message
	if msgIn == nil {
		t, ok := cacheType[env.Type]
		if !ok {
			return nil, fmt.Errorf("Unsupported message type %s", env.Type)
		}

		msg = reflect.New(t).Interface().(Message)
	} else {
		msg = msgIn
	}

	if err := json.Unmarshal(env.Payload, &msg); err != nil {
		return nil, fmt.Errorf("failed to unpack: %w", err)
	}
	return msg, nil
}

func UnpackInto(buffer []byte, msg Message) error {
	if _, err := unpack(buffer, msg); err != nil {
		return fmt.Errorf("failed to unpack: %w", err)
	}
	return nil
}

func Unpack(buffer []byte) (Message, error) {
	msg, err := unpack(buffer, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack: %w", err)
	}
	return msg, nil
}

func Pack(payload interface{}) ([]byte, error) {
	body, err := json.Marshal(struct {
		Type    string
		Payload interface{}
	}{
		Type:    reflect.TypeOf(payload).Elem().Name(),
		Payload: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to pack: %w", err)
	}
	return body, nil
}
