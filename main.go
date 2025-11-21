// tsnake
// Play the classic game snake in the terminal

package main

import (
	"flag"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func main() {
	fps := flag.Float64("fps", 60, "set fps")
	halfFlag := flag.Bool("square", false, "use half height blocks so the snake's body is more square")
	flag.Parse()
	draw := drawFull
	if *halfFlag {
		draw = drawHalf
	}
	ap := ansipixels.NewAnsiPixels(*fps)
	err := ap.Open()
	if err != nil {
		panic("error opening terminal")
	}
	defer func() {
		ap.Restore()
		ap.ShowCursor()
		ap.MoveCursor(0, 0)
	}()
	ap.SyncBackgroundColor()
	ap.HideCursor()
	ap.ClearScreen()
	h := ap.H
	if *halfFlag {
		h = ap.H * 2
	}
	s := newSnake(ap.W, h)
	_ = ap.FPSTicks(func() bool {
		if len(ap.Data) > 0 && ap.Data[0] == 'q' {
			return false
		}
		if len(ap.Data) >= 3 {
			handleInput(s, ap.Data[2])
		}
		ap.ClearScreen()
		if !s.next() {
			return false
		}
		draw(ap, s)
		return true
	})
}

func handleInput(s *snake, dataValue byte) {
	switch dataValue {
	case 65:
		if s.dir == R || s.dir == L {
			s.dir = U
		}
	case 66:
		if s.dir == R || s.dir == L {
			s.dir = D
		}
	case 67:
		if s.dir == U || s.dir == D {
			s.dir = R
		}
	case 68:
		if s.dir == U || s.dir == D {
			s.dir = L
		}
	}
}

func drawFull(ap *ansipixels.AnsiPixels, s *snake) {
	mouthCoords := s.snake[len(s.snake)-1]
	ap.WriteAt(mouthCoords.X, mouthCoords.Y, "%s ", ap.ColorOutput.Background(tcolor.Red.Color()))
	foodCoords := s.food
	ap.WriteAt(foodCoords.X, foodCoords.Y, "%s ", ap.ColorOutput.Background(tcolor.Green.Color()))
	for _, coords := range s.snake[:len(s.snake)-1] {
		ap.WriteAt(coords.X, coords.Y, "%s ", ap.ColorOutput.Background(tcolor.White.Color()))
	}
	ap.WriteString(tcolor.Reset)
}

type pixel struct {
	top, bottom           bool
	topColor, bottomColor tcolor.Color
}

func drawHalf(ap *ansipixels.AnsiPixels, s *snake) {
	pix := make(map[coords]*pixel)
	color := tcolor.White.Color()
	l := len(s.snake)
	for i, coords := range s.snake {
		if i == l-1 {
			color = tcolor.Red.Color()
		}
		if coords.Y%2 == 0 {
			coords.Y /= 2
			if pix[coords] == nil {
				pix[coords] = &pixel{}
			}
			pix[coords].top = true
			pix[coords].topColor = color
		} else {
			coords.Y /= 2
			if pix[coords] == nil {
				pix[coords] = &pixel{}
			}
			pix[coords].bottom = true
			pix[coords].bottomColor = color
		}
	}
	fy := coords{s.food.X, s.food.Y / 2}
	if pix[fy] == nil {
		pix[fy] = &pixel{}
	}
	if s.food.Y%2 == 0 {
		pix[fy].top = true
		pix[fy].topColor = tcolor.Green.Color()
	} else {
		pix[fy].bottom = true
		pix[fy].bottomColor = tcolor.Green.Color()
	}
	drawPixels(ap, pix)
	ap.WriteString(tcolor.Reset)
}

func drawPixels(ap *ansipixels.AnsiPixels, pix map[coords]*pixel) {
	var char rune
	var bg, fg tcolor.Color
	for coords, pixel := range pix {
		switch {
		case pixel.top && pixel.bottom:
			if pixel.topColor == pixel.bottomColor {
				char = ' '
				bg = pixel.topColor
				fg = pixel.topColor
			} else {
				char = ansipixels.BottomHalfPixel
				bg = pixel.topColor
				fg = pixel.bottomColor
			}
		case pixel.top:
			char = ansipixels.TopHalfPixel
			bg = ap.Background.Color()
			fg = pixel.topColor
		case pixel.bottom:
			char = ansipixels.BottomHalfPixel
			bg = ap.Background.Color()
			fg = pixel.bottomColor
		default:
			continue
		}
		ap.MoveCursor(coords.X, coords.Y)
		ap.WriteBg(bg)
		ap.WriteFg(fg)
		ap.WriteRune(char)
	}
}
