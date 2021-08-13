package common

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mrobinsn/go-rtorrent/rtorrent"
)

type RtorrentClient struct {
	*rtorrent.RTorrent
}

func NewRtorrentClient(address string) (*RtorrentClient, error) {
	c := rtorrent.New(address, false)
	if _, err := c.Name(); err != nil {
		return nil, err
	}
	return &RtorrentClient{c}, nil
}

func (rc *RtorrentClient) AddFromFile(path, destination string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return rc.addContent(content, path, destination)
}

func (rc *RtorrentClient) AddFromUrl(url, destination string) error {
	url, content, err := GetUrlContent(url)
	if err != nil {
		return err
	}
	return rc.addContent(content, url, destination)
}

func (rc *RtorrentClient) AddContent(content []byte, type_ string, destination string) error {
	return rc.addContent(content, fmt.Sprintf(".%v", type_), destination)
}

func (rc *RtorrentClient) addContent(content []byte, path, destination string) error {
	if filepath.Ext(path) == ".magnet" {
		return rc.addMagnet(content, destination)
	}
	if filepath.Ext(path) == ".torrent" {
		return rc.addTorrent(content, destination)
	}
	return fmt.Errorf("must be torrent or magnet, %q given", filepath.Ext(path))
}

func (rc *RtorrentClient) addMagnet(content []byte, destination string) error {
	return rc.Add(string(content), &rtorrent.FieldValue{Field: "d.directory", Value: destination})
}

func (rc *RtorrentClient) addTorrent(content []byte, destination string) error {
	return rc.AddTorrent(content, rtorrent.DBasePath.SetValue(destination))
}
