package kiya

import (
	"io/ioutil"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// GetAuthCredentials returns a credentials client option
func GetAuthCredentials(authLocation string) option.ClientOption {
	if len(authLocation) > 0 {
		// Your credentials should be obtained from the Google
		// Developer Console (https://console.developers.google.com).
		// Navigate to your project, then see the "Credentials" page
		// under "APIs & Auth".
		// To create a service account client, click "Create new Client ID",
		// select "Service Account", and click "Create Client ID". A JSON
		// key file will then be downloaded to your computer.
		data, err := ioutil.ReadFile(authLocation)
		if err != nil {
			log.Fatal("unable to read JSON key file", err)
		}
		creds, err := google.CredentialsFromJSON(context.Background(), data, cloudkms.DefaultAuthScopes()...)
		if err != nil {
			log.Fatal(err)
		}
		return option.WithCredentials(creds)
	}
	// Authorize the client using Aplication Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	creds, err := google.FindDefaultCredentials(context.Background(), cloudkms.DefaultAuthScopes()...)
	if err != nil {
		log.Fatal(err)
	}
	return option.WithCredentials(creds)
}
