package main

import (
	"fmt"
	"os"

	"github.com/daveadams/aws-ip-lookup/ipranges"
)

func main() {
	ipRanges, err := ipranges.GetIpRanges()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	var looper StringIterator
	if len(os.Args) == 1 {
		looper = NewStdinIterator()
	} else {
		looper = NewSliceIterator(os.Args[1:])
	}

	found := 0
	matcher := ipRanges.GetMatcher()
	for {
		ip, ok := looper.Next()
		if !ok {
			break
		}

		match, err := matcher.Match(ip)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s\n", err)
			continue
		}

		if match == nil {
			fmt.Fprintf(os.Stderr, "%s NOT FOUND\n", ip)
			continue
		}

		fmt.Printf("%s %s %s %s\n", ip, match.CidrString, match.Region, match.Service)
		found++
	}

	if found == 0 {
		os.Exit(1)
	}
}
