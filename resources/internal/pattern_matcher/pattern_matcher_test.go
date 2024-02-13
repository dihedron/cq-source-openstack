package pattern_matcher

import (
	"testing"
)

func TestPatternMatcher(t *testing.T) {
	type fields struct {
		include []string
		exclude []string
	}
	tests := []struct {
		name   string
		fields fields
		value  string
		want   bool
	}{
		{
			"test1",
			fields{
				include: []string{"*"},
				exclude: []string{},
			},
			"hello_world",
			true,
		},
		{
			"test2",
			fields{
				include: []string{"*"},
				exclude: []string{"hello*"},
			},
			"hello_world",
			false,
		},
		{
			"test3",
			fields{
				include: []string{"*"},
				exclude: []string{"hello*"},
			},
			"ciao_world",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PatternMatcher{
				include: tt.fields.include,
				exclude: tt.fields.exclude,
			}
			if got := pm.Match(tt.value); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
