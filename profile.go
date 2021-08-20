package kiya

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/kramphub/kiya/backend"
)

// Profiles is a collection of profiles as described in the .kiya configuration
var Profiles map[string]backend.Profile

func load(configFile string) (profs map[string]backend.Profile, err error) {
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
