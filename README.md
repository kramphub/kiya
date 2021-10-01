# Kiya #

<img align="right" src="kea.jpg">

Kiya is a tool to access secrets stored in Google Secret Manager(GSM) or a Google Bucket and encrypted by Google Key Management Service (KMS).

### Introduction

Developing and deploying applications to execution environments (dev,staging,production) requires all kinds of secrets.
Both continuous development enviroment and production environment require credentials to access other resources.
Examples are passwords, service accounts, TLS certificates, API tokens and Encryption Keys. These secrets should be
managed with great care. This means secrets must be stored encrypted on reliable shared storage and its access must
controlled by AAA (authentication, authorisation and auditing).

**Kiya** is a simple tool that eases the access to the secrets stored in GSM or KMS. It requires an authenticated Google account and permissions for that account to read secrets and perform
encryption and decryption.

#### Labeled secrets

A secret must have a label and a plain text representation of its value. A label is typically composed of a domain or
site or application (the parent key) and a secret key, e.g. google.gmail/info@mars.planets. A label must have at least
one parent key (lowercase with or without dots). The value must be a string which has a maximum length of 64Kb.

### Prerequisites

Kiya uses your authenticated Google account to access the Secret Manager / Storage Bucket, KMS and Audit Logging. The bucket stores the
encrypted secret value using the label as the storage key.

	gcloud auth application-default login

## Install

	go install github.com/kramphub/kiya/cmd/kiya@latest

## Usage

Read `setup.md` for detailed instructions how to setup the basic prerequisites.

### Configuration

Create a file name `.kiya` in your home directory with the content for a shareable secrets profile. You can have
multiple profiles for different usages. Each profile should mention a either `kms` or `gsm` to be used as the `backend`.
If no value is defined for a profile's `backend`, `kms` will be used as a default.

```json
{
	"teamF1-on-kms": {
		"backend": "kms",
		"projectID": "your-gcp-project",
		"location": "global",
		"keyring": "your-kiya-secrets-keyring",
		"cryptoKey": "your-kiya-secrets-cryptokey",
		"bucket": "your-kiya-secrets"
	},
	"teamF2-on-gsm": {
		"backend": "gsm",
		"projectID": "another-gcp-project"
	}
}

```

You should define `location`, `keyring`, `cryptoKey` and `bucket` for KMS based profiles.

For Google Secret Manager based profiles a `projectID` is only enough. 

### Store a password, _put_

	kiya teamF1 put concourse/cd-pipeline mySecretPassword

In this example, `teamF1` refers to the profile in your configuration. `concourse` refers to the site or
domain. `cd-pipeline` is the username which can be an email address too. `mySecretPassword` is the plain text password.

If a password was already stored then you will be warned about overwriting it. The -quiet flag can be used to skip the
confirmation prompt:

	kiya -quiet teamF1 put concourse/cd-pipeline myNewSecretPassword

_Note: this will put a secret in your command history; better use paste, see below._

### Generate a password, _generate_

	kiya teamF1 generate concourse/cd-pipeline 25

Generate a secret with length 25 store it as secret `concourse/cd-pipeline` and copy its value to the OS clipboard.

### Retrieve a password, _get_

	kiya teamF1 get concourse/cd-pipeline

_Note: this will put a secret in your command history; better use copy, see below._

### List labels of stored secrets, _list_

	kiya teamF1 list [|filter]

Specifying a filter argument will hide any keys that don't contain the filter string.

The list command is also used when the command is unknown, e.g. `kiya teamF1 list redbull` shows the same results
as `kiya teamF1 redbull`.

### Fill a template, _template_

    kiya teamF1 template template-file

Output will be written to stdout.

Example contents of `template-file`:

    bitbucket-password={{kiya "key-to-bitbucket-password"}}

Kiya also provides a builtin function for base64 encoding:

    artifatory-hashed-password={{base64 (kiya "key-to-artifatory-password")}}

For accessing OS environment values:

    gcp-project={{env "PROJECT"}}

### Write a secret to clipboard, _copy_

    kiya teamF1 copy concourse/cd-pipeline

### Create secret from clipboard, _paste_

    kiya teamF1 paste google/accounts/someone@gmail.com

### Move a secret from one profile to another, _move_

    kiya teamF1 move bitbucket.org/johndoe teamF2

## Troubleshooting

### 1. Error

	2017/06/24 22:14:24 google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.

Run

	gcloud auth application-default login

### 2. Error

	googleapi: Error 403: Caller does not have storage.objects.list access to bucket <some-bucket-name>., forbidden

You do not have access to encrypted secrets from `some-bucket-name`.

&copy; 2017 kramphub.com. Apache License v2.
