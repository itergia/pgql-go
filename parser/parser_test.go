package parser

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	des, err := os.ReadDir("./testdata/spec")
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	for _, de := range des {
		if !strings.HasSuffix(de.Name(), ".pgql") {
			continue
		}

		de := de
		t.Run(de.Name(), func(t *testing.T) {
			t.Parallel()

			bs, err := os.ReadFile(filepath.Join("./testdata/spec", de.Name()))
			if err != nil {
				t.Fatalf("ReadFile failed: %v", err)
			}
			bs = append(bs, ';', '\n')

			if err := Parse(bytes.NewReader(bs)); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
		})
	}
}
