package income

// Advertisement ...
type Advertisement struct {
	AdName     string
	CPC        int
	NoOfClicks int
}

// implement Income interface
func (a Advertisement) Calculate() int {
	return a.CPC * a.NoOfClicks
}

func (a Advertisement) Source() string {
	return a.AdName
}
