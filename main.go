package main

import (
	"github.com/mach-composer/mach-composer-plugin-honeycomb/internal"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
)

func main() {
	p := internal.NewHoneycombPlugin()
	plugin.ServePlugin(p)
}
