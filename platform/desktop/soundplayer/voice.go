//go:build !tinygo

package soundplayer

import (
	"encoding/binary"
	"io"
	"math"
	"time"
)

// I understand some of this, but not 100% on the filtering.
// Voice generates saw waves to imitate the signal going into a buzzer.
// The sawwave has super high harmonics removed via polyBlep which basically smooths the falling edge
// We use two bandpass filters which are multiplied together to imitate the buzzer's physical body (I think?)
// And these filters basically make up what a buzzer may be manufactured to resonate at.
// But... ultimately, it sounds like a buzzer, so I'm not really touching it.

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

func newBuzzerFilters() []*biquad {
	return []*biquad{newBandpass(3442, 1.63), newBandpass(5411, 9.75)}
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
