package yakgrpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aymanbagabas/go-pty"
	"github.com/google/shlex"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
	"golang.org/x/term"
)

func getShellCommand() (string, string, error) {
	switch os := runtime.GOOS; os {
	case "windows":
		return "C:\\\\Windows\\\\system32\\\\cmd.exe /k", "\r\n", nil
	case "linux", "darwin":
		var (
			finErr error
			shell  string
		)
		for _, shellName := range []string{"bash", "sh"} {
			cmd := exec.Command("which", shellName)
			shellBytes, err := cmd.CombinedOutput()
			if err == nil {
				shell = strings.TrimSpace(string(shellBytes))
				break
			} else {
				finErr = err
			}
		}

		if shell == "" && finErr != nil {
			return "", "", utils.Errorf("failed to find shell: %s", finErr)
		}
		return shell, "\n", nil
	default:
		return "", "", utils.Errorf("unsupported os: %s", os)
	}
}

func (s *Server) YaklangTerminal(inputStream ypb.Yak_YaklangTerminalServer) error {
	ctx, cancel := context.WithCancel(inputStream.Context())
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			cancel()
			return
		}
	}()

	firstInput, err := inputStream.Recv()
	if err != nil {
		return err
	}

	// exec
	shell, eol, err := getShellCommand()
	if err != nil {
		return err
	}

	streamerRWC := &OpenPortServerStreamerHelperRWC{
		stream: inputStream,
	}
	commands, _ := shlex.Split(shell)

	ptmx, err := pty.New()
	if err != nil {
		// fallback
		cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		if firstInput.GetPath() != "" {
			cmd.Dir = firstInput.GetPath()
		} else {
			path, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			cmd.Dir = path
		}
		cmd.Start()

		terminal := term.NewTerminal(streamerRWC, "")
		go io.Copy(terminal, stdout)
		go io.Copy(terminal, stderr)
		for {
			line, err := terminal.ReadLine()
			if errors.Is(err, io.EOF) {
				continue
			}
			if err != nil {
				return err
			}
			if line == "" {
				continue
			}
			stdin.Write([]byte(line + eol))
		}
	} else {
		defer ptmx.Close()

		go io.Copy(ptmx, streamerRWC) // stdin
		go func() {
			if runtime.GOOS == "windows" {
				// split the first output
				buf := make([]byte, 4096)
				n, err := ptmx.Read(buf)
				if err != nil {
					return
				}
				buf = buf[:n]
				_, after, ok := bytes.Cut(buf, []byte{0x1b, 0x5b, 0x48})
				if ok {
					buf = after
					before, _, ok := bytes.Cut(buf, []byte{0x1b, 0x5d, 0x30})
					if ok {
						buf = before
					}
				}
				streamerRWC.Write(buf)
			}

			io.Copy(streamerRWC, ptmx) // stdout
		}()

		defer func() {
			inputStream.Send(&ypb.Output{
				Control: true,
				Closed:  true,
			})
		}()

		cmd := ptmx.CommandContext(ctx, commands[0], commands[1:]...)
		return cmd.Run()
	}
}
