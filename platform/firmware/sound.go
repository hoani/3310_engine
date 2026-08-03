//go:build tinygo

package firmware

import (
	"time"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/platform/firmware/board"
)

const sampleRate = 44100

const maxInt16 float64 = 32767.0

type SoundPlayer struct {
	active      engine.Note
	sound       engine.Sound
	frame       uint16
	pwm         board.PwmGroup
	ch          uint8
	durPerFrame time.Duration
}

func NewSoundPlayer(pwm board.PwmGroup, ch uint8, fps int) *SoundPlayer {
	p := &SoundPlayer{
		active:      engine.Note{Index: note.None},
		sound:       nil,
		frame:       0,
		pwm:         pwm,
		ch:          ch,
		durPerFrame: time.Second / time.Duration(fps),
	}

	return p
}

func (p *SoundPlayer) Update() {
	if p.sound == nil {
		return
	}
	delta := p.active.Duration - (time.Duration(p.frame) * p.durPerFrame)
	if delta <= p.durPerFrame/2 { // Round it out.
		if p.sound.Done() {
			p.Stop()
			return
		}
		p.playNext()
	}
	p.frame++
}

func (p *SoundPlayer) Play(s engine.Sound) {
	p.Stop()
	s.Reset()

	if s.Done() {
		return // Hmm... nothing to do here.
	}

	p.sound = s
	p.playNext()
}

func (p *SoundPlayer) playNext() {
	p.frame = 0
	n := p.sound.Next()
	p.active = n
	p.pwm.Enable(true)
	period := uint64(float32(time.Second) / note.Frequency(n.Index))
	p.pwm.SetPeriod(period)
	duty := (p.pwm.Top() * uint32(n.Amplitude)) / (2 * 0xFF)
	p.pwm.Set(p.ch, duty)
}

func (p *SoundPlayer) Stop() {
	p.sound = nil
	p.pwm.Enable(false)
}
