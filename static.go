package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

/*
We handle everything pertaining to our "static/" folder here.
All that really is is just building TailwindCSS
and copying over everything else.
*/

// Use `os/exec` to build TailwindCSS.
func buildTailwind(outDir string) error {
	// TODO: Make these options/user-changeable later
	in := filepath.Join("static", "styles.css")
	out := filepath.Join(outDir, "static", "styles.css")

	// Assumes the user has `tailwindcss` available
	cmd := exec.Command("tailwindcss", "--input", in, "--output", out, "--minify", "--content", "./**/*.html")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Walk over and copy everything inside `src` into `dst`, excluding `styles.css`.
func copyStaticContents(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		// Check for traversal errors
		if err != nil {
			return err
		}

		// We run this _after_ generating our CSS, so kkip to avoid replacing generated CSS
		if d.Name() == "styles.css" {
			return nil
		}

		// Calculate relative path for destination construction
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Construct the full destination path
		dstPath := filepath.Join(dst, relPath)

		// Make directory if not already exists
		if d.IsDir() {
			return os.MkdirAll(dstPath, d.Type().Perm())
		}

		// All good, so copy file
		return copyFile(path, dstPath)
	})
}

// Copy file _and_ path from `src` to `dst`.
func copyFile(src, dst string) error {
	// Get file info
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check regular file, we don't support links yet
	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// Open and read file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create dstFile file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file over
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Ensure written
	return dstFile.Sync()
}
