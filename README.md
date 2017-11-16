# Kiya #

<img align="right" src="kea.jpg">

Kiya is a tool to access secrets stored in a Google Bucket and encrypted by Google Key Management Service (KMS).

A secret must have a label and a plain text representation of its value.
A label is typically composed of a domain or site or application (the parent key) and a secret key, e.g. google.gmail/info@mars.planets.
A label must have at least one parent key (lowercase with or without dots).
The value must be a string which has a maximum length of 64Kb.

### Prerequisites
Kiya uses your authenticated Google account to access the Storage Bucket, KMS and Audit Logging.
The bucket stores the encrypted secret value using the label as the storage key.

	gcloud auth application-default login
	
## Usage

### Configuration

Create a file name `.kiya` in your home directory with the content for a shareable secrets profile. You can have multiple profiles for different usages.

	{
		"shared": {
			"projectID": "your-gcp-project",
			"location": "global",
			"keyring": "your-kiya-secrets-keyring",
			"cryptoKey": "your-kiya-secrets-cryptokey",
			"bucket": "your-kiya-secrets"
		}
	}

Read `setup.md` for detailed instructions how to create the bucket, encryption ring and key and set the permissions.

### Store a password, _put_

	kiya shared put concourse/cd-pipeline mySecretPassword
	
In this example, `shared` refers to the profile in your configuration. `concourse` refers to the site or domain. `cd-pipeline` is the username which can be an email address too. `mySecretPassword` is the plain text password.

If a password was already stored then you will be warned about overwriting it.

### Retrieve a password, _get_

	kiya shared get concourse/cd-pipeline

### List labels of stored secrets, _list_

	kiya shared list

### Fill a template, _template_

    kiya shared template template-file

Output will be written to stdout.

Example contents of `template-file`:

    bitbucket-password={{kiya "key-to-bitbucket-password"}}
    
Kiya provides a builtin function for base64 encoding

    artifatory-hashed-password={{base64 (kiya "key-to-artifatory-password")}}

### Write a secret to clipboard, _copy_

	kiya shared copy concourse/cd-pipeline

### Create secret from clipboard, _paste_

	kiya shared paste google/accounts/someone@gmail.com

## Troubleshooting

### 1. Error

	2017/06/24 22:14:24 google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.

Run

	gcloud auth application-default login

### 2. Error

	googleapi: Error 403: Caller does not have storage.objects.list access to bucket <some-bucket-name>., forbidden

You do not have access to encrypted secrets from `some-bucket-name`.

(c) 2017 kramphub.com. Apache License v2.