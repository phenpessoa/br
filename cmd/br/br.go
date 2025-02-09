package main

import (
	"fmt"
	"os"

	"github.com/phenpessoa/br/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run br: %s\n", err)
		os.Exit(1)
	}
}
