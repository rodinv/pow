package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/rodinv/errors"
	"github.com/rs/zerolog/log"
)

type HandleFunc func(input, sender string, w *bufio.Writer)

// Server is a PoW server
type Server struct {
	host     string
	port     string
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup

	handlers map[string]HandleFunc
}

func New(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
		quit: make(chan struct{}),

		handlers: make(map[string]HandleFunc),
	}
}

// Run runs the server on the given host and port
func (s *Server) Run() error {
	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return errors.Wrap(err, "listening tcp")
	}

	log.Info().Str("host", s.host).Str("port", s.port).Msg("Server starts listening")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					log.Err(err).Msg("accepting connection")

					return
				}
			}

			log.Info().Str("remote", conn.RemoteAddr().String()).Msg("connection accepted")
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()

				s.handleRequest(conn)
			}()
		}
	}()

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

// AddHandler adds handler function
func (s *Server) AddHandler(name string, handler HandleFunc) {
	s.handlers[name] = handler
}

func (s *Server) getHandler(input string) (HandleFunc, string, error) {
	pos := strings.IndexByte(input, ' ')
	if pos == -1 {
		return nil, "", errors.Errorf("wrong request format %s", input)
	}

	h, ok := s.handlers[input[:pos]]
	if !ok {
		return nil, "", errors.Errorf("unknown handler %s", input[:pos])
	}

	return h, strings.TrimSuffix(input[pos+1:], " \n"), nil
}

func (s *Server) handleRequest(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()

	var err error
	defer func() {
		if err != nil {
			log.Err(err).Str("remote", remoteAddr).Msg("handle request")
		}
	}()

	reader := bufio.NewReader(conn)
	defer func() {
		if errClose := conn.Close(); errClose != nil {
			errors.Combine(err, errClose)
		}
	}()

	for {
		select {
		case <-s.quit:
			return
		default:
			var (
				row   string
				value string
				h     HandleFunc
			)
			row, err = reader.ReadString('\n')
			switch {
			case errors.Is(err, io.EOF):
				return
			case err != nil:
				err = errors.Wrap(err, "reading string")
				return
			}

			log.Info().Str("remote", remoteAddr).Str("value", row).Msg("reading from conn")

			h, value, err = s.getHandler(row)
			if err != nil {
				err = errors.Wrap(err, "getting handler")
				return
			}

			h(value, remoteAddr, bufio.NewWriter(conn))
		}
	}
}
