package register_bot

import (
	"fmt"
	"github.com/xlab/treeprint"
)

func (b *Bot) PrintRegisteredCommands() {
	tree := treeprint.New()

	botNode := tree.AddBranch(fmt.Sprintf("Bot with token: %s", b.Token))
	commandsNode := botNode.AddBranch("Registered Commands")
	buttonsNode := botNode.AddBranch("Registered Buttons")
	basicAuthNode := botNode.AddBranch("Registered BasicAuth commands")
	privateCommandsNode := botNode.AddBranch("Registered privateCommands commands")
	allowedUsernamesNode := botNode.AddBranch("Registered allowedUsernames")

	for command := range b.commandHandlers {
		commandsNode.AddNode(fmt.Sprintf("- %s", command))
	}

	for button := range b.buttons {
		buttonsNode.AddNode(fmt.Sprintf("- %s", button))
	}

	for basic := range b.registerBasicAuth {
		basicAuthNode.AddNode(fmt.Sprintf("- %s", basic))
	}

	for private := range b.privateCommands {
		privateCommandsNode.AddNode(fmt.Sprintf("- %s", private))
	}

	for user := range b.allowedUsernames {
		allowedUsernamesNode.AddNode(fmt.Sprintf("- %s", user))
	}

	fmt.Println(tree.String())
}
