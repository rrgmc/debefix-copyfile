package copyfile

import (
	"bytes"
	"io"
)

type parser struct {
	r *bytes.Reader
}

func newParser(str string) parser {
	return parser{r: bytes.NewReader([]byte(str))}
}

func (p *parser) next() (rune, bool) {
	ch, _, err := p.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
		return 0, false
	}
	return ch, true
}

func (p *parser) unread() {
	err := p.r.UnreadRune()
	if err != nil {
		panic(err)
	}
}
