package hashcash

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

type someMockReader struct{}

func (r someMockReader) Read(p []byte) (n int, err error) {
	copy(p, []byte{1, 2, 3, 4, 5, 6, 7, 8})
	return 8, nil
}

type failingReader struct{}

func (r failingReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read random bytes")
}

func TestHashCash_Challenge(t *testing.T) {
	h := &HashCash{r: someMockReader{}, sizeChallenge: 8}

	// Testing that the function returns an 8-byte array
	challenge, err := h.Challenge()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(challenge) != 8 {
		t.Errorf("expected 8-byte array, got %d bytes", len(challenge))
	}

	// Testing that the function handles the error case when failing to read random bytes
	h.r = failingReader{}
	_, err = h.Challenge()
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	h = New(16)
	challenge, err = h.Challenge()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(challenge) != 16 {
		t.Errorf("expected 8-byte array, got %d bytes", len(challenge))
	}
}

func TestHashCash_Calculate(t *testing.T) {
	type fields struct {
		r             io.Reader
		sizeChallenge int
	}
	type args struct {
		challenge  []byte
		difficulty int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "test_pass_1",
			fields: fields{
				r:             rand.Reader,
				sizeChallenge: 8,
			},
			args: args{
				challenge:  []byte("examplechallenge"),
				difficulty: 1,
			},
			want:    61,
			wantErr: false,
		},
		{
			name: "test_pass_2",
			fields: fields{
				r:             rand.Reader,
				sizeChallenge: 8,
			},
			args: args{
				challenge:  []byte("examplechallenge"),
				difficulty: 2,
			},
			want:    5619,
			wantErr: false,
		},
		{
			name: "test_pass_3",
			fields: fields{
				r:             rand.Reader,
				sizeChallenge: 8,
			},
			args: args{
				challenge:  []byte("examplechallenge"),
				difficulty: 3,
			},
			want:    5078300,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HashCash{
				r:             tt.fields.r,
				sizeChallenge: tt.fields.sizeChallenge,
			}
			got, err := h.Calculate(tt.args.challenge, tt.args.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkChallenge(t *testing.T) {
	type args struct {
		challenge []byte
		nonce     uint64
		target    []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_pass",
			args: args{
				challenge: []byte("examplechallenge"),
				nonce:     61,
				target:    []byte{0},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkChallenge(tt.args.challenge, tt.args.nonce, tt.args.target); got != tt.want {
				t.Errorf("checkChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkHashCash_Verify(b *testing.B) {
	h := New(8)
	for i := 0; i < b.N; i++ {
		_ = h.Verify([]byte("examplechallenge"), 61, 1)
	}
}
