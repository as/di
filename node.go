package main

import (
	"image"
	"image/color"
	"image/draw"
)
import "github.com/as/frame/win"
type Node interface {
	Bounds() image.Rectangle
	Draw(dst draw.Image)
	Move(sp image.Point)
	Kid() []Node
}

type Di struct {
	dx, dy int
	bg, fg image.Image
	Node   Node
}

type Rect struct {
	bg, fg image.Image
	image.Rectangle
	fill   int
	border int
	Node   Node
}

type Circ struct {
	bg, fg     image.Image
	a, b       int
	c          image.Point
	alpha, phi int
	Node       Node
}

type Line struct {
	bg, fg image.Image
	p0, p1 image.Point
	thick  int
}

type Text struct {
	dy     int
	bg, fg color.Color
	sp     image.Point
	s      *win.Win
}

type Link struct {
	id    []int
	arrow int
}

type Dot struct {
	sp image.Point
	id int
}

type Nest struct {
	m []Node
}

var Op = draw.Src

func (n *Di) Draw(dst draw.Image) {
	r := image.Rect(0, 0, n.dx, n.dy)
	draw.Draw(dst, r, n.bg, image.ZP, Op)
	n.Node.Draw(dst)
}
func (n *Rect) Draw(dst draw.Image) {
	draw.Draw(dst, n.Rectangle, n.bg, image.ZP, Op)
	if n.Node == nil {
		return
	}
	n.Node.Draw(dst)
}
func (n *Circ) Draw(dst draw.Image) {
	Ellipse(dst, n.c, n.bg, n.a, n.b, 1, image.ZP, 0, 0)
}
func (n *Text) Draw(dst draw.Image) {
	n.s.Refresh()
	draw.Draw(dst, n.s.Loc(), n.s.RGBA(), image.ZP, draw.Src)
}
func (n *Link) Draw(dst draw.Image) {}
func (n *Dot) Draw(dst draw.Image)  {}
func (n *Nest) Draw(dst draw.Image) {
	for _, v := range n.m {
		if v != nil {
			v.Draw(dst)
		}
	}
}
func (n *Line) Draw(dst draw.Image) {
	line(dst, n.p0, n.p1, n.bg, image.ZP, n.thick)
}

func (n *Di) Move(sp image.Point) {
	if n.Node == nil {
		return
	}
	return
	n.Node.Move(sp)
}
func (n *Rect) Move(sp image.Point) {
	r := n.Rectangle
	size := image.Pt(r.Dx(), r.Dy())
	n.Rectangle = image.Rectangle{sp, sp.Add(size)}
	return
	if n.Node == nil {
		return
	}
	n.Node.Move(sp)
}
func (n *Line) Move(sp image.Point) {
	r := image.Rectangle{n.p0, n.p1}
	size := image.Pt(r.Dx(), r.Dy())
	n.p0 = sp
	n.p1 = sp.Add(size)
}

func (n *Circ) Move(sp image.Point) {
	n.c = sp
}
func (n *Text) Move(sp image.Point) {
	n.sp = sp
	n.s.Move(sp)
}
func (n *Link) Move(sp image.Point) {
}
func (n *Dot) Move(sp image.Point) {}
func (n *Nest) Move(sp image.Point) {
	return
	for _, n := range n.m {
		if n != nil {
			n.Move(sp)
		}
	}
}

func Shift(n Node, dp image.Point){
	if n == nil{
		return
	}
	n.Move(n.Bounds().Min.Add(dp))
	for _,n := range n.Kid(){
		Shift(n, dp)
	}
}

func Move(n Node, sp image.Point){
	if n == nil{
		return
	}
	n.Move(sp)
	for _,n := range n.Kid(){
		Move(n, sp)
	}
}
func Handle(n Node, e interface{}){
	if n == nil{
		return
	}
	switch n := n.(type){
	case Handler:
		n.Handle(e)
		return
	}
	for _, n := range n.Kid() {
		switch n := n.(type){
		case Handler:
			n.Handle(e)
			return
		}		
	}
}
func (n *Di) Bounds() image.Rectangle {
	return image.Rect(0, 0, n.dx, n.dy)
}
func (n *Rect) Bounds() image.Rectangle {
	return n.Rectangle
}
func (n *Circ) Bounds() image.Rectangle {
	a, b := n.a, n.b
	return image.Rect(-a, -b, a, b).Add(n.c)
}
func (n *Text) Bounds() image.Rectangle { return n.s.Loc() }
func (n *Link) Bounds() image.Rectangle { return image.ZR }
func (n *Dot) Bounds() image.Rectangle  { return image.ZR }

func (n *Nest) Bounds() image.Rectangle {
	r := image.ZR
	for _, n := range n.m {
		if n == nil {
			continue
		}
		r0 := n.Bounds()
		r.Min.X = min(r.Min.X, r0.Min.X)
		r.Min.Y = min(r.Min.Y, r0.Min.Y)
		r.Max.X = max(r.Max.X, r0.Max.X)
		r.Max.Y = max(r.Max.Y, r0.Max.Y)
	}
	return r
}
func (n *Line) Bounds() image.Rectangle {
	return image.Rectangle{n.p0, n.p1}
}

//func (n *Rect) Handle(e interface{}){
//	Handle(n, e)
//}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
