package driver

import (
	"encoding/json"
	"fmt"

	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/rotate"
)

const rotateValidationPadding = 5

type rotateDriver struct {
	captcha rotate.Captcha
}

type rotateTarget struct {
	Angle int `json:"angle"`
}

func newRotateDriver() (*rotateDriver, error) {
	images, err := imagesv2.GetImages()
	if err != nil {
		return nil, fmt.Errorf("enhanced-verification: 加载旋转校验资源失败: %w", err)
	}
	builder := rotate.NewBuilder()
	builder.SetResources(rotate.WithImages(images))
	return &rotateDriver{captcha: builder.Make()}, nil
}

func (d *rotateDriver) Mode() string {
	return ModeRotate
}

func (d *rotateDriver) Generate() (*Data, error) {
	data, err := d.captcha.Generate()
	if err != nil {
		return nil, err
	}
	block := data.GetData()
	payload, err := json.Marshal(rotateTarget{Angle: block.Angle})
	if err != nil {
		return nil, err
	}
	imageData, err := data.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}
	thumb, err := data.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}
	return &Data{
		Image:          imageData,
		Thumb:          thumb,
		RequiredPoints: 1,
		Payload:        payload,
	}, nil
}

func (d *rotateDriver) Validate(payload json.RawMessage, points []Point) (bool, error) {
	if len(points) != 1 {
		return false, nil
	}
	var target rotateTarget
	if err := json.Unmarshal(payload, &target); err != nil {
		return false, err
	}
	return rotate.Validate(points[0].X, target.Angle, rotateValidationPadding), nil
}
