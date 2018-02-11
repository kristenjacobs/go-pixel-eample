package main

import (
	"image/color"
	"os/exec"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"
)

const (
	windowWidthPixels  = 1024
	windowHeightPixels = 768
	spriteSizePixels   = 16
	numX               = windowWidthPixels / spriteSizePixels
	numY               = windowHeightPixels / spriteSizePixels
)

func playSound() {
	go func() {
		cmd := exec.Command("paplay", "./resources/sample.wav")
		_, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
	}()
}

func draw(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int, color color.RGBA) {
	imd.Color = color
	x1 := float64(xPos * spriteSizePixels)
	y1 := float64(yPos * spriteSizePixels)
	x2 := x1 + spriteSizePixels
	y2 := y1 + spriteSizePixels
	imd.Push(pixel.V(x1, y1), pixel.V(x2, y2))
	imd.Rectangle(0)
	imd.Draw(win)
}

func move(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int,
	f func(xPos int, yPos int) (int, int)) (int, int) {
	draw(win, imd, xPos, yPos, colornames.Black)
	xPos, yPos = f(xPos, yPos)
	draw(win, imd, xPos, yPos, colornames.White)
	playSound()
	return xPos, yPos

}

func left(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int) (int, int) {
	return move(win, imd, xPos, yPos,
		func(xPos int, yPos int) (int, int) {
			xPos--
			if xPos < 0 {
				xPos = numX - 1
			}
			return xPos, yPos
		})
}

func right(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int) (int, int) {
	return move(win, imd, xPos, yPos,
		func(xPos int, yPos int) (int, int) {
			xPos++
			if xPos >= numX {
				xPos = 0
			}
			return xPos, yPos
		})
}

func up(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int) (int, int) {
	return move(win, imd, xPos, yPos,
		func(xPos int, yPos int) (int, int) {
			yPos++
			if yPos >= numY {
				yPos = 0
			}
			return xPos, yPos
		})
}

func down(win *pixelgl.Window, imd *imdraw.IMDraw, xPos int, yPos int) (int, int) {
	return move(win, imd, xPos, yPos,
		func(xPos int, yPos int) (int, int) {
			yPos--
			if yPos < 0 {
				yPos = numY - 1
			}
			return xPos, yPos
		})
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Go Pixel Example",
		Bounds: pixel.R(0, 0, windowWidthPixels, windowHeightPixels),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	xPos := 0
	yPos := 0

	imd := imdraw.New(nil)
	imd.Color = colornames.White

	draw(win, imd, xPos, yPos, colornames.White)

	for !win.Closed() {
		if win.Pressed(pixelgl.KeyLeft) {
			xPos, yPos = left(win, imd, xPos, yPos)
		}
		if win.Pressed(pixelgl.KeyRight) {
			xPos, yPos = right(win, imd, xPos, yPos)
		}
		if win.Pressed(pixelgl.KeyUp) {
			xPos, yPos = up(win, imd, xPos, yPos)
		}
		if win.Pressed(pixelgl.KeyDown) {
			xPos, yPos = down(win, imd, xPos, yPos)
		}
		win.Update()
		time.Sleep(25 * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}
