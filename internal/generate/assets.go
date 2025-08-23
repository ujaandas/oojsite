package generate

import (
	"oojsite/assets"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func WriteStatic(outDir string) error {
	staticFS, err := assets.Static()
	if err != nil {
		return err
	}
	return fs.WalkDir(staticFS, ".", func(relPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		dst := filepath.Join(outDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(dst, 0755)
		}

		in, err := staticFS.Open(relPath)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}
