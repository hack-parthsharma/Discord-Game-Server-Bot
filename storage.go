package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Asthetic/DiscordGameServerBot/network"
	"github.com/pkg/errors"
)

const filename = "data.json"

// WriteIP writes the current IP address to disk
func WriteIP(data network.Network) error {
	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "unable to marshal json data: %v", data)
	}

	if err = ioutil.WriteFile(filename, file, 0644); err != nil {
		return errors.Wrapf(err, "unable to write json data to file name: %v", filename)
	}

	return nil
}

// GetIP reads the current IP from disk
func GetIP() (string, error) {
	if !fileExists(filename) {
		if err := WriteIP(network.Network{}); err != nil {
			return "", errors.Wrapf(err, "error creating new file when it doesn't exist")
		}
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file: %v", filename)
	}

	data := network.Network{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal json: %v", filename)
	}

	return data.IP, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
