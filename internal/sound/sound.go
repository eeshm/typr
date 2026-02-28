package sound

import (
	"bytes"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	sampleRate = 44100
	channels   = 2     // stereo (oto v3 requires >= 2 on some backends)
	bitDepth   = 2     // 16-bit signed PCM
	frameSize  = channels * bitDepth
)

var (
	ctx      *oto.Context
	once     sync.Once
	initErr  error
	clickBuf []byte // pre-computed click waveform
	errorBuf []byte // pre-computed error waveform
)

// Init initializes the audio system. Safe to call multiple times;
// only the first call takes effect. Returns nil if audio is unavailable
// (the app should still work without sound).
func Init() error {
	once.Do(func() {
		op := &oto.NewContextOptions{
			SampleRate:   sampleRate,
			ChannelCount: channels,
			Format:       oto.FormatSignedInt16LE,
		}
		var ready chan struct{}
		ctx, ready, initErr = oto.NewContext(op)
		if initErr != nil {
			return
		}
		<-ready

		// Pre-generate waveforms so PlayClick/PlayError are allocation-light.
		clickBuf = generateClick()
		errorBuf = generateError()
	})
	return initErr
}

// PlayClick plays a short mechanical keyboard click sound (non-blocking).
func PlayClick() {
	playBuf(clickBuf)
}

// PlayError plays a short error buzz sound (non-blocking).
func PlayError() {
	playBuf(errorBuf)
}

func playBuf(buf []byte) {
	if ctx == nil || len(buf) == 0 {
		return
	}
	p := ctx.NewPlayer(bytes.NewReader(buf))
	p.Play()
	go func() {
		for p.IsPlaying() {
			time.Sleep(time.Millisecond)
		}
		_ = p.Close()
	}()
}

// generateClick creates a ~10ms click: a mix of mid+high frequencies
// with a sharp exponential decay — mimics a mechanical key switch.
func generateClick() []byte {
	duration := 0.010 // 10 ms
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*frameSize)

	for i := range samples {
		t := float64(i) / float64(sampleRate)

		// Sharp exponential decay envelope.
		envelope := math.Exp(-t * 600)

		// Mix of harmonics for a natural click feel.
		sample := 0.0
		sample += math.Sin(2*math.Pi*900*t) * 0.5  // fundamental
		sample += math.Sin(2*math.Pi*1800*t) * 0.3 // 1st harmonic
		sample += math.Sin(2*math.Pi*3600*t) * 0.2 // high end

		sample *= envelope * 0.25 // master volume

		val := int16(clamp(sample, -1, 1) * 32767)
		lo := byte(val)
		hi := byte(val >> 8)

		off := i * frameSize
		// Write same sample to both L and R channels.
		buf[off] = lo
		buf[off+1] = hi
		buf[off+2] = lo
		buf[off+3] = hi
	}
	return buf
}

// generateError creates a ~30ms low buzz for wrong keystrokes.
func generateError() []byte {
	duration := 0.030 // 30 ms
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*frameSize)

	for i := range samples {
		t := float64(i) / float64(sampleRate)

		envelope := math.Exp(-t * 200)

		// Low harsh buzz.
		sample := math.Sin(2*math.Pi*300*t)*0.6 + math.Sin(2*math.Pi*450*t)*0.4
		sample *= envelope * 0.20

		val := int16(clamp(sample, -1, 1) * 32767)
		lo := byte(val)
		hi := byte(val >> 8)

		off := i * frameSize
		buf[off] = lo
		buf[off+1] = hi
		buf[off+2] = lo
		buf[off+3] = hi
	}
	return buf
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
