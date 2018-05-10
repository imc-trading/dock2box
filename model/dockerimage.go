package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/pborman/uuid"
)

type DockerImage struct {
	UUID     string     `json:"uuid"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated,omitempty"`
	Registry string     `json:"registry,omitempty"`
	Repo     string     `json:"repo,omitempty"`
	Name     string     `json:"name"`
	Tag      string     `json:"tag,omitempty"`
}

type DockerImages []*DockerImage

func NewDockerImage(registry string, repo string, name string, tag string) *DockerImage {
	return &DockerImage{
		UUID:     uuid.New(),
		Created:  time.Now(),
		Registry: registry,
		Repo:     repo,
		Name:     name,
		Tag:      tag,
	}
}

func (ds *Datastore) QueryDockerImage(q *qry.Query) (DockerImages, error) {
	kvs, err := ds.Values("dockerImages")
	if err != nil {
		return nil, err
	}

	dockerImages := DockerImages{}
	if err := kvs.Decode(&dockerImages); err != nil {
		return nil, err
	}

	r, err := q.Eval(dockerImages)
	if err != nil {
		return nil, err
	}

	return r.(DockerImages), nil
}

func (ds *Datastore) OneDockerImage(uuid string) (*DockerImage, error) {
	kvs, err := ds.Values(fmt.Sprintf("dockerImages/%s", uuid))
	if err != nil {
		return nil, err
	}

	dockerImages := DockerImages{}
	if err := kvs.Decode(&dockerImages); err != nil {
		return nil, err
	}

	if len(dockerImages) > 0 {
		return dockerImages[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateDockerImage(dockerImage *DockerImage) error {
	return ds.Set(fmt.Sprintf("dockerImages/%s", dockerImage.UUID), dockerImage)
}

func (ds *Datastore) UpdateDockerImage(dockerImage *DockerImage) error {
	now := time.Now()
	dockerImage.Updated = &now
	return ds.Set(fmt.Sprintf("dockerImages/%s", dockerImage.UUID), dockerImage)
}
