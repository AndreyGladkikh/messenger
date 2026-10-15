package query

import (
	"context"
	"strings"
	"unicode"
	"uuid"
)

const SearchQueryName = "Search"

type Search struct {
	Str string
}

func (g Search) IsQuery() {}

func (g Search) Name() string {
	return SearchQueryName
}

type SearchResultItem struct {
	ID          uuid.UUID
	Type        string
	Handle      string
	DisplayName string
}

type Searcher interface {
	Search(ctx context.Context, terms []SearchTerm) ([]SearchResultItem, error)
}

type SearchHandler struct {
	searcher Searcher
}

func NewSearchHandler(
	reader Searcher,
) *SearchHandler {
	return &SearchHandler{
		searcher: reader,
	}
}

func (h *SearchHandler) Handle(ctx context.Context, qry Search) ([]SearchResultItem, error) {
	qryFields := strings.Fields(qry.Str)
	var terms []SearchTerm
	for _, f := range qryFields {
		terms = append(terms, NewSearchTerm(f))
	}

	result, err := h.searcher.Search(ctx, terms)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type SearchTerm struct {
	Val       string
	Translits []string
}

func NewSearchTerm(val string) SearchTerm {
	var tls []string
	switch detectScript(val) {
	case "Cyrillic":
		tls = translit(val)
	}

	return SearchTerm{
		Val:       val,
		Translits: tls,
	}
}

var scripts = map[string]*unicode.RangeTable{
	"Latin":    unicode.Latin,
	"Cyrillic": unicode.Cyrillic,
}

func detectScript(word string) string {
	var script string

	for _, r := range word {
		if !unicode.IsLetter(r) {
			continue
		}
		for name, rangeTable := range scripts {
			if unicode.Is(rangeTable, r) {
				if script == "" {
					script = name
				}
				if script != name {
					return "Mixed"
				}
			}
		}
	}

	return script
}

func translit(word string) []string {
	var tls []string
	translitsCh := make(chan string)
	go func() {
		for tl := range translitsCh {
			tls = append(tls, tl)
		}
	}()

	backtrack([]rune(word), 0, nil, translitsCh)
	close(translitsCh)

	return tls
}

func backtrack(word []rune, idx int, seq []string, translitsCh chan string) {
	if idx == len(word) {
		translitsCh <- strings.Join(seq, "")
		return
	}

	charTrnslts := tlMap[word[idx]]

	for _, ct := range charTrnslts {
		seq = append(seq, ct)
		backtrack(word, idx+1, seq, translitsCh)
		seq = seq[:len(seq)-1]
	}
}

var tlMap = map[rune][]string{
	'а': {"a"},
	'б': {"b"},
	'в': {"v"},
	'г': {"g"},
	'д': {"d"},
	'е': {"e", "ye"},
	'ё': {"yo", "jo", "e"},
	'ж': {"zh"},
	'з': {"z"},
	'и': {"i"},
	'й': {"j", "i", "y"},
	'к': {"k"},
	'л': {"l"},
	'м': {"m"},
	'н': {"n"},
	'о': {"o"},
	'п': {"p"},
	'р': {"r"},
	'с': {"s"},
	'т': {"t"},
	'у': {"u"},
	'ф': {"f"},
	'х': {"x", "kh"},
	'ц': {"c", "ts", "cz"},
	'ч': {"ch"},
	'ш': {"sh"},
	'щ': {"shh", "sch"},
	'ъ': {""},
	'ы': {"y", "w"},
	'ь': {"'", ""},
	'э': {"e"},
	'ю': {"yu", "ju"},
	'я': {"ya", "ja"},
}
