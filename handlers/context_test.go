package handlers

import (
	"strings"
	"testing"
)

func TestParseMounts(t *testing.T) {
	useCases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name: "deeper submount deduplication",
			input: `
/dev/sda1 /data ext4 rw 0 0
/dev/sda1 /data/sub ext4 rw 0 0
/dev/sda1 /data/sub/deep ext4 rw 0 0
`,
			want: []string{"/data"},
		},
		{
			name: "siblings kept",
			input: `
/dev/sda1 /projects/foo ext4 rw 0 0
/dev/sda1 /projects/bar ext4 rw 0 0
`,
			want: []string{"/projects/bar", "/projects/foo"},
		},
	}

	for _, u := range useCases {
		t.Run(u.name, func(t *testing.T) {
			got, err := parseMounts(strings.NewReader(u.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(u.want) {
				t.Fatalf("got %v, want %v", got, u.want)
			}
			for i := range got {
				if got[i] != u.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], u.want[i])
				}
			}
		})
	}
}
