package usecase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLooksLikeText(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{name: "plain text", data: []byte("package main\n"), want: true},
		{name: "empty", data: []byte{}, want: true},
		{name: "null byte", data: []byte("abc\x00def"), want: false},
		{name: "invalid utf8", data: []byte{0xff, 0xfe, 0xfd}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := looksLikeText(tt.data)
			if got != tt.want {
				t.Errorf("looksLikeText(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestReadNodeContent(t *testing.T) {
	dir := t.TempDir()

	textPath := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(textPath, []byte("env: local\n"), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	binPath := filepath.Join(dir, "data.bin")
	err = os.WriteFile(binPath, []byte{0x00, 0x01, 0x02}, 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	oversizedPath := filepath.Join(dir, "big.txt")
	err = os.WriteFile(oversizedPath, make([]byte, maxContentSize+1), 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	emptyPath := filepath.Join(dir, "empty.txt")
	err = os.WriteFile(emptyPath, []byte{}, 0600)
	if err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read dir: %v", err)
	}

	byName := make(map[string]os.DirEntry)
	for _, e := range entries {
		byName[e.Name()] = e
	}

	tests := []struct {
		file string
		want string
	}{
		{file: "config.yaml", want: "env: local\n"},
		{file: "data.bin", want: ""},
		{file: "big.txt", want: ""},
		{file: "empty.txt", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got := readNodeContent(filepath.Join(dir, tt.file), byName[tt.file])
			if got != tt.want {
				t.Errorf("readNodeContent(%s) = %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}
