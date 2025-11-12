package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func main() {

	if len(os.Args) < 2 {
		log.Println("| Usage: ./main https://github.com/rep_name")
		os.Exit(1)
	}

	gitHubSource := os.Args[1]
	gitHubSourceRaw := strings.Replace(gitHubSource, "github.com", "raw.githubusercontent.com", 1)
	gitHubSourceRaw += "/main/go.mod"

	client := &http.Client{}
	request, err := http.NewRequest("GET", gitHubSourceRaw, nil)
	if err != nil {
		log.Println("| Error:", err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("| Error:", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("| Error:", err)
		return
	}

	bodyStruct, err := modfile.Parse("go.mod", body, nil)
	if err != nil {
		log.Println("| Error:", err)
		return
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
			log.Println("| Error:", err)
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
}
