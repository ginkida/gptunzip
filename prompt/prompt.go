package prompt

import (
	"archive/zip"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type GitFile struct {
	Path     string `json:"path"`
	Tokens   int64  `json:"tokens"`
	Contents string `json:"contents"`
}

type GitRepo struct {
	TotalTokens int64     `json:"total_tokens"`
	Files       []GitFile `json:"files"`
	FileCount   int       `json:"file_count"`
}

func FindSingleSubdir(dirPath string) (string, error) {
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}
	for _, d := range dirs {
		if d.IsDir() {
			return filepath.Join(dirPath, d.Name()), nil
		}
	}
	return "", fmt.Errorf("no subdirectories found in: %s", dirPath)
}

func CreateTempDirectory() (string, error) {
	return os.MkdirTemp(".", "repo")
}

func UnzipFile(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func OutputGitRepo(repo *GitRepo) (string, error) {
	var repoBuilder strings.Builder
	for _, file := range repo.Files {
		repoBuilder.WriteString("----\n")
		repoBuilder.WriteString(fmt.Sprintf("%s\n", file.Path))
		repoBuilder.WriteString(fmt.Sprintf("%s\n", file.Contents))
	}
	repoBuilder.WriteString("--END--")
	output := repoBuilder.String()
	repo.TotalTokens = EstimateTokens(output)
	return output, nil
}

func SaveTextAsOneFile(path, text string) error {
	baseName := filepath.Base(path)
	fileName := fmt.Sprintf("%s.txt", strings.TrimSuffix(baseName, filepath.Ext(baseName)))
	if err := os.WriteFile(fileName, []byte(text), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	fmt.Printf("File %s has been saved successfully.\n", fileName)
	return nil
}

func SaveTextAsParts(path, text string) error {
	const maxFileSize = 10 * 1024 * 1024
	baseName := filepath.Base(path)
	fileBaseName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	textBytes := []byte(text)
	totalLength := len(textBytes)
	numFiles := (totalLength-1)/maxFileSize + 1
	for i := 0; i < numFiles; i++ {
		start := i * maxFileSize
		end := start + maxFileSize
		if end > totalLength {
			end = totalLength
		}
		fileName := fmt.Sprintf("%s_%d.txt", fileBaseName, i+1)
		if err := os.WriteFile(fileName, textBytes[start:end], 0644); err != nil {
			return fmt.Errorf("failed to write to file %s: %w", fileName, err)
		}
		fmt.Printf("File %s has been saved successfully.\n", fileName)
	}
	return nil
}

func ProcessGitRepo(repoPath string) (*GitRepo, error) {
	var repo GitRepo
	if err := processRepository(repoPath, &repo); err != nil {
		return nil, fmt.Errorf("error processing repository: %w", err)
	}
	return &repo, nil
}

func processRepository(filePath string, repo *GitRepo) error {
	return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			contents, err := os.ReadFile(path)

			if err != nil || !utf8.Valid(contents) {
				return err
			}

			relPath, err := filepath.Rel(filePath, path)

			if err != nil {
				return err
			}

			file := GitFile{
				Path:     filepath.ToSlash(relPath),
				Contents: string(contents),
				Tokens:   EstimateTokens(string(contents)),
			}
			repo.Files = append(repo.Files, file)
		}
		return nil
	})
}

func EstimateTokens(output string) int64 {
	return int64(math.Ceil(float64(len(output)) / 3.5))
}
