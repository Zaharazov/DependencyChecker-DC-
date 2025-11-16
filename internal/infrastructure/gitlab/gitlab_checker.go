package gitlab

import (
	"DependencyChecker-DC-/internal/infrastructure"
	"net/url"
)

type GitLabChecker struct{}

func (g GitLabChecker) FindGoModFiles(source string) ([]string, error) {

	group, project, err := infrastructure.ParseURL(source)
	if err != nil {
		return nil, err
	}

	projectPath := url.QueryEscape(group + "/" + project)
	apiUrl := "https://gitlab.com/api/v4/projects/" + projectPath + "/repository/tree?recursive=true"

	return infrastructure.Finder(apiUrl)
}

func (g GitLabChecker) GetGoModFile(source, path string) ([]byte, error) {

	gitLabSourceRaw := source + "/-/raw/main/" + path
	return infrastructure.GetBySourceRaw(gitLabSourceRaw)
}
