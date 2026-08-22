//go:build !tinygo

package soundplayer

import (
	"bytes"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/sound/note"
)

const sampleRate = 44100

const maxInt16 float64 = 32767.0

type Player struct {
	// Oto player
	ctx    *oto.Context
	player *oto.Player
	// Source
	src   *source
	track sourceData
	// Notes
	filters []*biquad
	notes   [note.Total]*voice
}

func New() (*Player, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   20 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	<-ready

	p := &Player{
		ctx:     ctx,
		filters: newBuzzerFilters(),
		src:     &source{},
	}
	p.src.reload = p.reload

	for i := range p.notes {
		freq := note.Frequency(note.Index(i))
		p.notes[i] = newVoice(float64(freq), p.filters)
	}

	p.player = ctx.NewPlayer(p.src)
	p.player.Play()

	return p, nil
}

func (p *Player) Track(s engine.Sound, loop bool) {
	p.load(s, true, loop)
	p.track = p.src.get()
}

func (p *Player) Play(s engine.Sound) {
	current := p.src.get()
	if current.music {
		p.track.pos = current.pos
	}
	p.load(s, false, false)
}

func (p *Player) reload() sourceData {
	if p.track.music {
		return p.track
	}
	return sourceData{}
}

func (p *Player) load(s engine.Sound, music, loop bool) {
	var buf bytes.Buffer
	p.player.Pause()
	s.Reset()
	for _, filter := range p.filters {
		filter.reset()
	}
	for !s.Done() {
		n := s.Next()
		if n.Index() == note.Custom {
			v := newVoice(float64(n.Frequency()), p.filters)
			v.Generate(n.Duration(), n.Amplitude(), &buf)
		} else {
			p.notes[n.Index()].Generate(n.Duration(), n.Amplitude(), &buf)
		}
	}
	p.src.set(buf.Bytes(), music, loop)
	p.player = p.ctx.NewPlayer(p.src)
	p.player.Play()
}

func (p *Player) Stop() {
	p.player.Pause()
	p.src.set(nil, false, false)
	p.track = sourceData{}
}

// Source plays sound continuously, this emulates a buzzer that we can feed generated sound into.
type sourceData struct {
	pcm   []byte
	pos   int
	loop  bool
	music bool
}

type source struct {
	mu sync.Mutex
	sourceData
	reload func() sourceData
}

func (s *source) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	if s.pos < len(s.pcm) {
		n = copy(p, s.pcm[s.pos:])
		n -= n % 2 // stay on 16-bit frame boundaries
		s.pos += n

		if s.pos >= len(s.pcm) {
			if s.loop {
				s.pos = 0
				return n, nil
			} else if !s.music {
				s.sourceData = s.reload()
			}
		}
	}

	// Generate silence if we are out of things to play
	if n == 0 {
		for i := n; i < len(p); i++ {
			p[i] = 0
		}
	}
	return len(p), nil
}

func (s *source) set(pcm []byte, music, loop bool) {
	s.mu.Lock()
	s.pcm, s.pos, s.music, s.loop = pcm, 0, music, loop
	s.mu.Unlock()
}

func (s *source) get() sourceData {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sourceData
}
