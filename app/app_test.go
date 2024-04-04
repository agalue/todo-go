package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
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

func newMockApp(dbFail bool) (*App, *MockDB) {
	db := &MockDB{fail: dbFail}
	db.Init()
	a := &App{
		db:     db,
		router: http.NewServeMux(),
	}
	a.initRoutes()
	return a, db
}

func TestApp(t *testing.T) {
	srv, _ := newMockApp(false)
	port, err := getFreePort()
	assert.NilError(t, err)
	os.Setenv("API_LISTEN", fmt.Sprintf("localhost:%d", port))
	srv.Start(context.Background())
	srv.Shutdown()
}
