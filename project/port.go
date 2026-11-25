package project

import (
	"bufio"
	"errors"
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
	return getLastUsedPort(os.DirFS("cmd"))
}

func getLastUsedPort(fileSystem fs.FS) (int, error) {
	matchPaths := []string{}
	err := fs.WalkDir(fileSystem, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			if filePath == "." && errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		match, err := path.Match("*-service/app/service.go", filePath)
		if err != nil {
			return err
		}
		if !match {
			return nil
		}
		matchPaths = append(matchPaths, filePath)
		return nil
	})
	if err != nil {
		return 0, err
	}
	ports := make([]int, 0, len(matchPaths))
	for _, servicePath := range matchPaths {
		port, err := getPortFromFile(fileSystem, servicePath)
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

var portLineRe = regexp.MustCompile(`^\s*(?:const\s+)?port\s*=\s*(\d+)\b`)

func getPortFromFile(fileSystem fs.FS, filePath string) (int, error) {
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
