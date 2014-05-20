package vals

import (
	"reflect"
	"testing"
)

func TestLInsertBefore(t *testing.T) {
	cases := []struct {
		l        []string
		piv, val string
		exp, at  int64
	}{
		0: {nil, "", "", -1, -1},
		1: {[]string{}, "", "", -1, -1},
		2: {[]string{"a"}, "a", "z", 2, 0},
		3: {[]string{"a"}, "x", "z", -1, -1},
		4: {[]string{"a", "b", "c"}, "a", "z", 4, 0},
		5: {[]string{"a", "b", "c"}, "b", "z", 4, 1},
		6: {[]string{"a", "b", "c"}, "c", "z", 4, 2},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LInsertBefore(c.piv, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.at >= 0 {
			if l[c.at] != c.val {
				t.Errorf("%d: value %q should be at index %d, got %q", i, c.val, c.at, l[c.at])
			}
		}
		t.Logf("%d: %v", i, l)
	}
}

func TestLInsertAfter(t *testing.T) {
	cases := []struct {
		l        []string
		piv, val string
		exp, at  int64
	}{
		0: {nil, "", "", -1, -1},
		1: {[]string{}, "", "", -1, -1},
		2: {[]string{"a"}, "a", "z", 2, 1},
		3: {[]string{"a"}, "x", "z", -1, -1},
		4: {[]string{"a", "b", "c"}, "a", "z", 4, 1},
		5: {[]string{"a", "b", "c"}, "b", "z", 4, 2},
		6: {[]string{"a", "b", "c"}, "c", "z", 4, 3},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LInsertAfter(c.piv, c.val)
		if got != c.exp {
			t.Errorf("%d: expected %d, got %d", i, c.exp, got)
		}
		if c.at >= 0 {
			if l[c.at] != c.val {
				t.Errorf("%d: value %q should be at index %d, got %q", i, c.val, c.at, l[c.at])
			}
		}
		t.Logf("%d: %v", i, l)
	}
}

func TestLRange(t *testing.T) {
	cases := []struct {
		l           []string
		start, stop int64
		exp         []string
	}{
		0: {nil, 0, 1, []string{}},
		1: {[]string{}, 0, 2, []string{}},
		2: {[]string{"a"}, 0, 2, []string{"a"}},
		3: {[]string{"a", "b", "c"}, 1, 2, []string{"b", "c"}},
		4: {[]string{"a", "b", "c"}, -3, 2, []string{"a", "b", "c"}},
		5: {[]string{"a", "b", "c"}, 1, 222, []string{"b", "c"}},
		6: {[]string{"a", "b", "c"}, -123, -2, []string{"a", "b"}},
		7: {[]string{"a", "b", "c"}, -123, -5, []string{}},
		8: {[]string{"a", "b", "c"}, 17, -1, []string{}},
		9: {[]string{"a", "b", "c"}, 17, -18, []string{}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LRange(c.start, c.stop)
		if !reflect.DeepEqual(got, c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, got)
		}
		t.Logf("%d: %v", i, got)
	}
}

func TestLRem(t *testing.T) {
	cases := []struct {
		l      []string
		val    string
		cnt, n int64
		exp    []string
	}{
		0:  {nil, "", 0, 0, nil},
		1:  {[]string{}, "", 0, 0, []string{}},
		2:  {[]string{"a", "b", "c"}, "z", 0, 0, []string{"a", "b", "c"}},
		3:  {[]string{"a", "b", "c"}, "z", 2, 0, []string{"a", "b", "c"}},
		4:  {[]string{"a", "b", "c"}, "z", -1, 0, []string{"a", "b", "c"}},
		5:  {[]string{"a", "z", "c", "z"}, "z", 0, 2, []string{"a", "c"}},
		6:  {[]string{"a", "z", "c", "z"}, "z", 1, 1, []string{"a", "c", "z"}},
		7:  {[]string{"a", "z", "c", "z"}, "z", 3, 2, []string{"a", "c"}},
		8:  {[]string{"a", "z", "c", "z"}, "z", -1, 1, []string{"a", "z", "c"}},
		9:  {[]string{"a", "z", "c", "z"}, "z", -4, 2, []string{"a", "c"}},
		10: {[]string{"a", "z", "c", "z"}, "a", -4, 1, []string{"z", "c", "z"}},
	}
	for i, c := range cases {
		l := list(c.l)
		got := l.LRem(c.cnt, c.val)
		if got != c.n {
			t.Errorf("%d: expected %d elements removed, got %d", i, c.n, got)
		}
		if !reflect.DeepEqual([]string(l), c.exp) {
			t.Errorf("%d: expected %v, got %v", i, c.exp, l)
		}
		t.Logf("%d: %v", i, l)
	}
}