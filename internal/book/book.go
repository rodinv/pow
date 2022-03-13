package book

import (
	"bytes"
	_ "embed"
	"math/rand"
	"time"
)

//go:embed quotes.txt
var quotesFile []byte

// Quotes is a set of quotes
type Quotes struct {
	data []string
}

func New() *Quotes {
	b := bytes.Split(quotesFile, []byte("\n"))
	d := make([]string, 0, len(b))

	for _, v := range b {
		d = append(d, string(v))
	}

	return &Quotes{data: d}
}

// Get gets random quote
func (q *Quotes) Get() string {
	rand.Seed(time.Now().Unix())
	return q.data[rand.Intn(len(q.data))]
}
