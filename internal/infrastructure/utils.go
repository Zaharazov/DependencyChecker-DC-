package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Finder(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var goModFiles []string

	var treeH struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}

	err = json.Unmarshal(body, &treeH)
	if err == nil {
		for _, t := range treeH.Tree {
			if t.Type == "blob" && strings.HasSuffix(t.Path, "go.mod") {
				goModFiles = append(goModFiles, t.Path)
			}
		}
		return goModFiles, nil
	}

	var treeL []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	}

	err = json.Unmarshal(body, &treeL)
	if err == nil {
		for _, t := range treeL {
			if t.Type == "blob" && strings.HasSuffix(t.Path, "go.mod") {
				goModFiles = append(goModFiles, t.Path)
			}
		}
		return goModFiles, nil
	}

	return nil, err

}

func ParseURL(source string) (string, string, error) {
	u, err := url.Parse(source)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")

	return parts[0], parts[1], nil
}

func GetBySourceRaw(raw string) ([]byte, error) {

	resp, err := http.Get(raw)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
