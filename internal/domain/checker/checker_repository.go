package checker

type RepoChecker interface {
	FindGoModFiles(source string) ([]string, error)
	GetGoModFile(goFileSource, path string) ([]byte, error)
}
