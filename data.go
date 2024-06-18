package copyfile

type FileDataList struct {
	Fields map[string]FileData `json:"fields"`
}

type FileData struct {
	ID          string  `yaml:"id"`
	Value       *string `yaml:"value"`
	Source      string  `yaml:"source"`
	Destination string  `yaml:"destination"`
}
