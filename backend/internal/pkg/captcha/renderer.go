package captcha

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math/rand"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	base64Captcha "github.com/mojocn/base64Captcha"
	"golang.org/x/image/font"
)

var (
	charBgColor    = color.RGBA{248, 249, 251, 255}
	charTextColors = []color.RGBA{
		{21, 101, 192, 255},
		{198, 40, 40, 255},
		{46, 125, 50, 255},
		{109, 58, 170, 255},
		{200, 101, 21, 255},
		{0, 131, 143, 255},
	}
	charNoiseColor = color.RGBA{180, 186, 196, 255}
	charLineColors = []color.RGBA{
		{100, 149, 237, 255},
		{229, 115, 115, 255},
		{129, 199, 132, 255},
		{186, 104, 200, 255},
		{255, 183, 77, 255},
		{77, 208, 225, 255},
	}

	charFonts = []*truetype.Font{
		base64Captcha.DefaultEmbeddedFonts.LoadFontByName("fonts/RitaSmith.ttf"),
		base64Captcha.DefaultEmbeddedFonts.LoadFontByName("fonts/actionj.ttf"),
		base64Captcha.DefaultEmbeddedFonts.LoadFontByName("fonts/chromohv.ttf"),
	}
)

type charDriver struct {
	height int
	width  int
	length int
	source string
	math   bool
	fonts  []*truetype.Font
}

func newCharDriver(height, width, length int, source string, mathMode bool) *charDriver {
	return &charDriver{height: height, width: width, length: length, source: source, math: mathMode, fonts: charFonts}
}

func (d *charDriver) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = base64Captcha.RandomId()
	if d.math {
		content, answer = randomMathQuestion()
		return id, content, answer
	}
	content = base64Captcha.RandText(d.length, d.source)
	return id, content, content
}

func (d *charDriver) DrawCaptcha(content string) (base64Captcha.Item, error) {
	img := image.NewNRGBA(image.Rect(0, 0, d.width, d.height))
	draw.Draw(img, img.Bounds(), image.NewUniform(charBgColor), image.Point{}, draw.Src)

	drawNoiseDots(img, d.width, d.height)
	drawInterferenceLines(img, d.width, d.height)
	if err := drawCaptchaText(img, content, d.fonts, d.width, d.height); err != nil {
		return nil, err
	}
	return &charItem{img: img}, nil
}

func randomMathQuestion() (question, answer string) {
	switch rand.Intn(3) {
	case 0:
		a, b := rand.Intn(20), rand.Intn(20)
		return fmt.Sprintf("%d+%d=?", a, b), fmt.Sprintf("%d", a+b)
	case 1:
		a, b := rand.Intn(10), rand.Intn(10)
		return fmt.Sprintf("%dx%d=?", a, b), fmt.Sprintf("%d", a*b)
	default:
		a := rand.Intn(89) + 10
		b := rand.Intn(a)
		return fmt.Sprintf("%d-%d=?", a, b), fmt.Sprintf("%d", a-b)
	}
}

func drawCaptchaText(img *image.NRGBA, content string, fonts []*truetype.Font, w, h int) error {
	runes := []rune(content)
	if len(runes) == 0 {
		return errors.New("captcha: empty content")
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetHinting(font.HintingFull)

	step := w / len(runes)
	fontSize := float64(h) * 0.58
	for i, r := range runes {
		c.SetFont(fonts[rand.Intn(len(fonts))])
		c.SetFontSize(fontSize)
		c.SetSrc(image.NewUniform(charTextColors[rand.Intn(len(charTextColors))]))
		x := step*i + step/6 + rand.Intn(step/8)
		y := int(float64(h)*0.72) + rand.Intn(h/12) - h/24
		if _, err := c.DrawString(string(r), freetype.Pt(x, y)); err != nil {
			return err
		}
	}
	return nil
}

func drawNoiseDots(img *image.NRGBA, w, h int) {
	count := w * h / 160
	for range count {
		fillCircle(img, rand.Intn(w), rand.Intn(h), rand.Intn(2)+1, charNoiseColor)
	}
}

func drawInterferenceLines(img *image.NRGBA, w, h int) {
	for range 3 {
		drawLine(img, rand.Intn(w/3), rand.Intn(h), w-1-rand.Intn(w/3), rand.Intn(h),
			charLineColors[rand.Intn(len(charLineColors))])
	}
}

func fillCircle(img *image.NRGBA, cx, cy, r int, c color.RGBA) {
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				img.Set(cx+dx, cy+dy, c)
			}
		}
	}
}

func drawLine(img *image.NRGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := absInt(x2 - x1)
	dy := absInt(y2 - y1)
	steps := max(dx, dy)
	if steps == 0 {
		return
	}
	for i := 0; i <= steps; i++ {
		x := x1 + (x2-x1)*i/steps
		y := y1 + (y2-y1)*i/steps
		img.Set(x, y, c)
		img.Set(x+1, y, c)
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type charItem struct {
	img *image.NRGBA
}

func (it *charItem) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, it.img); err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

func (it *charItem) EncodeB64string() string {
	var buf bytes.Buffer
	if err := png.Encode(&buf, it.img); err != nil {
		return ""
	}
	return fmt.Sprintf("data:%s;base64,%s", base64Captcha.MimeTypeImage, base64.StdEncoding.EncodeToString(buf.Bytes()))
}
