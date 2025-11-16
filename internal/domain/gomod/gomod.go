package gomod

import "DependencyChecker-DC-/internal/domain/dependency"

type GoModInfo struct {
	Path         string
	GoVersion    string
	Dependencies []dependency.DependencyInfo
}
