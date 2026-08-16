# Surface

Surfaces are a memory intensive way of holding and dumping graphics to a canvas.

The should typically be used for static images or forms which shouldn't need to be re-rendered every frame.

Using a Pico, a full screen surface rendering some rectangles with dithering will reduce cpu usage from 114% down to 45%.

