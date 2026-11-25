package project

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func writeTestFile(t *testing.T, root, relPath, content string) {
	t.Helper()
	fullPath := filepath.Join(root, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", fullPath, err)
	}
}

func TestGetLastUsedPort_NoCmdDir(t *testing.T) {
	t.Chdir(t.TempDir())

	got, err := GetLastUsedPort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != DefaultStartingPort {
		t.Fatalf("expected %d, got %d", DefaultStartingPort, got)
	}
}

func TestGetLastUsedPort_NoServices(t *testing.T) {
	t.Run("empty cmd directory", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "cmd"), 0o755); err != nil {
			t.Fatalf("mkdir cmd: %v", err)
		}
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != DefaultStartingPort {
			t.Fatalf("expected %d, got %d", DefaultStartingPort, got)
		}
	})

	t.Run("non-service files only", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "cmd/local-proxy/main.go", "package main\nconst port = 9999\n")
		writeTestFile(t, dir, "cmd/script-service/main.go", "package main\nconst port = 9998\n")
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != DefaultStartingPort {
			t.Fatalf("expected %d, got %d", DefaultStartingPort, got)
		}
	})
}

func TestGetLastUsedPort_SingleService(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "cmd/script-service/app/service.go", `package app

const (
	serviceName = "script"
	port        = 4001
)
`)
	t.Chdir(dir)

	got, err := GetLastUsedPort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 4001 {
		t.Fatalf("expected 4001, got %d", got)
	}
}

func TestGetLastUsedPort_MultipleServices(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, "cmd/alpha-service/app/service.go", `package app

const (
	serviceName = "alpha"
	port        = 4001
)
`)
	writeTestFile(t, dir, "cmd/beta-service/app/service.go", `package app

const (
	serviceName = "beta"
	port        = 8085
)
`)
	writeTestFile(t, dir, "cmd/gamma-service/app/service.go", `package app

const (
	serviceName = "gamma"
	port        = 8082
)
`)
	t.Chdir(dir)

	got, err := GetLastUsedPort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 8085 {
		t.Fatalf("expected 8085, got %d", got)
	}
}

func TestGetLastUsedPort_MissingPortLine(t *testing.T) {
	t.Run("only service has no port line", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "cmd/script-service/app/service.go", `package app

const (
	serviceName = "script"
)
`)
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != DefaultStartingPort {
			t.Fatalf("expected %d, got %d", DefaultStartingPort, got)
		}
	})

	t.Run("falls back to highest other service port", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "cmd/empty-service/app/service.go", `package app

const (
	serviceName = "empty"
)
`)
		writeTestFile(t, dir, "cmd/script-service/app/service.go", `package app

const (
	serviceName = "script"
	port        = 4001
)
`)
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 4001 {
			t.Fatalf("expected 4001, got %d", got)
		}
	})
}

func TestGetLastUsedPort_PortLineMatching(t *testing.T) {
	t.Run("ignores comments and other identifiers", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "cmd/script-service/app/service.go", `package app

// port = 9999
const (
	defaultPort = 9998
	export      = 9997
	serviceName = "script"
	port        = 4001 // local dev port
)
`)
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 4001 {
			t.Fatalf("expected 4001, got %d", got)
		}
	})

	t.Run("matches single-line const declaration", func(t *testing.T) {
		dir := t.TempDir()
		writeTestFile(t, dir, "cmd/script-service/app/service.go", `package app

const port = 4002
`)
		t.Chdir(dir)

		got, err := GetLastUsedPort()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 4002 {
			t.Fatalf("expected 4002, got %d", got)
		}
	})
}

type failingFS struct {
	err error
}

func (f failingFS) Open(string) (fs.File, error) {
	return nil, f.err
}

func TestGetLastUsedPort_WalkError(t *testing.T) {
	_, err := getLastUsedPort(failingFS{err: fs.ErrPermission})
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected permission error, got %v", err)
	}
}

func TestGetPortFromFile(t *testing.T) {
	fsys := fstest.MapFS{
		"script-service/app/service.go": &fstest.MapFile{
			Data: []byte("package app\n\nconst (\n\tport = 4001\n)\n"),
		},
		"empty-service/app/service.go": &fstest.MapFile{
			Data: []byte("package app\n"),
		},
	}

	got, err := getPortFromFile(fsys, "script-service/app/service.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 4001 {
		t.Fatalf("expected 4001, got %d", got)
	}

	got, err = getPortFromFile(fsys, "empty-service/app/service.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}

	if _, err := getPortFromFile(fsys, "missing/app/service.go"); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
