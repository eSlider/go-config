package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunConvert_TOMLToJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("[app]\nport = 8080\n")
	code := RunConvert([]string{"--from", "toml", "--to", "json", "--input", "-"}, stdin, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("stderr: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"port"`) || !strings.Contains(stdout.String(), "8080") {
		t.Fatalf("out: %s", stdout.String())
	}
}

func TestRunConvert_YAMLToJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("a: 1\n")
	code := RunConvert([]string{"--from", "yaml", "--to", "json", "--input", "-"}, stdin, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("stderr: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"a"`) {
		t.Fatalf("out: %s", stdout.String())
	}
}

func TestRunGet(t *testing.T) {
	var stdout, stderr bytes.Buffer
	in := "service:\n  name: x\n"
	code := RunGet([]string{"--from", "yaml", "--path", "service.name", "-"}, strings.NewReader(in), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("stderr: %s", stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "x" {
		t.Fatalf("got %q", stdout.String())
	}
}

func TestRunMerge_TwoYAMLFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.yaml")
	f2 := filepath.Join(dir, "b.yaml")
	if err := os.WriteFile(f1, []byte("k: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("k: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := RunMerge([]string{"--from", "yaml", "--to", "json", f1, f2}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("stderr: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"k"`) || !strings.Contains(stdout.String(), `2`) {
		t.Fatalf("out: %s", stdout.String())
	}
}
