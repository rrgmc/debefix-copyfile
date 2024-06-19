package copyfile

// FileData is the information of the !copyfile tag.
type FileData struct {
	ID          string  `yaml:"id"`
	Value       *string `yaml:"value"`
	Source      string  `yaml:"source"`
	Destination string  `yaml:"destination"`
}
