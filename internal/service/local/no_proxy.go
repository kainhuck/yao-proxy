package local

import (
	"regexp"
	"strings"
)

/*

	用户配置的国力规则必须满足以下格式：
		1. 纯 IPv4 地址: 192.168.1.1
		2. 带占位符的 IPv4地址: 192.178.1.x
		3. IPv4地址范围: 192.168.1.1-192.168.1.12 todo
		4. 域名
*/

const (
	IPv4Expr = `^(\d+?\.){3}\d+$`
	IPv4XExpr = `^(\d+?\.|[x]\.){3}\d+$`
	Ipv4RangeExpr = `^(\d+?\.){3}\d+?-(\d+?\.){3}\d+$`
)

var (
	IPv4Cpl = regexp.MustCompile(IPv4Expr)
	IPv4XCpl = regexp.MustCompile(IPv4XExpr)
	IPv4RangeCpl = regexp.MustCompile(Ipv4RangeExpr)
)

type Filter struct {
	re []*regexp.Regexp
}

func NewFilter(raw []string) *Filter {
	re := make([]*regexp.Regexp, 0)

	for _, r := range raw{
		if IPv4Cpl.MatchString(r){
			// 纯 IPv4 什么都不处理
		}else if IPv4XCpl.MatchString(r){
			// 将 x 替换为 \d+?
			r = strings.Replace(r, "x", `\d+?`, -1)
		}else if IPv4RangeCpl.MatchString(r){
			// todo 暂时不处理范围匹配
		}

		re = append(re, regexp.MustCompile(r))
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
