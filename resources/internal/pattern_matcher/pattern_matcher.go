package pattern_matcher

import "github.com/gobwas/glob"

type PatternMatcher struct {
	include []string
	exclude []string
}

// Option is the type for functional options.
type Option func(*PatternMatcher)

func New(options ...Option) *PatternMatcher {
	pm := &PatternMatcher{
		include: []string{"*"},
		exclude: []string{},
	}

	for _, option := range options {
		option(pm)
	}
	return pm
}

func WithInclude(list []string) Option {
	return func(c *PatternMatcher) {
		c.include = list
	}
}

func WithExclude(list []string) Option {
	return func(c *PatternMatcher) {
		c.exclude = list
	}
}

func WithIncludeFilter(value string) Option {
	return func(c *PatternMatcher) {
		c.include = append(c.include, value)
	}
}

func WithExcludeFilter(value string) Option {
	return func(c *PatternMatcher) {
		c.exclude = append(c.exclude, value)
	}
}

func (pm *PatternMatcher) Match(value string) bool {
	var g glob.Glob

	is_included := false
	// Check if the value matches any pattern in the pm.include list
	for _, pattern := range pm.include {
		g = glob.MustCompile(pattern)
		matched := g.Match(value)
		if matched {
			is_included = true
			break
		}
	}
	if !is_included {
		return false
	}
	// Check if the value matches any pattern in the excluded list
	for _, pattern := range pm.exclude {
		g = glob.MustCompile(pattern)
		matched := g.Match(value)
		if matched {
			return false
		}
	}
	// If the value matches at least one pattern in the included list
	// and does not match any pattern in the excluded list, return true
	return true
}
