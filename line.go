package main

import (
	"image"
	"image/draw"
)

func ants(dst draw.Image, q0, q1 image.Point, src image.Image, sp image.Point, thick int, stride int) {
	if q1.X < q0.X {
		q0, q1 = q1, q0
	}
	p0 := q0.Sub(q0)
	p1 := q1.Sub(q0)
	dy := p1.Y
	dx := p1.X
	s := 1
	if dy < 0 {
		dy = -dy
		s = -1
	}
	if thick == 0 {
		thick++
	}
	if thick > 0 {
		thick--
	}
	r := image.ZR
	//Ellipse(dst, q0, src, thick/2, thick/2, 1, q0, 0, 0)
	strideX := stride / 2
	X := 0
	drawstride := func() {
		if X < strideX {
			draw.Draw(dst, r, src, r.Min, draw.Src)
		}
		X++
		if X > stride {
			X = 0
		}
	}
	if dx > dy {
		d := 2*dy - dx
		for p0.X < dx {
			r = image.Rect(p0.X+q0.X, s*p0.Y+q0.Y-thick/2, p0.X+q0.X+1, s*p0.Y+q0.Y+1+thick/2)
			//draw.Draw(dst, r, src, r.Min, draw.Src)
			drawstride()
			if d > 0 {
				p0.Y++
				d -= dx
			}
			p0.X++
			d += dy
		}

	} else {
		d := 2*dx - dy
		for p0.Y < dy {
			r = image.Rect(p0.X+q0.X-thick/2, s*p0.Y+q0.Y, p0.X+q0.X+1+thick/2, s*p0.Y+q0.Y+1)
			drawstride() // draw.Draw(dst, r, src, r.Min, draw.Src)
			if d > 0 {
				p0.X++
				d -= dy
			}
			p0.Y++
			d += dx
		}
	}
	q1 = image.Pt(p0.X+q0.X, s*p0.Y+q0.Y)
	//Ellipse(dst, q1, src, thick/2, thick/2, 1, q1, 0, 0)
}

func CenterRect(r image.Rectangle, pt image.Point) image.Rectangle {
	return r.Sub(r.Min.Sub(pt.Sub(r.Size().Div(2))))
}

func controlPts(dst draw.Image, r, box image.Rectangle, src image.Image, sp image.Point) {
	dx2 := r.Dx() / 2
	dy2 := r.Dy() / 2
	if dx2 == 0 || dy2 == 0 {
		return
	}
	for x := r.Min.X; x <= r.Max.X; x += dx2 {
		for y := r.Min.Y; y <= r.Max.Y; y += dy2 {
			draw.Draw(dst, CenterRect(box, image.Pt(x, y)), src, boxr.Min, draw.Src)
		}
	}
}

func line(dst draw.Image, q0, q1 image.Point, src image.Image, sp image.Point, thick int) {
	if q1.X < q0.X {
		q0, q1 = q1, q0
	}
	p0 := q0.Sub(q0)
	p1 := q1.Sub(q0)
	dy := p1.Y
	dx := p1.X
	s := 1
	if dy < 0 {
		dy = -dy
		s = -1
	}
	if thick == 0 {
		thick++
	}
	if thick > 0 {
		thick--
	}
	Ellipse(dst, q0, src, thick/2, thick/2, 1, q0, 0, 0)
	if dx > dy {
		d := 2*dy - dx
		for p0.X < dx {
			r := image.Rect(p0.X+q0.X, s*p0.Y+q0.Y-thick/2, p0.X+q0.X+1, s*p0.Y+q0.Y+1+thick/2)
			draw.Draw(dst, r, src, r.Min, draw.Src)
			if d > 0 {
				p0.Y++
				d -= dx
			}
			p0.X++
			d += dy
		}

	} else {
		d := 2*dx - dy
		for p0.Y < dy {
			r := image.Rect(p0.X+q0.X-thick/2, s*p0.Y+q0.Y, p0.X+q0.X+1+thick/2, s*p0.Y+q0.Y+1)
			draw.Draw(dst, r, src, r.Min, draw.Src)
			if d > 0 {
				p0.X++
				d -= dy
			}
			p0.Y++
			d += dx
		}
	}
	q1 = image.Pt(p0.X+q0.X, s*p0.Y+q0.Y)
	Ellipse(dst, q1, src, thick/2, thick/2, 1, q1, 0, 0)
}

func polyline(dst draw.Image, p Polygon, src image.Image, sp image.Point, thick int) {
	pts := p.Points()
	if len(pts) < 2 {
		return
	}
	for i := 1; i < len(pts); i++ {
		line(dst, pts[i-1], pts[i], src, sp, thick)
	}
}
