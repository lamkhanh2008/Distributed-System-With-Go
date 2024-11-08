package main

type Employee struct {
	Name, Position string
	AnnualIncome   int
}

func NewEmployeeFactory(position string, annualIncome int) func(name string) *Employee {
	return func(name string) *Employee {
		return &Employee{name, position, annualIncome}
	}
}

func main() {
	// developerFactory := NewEmployeeFactory("developer", 60000)
	// managerFactory := NewEmployeeFactory("manager", 10000)
}
