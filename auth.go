package kiya

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
)

// NewAuthenticatedClient creates an authenticated google client
func NewAuthenticatedClient(authLocation string) *http.Client {
	var client *http.Client
	if len(authLocation) > 0 {
		// Your credentials should be obtained from the Google
		// Developer Console (https://console.developers.google.com).
		// Navigate to your project, then see the "Credentials" page
		// under "APIs & Auth".
		// To create a service account client, click "Create new Client ID",
		// select "Service Account", and click "Create Client ID". A JSON
		// key file will then be downloaded to your computer.
		data, err := os.ReadFile(authLocation)
		if err != nil {
			fatal("unable to read JSON key file", err)
		}
		conf, err := google.JWTConfigFromJSON(data, cloudkms.CloudPlatformScope)
		if err != nil {
			fatal("unable to parse JSON key file", err)
		}
		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(context.Background())
	} else {
		// Authorize the client using Aplication Default Credentials.
		// See https://g.co/dv/identity/protocols/application-default-credentials
		defaultClient, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
		if err != nil {
			fatal(err)
		}
		client = defaultClient
	}
	return client
}
