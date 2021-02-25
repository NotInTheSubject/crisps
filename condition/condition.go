package condition

import "context"

// ConditionFunc represent some condition
type Func func(context.Context) bool

func True() Func {
	return Func(func(context.Context) bool { return true })
}

func False() Func {
	return Func(func(context.Context) bool { return false })
}

func (rightCnd Func) And(leftCnd Func) Func {
	return Func(func(ctx context.Context) bool {
		return rightCnd(ctx) && leftCnd(ctx)
	})
}

func (rightCnd Func) Or(leftCnd Func) Func {
	return Func(func(ctx context.Context) bool {
		return rightCnd(ctx) || leftCnd(ctx)
	})
}
