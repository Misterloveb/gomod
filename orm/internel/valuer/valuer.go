package valuer

import (
	"database/sql"

	"github.com/Misterloveb/gomod/orm/model"
)

type Creator func(entity any, model *model.Model) ValuerFace

type ValuerFace interface {
	SetColumnVal(*sql.Rows) error
}
