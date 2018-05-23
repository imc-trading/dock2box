package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type Image struct {
	UUID     string     `json:"uuid"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated,omitempty"`
	Registry string     `json:"registry,omitempty"`
	Repo     string     `json:"repo,omitempty"`
	Name     string     `json:"name"`
	Tag      string     `json:"tag,omitempty"`
}

type Images []*Image

func NewImage(registry string, repo string, name string, tag string) *Image {
	return &Image{
		UUID:     uuid.New(),
		Created:  time.Now(),
		Registry: registry,
		Repo:     repo,
		Name:     name,
		Tag:      tag,
	}
}

func (ds *Datastore) AllImages() (Images, error) {
	kvs, err := ds.Values("images")
	if err != nil {
		return nil, err
	}

	images := Images{}
	if err := kvs.Decode(&images); err != nil {
		return nil, err
	}

	return images, nil
}

func (ds *Datastore) QueryImages(q *qry.Query) (Images, error) {
	images, err := ds.AllImages()
	if err != nil {
		return nil, err
	}

	filtered, err := q.Query(images)
	if err != nil {
		return nil, err
	}

	return filtered.(Images), nil
}

func (ds *Datastore) OneImage(uuid string) (*Image, error) {
	kvs, err := ds.Values(fmt.Sprintf("images/%s", uuid))
	if err != nil {
		return nil, err
	}

	images := Images{}
	if err := kvs.Decode(&images); err != nil {
		return nil, err
	}

	if len(images) > 0 {
		return images[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateImage(image *Image) error {
	return ds.Set(fmt.Sprintf("images/%s", image.UUID), image)
}

func (ds *Datastore) UpdateImage(image *Image) error {
	now := time.Now()
	image.Updated = &now
	return ds.Set(fmt.Sprintf("images/%s", image.UUID), image)
}

func (ds *Datastore) DeleteImage(uuid string) error {
	if err := ds.Delete(fmt.Sprintf("images/%s", uuid)); err != nil {
		return err
	}
	return nil
}
