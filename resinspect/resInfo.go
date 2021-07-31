package resinspect

type ResInfo struct {
	Name string
	Path string
}

func newResInfo(name, path string) *ResInfo {
	return &ResInfo{
		Name: name,
		Path: path,
	}
}