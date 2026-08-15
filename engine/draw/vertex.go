package draw

type triangleStrip struct {
	points []Point
}

func NewTriangleStrip(points ...Point) ShapeDrawer {
	return &triangleStrip{
		points: points,
	}
}

func (s triangleStrip) Fill(drawPixel drawPixel) {

	for i := range s.points {
		if (len(s.points) - i) < 3 {
			break
		}
		filledTriangle(drawPixel, s.points[i], s.points[i+1], s.points[i+2])
	}
	s.Outline(drawPixel)
}

func (s triangleStrip) Outline(drawPixel drawPixel) {
	n := len(s.points)
	if n < 3 {
		return
	}
	if n == 3 {
		line(drawPixel, s.points[0], s.points[1])
		line(drawPixel, s.points[0], s.points[2])
		line(drawPixel, s.points[1], s.points[2])
		return
	}
	if n > 3 {
		if !compareVertices(s.points[0], s.points[1], s.points[n-2], s.points[n-1]) {
			line(drawPixel, s.points[0], s.points[1])
			line(drawPixel, s.points[n-2], s.points[n-1])
		}
	}
	for i := range s.points {
		if (len(s.points) - i) < 4 {
			break
		}
		line(drawPixel, s.points[i], s.points[i+2])
		line(drawPixel, s.points[i+1], s.points[i+3])
	}
}

func compareVertices(p0, p1, p2, p3 Point) bool {
	return p0 == p2 && p1 == p3
}
