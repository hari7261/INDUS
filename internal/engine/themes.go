package engine

var themes = map[string]Theme{
	"saffron": {Name: "saffron", Prompt: "\033[38;5;208m"},
	"cobalt":  {Name: "cobalt", Prompt: "\033[34m"},
	"mint":    {Name: "mint", Prompt: "\033[32m"},
	"ember":   {Name: "ember", Prompt: "\033[31m"},
	"slate":   {Name: "slate", Prompt: "\033[36m"},
	"linen":   {Name: "linen", Prompt: "\033[97m"},
}

func defaultTheme() Theme {
	return themes["saffron"]
}

func themeByName(name string) (Theme, bool) {
	theme, ok := themes[name]
	return theme, ok
}

func themeNames() []string {
	return []string{"saffron", "cobalt", "mint", "ember", "slate", "linen"}
}
