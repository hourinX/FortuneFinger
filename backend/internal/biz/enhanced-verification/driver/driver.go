package driver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	ModeClickText   = "click-text"
	ModeClickIcon   = "click-icon"
	ModeRotate      = "rotate"
	ModeBehavior    = "behavior"
	ModeImageSelect = "image-select"
	ModeSlide       = "slide"
)

type Point struct {
	X int   `json:"x"`
	Y int   `json:"y"`
	T int64 `json:"t,omitempty"`
}

type Data struct {
	Image          string
	Prompt         string
	Thumb          string
	ThumbX         int
	ThumbY         int
	ThumbWidth     int
	ThumbHeight    int
	RequiredPoints int
	Payload        json.RawMessage
}

type Driver interface {
	Mode() string
	Generate() (*Data, error)
	Validate(payload json.RawMessage, points []Point) (bool, error)
}

func New(mode string) (Driver, error) {
	switch mode {
	case ModeClickText:
		return newClickTextDriver()
	case ModeClickIcon:
		return newClickIconDriver()
	case ModeRotate:
		return newRotateDriver()
	case ModeBehavior:
		return newBehaviorDriver(), nil
	case ModeImageSelect:
		return newImageSelectDriver()
	case ModeSlide:
		return newSlideDriver()
	default:
		return nil, fmt.Errorf("enhanced-verification: 不支持的验证模式「%s」", mode)
	}
}

func randomValue(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
