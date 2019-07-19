# Setup infrastructure

## Google Project
In order to use kiya, you need a Google Cloud Platform project. See https://console.cloud.google.com.
You will need the project id for a new kiya profile.

## Google KMS
kiya uses Google KMS (Key Management Service) to encrypt and decrypt secrets.

### create a keyring
In the console, goto IAM & Admin, Encryption keys and create a keyring, e.g. `kiya-keyring`.

### create a crypto key
Create a crypto key for your keyring, e.g. `kiya-cryptokey`.

### set permissions
Select your crypto key and edit permissions. 
Add a member with your google account and role `Cloud KMS CryptoKey Encrypter/Decrypter`.

## Google Storage
kiya using Google Storage (Buckets) to store encrypted secrets using a label.

### create a bucket
In the console, goto Storage Browser and create a bucket, e.g. `your-name-kiya-secrets`.
Must be unique across Cloud Storage.

### set permissions
Select your bucket and edit permissions.
Add a member with your google account and role `Storage Object Admin` for creating and listing secrets.

## kiya profile
In your .kiya file, add a new personal profile.

    "private": {
        "projectID": "your-project-id",
        "location": "global",
        "keyring": "kiya-keyring",
        "cryptoKey": "kiya-cryptokey",
        "bucket": "your-name-kiya-secrets",
        "secretChars": "optionally-specify-secret-generation-chars"
    }	

### Test
Authenticate your account with Google.

    gcloud auth application-default login

Show the list of keys (which should be empty)

    kiya private list    