package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type GitHubChecker struct{}

func parseGitHubURL(url string) (string, string) {
	parts := strings.Split(url, "/")
	return parts[3], parts[4]
}

func (g GitHubChecker) FindGoModFiles(source string) ([]string, error) {

	owner, repo := parseGitHubURL(source)
	api := "https://api.github.com/repos/" + owner + "/" + repo + "/git/trees/main?recursive=1"

	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tree struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}

	err = json.Unmarshal(body, &tree)
	if err != nil {
		return nil, err
	}

	var goModFiles []string
	for _, t := range tree.Tree {
		if t.Type == "blob" && strings.HasSuffix(t.Path, "go.mod") {
			goModFiles = append(goModFiles, t.Path)
		}
	}

	return goModFiles, nil
}

func (g GitHubChecker) GetGoModFile(goFileSource, path string) ([]byte, error) {
	gitHubSourceRaw := strings.Replace(goFileSource, "github.com", "raw.githubusercontent.com", 1)
	gitHubSourceRaw = gitHubSourceRaw + "/main/" + path

	resp, err := http.Get(gitHubSourceRaw)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
