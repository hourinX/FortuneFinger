package driver

type imageSelectDriver struct {
	*clickIconDriver
}

func newImageSelectDriver() (*imageSelectDriver, error) {
	driver, err := newIconClickDriver(ModeImageSelect)
	if err != nil {
		return nil, err
	}
	return &imageSelectDriver{clickIconDriver: driver}, nil
}
