package main

import "image"

type Pad struct {
	dir int
	image.Point
}

type Polygon interface {
	Points() []image.Point
	Add(...image.Point) Polygon
}

func (p poly) Add(qq ...image.Point) Polygon {
	pp := p.Points()
	rr := make(poly, 0, len(pp))
	if len(qq) == 0 {
		return append(poly{})
	}
	for i := range pp {
		rr = append(rr, p[i].Add(qq[i%(len(qq))]))
	}
	return rr
}

func (p poly) Points() []image.Point {
	return []image.Point(p)
}

type poly []image.Point

const Smacks = 20
const Offset = 4
const (
	L = 0
	R = 1
	B = 2
	T = 3
)

var dirs = [4][4]int{
	{T, B, L, R}, // Left
	{B, T, R, L}, // Right
	{R, L, T, B}, // Bottom
	{L, R, B, T}, // Top
}

func Connect2(r0, r1 image.Rectangle, out, in Pad) Polygon {
	in.Point = transform(out, in.Point)
	pts := connect(out, in)
	p := pts.Points()
	for i, v := range p {
		p[i] = invert(out, v)
	}
	pts = correct(r0, poly(p))
	pts = correct(r1, pts)
	pts = correct(r0, pts)
	return pts
}

func Connect(out, in Pad) Polygon {
	in.Point = transform(out, in.Point)
	p1 := connect(out, in)
	p := make(poly, len(p1.Points()))
	for i, v := range p1.Points() {
		p[i] = invert(out, v)
	}
	return p
}
func connect(out, in Pad) Polygon {
	switch dirs[out.dir][in.dir] {
	case L:
		if in.Y > 0 {
			if in.X > 0 {
				println("Left Shape3P")
				return Shape3P(in.Point)
			}
			println("Left Shape5B")
			return Shape5PB(in.Point, -1)
		}
		println("Left Shape5B (TODO)")
		return Shape5PB(in.Point, -1)
		// TODO
	case R:
		if in.Y > 0 {
			if in.X > 0 {
				return Shape5PB(in.Point, 1)
			}
			return Shape3P(in.Point)
		}
		return Shape5PB(in.Point, 1)
	case B:
		if in.Y > 0 {
			return Shape4PD(in.Point)
		}
		return Shape6P(in.Point, -1)
	case T:
		if in.X < -Smacks || in.X > Smacks {
			return Shape4PU(in.Point)
		}
		return Shape6P(in.Point, 1) // TODO
	}
	panic("connect: never happens")
}

func transform(o Pad, p image.Point) image.Point {
	switch o.dir {
	case L:
		return image.Pt(p.Y-o.Y, o.X-p.X)
	case R:
		return image.Pt(o.Y-p.Y, p.X-o.X)
	case B:
		return image.Pt(o.X-p.X, o.Y-p.Y)
	case T:
		return image.Pt(p.X-o.X, p.Y-o.Y)
	}
	panic("transform: never happens")
}
func invert(o Pad, p image.Point) image.Point {
	switch o.dir {
	case L:
		return image.Pt(o.X-p.Y, o.Y+p.X)
	case R:
		return image.Pt(o.X+p.Y, o.Y-p.X)
	case B:
		return image.Pt(o.X-p.X, o.Y-p.Y)
	case T:
		return image.Pt(o.X+p.X, o.Y+p.Y)
	}
	panic("invert: never happens")
}
func intersect(r image.Rectangle, p0, p1 image.Point) bool {
	return code(r, p0)&code(r, p1) == 0
}
func correct(r image.Rectangle, p Polygon) Polygon {
	pts := p.Points()
	for i := 2; i+1 < len(pts); i++ {
		a, b, c, d := pts[i-2], pts[i-1], pts[i], pts[i+1]
		if intersect(r, b, c) {

			if b.X == c.X {
				b.X, c.X = adjust(r.Min.X, r.Max.X, a.X, b.X, c.X, d.X)
			} else {
				b.Y, c.Y = adjust(r.Min.Y, r.Max.Y, a.Y, b.Y, c.Y, d.Y)
			}
		}
		pts[i-1], pts[i] = b, c
	}
	return poly(pts)
}

func eqsign(a, b, c, d int) bool {
	return (b-a)*(d-c) > 0
}
func adjust(min, max, a, b, c, d int) (int, int) {
	if eqsign(a, b, c, d) {
		if b < a {
			return min - Offset + 2, min - Offset + 2
		}
		return max + Offset - 3, max + Offset - 3
	}
	if c < d {
		return min - Offset + 2, min - Offset + 2
	}
	return max + Offset - 3, max + Offset - 3
}
func in(r image.Rectangle, p0, p1 image.Point) bool {
	return code(r, p0)&code(r, p1) == 0
}
func code(r image.Rectangle, p image.Point) int {
	c := 0
	if p.X < r.Min.X {
		c = 1
	} else if p.X >= r.Max.X {
		c = 2
	}
	if p.Y < r.Min.Y {
		c |= 4
	} else if p.Y >= r.Max.Y {
		c |= 8
	}
	return c
}

func Shape3P(sp image.Point) Polygon {
	return poly{
		image.ZP,
		image.Pt(0, sp.Y),
		sp,
	}
}
func Shape4PD(sp image.Point) Polygon {
	p := make(poly, 4)
	p[1] = image.Pt(0, Offset)
	p[2] = image.Pt(sp.X, Offset)
	p[3] = sp
	return p
}
func Shape4PU(sp image.Point) Polygon {
	p := make(poly, 4)
	p[1] = image.Pt(0, max(0, sp.Y)+Offset)
	p[2] = image.Pt(sp.X, p[1].Y)
	p[3] = sp
	return p
}
func Shape5PB(sp image.Point, s int) Polygon {
	p := make(poly, 5)
	p[1] = image.Pt(0, Offset)
	if s*sp.X < -Smacks {
		p[2] = image.Pt(sp.X+s+Offset, Offset)
	} else {
		p[2] = image.Pt(s+Offset, Offset)
		if s*sp.X >= 0 {
			p[2].X += sp.X
		}
	}
	p[3] = image.Pt(p[2].X, sp.Y)
	p[4] = sp
	return p

}
func Shape5PT(sp image.Point, s int) Polygon {
	p := make(poly, 5)
	p[1] = image.Pt(0, Offset)
	p[2] = image.Pt(sp.X+s*Offset, Offset)
	p[3] = image.Pt(p[2].X, sp.Y)
	p[4] = sp
	return p
}
func Shape6P(sp image.Point, s int) Polygon {
	p := make(poly, 6)
	m := sp.X / 2
	p[1] = image.Pt(0, Offset)
	if m > Smacks {
		p[2] = image.Pt(m, Offset)
	} else {
		s2 := 1
		if sp.Y < 0 {
			s2 = -s2
		}
		if sp.X > 0 {
			s2 = -s2
		}
		p[2] = image.Pt(s2*Smacks, Offset)
	}
	p[3] = image.Pt(p[2].X, sp.Y+s*Offset)
	p[4] = image.Pt(sp.X, p[3].Y)
	p[5] = sp
	return p
}

func pick(c bool, a, b int) int {
	if c {
		return a
	}
	return b
}
