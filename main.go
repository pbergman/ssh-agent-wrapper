package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"os/exec"
	"syscall"

	"github.com/ogier/pflag"
)

var (
	fileName string
	pattern  *regexp.Regexp = regexp.MustCompile(`(SSH_[A-Z|_]+)=([^;]+);`)
)

func getDefaultFile() string {
	if current, err := user.Current(); err != nil {
		panic(err)
	} else {
		return filepath.Join(current.HomeDir, ".ssh", "agent_env")
	}
}

func usage() string {
	return `
ssh-agent-wrapper is a simple wrapper for starting the ssh agent and
adding keys. It will check if the "ssh agent env" file exist. If not
then it will start the "ssh-agent" process and save the output to the
ssh agent env file. That can late be used to determine pid etc. of
the ssh-agent process.

this script can be used with startup of a shell by adding following to
~/.bashrc or ~/.profile or any other startup script.

	eval $(ssh-agent-wrapper)

usage: ssh-agent-wrapper [-f <SSH_AGENT_ENV_FILE>]

options:
  -f,--file    The location of the ssh agent env file [default: ~/.ssh/agent_env]
`
}

func init() {
	pflag.StringVarP(&fileName, "file", "f", getDefaultFile(), "")
	pflag.Usage = func() {
		fmt.Println(usage())
		os.Exit(0)
	}

}

func main() {
	pflag.Parse()

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("echo \"Failed to start/check ss-agent. %s\n\"", err)
		}
	}()

	if content, err := readFile(fileName); err != nil {
		if os.IsNotExist(err) {
			if err := startAgent(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		pid, _ := parseEnvContent(string(content))

		if process, err := os.FindProcess(pid); err != nil {
			panic(err)
		} else {
			if err = process.Signal(syscall.Signal(0)); err != nil {
				if err := startAgent(); err != nil {
					panic(err)
				}
			}
		}
	}
}

// readFile will try to open given file and read the content of that file.
func readFile(fileName string) ([]byte, error) {
	if file, err := os.Open(fileName); err != nil {
		return nil, err
	} else {

		defer file.Close()

		buf, content := make([]byte, 1024), make([]byte, 0)

		for {
			n, err := file.Read(buf)
			content = append(content, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return content, err
				}
			}

		}

		return content, nil
	}
}

// parseEnvContent will parse the content of the ssh env file,
// set the os envs and return the socket location and pid.
func parseEnvContent(content string) (pid int, socket string) {
	for _, env := range pattern.FindAllStringSubmatch(content, -1) {
		os.Setenv(env[1], env[2])
		fmt.Printf("export %s=%s;\n", env[1], env[2])
		switch env[1] {
		case "SSH_AUTH_SOCK":
			socket = env[2]
		case "SSH_AGENT_PID":
			i, _ := strconv.Atoi(env[2])
			pid = i
		}
	}
	return
}

// startAgent will start the ssh-agent and run the ssh-add command to
// add user keys.
func startAgent() error {
	if out, err := exec.Command("ssh-agent").Output(); err != nil {
		return err
	} else {
		if file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
			return err
		} else {
			if _, err := file.Write(out); err != nil {
				return err
			}
		}
		// set the envs
		parseEnvContent(string(out))
	}

	cmd := exec.Command("ssh-add")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
