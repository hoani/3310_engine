


# nRF Support

Download nrfutil from https://www.nordicsemi.com/Products/Development-tools/nrf-util

Place it somewhere on the system path.

Hold SW2 and press reset on the dongle to put it in bootloader mode prior to flashing.


Figure out where the usb is located:
```
ls /dev/tty.usb*
```

Then we have to manually package and flash:

```
tinygo build -o bin/out.hex -target=pca10059 ./example/fonts              
nrfutil pkg generate --hw-version 52 --sd-req 0x00 --application bin/out.hex --application-version 1 bin/out.zip
nrfutil dfu usb-serial -pkg out.zip -p  /dev/tty.usbmodemFADEA8C7F62D1
```

In my experience, the dongle runs around 3x slower than the nrf - the fonts example is enough to make it top out ~100% at 60FPS - could be fine if running something slower

# Teensy

Need to install the Teensy Loader CLI to your path. This needs to be built from source: https://github.com/PaulStoffregen/teensy_loader_cli

On MacOS, this works:
```
make OS=MACOSX SDK=$(xcrun --sdk macosx --show-sdk-path)
mv teensy_loader_cli /usr/local/bin/teensy_loader_cli
```

Teensy is kind of annoying - we cannot monitor it via usb - need to connect a seperate USB device like so:

```
└─ $ ▶ tinygo flash -target=teensy40 -serial=uart ./example/fonts
Teensy Loader, Command Line, Version 2.3
Read "/var/folders/xj/tnxd4dgn5db1ts_fxvgknfb00000gn/T/tinygo2974094189/main.hex": 172472 bytes, 8.5% usage
Waiting for Teensy device...
 (hint: press the reset button)
Found HalfKay Bootloader
Read "/var/folders/xj/tnxd4dgn5db1ts_fxvgknfb00000gn/T/tinygo2974094189/main.hex": 172472 bytes, 8.5% usage
Programming...................................................................................................................................................................
Booting
└─ $ ▶ tinygo monitor --port=/dev/cu.usbmodem5A980590091
```

The seperate USB device must be connected to UART1. 

Teensy is very fast - for the font example we see:
```
cpu 6%, draw 0% lcd 100% mem 17872 63216/524288
```

I think the cpu 6% is entirely just waiting on SPI comms.

The font drawing basically takes no time at all (compared to nrf which is at ~100% CPU).
