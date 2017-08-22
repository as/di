package main

import (
	draw2 "golang.org/x/image/draw"
	"image"
	"image/draw"
	_ "image/png"
	"os"
)

type Image struct {
	*Rect
	img image.Image
}

var kern = draw2.Interpolator(draw2.CatmullRom)

func (im *Image) Draw(dst draw.Image) {
	draw.Draw(dst, im.Bounds(), im.img, im.img.Bounds().Min, draw.Over)
	if im.Rect.Node != nil {
		im.Node.Draw(dst)
	}
}

func NewImage(path string, sp image.Point) (*Image, error) {
	img, err := readimage(path)
	if err != nil {
		return nil, err
	}
	return &Image{
		&Rect{
			bg:        Cyan,
			Rectangle: img.Bounds().Add(sp),
		},
		img,
	}, nil
}

func NewImageScaled(path string, r image.Rectangle) (*Image, error) {
	img, err := readimage(path)
	if err != nil {
		return nil, err
	}
	img2 := Scale(img, r)
	return &Image{
		&Rect{
			bg:        Cyan,
			Rectangle: r,
		},
		img2,
	}, nil
}

func NewImageEX(path string, r image.Rectangle) (*Image, error) {
	img, err := readimage(path)
	if err != nil {
		return nil, err
	}
	k := kern
	kern = draw2.NearestNeighbor
	pixel := Scale(img, image.Rect(0, 0, 32, 32))
	img2 := Scale(pixel, r)
	kern = k
	return &Image{
		&Rect{
			bg:        Cyan,
			Rectangle: r,
		},
		img2,
	}, nil
}

func readimage(path string) (image.Image, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	img, _, err := image.Decode(fd)
	return img, err
}

func Scale(img image.Image, r image.Rectangle) draw.Image {
	dst := image.NewRGBA(r)
	kern.Scale(dst, dst.Bounds(), img, img.Bounds(), draw2.Src, nil)
	return dst
}
