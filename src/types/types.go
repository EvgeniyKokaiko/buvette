package types

type Application struct {
	Name     string
	Usage    string
	Author   string
	Version  string
	HelpInfo string
}

type Command struct {
	Args   string
	Config map[string]string
}

type Config struct {
	PORT string
	PATH string
}
