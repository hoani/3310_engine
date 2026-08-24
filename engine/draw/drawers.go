package draw

func NewRectangleDrawer(p0, p1 Point) *RectangleBuilder {
	return &RectangleBuilder{
		p0: p0,
		p1: p1,
	}
}

func NewLineDrawer(p0, p1 Point) *LineBuilder {
	return &LineBuilder{
		p0: p0,
		p1: p1,
	}
}

func NewCircleDrawer(center Point, diameter uint16) *CircleBuilder {
	return &CircleBuilder{
		center:   center,
		diameter: diameter,
	}
}

func NewOvalDrawer(center Point, width, height uint16) *OvalBuilder {
	return &OvalBuilder{
		center: center,
		width:  width,
		height: height,
	}
}

func NewTriangleDrawer(p0, p1, p2 Point) *TriangleBuilder {
	return &TriangleBuilder{
		p0: p0,
		p1: p1,
		p2: p2,
	}
}
