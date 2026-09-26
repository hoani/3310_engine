//go:build tinygo

package firmware

import (
	"fmt"
	"machine"
	"runtime"
	"strings"
	"time"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/platform/firmware/board"
	"tinygo.org/x/drivers/pcd8544"
)

type Platform struct {
	game       engine.Game
	config     engine.Config
	canvas     *canvas
	display    *Display
	snd        *SoundPlayer
	keypad     *Keypad
	cmd        *command.CommandImpl[engine.Key]
	extensions []Extension
	debug      Debug
	period     time.Duration
}

type Debug struct {
	p *Platform
}

type Display struct {
	lcd     *pcd8544.Device
	rst     machine.Pin
	enabled bool
}

func New(game engine.Game, display *Display, snd *SoundPlayer, keypad *Keypad, cmd *command.CommandImpl[engine.Key], extensions ...Extension) *Platform {
	p := &Platform{
		game:       game,
		canvas:     NewCanvas(display.lcd),
		snd:        snd,
		keypad:     keypad,
		cmd:        cmd,
		display:    display,
		extensions: extensions,
		period:     time.Second / time.Duration(60),
	}
	p.debug = Debug{p: p}
	return p
}

func (p *Platform) Config(c engine.Config) {
	p.config = c
	p.snd.SetFps(p.config.Fps)
	p.period = time.Second / time.Duration(p.config.Fps)
}

func (p *Platform) Cmd() command.Command[engine.Key] {
	return p.cmd
}

func (p *Platform) Snd() engine.SoundPlayer {
	return p.snd
}

func (p *Platform) Debug() engine.Debug {
	return &p.debug
}

func (d *Debug) Enabled() bool {
	return d.p.config.Debug
}

func (d *Debug) Console(format string, args ...any) {
	if len(args) == 0 {
		fmt.Printf(format)
	} else {
		fmt.Printf(format, args)
	}

	if !strings.HasSuffix(format, "\n") {
		fmt.Printf("\n")
	}
}

func (p *Platform) Display() engine.Display {
	return p.display
}

func (d *Display) Enable(e bool) {
	if e == d.enabled {
		return
	}
	d.enabled = e

	if d.enabled {
		d.lcd.Configure(pcd8544.Config{
			Width:  84,
			Height: 48,
		})
	} else {
		// Reset pin puts the LCD into power-down mode until we reconfigure it.
		d.rst.Low()
		time.Sleep(100 * time.Microsecond)
		d.rst.High()
	}
}

func (p *Platform) Run() error {
	var m runtime.MemStats
	count := 0
	memFloor := uint64(0)
	memLast := uint64(0)

	for {
		count++
		start := time.Now()
		for _, extension := range p.extensions {
			if err := extension.Update(); err != nil {
				return err
			}
		}
		p.keypad.Update()
		p.cmd.Update()
		if err := p.game.Update(); err != nil {
			return err
		}
		p.snd.Update()
		dStart := time.Now()

		// Only draw if the display is enabled, otherwise we are better off sleeping.
		if p.display.enabled {
			if err := p.game.Draw(p.canvas); err != nil {
				return err
			}
		}

		lcdStart := time.Now()
		if err := p.Draw(); err != nil {
			return err
		}

		if p.config.Debug {
			if (count % (5 * p.config.Fps)) == 0 {
				dDur := lcdStart.Sub(dStart)
				lcdDur := time.Since(lcdStart)
				dur := time.Since(start)

				fmt.Printf("cpu %d%%, draw %d%% lcd %d%% mem %d %d/%d\n", 100*dur/p.period, 100*dDur/dur, 100*lcdDur/dur, memFloor, m.Alloc, m.Sys)
			}
			runtime.ReadMemStats(&m)
			if memLast > m.Alloc || memLast == 0 {
				memFloor = m.Alloc
			}
			memLast = m.Alloc
		}

		rem := p.period - time.Since(start)
		time.Sleep(rem)
	}
}

func (p *Platform) Draw() error {
	// Only draw if the display is enabled, otherwise we are better off sleeping.
	if !p.display.enabled {
		return nil
	}
	return p.canvas.Display()
}

func setupPcd(def *board.Pcd) *Display {
	def.Spi.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       def.SckPin,
		SDO:       def.SdoPin,
		SDI:       def.SdiPin,
	})

	dcPin := def.DcPin
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dcPin.High()

	rstPin := def.RstPin
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin.High()

	scePin := def.ScePin
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin.High() // CS inactive

	d := pcd8544.New(def.Spi, dcPin, rstPin, scePin)

	d.Configure(pcd8544.Config{
		Width:  84,
		Height: 48,
	})

	return &Display{
		lcd:     d,
		rst:     rstPin,
		enabled: true,
	}
}

func setupBuzzer(def *board.Buzzer) (*SoundPlayer, error) {
	def.Pwm.Configure(machine.PWMConfig{
		Period: 50 * 1e6,
	})
	ch, err := def.Pwm.Channel(def.Pin)
	if err != nil {
		return nil, err
	}
	def.Pwm.Enable(false)

	return NewSoundPlayer(def.Pwm, ch), nil
}

func handleErr(info string, err error) {
	if err == nil {
		return
	}
	for {
		fmt.Printf("%s: %s\n", info, err.Error())
		time.Sleep(5 * time.Second)
	}
}

func Run(game engine.Game, extensions ...Extension) {

	def, err := board.Get()
	handleErr("Board def", err)

	display := setupPcd(def.Pcd())

	buzzer, err := setupBuzzer(def.Buzzer())
	handleErr("Audio Setup", err)

	keypad, cmd := NewKeypad(def.Keypad().Col, def.Keypad().Row)

	def.Initialize()

	p := New(game, display, buzzer, keypad, cmd, extensions...)

	game.Setup(p)

	for i, extension := range extensions {
		handleErr(fmt.Sprintf("Extension %d Setup", i), extension.Setup())
	}

	err = p.Run()
	handleErr("Game Crash", err)

	for {
		time.Sleep(time.Second)
	}
}
