package analysis

import (
	"errors"
	"os/exec"
	"strings"
)

type ProjectDependency struct {
	*PackageInfo
	UsedBy []*PackageInfo
}

func NewPackageInfoFromString(s string) *PackageInfo {
	tokens := strings.Split(s, "@")

	var name, version string

	if len(tokens) == 1 {
		name = tokens[0]
	} else if len(tokens) == 2 {
		name, version = tokens[0], tokens[1]
	} else {
		return nil
	}

	return &PackageInfo{
		Name: name,
		Version: version,
	}
}

func GetProjectDependencies(projectPath string) ([]*ProjectDependency, error) {
	deps := make([]*ProjectDependency, 0, 10)

	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = projectPath
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	str := string(data)

	if str == "" {
		return nil, errors.New("bruh")
	}
	if str == "go: cannot find main module; see 'go help modules'" {
		return nil, errors.New("not in go project")
	}

	lines := strings.Split(str, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")

		if len(tokens) != 2 {
			continue
		}

		dependentPkg := NewPackageInfoFromString(tokens[0])
		dependencyPkg := NewPackageInfoFromString(tokens[1])

		var dependency *ProjectDependency
		for _, dep := range deps {
			if dep.Name == dependencyPkg.Name && dep.Version == dependencyPkg.Version {
				dependency = dep
			}
		}

		if dependency != nil {
			dependency.UsedBy = append(dependency.UsedBy, dependentPkg)
		} else {
			dependency = &ProjectDependency{
				PackageInfo: dependencyPkg,
				UsedBy: []*PackageInfo { dependentPkg },
			}
			deps = append(deps, dependency)
		}
	}

	return deps, nil
}
