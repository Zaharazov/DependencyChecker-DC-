package main

import (
	checker "DependencyChecker-DC-/internal/domain/checker"
	infrastructure "DependencyChecker-DC-/internal/infrastructure/github"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func CheckGoModByPath(source, f string, checker checker.RepoChecker) error {
	data, err := checker.GetGoModFile(source, f)

	if err != nil {
		return err
	}

	bodyStruct, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("| Module Name:", bodyStruct.Module.Mod.Path)
	fmt.Println("|-------------------------------------------------------")
	fmt.Println("| Go Version:", bodyStruct.Go.Version)
	fmt.Println("|-------------------------------------------------------")
	canBeUpdated := false
	for _, req := range bodyStruct.Require {

		var actual struct {
			Version string `json:"Version"`
		}

		actualData := exec.Command("go", "list", "-m", "-json", req.Mod.Path+"@latest")
		output, err := actualData.Output()
		if err != nil {
			fmt.Println("| Error:", err)
			continue
		}
		err = json.Unmarshal(output, &actual)
		if err != nil {
			fmt.Println("| Error:", err)
			continue
		}

		if semver.Compare(req.Mod.Version, actual.Version) == -1 {
			fmt.Println("| Can Be Updated:", req.Mod.Path, "(", req.Mod.Version, "->", actual.Version, ")")
			canBeUpdated = true
		}
	}
	if !canBeUpdated {
		fmt.Println("| There Are No Possible Updates!")
	}

	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println()
		fmt.Println("| Usage (example): ./main https://github.com/{owner}/{repo}")
		return
	}

	source := os.Args[1]

	parts := strings.Split(source, "/")

	if len(parts) < 5 {
		fmt.Println()
		fmt.Println("| Usage (example): ./main https://github.com/{owner}/{repo}")
		return
	}

	var checker checker.RepoChecker
	if parts[2] == "github.com" {
		checker = infrastructure.GitHubChecker{}
	} else {
		fmt.Println("| Will Be Available In The Future!")
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
		err = CheckGoModByPath(source, p, checker)
		if err != nil {
			fmt.Println("| Error:", err)
		}
	}

}
