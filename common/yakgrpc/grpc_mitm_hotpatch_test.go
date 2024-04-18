package yakgrpc

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/lowhttp"
	"github.com/yaklang/yaklang/common/yak"
	"github.com/yaklang/yaklang/common/yakgrpc/yakit"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
)

func TestGRPCMUSTPASS_MITM_HotPatch_Drop(t *testing.T) {
	ctx, cancel := context.WithCancel(utils.TimeoutContextSeconds(5))
	defer cancel()

	mockHost, mockPort := utils.DebugMockHTTPHandlerFuncContext(ctx, func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})

	mitmPort := utils.GetRandomAvailableTCPPort()
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.MITM(ctx)
	if err != nil {
		t.Fatal(err)
	}
	stream.Send(&ypb.MITMRequest{
		Host: "127.0.0.1",
		Port: uint32(mitmPort),
	})
	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}
		if data.GetMessage().GetIsMessage() {
			msg := string(data.GetMessage().GetMessage())
			if !strings.Contains(msg, "starting mitm server") {
				continue
			}
			// load hot-patch mitm plugin
			stream.Send(&ypb.MITMRequest{
				SetYakScript:     true,
				YakScriptContent: `hijackHTTPResponseEx = func(isHttps, url, req, rsp, forward, drop) { drop() }`,
			})
		} else if data.GetCurrentHook && len(data.GetHooks()) > 0 {
			// send packet
			packet := `GET / HTTP/1.1
Host: ` + utils.HostPort(mockHost, mockPort) + `

`
			packetBytes := lowhttp.FixHTTPRequest([]byte(packet))
			_, err := yak.Execute(`
rsp, req, err = poc.HTTPEx(packet, poc.proxy(mitmProxy))
assert rsp.RawPacket.Contains("响应被用户丢弃")
`, map[string]any{
				"packet":    string(packetBytes),
				"mitmProxy": `http://` + utils.HostPort("127.0.0.1", mitmPort),
			})
			if err != nil {
				t.Fatal(err)
			}
			cancel()
		}
	}
}

func TestGRPCMUSTPASS_MITM_HotPatch_Dangerous_FuzzTag(t *testing.T) {
	ctx, cancel := context.WithCancel(utils.TimeoutContextSeconds(5))
	defer cancel()

	// create a temporary file to test
	token1 := utils.RandStringBytes(16)
	fileName, err := utils.SaveTempFile(token1, "fuzztag-test-file")
	if err != nil {
		panic(err)
	}
	fileName = strings.ReplaceAll(fileName, "\\", "\\\\")
	// create a codec script to test
	token2 := utils.RandStringBytes(16)
	scriptName, err := yakit.CreateTemporaryYakScript("codec", fmt.Sprintf(`
	handle = func(origin)  {
		return "%s"
	}`, token2))
	if err != nil {
		t.Fatal(err)
	}
	defer yakit.DeleteYakScriptByName(consts.GetGormProjectDatabase(), scriptName)

	mockHost, mockPort := utils.DebugMockHTTPHandlerFuncContext(ctx, func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})

	mitmPort := utils.GetRandomAvailableTCPPort()
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.MITM(ctx)
	if err != nil {
		t.Fatal(err)
	}
	stream.Send(&ypb.MITMRequest{
		Host: "127.0.0.1",
		Port: uint32(mitmPort),
	})
	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}
		if data.GetMessage().GetIsMessage() {
			msg := string(data.GetMessage().GetMessage())
			if !strings.Contains(msg, "starting mitm server") {
				continue
			}
			// load hot-patch mitm plugin
			stream.Send(&ypb.MITMRequest{
				SetYakScript: true,
				YakScriptContent: fmt.Sprintf(`hijackHTTPResponseEx = func(isHttps, url, req, rsp, forward, drop) {
	token1, token2 = "%s", "%s"
	file_fuzztag = fuzz.Strings("{{file(%s)}}")
	codec_fuzztag = fuzz.Strings("{{codec(%s)}}")
	if file_fuzztag[0].Contains(token1) || codec_fuzztag[0].Contains(token2) {
		forward(poc.ReplaceBody(rsp, "no", false))
	} else {
		forward(poc.ReplaceBody(rsp, "yes", false))
	}
}`, token1, token2, fileName, scriptName),
			})
		} else if data.GetCurrentHook && len(data.GetHooks()) > 0 {
			// send packet
			packet := `GET / HTTP/1.1
Host: ` + utils.HostPort(mockHost, mockPort) + `

`
			packetBytes := lowhttp.FixHTTPRequest([]byte(packet))
			_, err := yak.Execute(`
rsp, req = poc.HTTPEx(packet, poc.proxy(mitmProxy))~
assert rsp.RawPacket.Contains("yes")
`, map[string]any{
				"packet":    string(packetBytes),
				"mitmProxy": `http://` + utils.HostPort("127.0.0.1", mitmPort),
			})
			if err != nil {
				t.Fatal(err)
			}
			cancel()
		}
	}
}

func TestGRPCMUSTPASS_MITM_HotPatch_BeforeRequest_AfterRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(utils.TimeoutContextSeconds(100000000))
	defer cancel()

	token1 := utils.RandStringBytes(16)
	token2 := utils.RandStringBytes(16)
	token3 := utils.RandStringBytes(16)
	token4 := utils.RandStringBytes(16)

	mockHost, mockPort := utils.DebugMockHTTPEx(func(req []byte) []byte {
		if !bytes.Contains(req, []byte(token2)) {
			panic("token2 not found")
		}
		return []byte("HTTP/1.1 200 OK\r\n\r\nyes")
	})

	mitmPort := utils.GetRandomAvailableTCPPort()
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.MITM(ctx)
	if err != nil {
		t.Fatal(err)
	}
	stream.Send(&ypb.MITMRequest{
		Host: "127.0.0.1",
		Port: uint32(mitmPort),
	})

	hotPatchScript := fmt.Sprintf(`hijackHTTPRequest = func(isHttps, url, req, forward , drop) {
    req = poc.ReplaceHTTPPacketBody(req,"%s")
    forward(req)
}

beforeRequest = func(req){
    return poc.ReplaceHTTPPacketBody(req, "%s")
}

hijackHTTPResponse = func(isHttps, url, rsp, forward, drop) {
    rsp = poc.ReplaceHTTPPacketBody(rsp,"%s")
    forward(rsp)
}

afterRequest = func(rsp){
    return poc.ReplaceHTTPPacketBody(rsp, "%s")
}



`, token1, token2, token3, token4)

	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}
		if data.GetMessage().GetIsMessage() {
			msg := string(data.GetMessage().GetMessage())
			if !strings.Contains(msg, "starting mitm server") {
				continue
			}
			// load hot-patch mitm plugin
			stream.Send(&ypb.MITMRequest{
				SetYakScript:     true,
				YakScriptContent: hotPatchScript,
			})
			stream.Send(&ypb.MITMRequest{
				SetAutoForward:   true,
				AutoForwardValue: true,
			})
		} else if data.GetCurrentHook && len(data.GetHooks()) > 0 {
			// send packet
			go func() {
				packet := `GET / HTTP/1.1
Host: ` + utils.HostPort(mockHost, mockPort) + `

`
				packetBytes := lowhttp.FixHTTPRequest([]byte(packet))
				_, err := yak.Execute(`
rsp, req = poc.HTTPEx(packet, poc.proxy(mitmProxy))~
dump(rsp.RawPacket)
assert rsp.RawPacket.Contains("`+token4+`")
`, map[string]any{
					"packet":    string(packetBytes),
					"mitmProxy": `http://` + utils.HostPort("127.0.0.1", mitmPort),
				})
				if err != nil {
					t.Fatal(err)
				}
				cancel()
			}()
		} else if data.Request != nil && !data.ForResponse {
			// send packet
			if !bytes.Contains(data.Request, []byte(token1)) {
				t.Fatal("token1 not found")
			}
			stream.Send(&ypb.MITMRequest{
				HijackResponse: true,
				Forward:        true,
			})
		} else if data.Response != nil {
			// send packet
			if !bytes.Contains(data.Response, []byte(token3)) {
				t.Fatal("token3 not found")
			}

			stream.Send(&ypb.MITMRequest{
				Forward: true,
			})
		}

	}
}
