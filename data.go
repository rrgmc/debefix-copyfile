package copyfile

type FileDataList struct {
	Fields map[string]FileData `json:"fields"`
}

type FileData struct {
	ID       string `yaml:"id"`
	SetValue bool   `yaml:"setValue"`
	Src      string `yaml:"src"`
	Dest     string `yaml:"dest"`
}
