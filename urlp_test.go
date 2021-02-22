package urlp

import (
	"reflect"
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
	H int    `urlp:"optional"`
}

type test2 struct {
	Inner `urlp:"inner,notinlined"`
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
		args values
		res  test
		err  sterr.Err
	}{
		{
			desc: "success",
			args: values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
				"H": {"10"},
			},
			res: test{10, Inner{10, true, []int{0, 1, 5}, 10}, "string", "string", 10},
		},
		{
			desc: "omit",
			args: values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
			},
			res: test{10, Inner{10, true, []int{0, 1, 5}, 10}, "string", "string", 0},
		},
		{
			desc: "omit with empty string",
			args: values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
				"h": {""},
			},
			res: test{10, Inner{10, true, []int{0, 1, 5}, 10}, "string", "string", 0},
		},
		{
			desc: "missing value",
			args: values{},
			err:  ErrMissingValue,
		},
		{
			desc: "parse fail",
			args: values{
				"A": {"10A"},
			},
			res: test{},
			err: ErrParseFail,
		},
		{
			desc: "bool parse fail",
			args: values{
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

func TestConfiguration(t *testing.T) {
	testCases := []struct {
		desc   string
		parser Parser
		args   values
		res    Inner
	}{
		{
			desc:   "none",
			parser: New(),
			args: values{
				"B": {"10"},
				"C": {"true"},
				"D": {"0", "1", "5"},
				"E": {"10"},
			},
			res: Inner{10, true, []int{0, 1, 5}, 10},
		},
		{
			desc:   "optional",
			parser: New(Optional),
			args:   values{},
		},
		{
			desc:   "lower case",
			parser: New(LowerCase),
			args: values{
				"b": {"10"},
				"c": {"true"},
				"d": {"0", "1", "5"},
				"e": {"10"},
			},
			res: Inner{10, true, []int{0, 1, 5}, 10},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var res test
			err := tC.parser.Parse(tC.args, &res.Inner)
			if err != nil {
				t.Error(err)
				return
			}

			if !res.cmp(&test{Inner: tC.res}) {
				t.Error(res, tC.res)
			}
		})
	}
}

func TestNotInlined(t *testing.T) {
	url := values{
		"inner.B": {"10"},
		"inner.C": {"true"},
		"inner.D": {"10", "20", "30"},
		"inner.E": {"10.3"},
	}

	res := test2{Inner{10, true, []int{10, 20, 30}, 10.3}}

	var r test2
	err := Parse(url, &r)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(res, r) {
		t.Error(r)
	}
}
