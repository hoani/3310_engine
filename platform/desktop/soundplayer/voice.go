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

		v := saw * 2.5 // Compensate filter attenuation.
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

// I came across these parameters by feeding some sample tones into Claude and asking it synthesize the sound.
// They seem pretty close to the real thing, but I've been unable to find any real documentation since
// Nokia's buzzer was a custom component.
// My understanding of the gains is as follows:
// * 3440 Hz emulates the resonant frequency of the buzzer
// * 5400 Hz emulatees the casing damping
func newBuzzerFilters() []*biquad {
	return []*biquad{newBandpass(3440, 1.00), newBandpass(5400, 10.0)}
}

// Second-order IIR filter in Direct form I (https://en.wikipedia.org/wiki/Digital_biquad_filter)
type biquad struct {
	b0, b1, b2, a1, a2 float64
	x1, x2, y1, y2     float64
}

// Create a new bandpass filter with center frequency fc and quality factor q.
func newBandpass(fc, q float64) *biquad {
	w0 := 2 * math.Pi * fc / sampleRate
	alpha := math.Sin(w0) / (2 * q)
	cosw0 := math.Cos(w0)
	a0 := 1 + alpha
	// Uses BPF (constant 0 dB peak gain)
	// See: https://webaudio.github.io/Audio-EQ-Cookbook/audio-eq-cookbook.html eq 20.
	// Predivides by a0 to avoid additional divisions when using the filter.
	return &biquad{
		b0: alpha / a0,
		b1: 0,
		b2: -alpha / a0,
		a1: -2 * cosw0 / a0,
		a2: (1 - alpha) / a0,
	}
}

// Apply IIR filter
func (f *biquad) process(x float64) float64 {
	y := f.b0*x + f.b1*f.x1 + f.b2*f.x2 - f.a1*f.y1 - f.a2*f.y2

	// Shift states
	f.x2, f.x1 = f.x1, x
	f.y2, f.y1 = f.y1, y
	return y
}

// Resets all states
func (f *biquad) reset() {
	f.x1, f.x2, f.y1, f.y2 = 0, 0, 0, 0
}
