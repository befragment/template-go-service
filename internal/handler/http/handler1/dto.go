// dto.go file is used to declare json schemas
// you can add mapping functions here too
package handler1

import (
	"github.com/befragment/template-go/internal/domain/entity1"
)

type someHTTPreq struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func fromDomainEntity1(e entity1.DomainEntity1) someHTTPreq {
	return someHTTPreq{Field1: e.Field1, Field2: e.Field2}
}

func toDomainEntity(r someHTTPreq) entity1.DomainEntity1 {
	return entity1.DomainEntity1{Field1: r.Field1, Field2: r.Field2}
}