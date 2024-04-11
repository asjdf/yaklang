package yakcmds

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/urfave/cli"
	"github.com/yaklang/yaklang/common/cybertunnel"
	"github.com/yaklang/yaklang/common/cybertunnel/tpb"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/mutate"
	"github.com/yaklang/yaklang/common/spec"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/filesys"
	"github.com/yaklang/yaklang/common/utils/omap"
	"github.com/yaklang/yaklang/common/yak/antlr4yak/dap"
	"github.com/yaklang/yaklang/common/yak/antlr4yak/yakast"
	"github.com/yaklang/yaklang/common/yak/yaklib"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"github.com/yaklang/yaklang/scannode"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var UtilsCommands = []*cli.Command{
	{
		Name:  "gzip",
		Usage: "gzip data or file",
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
	{
		Name: "hex",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "f,file",
				Usage: "input file",
			},
			cli.StringFlag{
				Name:  "d,data",
				Usage: "input data",
			},
		},
		Usage: "hex encode file or data to hex string",
		Action: func(c *cli.Context) {
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
	{
		Name:  "tag-stats",
		Usage: "Generate Tag Status(for Yakit)",
		Action: func(c *cli.Context) error {
			stats, err := yaklib.NewTagStat()
			if err != nil {
				return err
			}
			for _, v := range stats.All() {
				if v.Count <= 1 {
					continue
				}
				fmt.Printf("TAG:[%v]-%v\n", v.Name, v.Count)
			}
			return nil
		},
	},

	// dap
	{
		Name:  "dap",
		Usage: "Start a server based on the Debug Adapter Protocol (DAP) to debug scripts.",
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
		Usage: "Formatter for Yaklang Code",
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

	{
		Name:  "fuzz",
		Usage: "fuzztag short for fuzz tag, fuzz tag is a tool to generate fuzz string for fuzz testing",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "t,target",
				Usage: "Fuzztag Template, like: `{{int(1-10)}}`",
			},
		},
		Action: func(c *cli.Context) {
			for _, r := range mutate.MutateQuick(c.String("t")) {
				println(r)
			}
		},
	},

	// file tree size
	{
		Name:  "weight",
		Usage: "weight dir with depth",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "dir,d", Usage: "dir to weight"},
			cli.IntFlag{Name: "depth", Usage: "depth to weight", Value: 1},
			cli.BoolFlag{Name: "asc", Usage: "sort asc"},
			cli.StringFlag{Name: "blacklist,exclude", Usage: "ignore blacklist", Value: "*_test.go|.git*|*testdata*"},
			cli.StringFlag{Name: "show-exclude", Usage: "filter result", Value: "*.md|*.yak|*.DS_Store|*License|*.g4"},
			cli.IntFlag{Name: "show-min-size", Usage: "show min size", Value: 100000},
		},
		Action: func(c *cli.Context) error {
			m := omap.NewOrderedMap(map[string]int64{})
			err := filesys.Recursive(c.String("dir"), filesys.WithFileStat(func(pathname string, info os.FileInfo) error {
				if c.String("blacklist") != "" {
					if utils.MatchAnyOfGlob(pathname, utils.PrettifyListFromStringSplitEx(c.String("blacklist"), "|")...) {
						return nil
					}
				}
				log.Infof("path: %v, size: %v verbose: %v", pathname, info.Size(), utils.ByteSize(uint64(info.Size())))
				m.Set(pathname, info.Size())
				return nil
			}))
			if err != nil {
				return err
			}
			forest, err := utils.GeneratePathTrees(m.Keys()...)
			if err != nil {
				return err
			}

			results := omap.NewOrderedMap(make(map[string]int64))
			forest.Recursive(func(node2 *utils.PathNode) {
				if node2.GetDepth() > c.Int("depth") {
					return
				}
				count := int64(0)
				for _, child := range node2.AllChildren() {
					size, ok := m.Get(child.Path)
					if !ok {
						log.Warnf("path: %v, name: %v not found", child.Path, child.Name)
						continue
					}
					count += size
				}
				results.Set(node2.Path, count)
			})

			var desc []*sizeDescription
			results.ForEach(func(i string, v int64) bool {
				if c.String("show-exclude") != "" {
					if utils.MatchAnyOfGlob(i, utils.PrettifyListFromStringSplitEx(c.String("show-exclude"), "|")...) {
						return true
					}
				}
				desc = append(desc, &sizeDescription{Name: i, Size: uint64(v)})
				return true
			})

			sort.Slice(desc, func(i, j int) bool {
				if c.Bool("asc") {
					return desc[i].Size < desc[j].Size
				}
				return desc[i].Size > desc[j].Size
			})

			for _, i := range desc {
				fmt.Printf("[%6s]: %v\n", utils.ByteSize(i.Size), i.Name)
			}
			return nil
		},
	},
}

var DistributionCommands = []*cli.Command{
	&scannode.DistYakCommand,
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
		Name:  "tunnel",
		Usage: "Create Tunnel For CyberTunnel Service",
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
		Usage:   "Inspect Registered Tunnels",
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

type sizeDescription struct {
	Name string
	Size uint64
}
