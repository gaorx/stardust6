package sdblueprint

type Anns map[string][]string

func (la Anns) Get(lang string, langAliases ...string) []string {
	langs := append([]string{lang}, langAliases...)
	r := findLangValue(la, langs)
	if r != nil {
		return r
	}
	return nil
}
