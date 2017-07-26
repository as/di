package main

import (
	"image"

	"github.com/as/frame"
	"github.com/as/frame/win"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type aux struct {
	scr    screen.Screen
	events screen.Window
}

func NewText(a *aux, s string, dy int, sp, size, pad image.Point) *Text {
	cols := frame.Acme
	w := win.New(a.scr, frame.NewTTF(gomono.TTF, dy), a.events, sp, size, pad, cols)
	w.Insert([]byte(s), 0)
	w.Upload()
	return &Text{sp: sp, s: w}
}

type Handler interface {
	Handle(e interface{})
}

func (t *Text) Handle(e interface{}) {
	switch e := e.(type) {
	case mouse.Event:
		pt := image.Pt(int(e.X), int(e.Y)).Sub(t.s.Loc().Min)
		if e.Direction == 1 && e.Button == 1 {
			q := t.s.IndexOf(pt)
			t.s.Select(q,q)
		}
	case key.Event:
		if e.Direction == 2{
			return
		}
		q0, q1 := t.s.Dot()
		if e.Rune == '\x08' && q0 != 0 {
			q0--
		}
		if q0 != q1 {
			t.s.Delete(q0, q1)
		}
		if e.Rune != '\x08' {
			q0 += t.s.Insert([]byte{byte(e.Rune)},q0)
			t.s.Select(q0,q0)
		}
		t.s.Refresh()
	}
}
