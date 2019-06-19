package main

import (
	"fmt"
	"kuitest/additional/polymorphism/income"
	ic "kuitest/additional/polymorphism/income"
)

func calculateNetIncome(incomes []ic.Income) {
	var netincome int
	for _, i := range incomes {
		af, ok := i.(income.FixedBilling)
		if ok {
			fmt.Println("structure", af)
		}
		afp, ok := i.(*income.FixedBilling)
		if ok {
			fmt.Println("pointer", afp)
		}
		fmt.Printf("Income From %s = $%d\n", i.Source(), i.Calculate())
		netincome += i.Calculate()
	}
	fmt.Printf("Net income of organisation = $%d", netincome)
}

func show(showers []AdHocInterface) {
	fmt.Println(len(showers))
	for _, i := range showers {
		if i == nil {
			fmt.Printf("extend the object is nil: %d \n", i)
			continue
		}
		fmt.Printf("extend the object: %d \n", i.Show())
	}
}

func main() {
	project1 := &ic.FixedBilling{ProjectName: "Project 1", BiddedAmount: 5000}
	project2 := &ic.FixedBilling{ProjectName: "Project 2", BiddedAmount: 10000}
	project3 := &ic.TimeAndMaterial{ProjectName: "Project 3", NoOfHours: 160, HourlyRate: 25}
	bannerAd := &ic.Advertisement{AdName: "Banner Ad", CPC: 2, NoOfClicks: 500}
	popupAd := &ic.Advertisement{AdName: "Popup Ad", CPC: 5, NoOfClicks: 750}
	incomeStreams := []ic.Income{project1, project2, project3, bannerAd, popupAd}
	shows := []AdHocInterface{AFix(*project1), AFix(*project2), ATim(*project3), AAdv(*bannerAd), AAdv(*popupAd)}
	//vt := AFix(project1)
	// fmt.Println(vt.Show())
	// fmt.Println(vt.Calculate())
	// k := *project1
	// p, ok := k.(ic.Income)
	calculateNetIncome(incomeStreams)
	show(shows)
}
