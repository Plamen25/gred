package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type args struct {
	c, n        int
	net, addr   string
	wg          *sync.WaitGroup
	start, stop chan struct{}
	files       []jsonFile
}

func launchWorkers(a args) ([]*worker, error) {
	ws := make([]*worker, a.c)
	ix := 0

	// Launch all workers, assign files
	for i := 0; i < a.c; i++ {
		// Create new worker, which will try to establish connection to
		// the server.
		w, err := newWorker(i, a.net, a.addr)
		if err != nil {
			return nil, err
		}

		// Launch the worker, it will wait on the start chan signal
		w.files = append(w.files, a.files[ix])
		go w.work(a.n, a.wg, a.start, a.stop)
		ws[i] = w

		// Move on to next command file
		ix++
		if ix >= len(a.files) {
			ix = 0
		}
	}

	// If there are more files than workers, assign remaining files in order
	if len(a.files) > a.c {
		ix = 0
		for i := 0; i < len(a.files)-a.c; i++ {
			ws[ix].files = append(ws[ix].files, a.files[i])
			ix++
			if ix >= len(ws) {
				ix = 0
			}
		}
	}
	return ws, nil
}

type worker struct {
	id           int
	files        []jsonFile
	replies      *bytes.Buffer
	conn         redis.Conn
	ncmds, nerrs int
}

func newWorker(id int, net string, addr string) (*worker, error) {
	conn, err := redis.Dial(net, addr)
	if err != nil {
		return nil, err
	}

	return &worker{id: id, conn: conn, replies: bytes.NewBuffer(nil)}, nil
}

func (w *worker) work(n int, wg *sync.WaitGroup, start, stop chan struct{}) {
	// Wait for start signal
	<-start

loop:
	for i := 0; ; i++ {
		if n > 0 {
			if i >= n {
				break loop
			}
		}
		for _, f := range w.files {
			nc, ne := f.exec(w.id, w.conn, w.replies, stop)
			w.ncmds += nc
			w.nerrs += ne

			select {
			case <-stop:
				break loop
			default:
			}
		}
	}

	wg.Done()
}

func processResults(ws []*worker, out io.Writer) ([]byte, int, int) {
	ncmd, nerr := 0, 0

	var mw io.Writer
	h := sha256.New()
	mw = h
	if out != nil {
		mw = io.MultiWriter(h, out)
	}
	for _, w := range ws {
		io.Copy(mw, w.replies)
		ncmd += w.ncmds
		nerr += w.nerrs
	}
	return h.Sum(nil), ncmd, nerr
}
