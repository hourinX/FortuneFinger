package driver

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

type clickIconDriver struct {
	mode    string
	captcha click.Captcha
}

func newClickIconDriver() (*clickIconDriver, error) {
	return newIconClickDriver(ModeClickIcon)
}

func newIconClickDriver(mode string) (*clickIconDriver, error) {
	backgrounds, err := imagesv2.GetImages()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载图标点选背景资源失败: %w", err)
	}
	shapeImages, err := shapes.GetShapes()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载图标点选资源失败: %w", err)
	}
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 6, Max: 7}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 3, Max: 3}),
		click.WithUseShapeOriginalColor(true),
	)
	builder.SetResources(
		click.WithShapes(shapeImages),
		click.WithBackgrounds(backgrounds),
	)
	return &clickIconDriver{mode: mode, captcha: builder.MakeShape()}, nil
}

func (d *clickIconDriver) Mode() string {
	return d.mode
}

func (d *clickIconDriver) Generate() (*Data, error) {
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
		targets = append(targets, clickTarget{X: dot.X, Y: dot.Y, Width: dot.Width, Height: dot.Height})
	}
	payload, err := json.Marshal(targets)
	if err != nil {
		return nil, err
	}
	imageData, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}
	prompt, err := data.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}
	return &Data{
		Image:          imageData,
		Prompt:         prompt,
		RequiredPoints: len(targets),
		Payload:        payload,
	}, nil
}

func (d *clickIconDriver) Validate(payload json.RawMessage, points []Point) (bool, error) {
	var targets []clickTarget
	if err := json.Unmarshal(payload, &targets); err != nil {
		return false, err
	}
	if len(points) != len(targets) {
		return false, nil
	}
	for index, point := range points {
		target := targets[index]
		if !click.Validate(point.X, point.Y, target.X, target.Y, target.Width, target.Height, clickValidationPadding) {
			return false, nil
		}
	}
	return true, nil
}
