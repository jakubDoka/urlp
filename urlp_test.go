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
	E float64
}

type test struct {
	A int
	Inner
	F string
	G string `urlp:"g"`
	H string `urlp:"optional"`
}

func TestParse(t *testing.T) {
	testCases := []struct {
		desc string
		args url.Values
		res  test
		err  *sterr.Err
	}{
		{
			desc: "success",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
				"H": {"string"},
			},
			res: test{10, Inner{10, true, 10}, "string", "string", "string"},
		},
		{
			desc: "success omit",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"true"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
			},
			res: test{10, Inner{10, true, 10}, "string", "string", ""},
		},
		{
			desc: "success",
			args: url.Values{},
			err:  ErrMissingValue,
		},
		{
			desc: "success",
			args: url.Values{
				"A": {"10A"},
				"B": {"10"},
				"C": {"true"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
			},
			res: test{},
			err: ErrParseFail,
		},
		{
			desc: "success",
			args: url.Values{
				"A": {"10"},
				"B": {"10"},
				"C": {"trueu"},
				"E": {"10"},
				"F": {"string"},
				"g": {"string"},
			},
			res: test{A: 10, Inner: Inner{B: 10}},
			err: ErrParseFail,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var res test
			err := Parse(tC.args, &res)
			if res != tC.res {
				t.Error(res, tC.res)
			}

			if tC.err != nil && err != nil && !tC.err.SameSurface(err) {
				t.Error(err, "!=", tC.err)
			}
		})
	}
}
