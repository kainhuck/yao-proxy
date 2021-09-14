package local

import "regexp"

type Filter struct {
	re *regexp.Regexp
}

func NewFilter(raw []string) *Filter {
	return &Filter{}
}

func (f *Filter) Check(str string) bool {
	return false
}
