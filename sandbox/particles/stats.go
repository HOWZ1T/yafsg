package particles

type Stats struct {
	Empty    int
	Sand     int
	Water    int
	Concrete int
	Unknown  int
}

func (s Stats) TotalEmpty() int {
	return s.Empty
}

func (s Stats) TotalNonEmpty() int {
	return s.Sand + s.Water + s.Concrete + s.Unknown
}

func (s Stats) Total() int {
	return s.TotalEmpty() + s.TotalNonEmpty()
}
