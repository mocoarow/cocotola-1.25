package gin_test

import (
	"reflect"
	"testing"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
)

func TestSplitCommaSeparated_shouldReturnExpectedSlice_whenInputVaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "comma separated",
			input: "a,b",
			want:  []string{"a", "b"},
		},
		{
			name:  "spaces are trimmed",
			input: " a , b ",
			want:  []string{"a", "b"},
		},
		{
			name:  "single value",
			input: "only",
			want:  []string{"only"},
		},
		{
			name:  "wildcard remains single",
			input: "*",
			want:  []string{"*"},
		},
		{
			name:  "trailing comma ignores empty",
			input: "a,",
			want:  []string{"a"},
		},
		{
			name:  "all empty keeps original",
			input: "",
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := gin.SplitCommaSeparated(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SplitCommaSeparated(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
