package parsers

import (
	"FileCrawler/internal/ports"
	"FileCrawler/internal/spells"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type SpellParser struct {
}

func NewSpellHandler(allSpells *[]*spells.Spell, parser ports.Parser) func(*colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		spells, err := parser.Parse(e)
		if err != nil {
			fmt.Println("Error parsing spell:", err)
			return
		}
		*allSpells = append(*allSpells, spells...)
	}
}

func (p *SpellParser) Parse(e *colly.HTMLElement) ([]*spells.Spell, error) {
	spellList := make([]*spells.Spell, 0)
	var spell = spells.Spell{}

	spell.Name = e.ChildText("h1.page-header")
	var err error
	spell.Level, err = getSpellLevel(e)
	spell.School = getSchool(e)
	spell.Tags = getTags(e)
	spell.Classes = getClasses(e)
	spell.Components = getComponents(e)
	spell.Range = getRange(e)
	spell.Ritual = isRitual(e)
	spell.CastingTime = getCastingTime(e)
	spell.Duration = getDuration(e)
	spell.Target = getTarget(e)
	spell.SavingThrow = getSavingThrow(e)
	spell.Texts = getTexts(e)
	spell.Source = getSource(e)
	if err != nil {
		fmt.Println("Error getting spell level:", err)
	}
	// Print the whole spell struct to the console and append it to the spellList
	fmt.Printf("Spell: %+v\n", spell)
	spellList = append(spellList, &spell)
	return spellList, nil
}

// Helper functions to pull data from the HTML pages
func getSpellLevel(e *colly.HTMLElement) (int, error) {
	level := -1
	levelTxt := e.ChildText(".field--name-field-spell-level a")
	if levelTxt != "" {
		levelTxt = strings.TrimSpace(levelTxt)
		if strings.EqualFold(levelTxt, "Cantrip") {
			level = 0
			return level, nil
		}
		// Remove the last two characters from levelText
		if len(levelTxt) > 2 {
			levelTxt = levelTxt[:len(levelTxt)-2]
		}

		// Attempt to parse the remaining string into an integer
		level, err := strconv.Atoi(strings.TrimSpace(levelTxt))
		if err != nil {
			// Handle the error, for example, by returning a default value or forwarding the error
			return 0, err
		}

		return level, nil
	}
	err := fmt.Errorf("no level found")
	return level, err
}

func getSchool(e *colly.HTMLElement) string {
	return e.ChildText(".field--name-field-classical-spell-school a")
}

func getTags(e *colly.HTMLElement) []string {
	tags := make([]string, 0)
	e.ForEach(".field--name-field-spell-schools .field--item a", func(_ int, el *colly.HTMLElement) {
		tags = append(tags, el.Text)
	})
	return tags
}

func getClasses(e *colly.HTMLElement) []string {
	classes := make([]string, 0)
	e.ForEach(".field--name-field-spell-classes .field--item a", func(_ int, el *colly.HTMLElement) {
		classes = append(classes, el.Text)
	})
	return classes
}

func getComponents(e *colly.HTMLElement) []string {
	components := make([]string, 0)
	e.ForEach("#spell-components-display .component-value a", func(_ int, el *colly.HTMLElement) {
		components = append(components, el.Text)
	})
	e.ForEach("div.field.field--name-field-spellcomponent-description.field--type-string.field--label-hidden.field--item", func(_ int, el *colly.HTMLElement) {
		components = append(components, el.Text)
	})
	return components
}

func getRange(e *colly.HTMLElement) string {
	return e.ChildText(".field--name-field-spell-range .field--item a")
}

func isRitual(e *colly.HTMLElement) bool {
	return e.ChildText(".ritual-note .ritual-indicator") != ""
}

func getCastingTime(e *colly.HTMLElement) string {
	return e.ChildText(".field--name-field-spell-casting-time .field--item")
}

func getDuration(e *colly.HTMLElement) string {
	return e.ChildText("#duration .duration-value a")
}

func getTarget(e *colly.HTMLElement) string {
	return e.ChildText(".field--name-field-spell-target .field--item")
}

func getSavingThrow(e *colly.HTMLElement) string {
	return e.ChildText(".field.field--name-field-spell-saving-throw-desc .field--item")
}

func getTexts(e *colly.HTMLElement) []string {
	selectors := []string{
		"#spell-body .field.field--name-body.field--type-text-with-summary.field--label-hidden.field--item p",
		".field.field--name-field-spellcast-at-higher-levels .field--label",
		".field--name-field-spellcast-at-higher-levels .field--item p",
		".field--name-field-spell-rare-versions .field--label",
		".field--name-field-spell-rare-versions .field--item p",
	}
	var texts []string

	for _, selector := range selectors {
		e.ForEach(selector, func(_ int, elem *colly.HTMLElement) {
			// Clean up the text similar to the previous example
			text := strings.TrimSpace(elem.Text)
			text = strings.ReplaceAll(text, "\n", "")
			text = strings.ReplaceAll(text, "\r", "")
			text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
			texts = append(texts, text)
		})
	}

	return texts
}

func getSource(e *colly.HTMLElement) string {
	dirtySource := e.ChildText(".field--name-field-spell-source .field--item a")
	source := strings.TrimSpace(dirtySource)
	return removeHiddenChars(source)
}

func removeHiddenChars(input string) string {
	// Compile regular expression
	re := regexp.MustCompile("[^\x20-\x7E]")
	// Replace non-printable characters with empty string
	return re.ReplaceAllString(input, "")
}
