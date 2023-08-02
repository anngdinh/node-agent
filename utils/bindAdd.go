package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func ProcessExxxx() {
	// Command to get the bind address from the MySQL configuration file
	cmd := exec.Command("grep", "-oP", "(?i)bind-address\\s*=\\s*\\K\\S+", "/etc/mysql/my.cnf")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert the output to a string and remove any leading/trailing spaces or newlines
	bindAddress := strings.TrimSpace(string(output))

	if bindAddress == "" {
		fmt.Println("MySQL bind address is not set in the configuration file.")
	} else {
		fmt.Printf("MySQL is bound to address: %s\n", bindAddress)
	}
}
