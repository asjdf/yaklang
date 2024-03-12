package yakcmds

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/urfave/cli"
	"github.com/yaklang/yaklang/common/cybertunnel"
	"github.com/yaklang/yaklang/common/cybertunnel/tpb"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/spec"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/antlr4yak/dap"
	"github.com/yaklang/yaklang/common/yak/antlr4yak/yakast"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"github.com/yaklang/yaklang/scannode"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	for _, util := range UtilsCommands {
		util.Category = "Utils"
	}

	for _, util := range DistributionCommands {
		util.Category = "Distribution Deploy"
	}
}

var UtilsCommands = []cli.Command{
	{
		Name: "gzip",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "f,file",
				Usage: "input file",
			},
			cli.BoolFlag{Name: "d,decode"},
			cli.StringFlag{Name: "o,output"},
		}, Action: func(c *cli.Context) error {
			f := c.String("file")
			if utils.GetFirstExistedFile(f) == "" {
				return utils.Errorf("non-existed: %v", f)
			}
			originFp, err := os.Open(f)
			if err != nil {
				return err
			}
			defer originFp.Close()

			if c.Bool("decode") {
				outFile := c.String("output")
				if outFile == "" {
					return utils.Error("decode need output not empty")
				}
				log.Infof("start to d-gzip to %v", outFile)
				targetFp, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR, 0o666)
				if err != nil {
					return err
				}
				defer targetFp.Close()
				r, err := gzip.NewReader(originFp)
				if err != nil {
					return err
				}
				defer r.Close()
				io.Copy(targetFp, r)
				log.Infof("finished")
				return nil
			}

			gf := f + ".gzip"
			fp, err := os.OpenFile(gf, os.O_CREATE|os.O_RDWR, 0o666)
			if err != nil {
				return err
			}
			defer fp.Close()
			gzipWriter := gzip.NewWriter(fp)
			io.Copy(gzipWriter, originFp)
			gzipWriter.Flush()
			gzipWriter.Close()
			return nil
		}},
	{Name: "hex", Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "f,file",
			Usage: "input file",
		},
		cli.StringFlag{
			Name:  "d,data",
			Usage: "input data",
		},
	}, Action: func(c *cli.Context) {
		if c.String("file") != "" {
			raw, err := ioutil.ReadFile(c.String("file"))
			if err != nil {
				log.Error(err)
				return
			}
			println(codec.EncodeToHex(raw))
		}

		if c.String("data") != "" {
			println(codec.EncodeToHex(c.String("data")))
		}
	}},

	// dap
	{
		Name:  "dap",
		Usage: "启动基于调试适配器协议(dap)的服务器以调试脚本",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "host", Usage: "debugger adapter listen host"},
			cli.IntFlag{Name: "port", Usage: "debugger adapter listen port"},
			cli.BoolFlag{Name: "debug", Usage: "debug mode"},
			cli.BoolFlag{Name: "version,v", Usage: "show dap version"},
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.Int("port")
			debug := c.Bool("debug")
			versionFlag := c.Bool("version")
			if versionFlag {
				fmt.Printf("Debugger Adapter version: %v\n", dap.DAVersion)
				return nil
			}

			// 设置日志级别
			if debug {
				log.SetLevel(log.DebugLevel)
			}

			server, stopChan, err := dap.StartDAPServer(host, port)
			if err != nil {
				return err
			}
			defer server.Stop()

			forceStop := make(chan struct{})
			select {
			case <-stopChan:
			case <-forceStop:
			}

			return nil
		},
	},

	// fmt
	{
		Name:  "fmt",
		Usage: "格式化代码",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "version,v", Usage: "show formatter version"},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				fmt.Printf("Formatter version: %v\n", yakast.FormatterVersion)
				return nil
			}
			args := c.Args()
			file := args[0]
			if file != "" {
				var err error
				absFile := file
				if !filepath.IsAbs(absFile) {
					absFile, err = filepath.Abs(absFile)
					if err != nil {
						return utils.Errorf("fetch abs file path failed: %s", err)
					}
				}
				raw, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				vt := yakast.NewYakCompiler()
				vt.Compiler(string(raw))
				fmt.Printf("%s", vt.GetFormattedCode())
			} else {
				return utils.Errorf("empty yak file")
			}
			return nil
		},
	},
}

var DistributionCommands = []cli.Command{
	scannode.DistYakCommand,
	{
		Name:   "mq",
		Usage:  "distributed by private amqp application protocol, execute yak via rabbitmq",
		Before: nil,
		After:  nil,
		Action: func(c *cli.Context) error {
			config := spec.LoadAMQPConfigFromCliContext(c)
			node, err := scannode.NewScanNode(c.String("id"), c.String("server-port"), config)
			if err != nil {
				return err
			}
			node.Run()
			return nil
		},
		Flags: spec.GetCliBasicConfig("scannode"),
	},
	{
		Name: "tunnel",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "server", Value: "cybertunnel.run:64333"},
			cli.IntFlag{Name: "local-port", Value: 53},
			cli.StringFlag{Name: "local-host", Value: "127.0.0.1"},
			cli.IntFlag{Name: "remote-port", Value: 53},
			cli.StringFlag{Name: "secret", Value: ""},
			cli.StringFlag{Name: "network,proto", Value: "tcp"},
		},
		Action: func(c *cli.Context) error {
			return cybertunnel.MirrorLocalPortToRemoteEx(
				c.String("network"),
				c.String("local-host"),
				c.Int("local-port"),
				c.Int("remote-port"),
				"test-cli",
				c.String("server"),
				c.String("secret"),
				context.Background(),
			)
		},
	},
	{
		Name:    "inspect-tuns",
		Usage:   "查看注册 tunnels 信息",
		Aliases: []string{"lst"},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "server", Usage: "远程 Yak Bridge X 服务器", Value: "127.0.0.1:64333"},
			cli.StringFlag{Name: "secret", Usage: "远程 Yak Bridge X 服务器密码"},
			cli.StringFlag{Name: "secondary-password,x", Usage: "远程 Yak Bridge X 服务器的二级密码，避免别人查看注册管道"},
			cli.StringFlag{Name: "id", Usage: "指定 ID 查看 Tunnel 信息与认证"},
		},
		Action: func(c *cli.Context) error {
			ctx, client, _, err := cybertunnel.GetClient(context.Background(), c.String("server"), c.String("secret"))
			if err != nil {
				return err
			}

			showTunnel := func(tun *tpb.RegisterTunnelMeta) {
				withAuth, _ := client.GetRegisteredTunnelDescriptionByID(ctx, &tpb.GetRegisteredTunnelDescriptionByIDRequest{
					Id:                tun.GetId(),
					SecondaryPassword: c.String("secondary-password"),
				})
				fmt.Printf(`Tunnel: %v
	addr: %v
	note:
%v
	auth: 
%v
-----------------

`, tun.GetId(), utils.HostPort(tun.GetConnectHost(), tun.GetConnectPort()), tun.GetVerbose(), string(withAuth.GetAuth()))
			}

			id := c.String("id")
			if id != "" {
				rsp, err := client.GetRegisteredTunnelDescriptionByID(ctx, &tpb.GetRegisteredTunnelDescriptionByIDRequest{
					Id:                id,
					SecondaryPassword: c.String("secondary-password"),
				})
				if err != nil {
					return err
				}

				if len(rsp.GetAuth()) <= 0 {
					return utils.Errorf("cannot generate auth bytes for tun: %s", id)
				}

				showTunnel(rsp.GetInfo())
				println(string(rsp.GetAuth()))
				return nil
			}

			resp, err := client.GetAllRegisteredTunnel(ctx, &tpb.GetAllRegisteredTunnelRequest{
				SecondaryPassword: c.String("secondary-password"),
			})
			if err != nil {
				return err
			}
			for i := 0; i < len(resp.GetTunnels()); i++ {
				showTunnel(resp.Tunnels[i])
			}

			return nil
		},
	},
}
