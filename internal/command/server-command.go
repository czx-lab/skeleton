package command

import (
	"fmt"
	"log"
	"skeleton/internal/server"
	"skeleton/internal/variable"
	"skeleton/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type ServerCommand struct {
	Port string
	Mode string
}

func NewServerCommand() *ServerCommand {
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")
	if len(port) == 0 {
		port = ":8080"
	}
	if len(mode) == 0 {
		mode = gin.DebugMode
	}
	return &ServerCommand{
		Port: port,
		Mode: mode,
	}
}

func (s *ServerCommand) WithPort(port string) *ServerCommand {
	s.Port = port
	return s
}

func (s *ServerCommand) WithMode(mode string) *ServerCommand {
	s.Mode = mode
	return s
}

// Command implements Interface.
func (s *ServerCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "server:http",
		Short: "http server",
		Long: `Instructions:
  如果命令行参数-p[--port]或-m[--mode]有值，则会有优先取命令行的值，反之则取传入的值
  如果都没有默认值为 port=":8080"、mode="debug"`,
		Run: func(cmd *cobra.Command, args []string) {
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatalf("err: %s\n", err)
			}
			if len(port) > 0 {
				s.Port = port
			}
			mode, err := cmd.Flags().GetString("mode")
			if err != nil {
				log.Fatalf("err: %s\n", err)
			}
			if len(mode) > 0 {
				s.Mode = mode
			}
			httpStart(s.Port, s.Mode)
		},
	}
}

func httpStart(port, mode string) {
	http := server.New(
		server.WithPort(port),
		server.WithMode(mode),
		server.WithLogger(variable.Log),
	)
	fmt.Printf("Starting server at %s, Service model [%s]...\n", port, mode)
	http.SetRouters(router.New(http)).Run()
}

// Flags implements Interface.
func (s *ServerCommand) Flags(root *cobra.Command) {
	root.Flags().StringP("port", "p", "", "http服务端口")
	root.Flags().StringP("mode", "m", "", `http服务模式，取值范围如下：
	- debug
	- release
`)
}

var _ (Interface) = (*ServerCommand)(nil)
