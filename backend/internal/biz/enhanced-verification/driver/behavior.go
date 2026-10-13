package driver

import (
	"encoding/json"
	"math"
)

const (
	minimumBehaviorPoints   = 12
	minimumBehaviorDuration = 400
	maximumBehaviorDuration = 8000
	minimumBehaviorDistance = 120
)

type behaviorDriver struct{}

type behaviorTarget struct {
	Nonce string `json:"nonce"`
}

func newBehaviorDriver() *behaviorDriver {
	return &behaviorDriver{}
}

func (d *behaviorDriver) Mode() string {
	return ModeBehavior
}

func (d *behaviorDriver) Generate() (*Data, error) {
	nonce, err := randomValue(16)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(behaviorTarget{Nonce: nonce})
	if err != nil {
		return nil, err
	}
	return &Data{RequiredPoints: minimumBehaviorPoints, Payload: payload}, nil
}

func (d *behaviorDriver) Validate(_ json.RawMessage, points []Point) (bool, error) {
	if len(points) < minimumBehaviorPoints {
		return false, nil
	}
	duration := points[len(points)-1].T - points[0].T
	if duration < minimumBehaviorDuration || duration > maximumBehaviorDuration {
		return false, nil
	}
	distance := 0.0
	for index := 1; index < len(points); index++ {
		dx := float64(points[index].X - points[index-1].X)
		dy := float64(points[index].Y - points[index-1].Y)
		distance += math.Hypot(dx, dy)
	}
	return distance >= minimumBehaviorDistance, nil
}
