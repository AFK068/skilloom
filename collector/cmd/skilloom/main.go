package main

import "github.com/alecthomas/kong"

func main() {
	var cli CLI
	ctx := kong.Parse(&cli, kong.Name("skilloom"))
	ctx.FatalIfErrorf(ctx.Run())
}
