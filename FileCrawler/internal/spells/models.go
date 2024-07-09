package spells

type Spell struct {
	Name        string
	Level       int
	School      string
	CastingTime string
	Target      string
	Range       string
	Components  []string
	Duration    string
	Tags        []string
	Ritual      bool
	Classes     []string
	Source      string
	SavingThrow string
	Texts       []string
}
