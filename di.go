package main

import (
	//	"github.com/as/clip"
	//"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"fmt"
	"sync"
	
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	winSize = image.Pt(1024, 768)
	wg      sync.WaitGroup
)
var (
	Yellow   = image.NewUniform(color.RGBA{255, 255, 224, 255})
	Mauve    = image.NewUniform(color.RGBA{0x99, 0x99, 0xDD, 255})
	Blue     = image.NewUniform(color.RGBA{0, 0, 255, 255})
	EggShell = image.NewUniform(color.RGBA{128, 128, 128, 255})
	Black    = image.NewUniform(color.RGBA{0, 0, 0, 255})
	Gray     = image.NewUniform(color.RGBA{16, 16, 16, 255})
	White    = image.NewUniform(color.RGBA{255, 255, 255, 255})
	Cyan     = image.NewUniform(color.RGBA{234, 255, 255, 255})
)

type Mouse struct {
	bt   int
	down int
	mod  int
	dir  int
	pt   image.Point
}

var M Mouse

func parsemouse(e mouse.Event) {
	M.pt = image.Pt(int(e.X), int(e.Y))
	switch e.Direction {
	case mouse.DirRelease:
		M.dir = int(mouse.DirRelease)
		M.bt = int(e.Button)
		M.down &^= 1 << uint(M.bt)
	case mouse.DirPress:
		M.dir = int(mouse.DirPress)
		M.bt = int(e.Button)
		M.down |= 1 << uint(e.Button)
	default:
		M.bt = int(mouse.ButtonNone)
		M.dir = int(mouse.DirNone)
	}
}

func (i Item) Canon() Item {
	i.Rectangle = i.Rectangle.Canon()
	return i
}

type Item struct {
	image.Rectangle
	dirty bool
}

func attachline(n Node, pt0, pt1 image.Point) {
	switch t := n.(type) {
	case *Nest:
		t.m = append(t.m, &Line{bg: Yellow, p0: pt0, p1: pt1, thick: 30})
	}
}

var boxr = image.Rect(0, 0, 50, 50)
var a *aux

func node(name string, sp image.Point) *Rect {
	r := boxr.Add(sp)
	return &Rect{
		bg:        Cyan,
		Rectangle: r,
		Node:      NewText(a, name, 15, r.Min.Add(image.Pt(0, r.Dy()+2)), image.Pt(r.Dx(), 25), image.Pt(2, 2)),
	}
}

func nodeCPU(name string, sp image.Point) Node {
	im, err := NewImageScaled(`cpu.png`, 
		image.Rectangle{sp, sp.Add(image.Pt(55,55))})
	if err != nil {
		log.Printf("nodeCPU: %s\n", err)
		return nil
	}
	r := im.Rect.Bounds()
	im.Rect.Node = NewText(a, name, 15, r.Min.Add(image.Pt(0, r.Dy()+2)), image.Pt(r.Dx(), 25), image.Pt(2, 2))
	return im
}

func root(sp image.Point) Node{
	im, err := NewImageEX(`root.png`, 
		image.Rectangle{sp, sp.Add(image.Pt(320,240))})
	if err != nil {
		log.Printf("nodeCPU: %s\n", err)
		return nil
	}
	//r := im.Rect.Bounds()
	//im.Rect.Node = NewText(a, "root", 15, r.Min.Add(image.Pt(0, r.Dy()+2)), image.Pt(r.Dx(), 25), image.Pt(2, 2))
	return im
}

func main() {
	driver.Main(func(src screen.Screen) {
		win, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y, "di"})
		focused := false
		focused = focused
		//		repaint := true
		buf, _ := src.NewBuffer(winSize)
		draw.Draw(buf.RGBA(), buf.Bounds(), Black, image.ZP, draw.Src)
		rt := root(image.Pt(100,100)) 
		dirty := true
		a = &aux{src, win}
		cnt := 25
		RPs := []Node{
			nodeCPU("W21", image.Pt(100, 0)),
			nodeCPU("W22", image.Pt(200, 0)),
			nodeCPU("W23", image.Pt(300, 0)),
			nodeCPU("W24", image.Pt(400, 0)),
		}
		nest1 := &Nest{
			[]Node{
				&Circ{bg: Black, a: 100, b: 100, c: image.Pt(50, 50)},
				&Rect{bg: Yellow, Rectangle: image.Rect(0, 0, 64, 64)},
				&Line{bg: Yellow, p0: image.Pt(0, 50), p1: image.Pt(100, 49)},
			},
		}
		nest1.m = append(nest1.m, rt)
		nest1.m = append(nest1.m, RPs...)

		di := &Di{
			dx: 800, dy: 600, bg: Black,
			Node: &Nest{[]Node{
				&Rect{bg: EggShell, Rectangle: image.Rect(0, 0, 1600, 1200),
					Node: nest1,
				},
				//&Rect{bg: Blue,Rectangle: image.Rect(150, 100, 250, 200)},
			},
			},
		}
		links := make(map[Node][]Node)
		links[rt] = RPs
		var lpt0, lpt1 image.Point
		act := Node(di)
		addplat := func(sp image.Point){
			plat := nodeCPU(fmt.Sprintf("W%02d", cnt), sp)
			nest1.m = append(nest1.m, plat)
			cnt++
			links[act] = append(links[act], plat)
		}
		var s0,s1 image.Point
		sr := image.NewUniform(color.RGBA{255,32,32,255})
		for {
			switch e := win.NextEvent().(type) {
			case key.Event:
				Handle(act, e);
				if dirty{
					win.Send(paint.Event{})
				}
			case mouse.Event:
				parsemouse(e)
				switch{
				case M.bt == 1<<3 || M.down == 1<<3:
					if M.dir == 1{
						s0 = M.pt
					} else if M.dir == 2{
						s0,s1 = image.ZP, image.ZP
						
					} else {
						s1 = M.pt
					}
					dirty=true
				case M.bt == 2 && M.dir == 1:
					addplat(M.pt)
					dirty=true
				case M.bt == 1 || M.down == 2:
					switch M.dir {
					case 1:
						act = Select(M.pt, di)
						lpt0 = M.pt
					case 2:
						//attachline(nest1, lpt0, lpt1)
					case 0:
						if act != nil {
							lpt1 = M.pt
							d := lpt1.Sub(lpt0)
							Shift(act, d)
							lpt0 = lpt1
						}
					}
					if act != nil && M.pt.In(act.Bounds()){
						Handle(act, e)
					}
					dirty=true
				}
				if dirty{
					win.Send(paint.Event{})
				}
				//repaint = true
			case size.Event:
			case paint.Event:
				di.Draw(buf.RGBA())
				if act != nil {
					act.Draw(buf.RGBA())
				}
				for k, v := range links {
					p0 := Pad{B, k.Bounds().Min.Add(image.Pt(k.Bounds().Dx()/2, 0))}
					for _, v := range v {
						p1 := Pad{B, v.Bounds().Min.Add(image.Pt(v.Bounds().Dx()/2, 0))}
						p := Connect2(v.Bounds(), p0, p1)
						polyline(buf.RGBA(), p, Mauve, image.ZP, 3)
					}
				}
				drawBorder(buf.RGBA(), image.Rect(s0.X,s0.Y,s1.X,s1.Y), sr, image.ZP, 4)
			
				win.Upload(image.ZP, buf, buf.Bounds())
				wg.Wait()
				win.Publish()
				dirty=false
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				// NT doesn't repaint the window if another window covers it
				if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff {
					focused = false
				} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn {
					focused = true
				}
			}
		}
	})
}
func drawBorder(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, thick int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}

//polyline(buf.RGBA(), poly{{0,0},{2,2},{5,5},{100,500},{777,787}}, Yellow, image.ZP, 3)
//polyline(buf.RGBA(), Connect(Pad{T,image.Pt(200,200)}, Pad{R,M.pt}), Yellow, image.ZP, 3)
//polyline(buf.RGBA(), Connect(Pad{T,image.Pt(500,500)}, Pad{L,M.pt}), Blue, image.ZP, 3)
//polyline(buf.RGBA(), Connect(Pad{T,image.Pt(500,100)}, Pad{B,M.pt}), Yellow, image.ZP, 3)
//polyline(buf.RGBA(), correct(br.Add(pt), Connect(Pad{T,image.Pt(300,100)}, Pad{T,M.pt})), Mauve, image.ZP, 1)
