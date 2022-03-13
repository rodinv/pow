package delivery

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Handler is a set of request handlers
type Handler struct {
	quotes QuotesProvider
	pow    ProofOfWorkProvider
}

// QuotesProvider provides work with a set of quotes
type QuotesProvider interface {
	Get() string
}

// ProofOfWorkProvider provides work with PoW algorithm
type ProofOfWorkProvider interface {
	GetHeader(resource string) (string, error)
	Compute(header string) string
	Verify(data string) error
}

func New(q QuotesProvider, p ProofOfWorkProvider) *Handler {
	return &Handler{
		quotes: q,
		pow:    p,
	}
}

// Challenge responses a new challenge
func (h *Handler) Challenge(_, sender string, w *bufio.Writer) {
	header, err := h.pow.GetHeader(sender)
	if err != nil {
		log.Err(err).Str("remote", sender).Msg("new challenge creating failed")
		err = writeResponse(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Err(err).Str("remote", sender).Msg("writing response failed")
		}
		return
	}
	log.Info().Str("remote", sender).Str("header", header).Msg("new challenge created")

	err = writeResponse(w, http.StatusOK, header)
	if err != nil {
		log.Err(err).Str("remote", sender).Msg("writing response failed")
	}
}

// GetQuote verifies hash and responses random quote
func (h *Handler) GetQuote(in, sender string, w *bufio.Writer) {
	err := h.pow.Verify(in)
	if err != nil {
		log.Err(err).Str("remote", sender).Msg("hash verify failed")
		err = writeResponse(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Err(err).Str("remote", sender).Msg("writing response failed")
		}
		return
	}

	log.Info().Str("remote", sender).Msg("hash verify success")
	q := h.quotes.Get()

	err = writeResponse(w, http.StatusOK, q)
	if err != nil {
		log.Err(err).Str("remote", sender).Msg("writing response failed")
	}
}

func writeResponse(w *bufio.Writer, code int, data string) error {
	_, err := w.Write([]byte(fmt.Sprintf("%d %s\n", code, data)))
	if err != nil {
		return err
	}

	return w.Flush()
}
