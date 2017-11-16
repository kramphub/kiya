package main

import "flag"

var (
	oForce          = flag.Bool("f", false, "overwrite existing secret values")
	oVerbose        = flag.Bool("v", false, "set verbose mode")
	oConfigFilename = flag.String("c", "", "location of the configuration file. If empty then expect .kiya in $HOME.")
	oAuthLocation   = flag.String("a", "", "location of the JSON key credentials file. If empty then use the Google Application Defaults.")
)
