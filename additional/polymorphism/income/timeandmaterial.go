package income

// TimeAndMaterial ...
type TimeAndMaterial struct {
	ProjectName string
	NoOfHours   int
	HourlyRate  int
}

// implement Income interface
func (tm TimeAndMaterial) Calculate() int {
	return tm.NoOfHours * tm.HourlyRate
}

func (tm TimeAndMaterial) Source() string {
	return tm.ProjectName
}
