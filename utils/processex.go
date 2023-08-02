package utils

import (
	"fmt"
	"os"

	"github.com/biter777/processex"
)

func Processex() {
	processName := "mysql"
	process, _, err := processex.FindByName(processName)
	if err == processex.ErrNotFound {
		fmt.Printf("Process %v not running", processName)
		os.Exit(0)
	}
	if err != nil {
		fmt.Printf("Process %v find error: %v", processName, err)
		os.Exit(1)
	}
	fmt.Printf("Process %v PID: %v", processName, process)
}
