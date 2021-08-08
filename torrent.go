package common

import "fmt"

type TorrentClient interface {
	// path: filepath to .torrent or a .magnet file
	// destination: path to dir where to download
	Download(filepath, destination string) error
}

type TorrentClientConfig struct {
	Address string `yaml:"address,omitempty"`
	User    string `yaml:"user,omitempty"`
	Pass    string `yaml:"pass,omitempty"`
	Type    string `yaml:"type,omitempty"`
}

func (tcc TorrentClientConfig) NewTorrentClientByType() (TorrentClient, error) {
	switch tcc.Type {
	// todo implement more!
	case "transmission":
		return NewTransmissionClient(tcc.Address, tcc.User, tcc.Pass)
	case "rtorrent":
		return NewRtorrentClient(tcc.Address)
	}
	return nil, fmt.Errorf("no torrent client of type %v available", tcc.Type)
}
