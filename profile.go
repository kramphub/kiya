package kiya

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

// Profiles is a collection of profiles as described in the .kiya configuration file
var Profiles map[string]Profile

// Profile represents a single profile within the .kiya configuration file
type Profile struct {
	Label       string
	ProjectID   string
	Location    string
	Keyring     string
	CryptoKey   string
	Bucket      string
	SecretRunes []rune
}

func load(configFilename string) (profs map[string]Profile, err error) {
	reader, err := os.Open(configLocation(configFilename))
	defer reader.Close()
	if err != nil {
		return
	}
	err = json.NewDecoder(reader).Decode(&profs)
	// ensure profile knows label
	for l, p := range profs {
		each := p
		each.Label = l
		profs[l] = each
	}
	return
}

func configLocation(configFilename string) string {
	location := configFilename
	if len(location) == 0 {
		location = path.Join(os.Getenv("HOME"), ".kiya")
	}
	return location
}

// Load configuration from .kiya file
func LoadConfiguration(configFilename string) {
	profs, err := load(configFilename)
	if err != nil {
		log.Fatal("unable to read/parse kiya configration file ("+configLocation(configFilename)+")", err)
	}
	Profiles = profs
}
