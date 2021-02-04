package urlp

import (
	"net/url"
	"testing"

	"github.com/jakubDoka/sterr"
)

// Inner is part of a unit test, has no purpose other then that
type Inner struct {
	B uint
	C bool
	D []int
	E float64
}

type test struct {
	A int
	Inner
	F string
	G string `urlp:"g"`
	H string `urlp:"optional"`
}

func (t *test) cmp(b *test) bool {
	for i := range t.D {
		if t.D[i] != b.D[i] {
			return false
		}
	}

	return t.A == b.A && t.B == b.B && t.C == b.C && t.E == b.E && t.F == b.F && t.G == b.G && t.H == b.H
}

func TestParse(t *testing.T) {
	testCases := []struct {
		desc string
		args url.Values
		res  test
		err  sterr.Err
	}{
		{
			desc: "success",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
				"H": {"string"},
			},
			res: test{10, Inner{10, true, []int{0, 1, 5}, 10}, "string", "string", "string"},
		},
		{
			desc: "success omit",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
			},
			res: test{10, Inner{10, true, []int{0, 1, 5}, 10}, "string", "string", ""},
		},
		{
			desc: "missing value",
			args: url.Values{},
			err:  ErrMissingValue,
		},
		{
			desc: "parse fail",
			args: url.Values{
				"A": {"10A"},
			},
			res: test{},
			err: ErrParseFail,
		},
		{
			desc: "bool parse fail",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"trueu"},
			},
			res: test{A: 10, Inner: Inner{B: 10}},
			err: ErrParseFail,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var res test
			err := Parse(tC.args, &res)
			if !res.cmp(&tC.res) {
				t.Error(res, tC.res)
			}

			if !tC.err.SameSurface(err) {
				t.Error(err, "!=", tC.err)
			}
		})
	}
}
