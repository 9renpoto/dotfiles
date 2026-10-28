package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRunPackages(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Brewfile"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	var gotName string
	var gotArgs []string
	var stdout, stderr bytes.Buffer
	wantErr := errors.New("brew failed")
	err := run(context.Background(), []string{"packages"}, root, &stdout, &stderr,
		func(_ context.Context, name string, args ...string) error {
			gotName, gotArgs = name, args
			return wantErr
		})
	if !errors.Is(err, wantErr) {
		t.Fatalf("run() error = %v, want %v", err, wantErr)
	}
	if gotName != "brew" || !reflect.DeepEqual(gotArgs, []string{"bundle", "--file=" + filepath.Join(root, "Brewfile")}) {
		t.Fatalf("command = %q %q", gotName, gotArgs)
	}
	if stdout.String() != "Installing Homebrew packages from Brewfile.\n" || stderr.Len() == 0 {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunPackagesRejectsMissingManifest(t *testing.T) {
	err := run(context.Background(), []string{"packages"}, t.TempDir(), &bytes.Buffer{}, &bytes.Buffer{},
		func(context.Context, string, ...string) error { t.Fatal("unexpected command"); return nil })
	if err == nil {
		t.Fatal("expected missing manifest error")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := run(context.Background(), []string{"unknown"}, t.TempDir(), &bytes.Buffer{}, &bytes.Buffer{},
		func(context.Context, string, ...string) error { t.Fatal("unexpected command"); return nil })
	if err == nil {
		t.Fatal("expected unknown command error")
	}
}
