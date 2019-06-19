package income

// FixedBilling ...
type FixedBilling struct {
	ProjectName  string
	BiddedAmount int
}

// implement Income interface
func (fb FixedBilling) Calculate() int {
	return fb.BiddedAmount
}

func (fb FixedBilling) Source() string {
	return fb.ProjectName
}
