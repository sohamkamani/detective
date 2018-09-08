package detective

type DetectorFunc func() error

type Dependency struct {
	name     string
	detector DetectorFunc
	state    State
}

func noopDetectorFunc() DetectorFunc {
	return func() error {
		return nil
	}
}

func NewDependency(name string) *Dependency {
	return &Dependency{
		name:     name,
		detector: noopDetectorFunc(),
	}
}

func (d *Dependency) Detect(df DetectorFunc) {
	d.detector = df
}

func (d *Dependency) updateState() {
	d.state = d.getState()
}

func (d *Dependency) getState() State {
	err := d.detector()
	s := State{Name: d.name}
	if err != nil {
		return s.WithError(err)
	}
	return s.WithOk()
}
