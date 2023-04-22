package err

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("只支持指向结构体的一级指针")
)

var (
	ErrUnKnowColumn = func(col any) error {
		return fmt.Errorf("未知字段:%t", col)
	}
)