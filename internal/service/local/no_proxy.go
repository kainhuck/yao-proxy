package local

import "regexp"

type Filter struct {
	re []*regexp.Regexp
}

func NewFilter(raw []string) *Filter {
	re := make([]*regexp.Regexp, len(raw))

	for i, r := range raw{
		re[i] = regexp.MustCompile(r)
	}

	return &Filter{
		re: re,
	}
}

func (f *Filter) Check(str string) bool {
	for _, r := range f.re{
		if r.MatchString(str){
			return true
		}
	}

	return false
}
