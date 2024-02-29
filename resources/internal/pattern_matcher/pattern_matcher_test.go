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
			"test_1",
			fields{
				include: []string{"*"},
				exclude: []string{},
			},
			"hello_world",
			true,
		},
		{
			"test_2",
			fields{
				include: []string{"*"},
				exclude: []string{"hello*"},
			},
			"hello_world",
			false,
		},
		{
			"test_3",
			fields{
				include: []string{"*"},
				exclude: []string{"hello*"},
			},
			"hi_world",
			true,
		},
		{
			"test_4",
			fields{
				include: []string{"hello"},
				exclude: []string{"*"},
			},
			"hello",
			false,
		},
		{
			"test_5",
			fields{
				include: []string{"openstack_compute*", "openstack_blockstorage*"},
				exclude: []string{"openstack_blockstorage_attachment*"},
			},
			"openstack_blockstorage_attachments",
			false,
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
