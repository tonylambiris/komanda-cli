package command

import (
	"github.com/jroimartin/gocui"
	"github.com/mephux/komanda-cli/komanda/client"
	"github.com/mephux/komanda-cli/komanda/ui"
)

// JoinCmd strcut
type JoinCmd struct {
	*MetadataTmpl
}

// Metadata for join command
func (e *JoinCmd) Metadata() CommandMetadata {
	return e
}

// Exec join command
func (e *JoinCmd) Exec(args []string) error {
	Server.Exec(client.StatusChannel, func(c *client.Channel, g *gocui.Gui, v *gocui.View, s *client.Server) error {

		if !s.Client.Connected() {
			client.StatusMessage(v, "Not connected")
			return nil
		}

		if len(args) >= 2 && len(args[1]) > 0 {

			if channel, _, has := Server.HasChannel(args[1]); has {
				CurrentChannel = args[1]
				s.CurrentChannel = args[1]

				Server.Gui.SetViewOnTop(Server.CurrentChannel)

				if _, err := g.SetCurrentView(channel.Name); err != nil {
					return err
				}

				channel.Unread = false
				channel.Highlight = false

				if _, err := g.SetCurrentView("input"); err != nil {
					return err
				}

				ui.UpdateMenuView(g)

				return nil
			}

			s.Client.Join(args[1])
			CurrentChannel = args[1]
			s.CurrentChannel = args[1]

			return s.NewChannel(args[1], false)
		}

		return nil
	})

	return nil
}

func joinCmd() Command {
	return &JoinCmd{
		MetadataTmpl: &MetadataTmpl{
			name: "join",
			args: "<channel>",
			aliases: []string{
				"j",
			},
			description: "join irc channel",
		},
	}
}
