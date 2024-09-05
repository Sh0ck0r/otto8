package cli

import (
	"os"

	"github.com/gptscript-ai/cmd"
	"github.com/gptscript-ai/gptscript/pkg/mvl"
	"github.com/gptscript-ai/otto/pkg/api/client"
	"github.com/spf13/cobra"
)

type Otto struct {
	Debug  bool `usage:"Enable debug logging"`
	client *client.Client
}

func (a *Otto) PersistentPre(cmd *cobra.Command, args []string) error {
	if a.Debug {
		mvl.SetDebug()
	}
	return nil
}

func New() *cobra.Command {
	root := &Otto{
		client: &client.Client{
			BaseURL: "http://localhost:8080",
			Token:   os.Getenv("OTTO_TOKEN"),
		},
	}
	return cmd.Command(root,
		&Create{root: root},
		&Agents{root: root},
		&Update{root: root},
		&Delete{root: root},
		&Invoke{root: root},
		&Threads{root: root},
		cmd.Command(&Runs{root: root}, &Debug{root: root}),
		&Server{})
}

func (a *Otto) Run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
