package pattern_matcher

import "github.com/gobwas/glob"

type PatternMatcher struct {
	include []string
	exclude []string
	name    string
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
	is_excluded := false

	for _, in := range pm.include {
		g = glob.MustCompile(in)
		is_included = g.Match(value)
		if is_included {
			break
		}
	}
	if !is_included {
		return false
	}

	for _, ex := range pm.exclude {
		g = glob.MustCompile(ex)
		is_excluded = g.Match(value)
		if !is_excluded {
			return true
		}
	}

	if is_included && !is_excluded {
		return true
	}
	return false
}
