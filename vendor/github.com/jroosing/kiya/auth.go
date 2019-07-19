package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func newAuthenticatedClient() *http.Client {
	var client *http.Client
	if len(*oAuthLocation) > 0 {
		// Your credentials should be obtained from the Google
		// Developer Console (https://console.developers.google.com).
		// Navigate to your project, then see the "Credentials" page
		// under "APIs & Auth".
		// To create a service account client, click "Create new Client ID",
		// select "Service Account", and click "Create Client ID". A JSON
		// key file will then be downloaded to your computer.
		data, err := ioutil.ReadFile(*oAuthLocation)
		if err != nil {
			log.Fatal("unable to read JSON key file", err)
		}
		conf, err := google.JWTConfigFromJSON(data, cloudkms.CloudPlatformScope)
		if err != nil {
			log.Fatal("unable to parse JSON key file", err)
		}
		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(oauth2.NoContext)
	} else {
		// Authorize the client using Aplication Default Credentials.
		// See https://g.co/dv/identity/protocols/application-default-credentials
		defaultClient, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
		if err != nil {
			log.Fatal(err)
		}
		client = defaultClient
	}
	return client
}
