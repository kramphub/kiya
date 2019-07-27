package kiya

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

var profiles map[string]profile

type profile struct {
	Label       string
	ProjectID   string
	Location    string
	Keyring     string
	CryptoKey   string
	Bucket      string
	SecretRunes []rune
}

func load() (profs map[string]profile, err error) {
	reader, err := os.Open(configLocation())
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

func configLocation() string {
	location := *oConfigFilename
	if len(location) == 0 {
		location = path.Join(os.Getenv("HOME"), ".kiya")
	}
	return location
}

// LoadConfiguration loads the .kiya file
func LoadConfiguration() {
	profs, err := load()
	if err != nil {
		log.Fatal("unable to read/parse kiya configration file ("+configLocation()+")", err)
	}
	profiles = profs
}
