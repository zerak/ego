package config

type Section struct {
	data         map[string]string // key:value
	dataOrder    []string
	dataComments map[string][]string // key:comments
	Name         string
	comments     []string
	Comment      string
}
