package common

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/odwrtw/transmission"
)

type TransmissionConfig struct {
	Address string `yaml:"address,omitempty"`
	User    string `yaml:"user,omitempty"`
	Pass    string `yaml:"pass,omitempty"`
}

type TransmissionClient struct {
	*transmission.Client
}

func NewTransmissionClient(address, user, pass string) (*TransmissionClient, error) {
	tc, err := transmission.New(transmission.Config{Address: address, User: user, Password: pass})
	if err != nil {
		return nil, err
	}
	return &TransmissionClient{tc}, nil
}

func (tm *TransmissionClient) Download(path, destination string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if filepath.Ext(path) == ".magnet" {
		return tm.addMagnet(content, destination)
	}
	if filepath.Ext(path) == ".torrent" {
		return tm.addTorrent(content, destination)
	}
	return fmt.Errorf("Extension must be torrent or magnet, %q given", filepath.Ext(path))
}

func (tm *TransmissionClient) addMagnet(content []byte, destination string) error {
	_, err := tm.AddTorrent(transmission.AddTorrentArg{
		DownloadDir: destination,
		Filename:    string(content),
	})
	return err
}

func (tm *TransmissionClient) addTorrent(content []byte, destination string) error {
	_, err := tm.AddTorrent(transmission.AddTorrentArg{
		DownloadDir: destination,
		Metainfo:    base64.StdEncoding.EncodeToString(content),
	})
	return err
}
