package models

import (
	"github.com/d3v-friends/mango/mtype"
)

var All = []mtype.IfMigrateModel{
	Account{},
	System{},
}

type (
	IPager struct {
		Page int64
		Size int64
	}

	List[T any] struct {
		Page  int64
		Size  int64
		Total int64
		List  []T
	}
)

func (x *IPager) Skip() *int64 {
	skip := x.Page * x.Size
	return &skip
}

func (x *IPager) Limit() *int64 {
	return &x.Size
}
