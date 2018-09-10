package credentials

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ParseInput parses a git credential-like standard in into a pullParams data struct. The input
// is something like:
// protocol=https
// host=example.com
// path=foo.git
func ParseInput(r io.Reader) (*pullParams, error) {
	scn := bufio.NewScanner(r)
	var lines []string

	for scn.Scan() {
		line := scn.Text()
		lines = append(lines, line)
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("Invalid number of input lines")
	}

	firstLine := strings.Split(lines[0], "=")
	if firstLine[0] != "protocol" || len(firstLine) < 2 {
		return nil, fmt.Errorf("Invalid argument, protocol not found")
	}

	protocol := firstLine[1]
	secondLine := strings.Split(lines[1], "=")
	if secondLine[0] != "host" || len(secondLine) < 2 {
		return nil, fmt.Errorf("Invalid argument, host not found")
	}
	host := secondLine[1]
	return &pullParams{host, protocol}, nil
}

// FromJSON returns the matching credentials from a json formatted reader
func FromJSON(reader io.Reader, params *pullParams) (credentials, error) {
	var cc []credentials
	dec := json.NewDecoder(reader)
	for {
		if err := dec.Decode(&cc); err == io.EOF {
			break
		} else if err != nil {
			return credentials{}, err
		}
	}

	for _, c := range cc {
		if urlMatches(params.RepoURL, c.URL) {
			return c, nil
		}
	}
	return credentials{}, fmt.Errorf("Url not found")
}

func urlMatches(requested, toMatch string) bool {
	match, err := regexp.Match(toMatch, []byte(requested))
	return match && err == nil
}
