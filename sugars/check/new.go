package check

import "nhatp.com/go/sugar"

func New() sugar.Sugar {
	return &impl{}
}

type impl struct{}

var _ sugar.Sugar = (*impl)(nil)
