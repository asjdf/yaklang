package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"testing"
)

func TestHTTPRequestDumper_BodyIsLager(t *testing.T) {
	packet := `GET / HTTP/1.1` + CRLF +
		`Host: www.example.com` + CRLF +
		`Content-Length: 3` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBytes([]byte(packet))
	if err != nil {
		panic(err)
	}
	if req.ContentLength == 3 {
		t.Fatal("ContentLength should be 10")
	}
}

func TestHTTPRequestDumper_BodyIsSmall(t *testing.T) {
	packet := `GET / HTTP/1.1` + CRLF +
		`Host: www.example.com` + CRLF +
		`Content-Length: 13` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBytes([]byte(packet))
	if err != nil {
		panic(err)
	}
	if req.ContentLength == 13 {
		t.Fatal("ContentLength should be 10")
	}
}

func TestHTTPRequestDumper_Stream_BodyIsLager(t *testing.T) {
	packet := `GET / HTTP/1.1` + CRLF +
		`Host: www.example.com` + CRLF +
		`Content-Length: 3` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBufioReader(bufio.NewReader(bytes.NewBufferString(packet)))
	if err != nil {
		panic(err)
	}
	if req.ContentLength != 3 {
		t.Fatal("ContentLength should be 3")
	}
}

func TestHTTPRequestDumper_C1(t *testing.T) {
	packet := `GET https://example.com/bac HTTP/1.1` + CRLF +
		`Host: www.example.com` + CRLF +
		`Content-Length: 3` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBytes([]byte(packet))
	if err != nil {
		panic(err)
	}
	raw, _ := DumpHTTPRequest(req, true)
	fmt.Println(string(raw))
	if !bytes.HasPrefix(raw, []byte(`GET /bac HTTP/1.1`)) {
		t.Fatal("should be GET /bac HTTP/1.1")
	}
}

func TestHTTPRequestDumper_Stream_BodyIsSmall(t *testing.T) {
	packet := `GET / HTTP/1.1` + CRLF +
		`Host: www.example.com` + CRLF +
		`Content-Length: 13` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBufioReader(bufio.NewReader(bytes.NewBufferString(packet)))
	if err != nil {
		panic(err)
	}
	if req.ContentLength != 13 {
		t.Fatal("ContentLength should be 13")
	}
	raw, _ := io.ReadAll(req.Body)
	if string(raw) != "abccccddef   " && len(string(raw)) != 13 {
		spew.Dump(raw)
		t.Fatal("body should be abcccddef[SP][SP][SP]")
	}
}

func TestHTTPResponseDumper_BodyIsLager(t *testing.T) {
	packet := `HTTP/1.1 200 OK` + CRLF +
		`Server: Test-ABC` + CRLF +
		`Content-Length: 3` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBytes([]byte(packet))
	if err != nil {
		panic(err)
	}
	if req.ContentLength == 3 {
		t.Fatal("ContentLength should be 10")
	}
}

func TestHTTPResponseDumper_BodyIsSmall(t *testing.T) {
	packet := `HTTP/1.1 200 OK` + CRLF +
		`Server: Test-ABC` + CRLF +
		`Content-Length: 13` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBytes([]byte(packet))
	if err != nil {
		panic(err)
	}
	if req.ContentLength == 13 {
		t.Fatal("ContentLength should be 10")
	}
}

func TestHTTPResponseDumper_Stream_BodyIsLager(t *testing.T) {
	packet := `HTTP/1.1 200 OK` + CRLF +
		`Server: Test-ABC` + CRLF +
		`Content-Length: 3` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBufioReader(bufio.NewReader(bytes.NewBufferString(packet)))
	if err != nil {
		panic(err)
	}
	if req.ContentLength != 3 {
		t.Fatal("ContentLength should be 3")
	}
}

func TestHTTPResponseDumper_Stream_BodyIsSmall(t *testing.T) {
	packet := `HTTP/1.1 200 OK` + CRLF +
		`Server: Test-ABC` + CRLF +
		`Content-Length: 13` + CRLF + CRLF + "abccccddef"
	req, err := ReadHTTPRequestFromBufioReader(bufio.NewReader(bytes.NewBufferString(packet)))
	if err != nil {
		panic(err)
	}
	if req.ContentLength != 13 {
		t.Fatal("ContentLength should be 13")
	}
	raw, _ := io.ReadAll(req.Body)
	if string(raw) != "abccccddef   " && len(string(raw)) != 13 {
		spew.Dump(raw)
		t.Fatal("body should be abcccddef[SP][SP][SP]")
	}
}
