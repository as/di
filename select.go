package main

import "image"

func Select(pt image.Point, root *Di) Node {
	if root == nil {
		return nil
	}
	return select0(pt, root)
}

func select0(pt image.Point, x Node) (a Node) {
	if x == nil {
		return nil
	}
	switch x := x.(type) {
	case *Nest:
	case Node:
		if !pt.In(x.Bounds()) {
			return nil
		}
		a = x
	}
	if x.Kid() == nil {
		return x
	}
	for _, y := range x.Kid() {
		z := select0(pt, y)
		if z != nil {
			return z
		}
	}
	return a
}
func (n *Di) Kid() []Node {
	if n.Node == nil {
		return nil
	}
	return []Node{n.Node}
}
func (n *Rect) Kid() []Node {
	if n.Node == nil {
		return nil
	}
	return []Node{n.Node}
}
func (n *Nest) Kid() []Node {
	return n.m
}

func (n *Circ) Kid() []Node { return []Node{n.Node} }
func (n *Text) Kid() []Node { return nil}
func (n *Link) Kid() []Node { panic("!") }
func (n *Dot) Kid() []Node  { panic("!") }
func (n *Line) Kid() []Node { return nil }
