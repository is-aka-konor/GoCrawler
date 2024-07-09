package ports

import (
	"FileCrawler/internal/spells"

	"github.com/gocolly/colly/v2"
)

type Parser interface {
	Parse(e *colly.HTMLElement) ([]*spells.Spell, error)
}
