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

type Sound struct {
	sound engine.Sound
	frame uint16
	loop  bool
}

type SoundPlayer struct {
	active      engine.Note
	sfx         Sound
	track       Sound
	pwm         board.PwmGroup
	ch          uint8
	durPerFrame time.Duration
}

func NewSoundPlayer(pwm board.PwmGroup, ch uint8) *SoundPlayer {
	p := &SoundPlayer{
		active:      engine.NoneNote(0),
		sfx:         Sound{sound: nil, frame: 0, loop: false},
		track:       Sound{sound: nil, frame: 0, loop: false},
		pwm:         pwm,
		ch:          ch,
		durPerFrame: time.Second / time.Duration(60),
	}
	p.pwm.Enable(true)
	return p
}

func (p *SoundPlayer) SetFps(fps int) {
	if fps == 0 {
		return
	}
	p.durPerFrame = time.Second / time.Duration(fps)
}

func (p *SoundPlayer) Update() {
	if p.update(&p.sfx) {
		return
	}
	if p.update(&p.track) {
		return
	}
}

func (p *SoundPlayer) update(s *Sound) (ok bool) {
	if s.sound == nil {
		return false
	}
	delta := p.active.Duration() - (time.Duration(s.frame) * p.durPerFrame)
	if delta <= p.durPerFrame/2 { // Round it out.
		if s.sound.Done() {
			if s.loop {
				s.sound.Reset()
			} else {
				s.sound = nil
				p.pwm.Set(p.ch, 0)
				return false
			}
		}
		p.playNext(s)
	}
	s.frame++
	return true
}

func (p *SoundPlayer) Track(s engine.Sound, loop bool) {
	p.track.sound = nil

	s.Reset()

	if s.Done() {
		return // Hmm... nothing to do here.
	}

	p.track.loop = loop
	p.track.sound = s
	if p.sfx.sound == nil {
		p.playNext(&p.track)
	}
}

func (p *SoundPlayer) Play(s engine.Sound) {
	p.sfx.sound = nil
	s.Reset()

	if s.Done() {
		return // Hmm... nothing to do here.
	}

	p.sfx.loop = false
	p.sfx.sound = s
	p.playNext(&p.sfx)
}

func (p *SoundPlayer) playNext(s *Sound) {
	s.frame = 0
	n := s.sound.Next()

	p.active = n
	if n.Index() != note.None {
		period := uint64(float32(time.Second) / n.Frequency())
		p.pwm.SetPeriod(period)
		duty := (p.pwm.Top() * uint32(n.Amplitude())) / (2 * 0xFF)
		p.pwm.Set(p.ch, duty)
	} else {
		p.pwm.Set(p.ch, 0)
	}

}

func (p *SoundPlayer) Stop() {
	p.sfx.sound = nil
	p.track.sound = nil
	p.pwm.Set(p.ch, 0)
}
