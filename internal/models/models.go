package models

type SpellList struct {
	Name     string `json:"name"`
	SpellUrl string `json:"spellUrl"`
}

type Spell struct {
	Name        string
	Level       int
	School      string
	CastingTime string
	Range       string
	Components  []string
	Duration    string
	Tags        []string
	Ritual      bool
	Classes     []string
	Source      string
	Texts       []string
}

type ParserSettings struct {
	StartPoint int
	EndPoint   int
	QueryParam string
	BaseURL    string
	Index      []string
}
