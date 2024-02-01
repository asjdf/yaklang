package lowhttp

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/utils"
	"strconv"
	"strings"
	"testing"
)

func TestExtractURLFromHTTPRequest(t *testing.T) {
	const packet = `GET / HTTP/1.1
Host: asdfasd:123
 Cookie: 123
  d: 1
`
	u, err := ExtractURLFromHTTPRequestRaw([]byte(packet), false)
	if err != nil {
		panic(err)
	}
	spew.Dump(u.String())
	var a = FixHTTPRequest([]byte(packet))
	if !strings.Contains(string(a), "\r\n Cookie: 123\r\n  d: 1\r\n") {
		panic(1)
	}
}

func TestParseStringToHttpRequest2(t *testing.T) {
	req, err := ParseStringToHttpRequest(`
GET / HTTP/1.1
Host: www.baidu.com

teadfasdfasd
`)
	if err != nil {
		t.FailNow()
		return
	}
	_ = req
}

func TestSplitHTTPHeader(t *testing.T) {
	key, value := SplitHTTPHeader("abc")
	if !(key == "abc" && value == "") {
		panic("111")
	}

	key, value = SplitHTTPHeader("abc:111")
	if !(key == "abc" && value == "111") {
		panic("111")
	}

	key, value = SplitHTTPHeader("abc: 111")
	if !(key == "abc" && value == "111") {
		panic("111")
	}

	key, value = SplitHTTPHeader("abc: 111\r\n")
	if !(key == "abc" && value == "111") {
		panic("111")
	}

	key, value = SplitHTTPHeader("Abc: 111\r\n")
	if !(key == "Abc" && value == "111") {
		panic("111")
	}

	key, value = SplitHTTPHeader("Abc: 1::11\r\n")
	if !(key == "Abc" && value == "1::11") {
		panic("111")
	}
}

func TestParseStringToHttpRequest(t *testing.T) {
	test := assert.New(t)

	req, err := ParseStringToHttpRequest(`
GET / HTTP/1.1
Host: www.baidu.com
Connection: close
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36

`)
	if err != nil {
		test.FailNow(err.Error())
	}

	u, err := ExtractURLFromHTTPRequest(req, true)
	if err != nil {
		test.FailNow(err.Error())
		return
	}
	_ = u
}

func TestConvertContentToChunked(t *testing.T) {
	raws := fixInvalidHTTPHeaders([]byte(`
GET / HTTP/1.1
Host: www.baidu.com
Content-Length: 12

123123123123
`))
	println(string(raws))
}

func TestGetRedirectFromHTTPResponse(t *testing.T) {
	target := GetRedirectFromHTTPResponse([]byte(`HTTP/1.1 300 ...
Location: /target`), false)
	println(target)
	if target != "/target" {
		t.FailNow()
		return
	}
}

func TestRemoveZeroContentLengthHTTPHeader(t *testing.T) {
	target := RemoveZeroContentLengthHTTPHeader([]byte(`GET / HTTP/1.1
Host: www.baidu.com
Content-Length: 0

`))
	println(string(target))
	println(strconv.Quote(string(target)))
}

func TestConnectExtractedUrl(t *testing.T) {
	var testcases = []struct {
		url    string
		packet string
	}{
		{url: "http://baidu.com/a?c=1", packet: `POST a?c=1 HTTP/1.1
Host: baidu.com`},
		{url: "http://baidu.com:11/a?c=1", packet: `POST a?c=1 HTTP/1.1
Host: baidu.com:11`},
		{url: "http://baidu.com:11/./a?c=1", packet: "POST /./a?c=1 HTTP/1.1\r\nHost: baidu.com:11\n"},
		{url: "http://baidu.com:11//a?c=1", packet: "POST //a?c=1 HTTP/1.1\r\nHost: baidu.com:11\n"},

		{url: "http://baidu.com:11/?c=1", packet: `POST ?c=1 HTTP/1.1
Host: baidu.com:11`},
		{url: "http://baidu.com:11/", packet: `POST http://baidu.com:11/ HTTP/1.1`},
		{url: "http://baidu.com:11/", packet: `GET http://baidu.com:11/ HTTP/1.1`},
		{url: "http://baidu.com:11/c", packet: `GET http://baidu.com:11/c HTTP/1.1`},
		{url: "http://baidu.com:11", packet: `GET http://baidu.com:11 HTTP/1.1`},
		{url: "http://baidu.com:11/a?c=1", packet: `GET http://baidu.com:11/a?c=1 HTTP/1.1`},
		{url: "http://baidu.com:11", packet: `CONNECT http://baidu.com:11 HTTP/1.1`},
		{url: "http://baidu.com:11", packet: `CONNECT http://baidu.com:11 HTTP/1.1
Host: www.example.com`},
		{url: "http://baidu.com", packet: `CONNECT http://baidu.com`},
		{url: "https://baidu.com", packet: `CONNECT https://baidu.com`},
		{url: "https://baidu.com/ab", packet: `CONNECT https://baidu.com/ab`},
		{url: "https://baidu.com:1/ab?a=1", packet: `CONNECT https://baidu.com:1/ab?a=1`},
	}

	for _, testcase := range testcases {
		req, err := ParseStringToHttpRequest(testcase.packet)
		if err != nil {
			t.FailNow()
			return
		}
		u, err := ExtractURLFromHTTPRequest(req, false)
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
			return
		}
		if u.String() != testcase.url {
			fmt.Println(string(testcase.packet))
			t.Fatalf("url not match: %s != %s", u.String(), testcase.url)
		} else {
			t.Logf("url match: %v == %v", u.String(), testcase.url)
		}
	}
}

func TestParseResponseLine(t *testing.T) {
	testcases := []struct {
		line          string
		proto, status string
		code          int
	}{
		{
			line:   "HTTP/1.1 200 OK",
			proto:  "HTTP/1.1",
			code:   200,
			status: "OK",
		},
		{
			line:   "HTTP/1.1 200",
			proto:  "HTTP/1.1",
			code:   200,
			status: "",
		},
		{
			line:   "HTTP/1.1 301 Moved Permanently",
			proto:  "HTTP/1.1",
			code:   301,
			status: "Moved Permanently",
		},
	}

	for _, testcase := range testcases {
		proto, code, status, _ := utils.ParseHTTPResponseLine(testcase.line)
		if proto != testcase.proto {
			t.Fatalf("utils.ParseHTTPResponseLine error: %s(got) != %s(want)", proto, testcase.proto)
		}
		if code != testcase.code {
			t.Fatalf("utils.ParseHTTPResponseLine error: %d(got) != %d(want)", code, testcase.code)
		}
		if status != testcase.status {
			t.Fatalf("utils.ParseHTTPResponseLine error: %s(got) != %s(want)", status, testcase.status)
		}

	}
}

func TestGZIP_IN_REQUEST(t *testing.T) {
	raw, _ := utils.GzipCompress("abc")
	var packetResult []byte
	packetResult = ReplaceHTTPPacketBody([]byte(`POST / HTTP/1.1
Host: www.baidu.com
Content-Encoding: gzip

`), raw, false)
	packetResult = FixHTTPRequest(packetResult)
	fmt.Println(string(packetResult))

	if strings.Contains(string(packetResult), "abc") {
		panic("gzip must in request error")
	}

	var result = DeletePacketEncoding(packetResult)
	fmt.Println(string(result))
	if !strings.Contains(string(result), "abc") || strings.Contains(string(result), `-Encoding: gzip`) {
		panic("clear in request error")
	}
}

func TestSplitHTTPPacket_BlankCharacterBody(t *testing.T) {
	_, body := SplitHTTPHeadersAndBodyFromPacketEx([]byte("HTTP/1.1 200 OK\r\nServer: Apache-Coyote/1.1\r\nSet-Cookie: OASESSIONID=abc; Path=/defaultroot\r\nContent-Type: text/html;charset=UTF-8\r\nContent-Length: 6\r\nDate: Wed, 31 Jan 2024 07:44:24 GMT\r\n\r\n\r\n\r\n\r\n"), nil)
	if bytes.Compare(body, []byte("\r\n\r\n\r\n")) != 0 {
		t.Fatal("split body error ")
	}
}
