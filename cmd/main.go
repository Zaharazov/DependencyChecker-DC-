package main

import (
	checker "DependencyChecker-DC-/internal/domain/checker"
	"DependencyChecker-DC-/internal/domain/dependency"
	"DependencyChecker-DC-/internal/domain/gomod"
	"DependencyChecker-DC-/internal/infrastructure/github"
	"DependencyChecker-DC-/internal/infrastructure/gitlab"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func checkGoModByPath(source, f string, checker checker.RepoChecker) (*gomod.GoModInfo, int, error) {

	errors := 0

	data, err := checker.GetGoModFile(source, f)

	if err != nil {
		return nil, 0, err
	}

	bodyStruct, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, 0, err
	}

	fileInfo := &gomod.GoModInfo{
		Path:         bodyStruct.Module.Mod.Path,
		GoVersion:    bodyStruct.Go.Version,
		Dependencies: []dependency.DependencyInfo{},
	}

	for _, req := range bodyStruct.Require {

		var actual struct {
			Version string `json:"Version"`
		}

		actualData := exec.Command("go", "list", "-m", "-json", req.Mod.Path+"@latest")
		output, err := actualData.Output()
		if err != nil {
			errors++
			continue
		}
		err = json.Unmarshal(output, &actual)
		if err != nil {
			errors++
			continue
		}

		if semver.Compare(req.Mod.Version, actual.Version) == -1 {
			fileInfo.Dependencies = append(fileInfo.Dependencies, dependency.DependencyInfo{
				Name:        req.Mod.Path,
				CurVersion:  req.Mod.Version,
				LastVersion: actual.Version,
			})
		}
	}

	return fileInfo, errors, nil
}

func determineChecker(source string) (checker.RepoChecker, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	switch u.Host {
	case "github.com":
		return github.GitHubChecker{}, nil
	case "gitlab.com":
		return gitlab.GitLabChecker{}, nil
	default:
		return nil, fmt.Errorf("Host Not Supported: %s", u.Host)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println()
		fmt.Println("| Usage (example): ./main https://github.com/{owner}/{repo}")
		return
	}

	source := os.Args[1]
	checker, err := determineChecker(source)

	if err != nil {
		fmt.Println("| Error:", err)
		return
	}

	paths, err := checker.FindGoModFiles(source)
	if err != nil {
		fmt.Println("| Error:", err)
		return
	}

	if len(paths) == 0 {
		fmt.Println()
		fmt.Println("| No Matching Files!")
		return
	}

	for _, p := range paths {
		result, errors, err := checkGoModByPath(source, p, checker)
		if err != nil {
			fmt.Println("| Error:", err)
			continue
		}

		fmt.Println()
		fmt.Println("| Module Name:", result.Path)
		fmt.Println("|-------------------------------------------------------")
		fmt.Println("| Go Version:", result.GoVersion)
		fmt.Println("|-------------------------------------------------------")

		if len(result.Dependencies) == 0 {
			fmt.Println("| There Are No Possible Updates!")
		} else {
			for _, d := range result.Dependencies {
				fmt.Println("| Can Be Updated:", d.Name, "(", d.CurVersion, "->", d.LastVersion, ")")
			}
		}

		if errors > 0 {
			fmt.Println("| Not Processed Due To Errors:", errors)
		}
	}

}
