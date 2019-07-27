package kiya

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

// Profiles is a collection of profiles as described in the .kiya configuration
var Profiles map[string]Profile

// Profile describes a single profile in a .kiya configuration
type Profile struct {
	Label       string
	ProjectID   string
	Location    string
	Keyring     string
	CryptoKey   string
	Bucket      string
	SecretRunes []rune
}

func load(configFile string) (profs map[string]Profile, err error) {
	reader, err := os.Open(configLocation(configFile))
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

func configLocation(configFile string) string {
	location := configFile
	if len(location) == 0 {
		location = path.Join(os.Getenv("HOME"), ".kiya")
	}
	return location
}

// LoadConfiguration loads the .kiya file
func LoadConfiguration(configFile string) {
	profs, err := load(configFile)
	if err != nil {
		log.Fatal("unable to read/parse kiya configration file ("+configLocation(configFile)+")", err)
	}
	Profiles = profs
}
