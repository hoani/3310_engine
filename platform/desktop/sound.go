package desktop

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/sound/note"
)

const sampleRate = 44100

const maxInt16 float64 = 32767.0

type SoundPlayer struct {
	// Oto player
	ctx    *oto.Context
	player *oto.Player
	// Source
	src *source
	// Notes
	filters []*biquad
	notes   [note.Total]*voice
}

func NewSoundPlayer() (*SoundPlayer, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return nil, err
	}
	<-ready

	p := &SoundPlayer{
		ctx:     ctx,
		filters: []*biquad{newBandpass(3442, 1.63), newBandpass(5411, 9.75)},
		src:     &source{},
	}

	for i := range p.notes {
		freq := note.Frequency(note.Index(i))
		p.notes[i] = newVoice(float64(freq), p.filters)
	}

	p.player = ctx.NewPlayer(p.src)
	p.player.Play()

	return p, nil
}

func (p *SoundPlayer) Play(s engine.Sound) {
	var buf bytes.Buffer
	s.Reset()
	for _, filter := range p.filters {
		filter.reset()
	}
	for !s.Done() {
		n := s.Next()
		p.notes[n.Index].Generate(n.Duration, n.Amplitude, &buf)
	}
	p.src.set(buf.Bytes())
}

func (p *SoundPlayer) Stop() {
	p.src.set(nil)
}

type biquad struct {
	b0, b1, b2, a1, a2 float64
	x1, x2, y1, y2     float64
}

func newBandpass(fc, q float64) *biquad {
	w0 := 2 * math.Pi * fc / sampleRate
	al := math.Sin(w0) / (2 * q)
	c := math.Cos(w0)
	a0 := 1 + al
	return &biquad{b0: al / a0, b1: 0, b2: -al / a0,
		a1: -2 * c / a0, a2: (1 - al) / a0}
}

func (f *biquad) process(x float64) float64 {
	y := f.b0*x + f.b1*f.x1 + f.b2*f.x2 - f.a1*f.y1 - f.a2*f.y2
	f.x2, f.x1 = f.x1, x
	f.y2, f.y1 = f.y1, y
	return y
}

func (f *biquad) reset() { f.x1, f.x2, f.y1, f.y2 = 0, 0, 0, 0 }

type voice struct {
	freq    float64
	filters []*biquad
}

func newVoice(freq float64, filters []*biquad) *voice {
	return &voice{
		freq:    freq,
		filters: filters,
	}
}

// polyBlep returns the correction for a discontinuity, where t is the
// current phase (0..1) and dt is the phase increment per sample.
func polyBlep(t, dt float64) float64 {
	if t < dt {
		// first sample after the reset
		t /= dt
		return t + t - t*t - 1
	}
	if t > 1.0-dt {
		// last sample before the reset
		t = (t - 1.0) / dt
		return t*t + t + t + 1.0
	}
	return 0
}

func (s *voice) Generate(dur time.Duration, amplitude uint8, w io.Writer) error {
	N := int((sampleRate * dur) / time.Second)

	volume := (maxInt16 * float64(amplitude)) / float64(0xff)

	var phase float64 = 0

	dt := s.freq / sampleRate
	for range N {
		saw := 2.0*phase - 1.0
		saw -= polyBlep(phase, dt)

		v := saw * 2.7 // Compensate filter attenuation.
		for _, f := range s.filters {
			v = f.process(v)
		}

		// Clamp, just in case.
		v = math.Min(math.Max(v, -1.0), 1.0)

		out := int32(v * volume)
		if err := binary.Write(w, binary.LittleEndian, uint16(out)); err != nil {
			return err
		}
		phase += dt
		if phase >= 1.0 {
			phase -= 1.0
		}
	}
	return nil
}

// Source plays sound continuously, this emulates a buzzer that we can feed generated sound into.
type source struct {
	mu  sync.Mutex
	pcm []byte
	pos int
}

func (s *source) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	if s.pos < len(s.pcm) {
		n = copy(p, s.pcm[s.pos:])
		n -= n % 2 // stay on 16-bit frame boundaries
		s.pos += n
	}
	// Generate silence if we are out of things to play
	if n == 0 {
		for i := n; i < len(p); i++ {
			p[i] = 0
		}
	}
	return len(p), nil
}

func (s *source) set(pcm []byte) {
	s.mu.Lock()
	s.pcm, s.pos = pcm, 0
	s.mu.Unlock()
}
