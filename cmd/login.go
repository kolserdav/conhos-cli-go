package cmd

import (
	"fmt"

	"conhos-cli/pkg/connectors"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Login via browser",
	Run: func(cmd *cobra.Command, args []string) {
		// Create WS connection
		ws := connectors.NewWS(nil)
		if err := ws.Connect(); err != nil {
			fmt.Println("Unable to connect:", err)
			return
		}
		defer ws.Close()

		// Start listening for messages
		go func() {
			for {
				_, message, err := ws.Conn.ReadMessage()
				if err != nil {
					break
				}

				// TODO: Parse message type
				messageType := "UNKNOWN"

				// Process message types
				switch messageType {
				case "AUTH_START":
					fmt.Println("Received AUTH_START")
					// TODO: handle auth start
				case "AUTH_COMPLETE":
					fmt.Println("Received AUTH_COMPLETE")
					// TODO: handle auth complete
				default:
					fmt.Printf("Received unknown message type: %s\n", messageType)
				}
			}
		}()

		// Main command logic
		fmt.Println("Login command started. Waiting for messages...")

		// Block until interrupted
		select {}
	},
}

func init() {
	loginCmd.Flags().BoolP("crypt", "c", false, "encrypt session token with password")
	loginCmd.Flags().BoolP("remove", "r", false, "remove session token from this device")
}
