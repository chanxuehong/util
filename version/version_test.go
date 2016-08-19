package version

import "testing"

func TestString(t *testing.T) {
	v := Version{7, 8, 9}
	have := v.String()
	want := "7.8.9"
	if have != want {
		t.Errorf("Version{7, 8, 9}.String() failed, have %q, want %q", have, want)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		src string
		ver Version
		err error
	}{
		{
			"1",
			Version{1, 0, 0},
			nil,
		},
		{
			"1.2",
			Version{1, 2, 0},
			nil,
		},
		{
			"1.2.3",
			Version{1, 2, 3},
			nil,
		},
	}

	var v Version
	var err error
	for _, item := range tests {
		v, err = Parse(item.src)
		if v != item.ver || err != item.err {
			t.Errorf("Parse(%q) failed, have(%+v, %v), want(%+v, %v)", item.src, v, err, item.ver, item.err)
		}
	}

	tests2 := []string{
		"1.2.3.4",
		"1.2.3.",
		".2.3",
		"1..3",
		"1.2.",
	}
	for _, item := range tests2 {
		_, err = Parse(item)
		if err == nil {
			t.Errorf("Parse(%q) should have a error, but not", item)
		}
	}
}
