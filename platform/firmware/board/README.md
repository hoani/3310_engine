


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

