package analysis

import "fmt"

type PackageInfo struct {
	Name, Version string
}

func (p *PackageInfo) ToString() string {
	if p.Version == "" {
		return p.Name
	}

	return fmt.Sprintf("%s@%s", p.Name, p.Version)
}
