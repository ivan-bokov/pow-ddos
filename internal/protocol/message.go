package protocol

import (
	"encoding/json"
	"reflect"
)

var cacheType map[string]reflect.Type

func init() {
	cacheType = make(map[string]reflect.Type)

	t := func(obj interface{}) reflect.Type { return reflect.TypeOf(obj).Elem() }
	cacheType["Quit"] = t((*Quit)(nil))
	cacheType["Auth"] = t((*Auth)(nil))
	cacheType["Quote"] = t((*Quote)(nil))
	cacheType["ChallengeResponse"] = t((*ChallengeResponse)(nil))
	cacheType["ChallengeRequest"] = t((*ChallengeRequest)(nil))
}

type Message interface{}

type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Quit struct {
	Reason string `json:"reason"`
}

type ChallengeResponse struct {
	Nonce uint64 `json:"nonce"`
}
type ChallengeRequest struct {
	Challenge  []byte `json:"challenge"`
	Difficulty uint8  `json:"difficulty"`
}

type Auth struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type Quote struct {
	Quote string `json:"quote"`
}
