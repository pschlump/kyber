package test

import (
	"github.com/dedis/kyber"
	"github.com/dedis/kyber/util/random"
)

// A generic benchmark suite for kyber.groups.
type GroupBench struct {
	g kyber.Group

	// Random secrets and points for testing
	x, y kyber.Scalar
	X, Y kyber.Point
	xe   []byte // encoded Scalar
	Xe   []byte // encoded Point
}

func NewGroupBench(g kyber.Group) *GroupBench {
	var gb GroupBench
	gb.g = g
	gb.x = g.Scalar().Pick(random.Stream)
	gb.y = g.Scalar().Pick(random.Stream)
	gb.xe, _ = gb.x.MarshalBinary()
	gb.X, _ = g.Point().Pick(nil, random.Stream)
	gb.Y, _ = g.Point().Pick(nil, random.Stream)
	gb.Xe, _ = gb.X.MarshalBinary()
	return &gb
}

func (gb GroupBench) ScalarAdd(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Add(gb.x, gb.y)
	}
}

func (gb GroupBench) ScalarSub(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Sub(gb.x, gb.y)
	}
}

func (gb GroupBench) ScalarNeg(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Neg(gb.x)
	}
}

func (gb GroupBench) ScalarMul(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Mul(gb.x, gb.y)
	}
}

func (gb GroupBench) ScalarDiv(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Div(gb.x, gb.y)
	}
}

func (gb GroupBench) ScalarInv(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Inv(gb.x)
	}
}

func (gb GroupBench) ScalarPick(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.Pick(random.Stream)
	}
}

func (gb GroupBench) ScalarEncode(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.MarshalBinary()
	}
}

func (gb GroupBench) ScalarDecode(iters int) {
	for i := 1; i < iters; i++ {
		gb.x.UnmarshalBinary(gb.xe)
	}
}

func (gb GroupBench) PointAdd(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Add(gb.X, gb.Y)
	}
}

func (gb GroupBench) PointSub(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Sub(gb.X, gb.Y)
	}
}

func (gb GroupBench) PointNeg(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Neg(gb.X)
	}
}

func (gb GroupBench) PointMul(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Mul(gb.X, gb.y)
	}
}

func (gb GroupBench) PointBaseMul(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Mul(nil, gb.y)
	}
}

func (gb GroupBench) PointPick(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.Pick(nil, random.Stream)
	}
}

func (gb GroupBench) PointEncode(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.MarshalBinary()
	}
}

func (gb GroupBench) PointDecode(iters int) {
	for i := 1; i < iters; i++ {
		gb.X.UnmarshalBinary(gb.Xe)
	}
}
