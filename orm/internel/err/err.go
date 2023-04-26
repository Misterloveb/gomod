package err

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("只支持指向结构体的一级指针")
	ErrTagNoOrm    = errors.New("结构体tag格式错误,缺少orm")
	ErrTagNoDeng   = errors.New("结构体tag格式错误,缺少=")
	ErrNoRows      = errors.New("没有数据")
)

var (
	ErrUnKnowColumn = func(col any) error {
		return fmt.Errorf("未知字段:%t", col)
	}
)
