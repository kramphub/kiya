package main

import "flag"

var (
	oForce          = flag.Bool("f", false, "overwrite existing secret values")
	oVerbose        = flag.Bool("v", false, "set verbose mode")
	oConfigFilename = flag.String("c", "", "location of the configuration file. If empty then expect .kiya in $HOME.")
	oAuthLocation   = flag.String("a", "", "location of the JSON key credentials file. If empty then use the Google Application Defaults.")
	oVersion        = flag.Bool("version", false, "show the version of the tool")
	oOutputFilename = flag.String("o", "", "if not empty then write the secret to a file else write to stdout (get)")
	oCreatedAt      = flag.String("ca", "", "if not empty then use this formatted timestamp (RFC3339) as the creation date of the new secret")
	oQuiet          = flag.Bool("quiet", false, "don't prompt for confirmation on destructive actions")

	oEncrypt           = flag.Bool("encrypt", false, "encrypt backups")
	oEncryptDecryptKey = flag.String("key", "", "encrypt/decrypt key for backup and restore commands")
	oPath              = flag.String("path", "", "location for backup/restore file")
)
