package app

import (
	"context"
	"fmt"
	"net"
	"testing"

	"gotest.tools/v3/assert"
)

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func TestApp(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)
	port, err := getFreePort()
	assert.NilError(t, err)
	srv.Start(context.Background(), fmt.Sprintf("localhost:%d", port))
	srv.Shutdown()
}
