package ipranges

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func CachedIpRangesFilename() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return path.Join(cacheDir, "aws-ip-lookup", "ip-ranges.json"), nil
}

func isCachedIpRangesValid() bool {
	filename, err := CachedIpRangesFilename()
	if err != nil {
		return false
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}

	if stat.IsDir() || stat.Size() < 100 {
		return false
	}

	if time.Since(stat.ModTime()) < (24 * time.Hour) {
		return true
	}
	return false
}

func readCachedIpRangesJSON() ([]byte, error) {
	filename, err := CachedIpRangesFilename()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filename)
}

func writeCachedIpRangesJSON(body []byte) error {
	filename, err := CachedIpRangesFilename()
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return os.WriteFile(filename, body, 0644)
}

func fetchIpRangesJSON() ([]byte, error) {
	resp, err := http.Get("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func parseIpRangesJSON(body []byte) (*IpRanges, error) {
	var ipRanges IpRanges
	if err := json.Unmarshal(body, &ipRanges); err != nil {
		return nil, err
	}
	return &ipRanges, nil
}

func GetIpRanges() (*IpRanges, error) {
	if isCachedIpRangesValid() {
		body, err := readCachedIpRangesJSON()
		if err == nil {
			rv, err := parseIpRangesJSON(body)
			if err == nil {
				rv.IsFromCache = true
			}
			return rv, err
		}
	}

	body, err := fetchIpRangesJSON()
	if err != nil {
		return nil, err
	}

	rv, err := parseIpRangesJSON(body)
	if err != nil {
		return nil, err
	}

	if err := writeCachedIpRangesJSON(body); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Unable to write to local cache because '%s'", err)
	}

	return rv, nil
}
