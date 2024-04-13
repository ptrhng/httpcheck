package main

import (
	"context"
	"os"
)

func main() {
	cmd := NewCommand()
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
