package main

import "flag"

var (
	oConfigFilename = flag.String("c", "", "location of the configuration file. If empty then expect .kiya in $HOME.")
	oAuthLocation   = flag.String("a", "", "location of the JSON key credentials file. If empty then use the Google Application Defaults.")
	oVersion        = flag.Bool("version", false, "show the version of the tool")
	oOutputFilename = flag.String("o", "", "if not empty then write the secret to a file else write to stdout (get)")
	oQuiet          = flag.Bool("quiet", false, "don't prompt for confirmation on destructive actions")

	// Backup flags
	oEncryptBackup          = flag.Bool("encrypt-backup", false, "if true, the backup will be encrypted")
	oBackupKeyStore         = flag.String("backup-key-store", "file", "storage type for public key, 'store' or 'file'")
	oBackupKey              = flag.String("backup-key", "./kiya_backupkey_rsa", "key to encrypt/decrypt the backup")
	oBackupPath             = flag.String("backup-path", "./kiya_backup", "backup file path")
	oBackupRestoreOverwrite = flag.Bool("backup-restore-overwrite", false, "if true, the restore will overwrite existing secrets")
)
