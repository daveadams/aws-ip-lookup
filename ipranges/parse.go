package ipranges

import (
	"fmt"
	"os"
)

type IpPrefix struct {
	Ipv4Prefix string `json:"ip_prefix"`
	Ipv6Prefix string `json:"ipv6_prefix"`
	Region     string `json:"region"`
	Service    string `json:"service"`
}

func (p *IpPrefix) Prefix() string {
	if p.Ipv4Prefix != "" {
		return p.Ipv4Prefix
	} else {
		return p.Ipv6Prefix
	}
}

type IpRanges struct {
	IsFromCache bool       `json:"-"`
	SyncToken   string     `json:"syncToken"`
	CreateDate  string     `json:"createDate"`
	Prefixes    []IpPrefix `json:"prefixes"`
}

func (f *IpRanges) GetMatcher() *CidrMatcher {
	matches := make([]*CidrMatch, 0, len(f.Prefixes))

	for _, prefix := range f.Prefixes {
		match, err := NewCidrMatch(prefix.Prefix(), prefix.Region, prefix.Service)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Bad CIDR '%s': %s\n", prefix.Prefix(), err)
		} else {
			matches = append(matches, match)
		}
	}

	return &CidrMatcher{
		matches: matches,
	}
}
