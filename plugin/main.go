package plugin

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-honeycomb/internal"
)

// Serve serves the plugin
func Serve() {
	p := internal.NewHoneycombPlugin()
	plugin.ServePlugin(p)
}
