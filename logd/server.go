//a simple log server to receive remote log and write with specified handler
package main

import (
	"bufio"
	"encoding/binary"
	"github.com/siddontang/go-log/log"
	"io"
	"net"
	"strings"
)

//a log server for handling SocketHandler send log

type server struct {
	closed   bool
	listener net.Listener
	h        log.Handler
}

func newServer(addr string, h log.Handler) (*server, error) {
	s := new(server)

	s.closed = false

	var err error

	var protocol = "tcp"

	if strings.Contains(addr, "/") {
		protocol = "unix"
	}

	s.listener, err = net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}

	s.h = h

	return s, nil
}

func (s *server) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true

	s.h.Close()

	s.listener.Close()
	return nil
}

func (s *server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.onRead(conn)
	}
}

func (s *server) onRead(c net.Conn) {
	br := bufio.NewReaderSize(c, 1024)

	var bufLen uint32

	for {
		if err := binary.Read(br, binary.BigEndian, &bufLen); err != nil {
			c.Close()
			return
		}

		buf := make([]byte, bufLen, bufLen+1)

		if _, err := io.ReadFull(br, buf); err != nil && err != io.ErrUnexpectedEOF {
			c.Close()
			return
		} else {
			if len(buf) == 0 {
				continue
			}
			if buf[len(buf)-1] != '\n' {
				buf = append(buf, '\n')
			}

			s.h.Write(buf)
		}

	}
}
