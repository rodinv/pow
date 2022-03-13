package client

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/rodinv/errors"
	"github.com/rs/zerolog/log"

	"github.com/rodinv/pow/internal/hashcash"
)

// Client is a client for tcp server
type Client struct {
	host string
	port string
}

func New(host, port string) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

// response tcp server response
type response struct {
	Code int
	Data string
}

// GetQuote gets new random quote with hashcash algorithm computing
func (c *Client) GetQuote() (string, error) {
	log.Info().Msg("start getting quote")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.host, c.port))
	if err != nil {
		return "", errors.Wrapf(err, "connection failed %s:%s", c.host, c.port)
	}
	defer func() {
		if errClose := conn.Close(); errClose != nil {
			errors.Combine(err, errClose)
		}
	}()

	// get new challenge
	challenge, err := c.getNewChallenge(conn)
	if err != nil {
		return "", errors.Wrap(err, "getting new challenge")
	}
	log.Info().Str("challenge", challenge).Msg("new challenge received")

	// compute
	bits, err := hashcash.GetBitsFromHeader(challenge)
	if err != nil {
		return "", errors.Wrap(err, "getting bits")
	}

	pow, err := hashcash.New(bits, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating pow instance")
	}

	log.Info().Msg("computing the challenge...")
	result := pow.Compute(challenge)
	log.Info().Str("result", challenge).Msg("the result is already being computed")

	// request quote
	return c.getNewQuote(conn, result)
}

func (c *Client) getNewQuote(conn net.Conn, challenge string) (string, error) {
	_, err := conn.Write([]byte(fmt.Sprintf("get_quote %s \n", challenge)))
	if err != nil {
		return "", errors.Wrap(err, "writing get_quote request")
	}

	resp, err := readResponse(conn)
	if err != nil {
		return "", errors.Wrap(err, "reading get_quote response")
	}

	if resp.Code != http.StatusOK {
		return "", errors.New(resp.Data)
	}

	return resp.Data, nil
}

func (c *Client) getNewChallenge(conn net.Conn) (string, error) {
	_, err := conn.Write([]byte("challenge \n"))
	if err != nil {
		return "", errors.Wrap(err, "writing challenge request")
	}

	resp, err := readResponse(conn)
	if err != nil {
		return "", errors.Wrap(err, "reading challenge response")
	}

	if resp.Code != http.StatusOK {
		return "", errors.New(resp.Data)
	}

	return resp.Data, nil
}

func readResponse(conn net.Conn) (*response, error) {
	reader := bufio.NewReader(conn)

	row, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	row = strings.TrimSuffix(row, "\n")

	pos := strings.IndexByte(row, ' ')
	if pos == -1 {
		return nil, errors.Errorf("wrong response: %s", row)
	}

	r := new(response)
	r.Code, err = strconv.Atoi(row[:pos])
	if err != nil {
		return nil, errors.Wrapf(err, "parsing code %s", row[:pos])
	}

	r.Data = row[pos+1:]

	return r, nil
}
