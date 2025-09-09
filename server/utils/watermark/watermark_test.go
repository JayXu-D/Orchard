package watermark

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestAddWatermark_GeneratesOutput(t *testing.T) {
	name := "d6d68c372b80ce337ab9eff1ab688862_20250907154149.jpg"
	// name := "d6d68c372b80ce337ab9eff1ab688862_20250819214417.jpg"
	bases := [][]string{
		{"uploads", "file", name},
		{"uploads", "files", name},
		{"server", "uploads", "file", name},
		{"server", "uploads", "files", name},
	}
	prefixes := []string{"", "..", filepath.Join("..", ".."), filepath.Join("..", "..", "..")}
	t.Logf("cwd=%s", mustGetwd(t))
	t.Logf("probing %d prefixes x %d bases", len(prefixes), len(bases))
	var inputPath string
	for _, p := range prefixes {
		for _, b := range bases {
			cand := filepath.Join(append([]string{p}, b...)...)
			t.Logf("check: %s", cand)
			if fi, err := os.Stat(cand); err == nil && fi.Size() > 0 {
				t.Logf("found: %s (size=%d)", cand, fi.Size())
				inputPath = cand
				break
			}
		}
		if inputPath != "" {
			break
		}
	}
	if inputPath == "" {
		t.Skipf("test image not found in expected locations (cwd=%s)", mustGetwd(t))
	}

	ws := NewWatermarkService()
	// Output to current working directory
	ws.cacheDir = "."
	t.Logf("cacheDir=%s", ws.cacheDir)
	logText := "Admin"
	t.Logf("input=%s, text=%q", inputPath, logText)

	outPath, err := ws.AddWatermark(inputPath, logText)
	if err != nil {
		t.Fatalf("AddWatermark error: %v", err)
	}
	t.Logf("outPath=%s", outPath)
	fi, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output stat error: %v", err)
	}
	t.Logf("outSize=%d", fi.Size())
	if fi.Size() == 0 {
		t.Fatalf("output is empty")
	}
	// Try to decode (supports png/jpeg automatically)
	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer func() { _ = f.Close() }()
	if _, format, err := image.Decode(f); err != nil {
		t.Fatalf("decode output: %v", err)
	} else {
		t.Logf("decoded format=%s", format)
	}
}

func mustGetwd(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}
