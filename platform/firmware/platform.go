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
	game   engine.Game
	canvas *canvas
	lcdLed machine.Pin
	lcd    *pcd8544.Device
	snd    *SoundPlayer
	keypad *Keypad
	cmd *command.CommandImpl[engine.Key]
}

func New(game engine.Game, lcd *pcd8544.Device, led machine.Pin, snd *SoundPlayer, keypad *Keypad, cmd *command.CommandImpl[engine.Key]) *Platform {
	return &Platform{
		game:   game,
		canvas: NewCanvas(lcd),
		lcdLed: led,
		lcd:    lcd,
		snd:    snd,
		keypad: keypad,
		cmd: cmd,
	}
}

func (p *Platform) Console(format string, args ...any) {
	if len(args) == 0 {
		fmt.Printf(format)
	} else {
		fmt.Printf(format, args)
	}

	if !strings.HasSuffix(format, "\n") {
		fmt.Printf("\n")
	}

}

func (p *Platform) Run() error {
	var m runtime.MemStats
	count := 0
	info := p.game.Info()
	memFloor := uint64(0)
	memLast := uint64(0)
	period := time.Second / time.Duration(p.game.Info().Fps)
	illuminated := false
	for {
		count++
		start := time.Now()
		p.keypad.Update()
		p.cmd.Update()
		if err := p.game.Update(); err != nil {
			return err
		}
		p.snd.Update()
		dStart := time.Now()
		if err := p.game.Draw(p.canvas); err != nil {
			return err
		}
		lcdStart := time.Now()

		if err := p.Draw(); err != nil {
			return err
		}
		if info.Debug {
			if (count % (5 * 60)) == 0 {
				dDur := lcdStart.Sub(dStart)
				lcdDur := time.Since(lcdStart)
				dur := time.Since(start)

				fmt.Printf("cpu %d%%, draw %d%% lcd %d%% mem %d %d/%d\n", 100*dur/period, 100*dDur/dur, 100*lcdDur/dur, memFloor, m.Alloc, m.Sys)
			}
			runtime.ReadMemStats(&m)
			if memLast > m.Alloc || memLast == 0 {
				memFloor = m.Alloc
			}
			memLast = m.Alloc
		}
		if info.Illuminated != illuminated {
			p.lcdLed.Set(info.Illuminated)
			illuminated = info.Illuminated
		}
		
		rem := period - time.Since(start)
		time.Sleep(rem)
	}
}

func (p *Platform) Draw() error {
	return p.canvas.Display()
}

func setupPcd(def *board.Pcd) *pcd8544.Device {
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

	return d
}

func setupBuzzer(def *board.Buzzer, fps int) (*SoundPlayer, error) {
	def.Pwm.Configure(machine.PWMConfig{
		Period: 50 * 1e6,
	})
	ch, err := def.Pwm.Channel(def.Pin)
	if err != nil {
		return nil, err
	}
	def.Pwm.Enable(false)

	return NewSoundPlayer(def.Pwm, ch, fps), nil
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

func Run(game engine.Game) {

	def, err := board.Get()
	handleErr("Board def", err)

	// Configure SPI with a 1 MHz frequency.
	pcd := setupPcd(&def.Pcd)

	def.Pcd.LedPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	buzzer, err := setupBuzzer(&def.Buzzer, game.Info().Fps)
	handleErr("Audio Setup", err)

	keypad, cmd := NewKeypad(def.Keypad.Col, def.Keypad.Row)

	p := New(game, pcd, def.Pcd.LedPin, buzzer, keypad, cmd)

	game.Setup(cmd, buzzer, p)

	err = p.Run()
	handleErr("Game Crash", err)
}
