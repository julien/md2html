package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/julien/md2html/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: md2html <input> <output>")
		os.Exit(1)
	}

	if err := run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	src, dst := args[1], args[1]

	if len(args) >= 3 {
		dst = args[2]
	}

	return internal.Run(src, dst)
}
