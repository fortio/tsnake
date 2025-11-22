// tsnake
// Play the classic game snake in the terminal
package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func main() {
	os.Exit(Main())
}

type bufferType uint

const (
	EMPTY bufferType = iota
	WASD
	ARROWKEYS
)

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fCpuprofile := flag.String("profile-cpu", "", "write cpu profile to `file`")
	fMemprofile := flag.String("profile-mem", "", "write memory profile to `file`")
	fps := flag.Float64("fps", 10, "set fps")
	halfFlag := flag.Bool("square", false, "use half height blocks so the snake's body is more square")
	relativeMovementFlag := flag.Bool("relative", false,
		"move the snake with a/d or left/right so that it turns relative to its current direction")
	cli.Main()
	if *fCpuprofile != "" {
		f, err := os.Create(*fCpuprofile)
		if err != nil {
			return log.FErrf("can't open file for cpu profile: %v", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return log.FErrf("can't start cpu profile: %v", err)
		}
		log.Infof("Writing cpu profile to %s", *fCpuprofile)
		defer pprof.StopCPUProfile()
	}
	draw := drawFull
	if *halfFlag {
		draw = drawHalf
	}
	ap := ansipixels.NewAnsiPixels(*fps)
	ap.TrueColor = *fTrueColor

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
	var s *snake
	ap.OnResize = func() error {
		h := ap.H
		if *halfFlag {
			h = ap.H * 2
		}
		s = newSnake(ap.W, h, *halfFlag)
		return nil
	}
	handleWasd, handleArrowKeys := handleWasd, handleArrowKeys
	if *relativeMovementFlag {
		handleWasd, handleArrowKeys = handleRelativeAD, handleRelativeArrowKeys
	}
	_ = ap.OnResize()
	var buffer byte
	bType := EMPTY
	err = ap.FPSTicks(func() bool {
		if len(ap.Data) > 0 && ap.Data[0] == 'q' {
			return false
		}
		switch bType {
		case EMPTY:
			if len(ap.Data) == 1 {
				handleWasd(s, ap.Data[0], &buffer, &bType, false)
			}
			if len(ap.Data) >= 3 {
				handleArrowKeys(s, ap.Data[2], &buffer, &bType, false)
			}

		case WASD:
			handleWasd(s, buffer, &buffer, &bType, true)
		case ARROWKEYS:
			handleArrowKeys(s, buffer, &buffer, &bType, true)
		}
		ap.ClearScreen()
		if !s.next() {
			return false
		}
		draw(ap, s)
		ap.WriteAt(0, 0, "%v %v %v", s.dir, s.firstFrame, bType)
		return true
	})
	if *fMemprofile != "" {
		f, errMP := os.Create(*fMemprofile)
		if errMP != nil {
			return log.FErrf("can't open file for mem profile: %v", errMP)
		}
		errMP = pprof.WriteHeapProfile(f)
		if errMP != nil {
			return log.FErrf("can't write mem profile: %v", err)
		}
		log.Infof("Wrote memory profile to %s", *fMemprofile)
		_ = f.Close()
	}
	if err != nil {
		log.Infof("Exiting on %v", err)
		return 1
	}
	return 0
}

const (
	ARROWUP = 65 + iota
	ARROWDOWN
	ARROWRIGHT
	ARROWLEFT
)

func handleRelativeAD(s *snake, dataValue byte, buffer *byte, bufferType *bufferType, clearBuffer bool) {
	switch dataValue {
	case 'a', 'A':
		switch s.dir {
		case U:
			if s.firstFrame && !s.square {
				*buffer = ARROWLEFT
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		case D:
			if s.firstFrame && !s.square {
				*buffer = ARROWLEFT
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		case L:
			s.dir = D
			s.firstFrame = true
		case R:
			s.dir = U
			s.firstFrame = true
		}
	case 'd', 'D':
		switch s.dir {
		case U:
			if s.firstFrame && !s.square {
				*buffer = ARROWRIGHT
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		case D:
			if s.firstFrame && !s.square {
				*buffer = ARROWRIGHT
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		case L:
			s.dir = U
			s.firstFrame = true
		case R:
			s.dir = D
			s.firstFrame = true
		}
	}
	if clearBuffer {
		*bufferType = EMPTY
	}
}

func handleRelativeArrowKeys(s *snake, dataValue byte, buffer *byte, bufferType *bufferType, clearBuffer bool) {
	switch dataValue {
	case ARROWLEFT:
		switch s.dir {
		case U:
			if s.firstFrame && !s.square {
				*buffer = ARROWLEFT
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		case D:
			if s.firstFrame && !s.square {
				*buffer = ARROWLEFT
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		case L:
			s.dir = D
			s.firstFrame = true
		case R:
			s.dir = U
			s.firstFrame = true
		}
	case ARROWRIGHT:
		switch s.dir {
		case U:
			if s.firstFrame && !s.square {
				*buffer = ARROWRIGHT
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		case D:
			if s.firstFrame && !s.square {
				*buffer = ARROWRIGHT
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		case L:
			s.dir = U
			s.firstFrame = true
		case R:
			s.dir = D
			s.firstFrame = true
		}
	}
	if clearBuffer {
		*bufferType = EMPTY
	}
}

func handleArrowKeys(s *snake, dataValue byte, buffer *byte, bufferType *bufferType, clearBuffer bool) {
	switch dataValue {
	case ARROWUP:
		if s.dir == R || s.dir == L {
			s.dir = U
			s.firstFrame = true
		}
	case ARROWDOWN:
		if s.dir == R || s.dir == L {
			s.dir = D
			s.firstFrame = true
		}
	case ARROWRIGHT:
		if s.dir == U || s.dir == D {
			if s.firstFrame && !s.square {
				*buffer = 67
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		}
	case ARROWLEFT:
		if s.dir == U || s.dir == D {
			if s.firstFrame && !s.square {
				*buffer = 68
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		}
	}
	if clearBuffer {
		*bufferType = EMPTY
	}
}

func handleWasd(s *snake, dataValue byte, buffer *byte, bufferType *bufferType, clearBuffer bool) {
	switch dataValue {
	case 'w', 'W':
		if s.dir == R || s.dir == L {
			s.dir = U
			s.firstFrame = true
		}
	case 's', 'S':
		if s.dir == R || s.dir == L {
			s.dir = D
			s.firstFrame = true
		}
	case 'd', 'D':
		if s.dir == U || s.dir == D {
			if s.firstFrame && !s.square {
				*buffer = 67
				*bufferType = ARROWKEYS
			} else {
				s.dir = R
			}
		}
	case 'a', 'A':
		if s.dir == U || s.dir == D {
			if s.firstFrame && !s.square {
				*buffer = 68
				*bufferType = ARROWKEYS
			} else {
				s.dir = L
			}
		}
	}
	if clearBuffer {
		*bufferType = EMPTY
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
