package ipranges

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type CidrMatch struct {
	cidr        *net.IPNet
	CidrString  string
	Region      string
	Service     string
	NetworkSize int
}

func NewCidrMatch(cidrstr, region, service string) (*CidrMatch, error) {
	_, cidr, err := net.ParseCIDR(cidrstr)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(cidrstr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Could not parse network bit length from CIDR string '%s'", cidrstr)
	}
	networkSize, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Could not parse network bit length from CIDR string '%s' because: %s", cidrstr, err)
	}

	return &CidrMatch{
		cidr:        cidr,
		CidrString:  cidrstr,
		Region:      region,
		Service:     service,
		NetworkSize: networkSize,
	}, nil
}

type CidrMatcher struct {
	matches []*CidrMatch
}

// Match finds the most specific match from all matching CIDRs
func (m *CidrMatcher) Match(ipstr string) (*CidrMatch, error) {
	matches, err := m.AllMatches(ipstr)
	if err != nil {
		return nil, err
	}

	var best *CidrMatch
	for _, match := range matches {
		if best == nil ||
			match.NetworkSize > best.NetworkSize ||
			(match.NetworkSize == best.NetworkSize && match.Service != "AMAZON") {
			best = match
		}
	}

	return best, nil
}

func (m *CidrMatcher) AllMatches(ipstr string) ([]*CidrMatch, error) {
	var rv []*CidrMatch

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return nil, fmt.Errorf("Could not parse IP '%s'", ipstr)
	}

	for _, match := range m.matches {
		if match.cidr.Contains(ip) {
			rv = append(rv, match)
		}
	}

	return rv, nil
}
