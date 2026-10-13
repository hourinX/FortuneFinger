package driver

import "testing"

func TestConfiguredDriversGenerate(t *testing.T) {
	modes := []string{ModeClickText, ModeClickIcon, ModeRotate, ModeBehavior, ModeImageSelect, ModeSlide}
	for _, mode := range modes {
		driver, err := New(mode)
		if err != nil {
			t.Fatalf("new driver %s: %v", mode, err)
		}
		if _, err := driver.Generate(); err != nil {
			t.Fatalf("generate %s: %v", mode, err)
		}
	}
}
