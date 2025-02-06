package project

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
)

// DefaultStartingPort is the root port used by the local proxy
// and additional services will just increment by one
var DefaultStartingPort = 8080

// GetLastUsedPort returns the last used port for a service
func GetLastUsedPort() (int, error) {
	fileSystem := os.DirFS("cmd")
	matchPaths := []string{}
	fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		match, err := path.Match("*-service/main.go", filePath)
		if err != nil {
			return err
		}
		if !match {
			return nil
		}
		matchPaths = append(matchPaths, filePath)
		return nil
	})
	ports := make([]int, 0, len(matchPaths))
	for _, mainPath := range matchPaths {
		port, err := GetPortFromMain(fileSystem, mainPath)
		if err != nil {
			return 0, err
		}
		ports = append(ports, port)
	}
	if len(ports) == 0 {
		return DefaultStartingPort, nil
	}
	sort.Ints(ports)
	port := ports[len(ports)-1]
	if port == 0 {
		port = DefaultStartingPort
	}
	return port, nil
}

var portLineRe = regexp.MustCompile(`\s*port\s*=\s*(\d+)`)

func GetPortFromMain(fileSystem fs.FS, filePath string) (int, error) {
	f, err := fileSystem.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed opening %v: %w", filePath, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matches := portLineRe.FindStringSubmatch(scanner.Text())
		if matches != nil {
			port, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}
			return port, nil
		}
	}
	return 0, scanner.Err()
}
