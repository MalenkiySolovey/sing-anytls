package session

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/anytls/sing-anytls/padding"
	"github.com/sagernet/sing/common/atomic"
	"github.com/sagernet/sing/common/logger"
)

// A full-duplex owner fixture: reads consume peer settings while writes are
// accepted independently. It records no credentials and owns one close signal.
type concurrencyConn struct {
	net.Conn
	reader  *bytes.Reader
	drained chan struct{}
	closed  chan struct{}
	once    sync.Once
	drain   sync.Once
}

func (c *concurrencyConn) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	if err == io.EOF {
		c.drain.Do(func() { close(c.drained) })
		<-c.closed
	}
	return n, err
}
func (c *concurrencyConn) Write(p []byte) (int, error)      { return len(p), nil }
func (c *concurrencyConn) SetWriteDeadline(time.Time) error { return nil }
func (c *concurrencyConn) Close() error {
	c.once.Do(func() { close(c.closed) })
	return nil
}

func concurrencySession() (*Session, *concurrencyConn) {
	var pad atomic.TypedValue[*padding.PaddingFactory]
	padding.UpdatePaddingScheme(padding.DefaultPaddingScheme, &pad)
	// Actual version-2 server-settings frames, read by the real recvLoop.
	frame := []byte{cmdServerSettings, 0, 0, 0, 0, 0, 3, 'v', '=', '2'}
	conn := &concurrencyConn{reader: bytes.NewReader(bytes.Repeat(frame, 1024)),
		drained: make(chan struct{}), closed: make(chan struct{})}
	return NewClientSession(conn, &pad, logger.NOP()), conn
}

func TestPeerVersionNegotiationConcurrentOpenStream(t *testing.T) {
	session, conn := concurrencySession()
	defer session.Close()
	start := make(chan struct{})
	received := make(chan error, 1)
	opened := make(chan error, 1)
	go func() { <-start; received <- session.recvLoop() }()
	go func() {
		<-start
		for index := 0; index < 1024; index++ {
			stream, err := session.OpenStream()
			if err != nil {
				opened <- err
				return
			}
			_ = stream.Close()
		}
		opened <- nil
	}()
	close(start)
	select {
	case err := <-opened:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("stream creation did not finish")
	}
	select {
	case <-conn.drained:
	case <-time.After(5 * time.Second):
		t.Fatal("peer negotiation did not finish")
	}
	// Check negotiated behavior, without coupling the test to field storage.
	stream, err := session.OpenStream()
	if err != nil {
		t.Fatal(err)
	}
	if session.synDone == nil {
		t.Fatal("version-2 stream confirmation was not enabled")
	}
	_ = stream.Close()
	_ = session.Close()
	select {
	case <-received:
	case <-time.After(5 * time.Second):
		t.Fatal("receiver did not close")
	}
}

func TestStreamCloseConcurrentWrite(t *testing.T) {
	session, _ := concurrencySession()
	defer session.Close()
	stream, err := session.OpenStream()
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var done sync.WaitGroup
	done.Add(2)
	go func() {
		defer done.Done()
		<-start
		for index := 0; index < 1024; index++ {
			_, _ = stream.Write([]byte{1})
		}
	}()
	go func() { defer done.Done(); <-start; _ = stream.Close() }()
	close(start)
	done.Wait()
	if _, err := stream.Write([]byte{1}); err == nil {
		t.Fatal("closed stream still accepted a write")
	}
}
