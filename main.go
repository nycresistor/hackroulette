package main

import (
	"flag"
	"fmt"
	"io"
	"net"

	"github.com/golang/glog"
)

const listenAddr = ":4000"

var telnetOneChar []byte = []byte("\377\375\042\377\373\001")

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
		//case <-time.After(1 * time.Second):
		//	chat(Bot(), c)
	}
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}

func chat(a, b io.ReadWriteCloser) {
	defer a.Close()
	defer b.Close()
	//a.Write(telnetOneChar)
	//b.Write(telnetOneChar)
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		glog.Warningln(err)
	}
}

type socket struct {
	io.ReadCloser
	io.WriteCloser
}

func (s socket) Close() error {
	s.ReadCloser.Close()
	return s.WriteCloser.Close()
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		glog.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			glog.Fatal(err)
		}

		//r, w := io.Pipe()
		//go func() {
		//	_, err := io.Copy(io.MultiWriter(w, chain), c)
		//	w.CloseWithError(err)
		//	glog.Warningln(err)
		//}()
		//s := socket{r, c}
		go match(c)
	}
}
