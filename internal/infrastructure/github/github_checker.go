package github

import (
	"DependencyChecker-DC-/internal/infrastructure"
	"strings"
)

type GitHubChecker struct{}

func (g GitHubChecker) FindGoModFiles(source string) ([]string, error) {

	owner, repo, err := infrastructure.ParseURL(source)
	if err != nil {
		return nil, err
	}
	apiUrl := "https://api.github.com/repos/" + owner + "/" + repo + "/git/trees/main?recursive=1"

	return infrastructure.Finder(apiUrl)
}

func (g GitHubChecker) GetGoModFile(goFileSource, path string) ([]byte, error) {

	gitHubSourceRaw := strings.Replace(goFileSource, "github.com", "raw.githubusercontent.com", 1)
	gitHubSourceRaw = gitHubSourceRaw + "/main/" + path
	return infrastructure.GetBySourceRaw(gitHubSourceRaw)
}
