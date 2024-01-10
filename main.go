package main

import (
	"github.com/mach-composer/mach-composer-plugin-honeycomb/internal"
	"github.com/mach-composer/mach-composer-plugin-sdk/v2/plugin"
)

func main() {
	p := internal.NewHoneycombPlugin()
	plugin.ServePlugin(p)
}
