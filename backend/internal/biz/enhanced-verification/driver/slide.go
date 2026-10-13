package driver

import (
	"encoding/json"
	"fmt"

	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

const slideValidationPadding = 5

type slideDriver struct {
	captcha slide.Captcha
}

type slideTarget struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func newSlideDriver() (*slideDriver, error) {
	backgrounds, err := imagesv2.GetImages()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载拼图背景资源失败: %w", err)
	}
	tileImages, err := tiles.GetTiles()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载拼图资源失败: %w", err)
	}
	graphs := make([]*slide.GraphImage, 0, len(tileImages))
	for _, tile := range tileImages {
		graphs = append(graphs, &slide.GraphImage{
			OverlayImage: tile.OverlayImage,
			ShadowImage:  tile.ShadowImage,
			MaskImage:    tile.MaskImage,
		})
	}
	builder := slide.NewBuilder()
	builder.SetResources(
		slide.WithBackgrounds(backgrounds),
		slide.WithGraphImages(graphs),
	)
	return &slideDriver{captcha: builder.Make()}, nil
}

func (d *slideDriver) Mode() string {
	return ModeSlide
}

func (d *slideDriver) Generate() (*Data, error) {
	data, err := d.captcha.Generate()
	if err != nil {
		return nil, err
	}
	block := data.GetData()
	payload, err := json.Marshal(slideTarget{X: block.X, Y: block.Y})
	if err != nil {
		return nil, err
	}
	image, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}
	thumb, err := data.GetTileImage().ToBase64()
	if err != nil {
		return nil, err
	}
	return &Data{
		Image:          image,
		Thumb:          thumb,
		ThumbX:         block.DX,
		ThumbY:         block.DY,
		ThumbWidth:     block.Width,
		ThumbHeight:    block.Height,
		RequiredPoints: 1,
		Payload:        payload,
	}, nil
}

func (d *slideDriver) Validate(payload json.RawMessage, points []Point) (bool, error) {
	if len(points) != 1 {
		return false, nil
	}
	var target slideTarget
	if err := json.Unmarshal(payload, &target); err != nil {
		return false, err
	}
	return slide.Validate(
		target.X,
		target.Y,
		points[0].X,
		points[0].Y,
		slideValidationPadding,
	), nil
}
