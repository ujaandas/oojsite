package assets

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildTailwind(outDir, staticDir string) error {
	in := filepath.Join(staticDir, "styles.css")
	out := filepath.Join(outDir, "static", "styles.css")

	if _, err := os.Stat(in); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	cmd := exec.Command("tailwindcss", "--input", in, "--output", out, "--minify", "--content", fmt.Sprintf("%s/**/*.html", outDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func CopyStaticContents(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "styles.css" {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)
		if d.IsDir() {
			return os.MkdirAll(dstPath, d.Type().Perm())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}
