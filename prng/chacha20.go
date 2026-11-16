package prng

import (
	"context"
	cr "crypto/rand"
	"crypto/sha256"
	"errors"
	mr "math/rand/v2"
	"sync"
	"time"

	"golang.org/x/crypto/chacha20"
)

const (
	bufferWords     = 512                  // 512 * 8 = 4096 bytes (64 blocks of ChaCha20)
	reseedThreshold = 64 * 1024 * 1024 / 8 // default reseeding every 64 MB
	reseedPeriod    = 15 * time.Minute     // default reseeding every 15 minutes
	cyclingPeriod   = 50 * time.Millisecond
)

var (
	ErrUnmarshalWyRand = errors.New("invalid WyRand encoding")
)

type ChaCha20 struct {
	mux    sync.Mutex
	cipher *chacha20.Cipher
	words  [bufferWords]uint64
	index  int
	readed uint64 // readed uint64 words
	thres  uint64 // ressed threshold, in uint64
	period time.Duration
	last   time.Time
}

var _ mr.Source = (*ChaCha20)(nil)

// Option configures the generator's parameters.
type Option func(*ChaCha20)

// WithReseedThreshold option sets the threshold of readed uint64 values
// to call reseed after that, default value is set on 64MB.
//
// ATTENTION: after 256GB readed bytes without reseed
// ChaCha20 will panic (uint32 blocks counter overflow).
func WithReseedThreshold(t uint64) Option {
	return func(s *ChaCha20) {
		s.thres = t
	}
}

// WithReseedPeriod option sets the period to reseed ChaCha20,
// default value is set on 15 minutes.
func WithReseedPeriod(d time.Duration) Option {
	return func(s *ChaCha20) {
		s.period = d
	}
}

// NewChaCha20 creates a new generator with entropy from "crypto/rand".
// Usage sample:
//   var ctx, cancel = context.WithCancel(context.Background())
//   defer cancel()
//
//   var src = NewChaCha20()
//   src.StartBackgroundCycling(ctx, 0)
//
//   var rng = mr.New(src) // *Rand from "math/rand/v2"
//   fmt.Println("Stopping position:", rng.IntN(128))
func NewChaCha20(opts ...Option) *ChaCha20 {
	var s = &ChaCha20{
		index:  bufferWords,
		thres:  reseedThreshold,
		period: reseedPeriod,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.reseed()
	return s
}

// reseed collects new entropy, mixes it with the current state
// using SHA-256, and recreates the cipher.
func (s *ChaCha20) reseed() {
	var entropy [64]byte // 32 bytes of new entropy + 32 bytes of the current stream (to preserve state)

	cr.Read(entropy[:32])

	// If the cipher has already been initialized, we mix its output into the new entropy
	// (a NIST/GLI requirement).
	if s.cipher != nil {
		s.cipher.XORKeyStream(entropy[32:], entropy[32:])
	} else {
		cr.Read(entropy[32:])
	}

	// Entropy conditioning via SHA-256 to obtain a uniform 32-byte key
	var key = sha256.Sum256(entropy[:])

	var nonce [12]byte
	cr.Read(nonce[:])

	// Key and the nonce are of the valid size, so no errors should occur
	var cipher, _ = chacha20.NewUnauthenticatedCipher(key[:], nonce[:])

	s.cipher = cipher
	s.readed = 0
	s.last = time.Now()
	s.index = bufferWords // Flush the buffer
}

// Uint64 returns a random 64-bit number.
func (s *ChaCha20) Uint64() uint64 {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.index >= bufferWords {
		s.refill()
	}
	var val = s.words[s.index]
	s.index++
	return val
}

// StartBackgroundCycling runs a background goroutine to shift
// the state in accordance with GLI-19 requirements.
func (s *ChaCha20) StartBackgroundCycling(ctx context.Context, d time.Duration) {
	if d == 0 {
		d = cyclingPeriod
	}
	go func() {
		var ticker = time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.Uint64()
			}
		}
	}()
}
