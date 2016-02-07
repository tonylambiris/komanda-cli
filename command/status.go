package command

import (
	"github.com/jroimartin/gocui"
	"github.com/mephux/komanda-cli/client"
)

type StatusCmd struct {
	*MetadataTmpl
}

func (e *StatusCmd) Metadata() CommandMetadata {
	return e
}

func (e *StatusCmd) Exec(args []string) error {

	Server.Exec(client.StatusChannel, func(v *gocui.View, s *client.Server) error {
		if s.Client != nil && s.Client.Connected() {
			client.StatusMessage(v, "Status: Connected.")
		} else {
			client.StatusMessage(v, "Status: Disconnected.")
		}

		return nil
	})

	return nil
}

func statusCmd() Command {
	return &StatusCmd{
		MetadataTmpl: &MetadataTmpl{
			name:        "status",
			description: "status command",
		},
	}
}