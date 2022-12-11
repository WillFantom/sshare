package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willfantom/sshare/internal/ui"
	"github.com/willfantom/sshare/keys"
	"github.com/willfantom/sshare/pastebin"
	"golang.design/x/clipboard"
)

var (
	pastebinURL      string = pastebin.DefaultPastebinURL
	pastebinDevToken string = ""
	pastebinUserkey  string = ""

	pasteExpiration string = "1H"
	pasteVisibility string = "0"

	sshAgentPath string = os.Getenv("SSH_AUTH_SOCK")
)

var (
	rootCmd = &cobra.Command{
		Use:   "sshare",
		Short: "Easily share links to your SSH public keys",
		Long:  `Share your public SSH keys found in your ssh agent via curl-able Pastebin links.`,
	}

	generateUserKey = &cobra.Command{
		Use:   "userkey [username] [password]",
		Short: "Generate your pastebin userkey",
		Long: `Generate a reusable Pastebin userkey for your
						account`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			pb, err := pastebin.NewPastebin(pastebinDevToken)
			if err != nil {
				ui.Errorln(fmt.Sprintf("failed to create pastebin client: %s", err.Error()), true)
			}
			userkey, err := pb.GenerateUserKey(args[0], args[1])
			if err != nil {
				ui.Errorln(fmt.Sprintf("failed to generate user key: %s", err.Error()), true)
			}
			ui.Successln(fmt.Sprintf("Pastebin Userkey: %s\n", userkey))
		},
	}

	share = &cobra.Command{
		Use:   "share",
		Short: "Share SSH public keys found in the your SSH agent",
		Long:  `Create a shareable Pastebin link to the SSH keys currently in your SSH agent.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if pastebin.ExpirationFromString(pasteExpiration).String() == "" {
				return fmt.Errorf("paste expiration time is not valid")
			}
			if pastebin.VisibilityFromString(pasteVisibility).String() == "" {
				return fmt.Errorf("paste visibility is not valid")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			sshAgent, err := keys.NewSSHAgent(sshAgentPath)
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			authorizedKeys, err := sshAgent.GetKeys()
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			selectedKeys, err := ui.SelectKey(authorizedKeys)
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			keyText := keys.CreateAuthorizedKeys(selectedKeys)
			ui.Infoln(fmt.Sprintf("Generated authorized_keys:\n%s", keyText))
			opts := make([]pastebin.PastebinOpt, 0)
			if pastebinUserkey != "" {
				opts = append(opts, pastebin.PastebinLoginOpt(pastebinUserkey))
			}
			pb, err := pastebin.NewPastebin(pastebinDevToken, opts...)
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			pasteKey, err := pb.Post(keyText, pastebin.ExpirationFromString(pasteExpiration), pastebin.VisibilityFromString(pasteVisibility))
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			rawURL, err := pb.GenerateRawPasteURL(pasteKey)
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			ui.Successln(fmt.Sprintf("Generated share link: %s", rawURL))
			if err := clipboard.Init(); err != nil {
				ui.Warnln("clipboard is not accessible")
			} else {
				clipboard.Write(clipboard.FmtText, []byte(rawURL))
				ui.Infoln("URL coppied to clipboard")
			}
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		ui.Errorln(err.Error(), true)
		os.Exit(1)
	}
}

func init() {
	clipboard.Init()
	rootCmd.AddCommand(generateUserKey)
	rootCmd.AddCommand(share)

	rootCmd.PersistentFlags().StringVar(&pastebinURL, "pb-url", pastebinURL, "URL for the target Pastebin instance")
	rootCmd.PersistentFlags().StringVarP(&pastebinDevToken, "pb-token", "t", "", "Your developer token for the Pastebin API")

	share.PersistentFlags().StringVarP(&pastebinUserkey, "userkey", "k", pastebinUserkey, "Your Pastebin user login key")
	share.PersistentFlags().StringVarP(&pasteExpiration, "expiration", "e", pasteExpiration, "Time till paste link expires")
	share.PersistentFlags().StringVarP(&pasteVisibility, "visibility", "v", pasteVisibility, "Visibility of the created Pastebin paste")
	share.PersistentFlags().StringVarP(&sshAgentPath, "agent", "a", sshAgentPath, "Path to the target SSH Agent socket")
}
