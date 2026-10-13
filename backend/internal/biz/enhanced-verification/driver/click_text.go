package driver

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

const clickValidationPadding = 5

var clickTextChars = []string{
	"福", "运", "奖", "喜", "庆", "年", "会", "礼", "星", "彩",
	"春", "乐", "吉", "祥", "顺", "利", "金", "光", "梦", "圆",
}

type clickTextDriver struct {
	captcha click.Captcha
}

type clickTarget struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func newClickTextDriver() (*clickTextDriver, error) {
	backgrounds, err := imagesv2.GetImages()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载文字点选背景资源失败: %w", err)
	}
	font, err := fzshengsksjw.GetFont()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载文字点选字体资源失败: %w", err)
	}
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 6, Max: 7}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 3, Max: 3}),
	)
	builder.SetResources(
		click.WithChars(clickTextChars),
		click.WithFonts([]*truetype.Font{font}),
		click.WithBackgrounds(backgrounds),
	)
	return &clickTextDriver{captcha: builder.Make()}, nil
}

func (d *clickTextDriver) Mode() string {
	return ModeClickText
}

func (d *clickTextDriver) Generate() (*Data, error) {
	data, err := d.captcha.Generate()
	if err != nil {
		return nil, err
	}
	dots := data.GetData()
	indexes := make([]int, 0, len(dots))
	for index := range dots {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	targets := make([]clickTarget, 0, len(indexes))
	for _, index := range indexes {
		dot := dots[index]
		targets = append(targets, clickTarget{
			X:      dot.X,
			Y:      dot.Y,
			Width:  dot.Width,
			Height: dot.Height,
		})
	}
	payload, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	image, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}
	prompt, err := data.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}
	return &Data{
		Image:          image,
		Prompt:         prompt,
		RequiredPoints: len(targets),
		Payload:        payload,
	}, nil
}

func (d *clickTextDriver) Validate(payload json.RawMessage, points []Point) (bool, error) {
	var targets []clickTarget
	if err := json.Unmarshal(payload, &targets); err != nil {
		return false, err
	}
	if len(points) != len(targets) {
		return false, nil
	}
	for index, point := range points {
		target := targets[index]
		if !click.Validate(
			point.X,
			point.Y,
			target.X,
			target.Y,
			target.Width,
			target.Height,
			clickValidationPadding,
		) {
			return false, nil
		}
	}
	return true, nil
}
