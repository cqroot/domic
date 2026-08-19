package initcmd

import "testing"

func TestResolveRepo(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"username only", "alice", "https://github.com/alice/dotfiles"},
		{"owner/repo", "alice/dots", "https://github.com/alice/dots"},
		{"https URL", "https://github.com/alice/dots", "https://github.com/alice/dots"},
		{"ssh URL", "git@github.com:alice/dots.git", "git@github.com:alice/dots.git"},
		{"ssh scheme", "ssh://git@github.com/alice/dots", "ssh://git@github.com/alice/dots"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := resolveRepo(c.in)
			if got != c.want {
				t.Errorf("resolveRepo(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}