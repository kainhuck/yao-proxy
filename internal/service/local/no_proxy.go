package local

import (
	"encoding/binary"
	"net"
	"regexp"
	"strings"
)

/*
	用户配置的过滤规则必须满足以下格式：
		1. 纯 IPv4 地址: 192.168.1.1
		2. 带占位符的 IPv4地址: 192.178.1.x
		3. IPv4地址范围: 192.168.1.1-192.168.1.12 todo
		4. 域名
*/

const (
	IPv4Expr      = `^(\d+?\.){3}\d+$`
	IPv4XExpr     = `^(\d+?\.|[x]\.){3}\d+$`
	Ipv4RangeExpr = `^(\d+?\.){3}\d+?-(\d+?\.){3}\d+$`
)

var (
	IPv4Cpl      = regexp.MustCompile(IPv4Expr)
	IPv4XCpl     = regexp.MustCompile(IPv4XExpr)
	IPv4RangeCpl = regexp.MustCompile(Ipv4RangeExpr)

	prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255}
)

type Filter struct {
	re []*regexp.Regexp
	mp map[string]struct{}
}

func NewFilter(raw []string) *Filter {
	re := make([]*regexp.Regexp, 0)
	mp := make(map[string]struct{})
	for _, r := range raw {
		if IPv4Cpl.MatchString(r) {
			// 纯 IPv4 do nothing
			re = append(re, regexp.MustCompile(r))
		} else if IPv4XCpl.MatchString(r) {
			// 将 x 替换为 \d+?
			r = strings.Replace(r, "x", `\d+?`, -1)
			re = append(re, regexp.MustCompile(r))
		} else if IPv4RangeCpl.MatchString(r) {
			ips := strings.Split(r, "-")
			ip1 := binary.BigEndian.Uint32(net.ParseIP(ips[0])[12:])
			ip2 := binary.BigEndian.Uint32(net.ParseIP(ips[1])[12:])
			if ip2 < ip1 { // ip2 不可以小于 ip1
				continue
			}

			for i := ip1; i <= ip2; i++ {
				ipBts := make([]byte, 4)
				binary.BigEndian.PutUint32(ipBts, i)
				ip := append(prefix, ipBts...)

				a := net.IP(ip)
				mp[a.String()] = struct{}{}
			}

			continue
		}else{
			// 这种情况是域名，需要解析出真实IP地址 再保存一份其IP
			re = append(re, regexp.MustCompile(r))
			ip, err := net.ResolveIPAddr("ip", r)
			if err != nil {
				continue
			}
			re = append(re, regexp.MustCompile(ip.String()))
		}
	}

	return &Filter{
		re: re,
		mp: mp,
	}
}

func (f *Filter) Check(host string) bool {
	ip := strings.Split(host, ":")[0]
	for _, r := range f.re {
		if r.MatchString(ip) {
			return true
		}
	}

	_, ok := f.mp[ip]

	return ok
}
