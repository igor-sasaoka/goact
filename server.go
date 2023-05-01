package goact

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
    mu sync.Mutex
    connections map[*conn]struct{}
    actionHandler *ActionHandler
}

func (s *Server) listenAndServe(addr string) error {
    l, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    } 

    log.Printf("Listening to %s...\n", addr)
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Print(err)
        }

        log.Printf("connection accepted\n")
        go s.handleConnection(conn)
    }
}

func (s *Server) trackConn(c *conn, add bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if s.connections == nil {
        s.connections = make(map[*conn]struct{})
    }
    if add {
        s.connections[c] = struct{}{}
    } else {
        delete(s.connections, c)
    }
}

func (s *Server) newConn(c net.Conn) *conn {
    conn := new(conn).init(c, s)
    s.trackConn(conn, true)
    return conn
}

func (s *Server) handleConnection(c net.Conn) {
    conn := s.newConn(c)
    defer conn.Close()
    for {
        line, err := conn.r.ReadSlice('\n') 
        if err == io.EOF {
            log.Print("connection ended")
            return
        }
        if err != nil {
            log.Print("unexpected error:", err)
            return
        } 
        m, messageErr := DecodeMessage(line)
        if messageErr != nil {
            log.Print(messageErr)
            return
        }
        s.actionHandler.executeAction(m.Action, conn.w, m.Body)
    }
}

type conn struct {
    net.Conn
    server *Server
    r *bufio.Reader
    w MessageWriter 
}

func (c *conn) init(conn net.Conn, s *Server) *conn {
    c.Conn = conn
    c.server = s
    c.r = bufio.NewReader(c.Conn)
    c.w = bufio.NewWriter(c.Conn)

    return c
}

func (c *conn) Close() error {
    c.server.trackConn(c, false)
    return c.Conn.Close()
}

//Starts a new server listening to the given address and will execute actions
//based on incoming messages matching the pattern in the action handler
func ListenAndServe(addr string, handler *ActionHandler) error {
    s := &Server{actionHandler: handler}
    return s.listenAndServe(addr)
}
