package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"tapo/cloud/tapo"

	"github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	RunE:    runLogin,
	PreRun:  preRunLogin,
	PostRun: postRunLogin,
}

func init() {

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")

	rootCmd.AddCommand(loginCmd)
}

func preRunLogin(cmd *cobra.Command, args []string) {

	termID := viper.GetString("TERM_ID")

	if termID == "" {

		id := uuid.New()
		termID = hex.EncodeToString(id[:])

		viper.Set("TERM_ID", termID)
	}

	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	if username == "" {

		fmt.Print("Username:")

		username = prompt.Input(" ", func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		})

		cmd.Flags().Set("username", username)
	}

	if password == "" {

		fmt.Print("Password:")

		password = prompt.Input(" ", func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		})

		cmd.Flags().Set("password", password)
	}
}

func postRunLogin(cmd *cobra.Command, args []string) {

	viper.WriteConfig()

	fmt.Printf("TERM_ID=%s\n", viper.GetString("TERM_ID"))
	fmt.Printf("TOKEN=%s\n", viper.GetString("TOKEN"))
}

func runLogin(cmd *cobra.Command, args []string) error {

	appName := "TP-Link_Tapo_Android"
	appVersion := "3.0.536"

	termID := viper.GetString("TERM_ID")
	termName := viper.GetString("TERM_NAME")

	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	client := tapo.NewTpLinkCloudClient("https://n-wap-gw.tplinkcloud.com")

	params := &tapo.GenericTapoParams{
		AppName: appName,
		AppVer:  appVersion,

		NetType: "wifi",
		Locale:  viper.GetString("LOCALE"),

		OSPF:  viper.GetString("OSPF"),
		Brand: viper.GetString("BRAND"),
		Model: viper.GetString("MODEL"),

		TermID:   termID,
		TermName: termName,
		TermMeta: "1",
	}

	resp, err := client.AccountLogin(params, &tapo.AccountLoginRequest{

		AppType:    appName,
		AppVersion: appVersion,

		Platform: params.OSPF,

		TerminalMeta: params.TermMeta,
		TerminalName: params.TermName,
		TerminalUUID: params.TermID,

		CloudUserName: username,
		CloudPassword: password,
	})

	if err != nil {
		return err
	}

	if resp.Result.ErrorCode == "-20677" {

		if slices.Contains(resp.Result.SupportedMFATypes, tapo.MFATypePush) {

			_, err := client.GetPushVC4TerminalMFA(params, &tapo.GetPushVC4TerminalMFARequest{
				AppType:       appName,
				CloudUserName: username,
				CloudPassword: password,
				TerminalUUID:  termID,
			})

			if err != nil {
				return err
			}

		} else if slices.Contains(resp.Result.SupportedMFATypes, tapo.MFATypeEmail) {

			_, err := client.GetEmailVC4TerminalMFA(params, &tapo.GetEmailVC4TerminalMFARequest{
				AppType:       appName,
				CloudUserName: username,
				CloudPassword: password,
				TerminalUUID:  termID,
			})

			if err != nil {
				return err
			}

		} else {
			return errors.New("MFA not supported")
		}

		fmt.Print("MFA Code:")

		code := prompt.Input(" ", func(d prompt.Document) []prompt.Suggest {
			return []prompt.Suggest{}
		})

		resp, err := client.CheckMFACodeAndLogin(params, &tapo.CheckMFACodeAndLoginRequest{

			MFAProcessId: resp.Result.MFAProcessId,
			MFAType:      tapo.MFATypeEmail,

			AppType:             appName,
			CloudUserName:       username,
			Code:                code,
			TerminalBindEnabled: true,
		})

		if err != nil {
			return err
		}

		viper.Set("TOKEN", resp.Result.Token)

	} else {
		viper.Set("TOKEN", resp.Result.Token)
	}

	return nil
}
