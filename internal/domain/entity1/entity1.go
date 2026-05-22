package entity1

type DomainEntity1 struct {
	Field1 string
	Field2 int
}


func NewDomainEntity1() *DomainEntity1 { // add params and logic if needed
	return &DomainEntity1{}
}

func (e *DomainEntity1) SomeMethod() {} // add params and logic if needed