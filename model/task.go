package model

import (
	"fmt"
	"time"

	"github.com/mickep76/qry"
	"github.com/mickep76/runcmd"
	"github.com/pborman/uuid"
)

type Task struct {
	UUID         string            `json:"uuid"`
	Created      time.Time         `json:"created"`
	Updated      *time.Time        `json:"updated,omitempty"`
	HostUUID     string            `json:"hostUUID,omitempty"`
	Host         *Host             `json:"host,omitempty"`
	TaskDefUUID  string            `json:"taskDefUUID,omitempty"`
	TaskDef      *TaskDef          `json:"taskDef,omitempty"`
	Env          map[string]string `json:"env,omitempty" toml:"env"`
	Args         []string          `json:"args,omitempty" toml:"args"`
	Tail         []string          `json:"tail,omitempty"`
	Orphaned     bool              `json:"orphaned,omitempty"`
	Progress     float64           `json:"progress"`
	DurationLeft time.Duration     `json:"durationLeft"`

	*runcmd.Status
}

type Tasks []*Task

func NewTask() *Task {
	return &Task{
		UUID:    uuid.New(),
		Created: time.Now(),
	}
}

func (ds *Datastore) AllTasks() (Tasks, error) {
	kvs, err := ds.Values("tasks")
	if err != nil {
		return nil, err
	}

	tasks := Tasks{}
	if err := kvs.Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (ds *Datastore) QueryTasks(q *qry.Query) (Tasks, error) {
	tasks, err := ds.AllTasks()
	if err != nil {
		return nil, err
	}

	filtered, err := q.Query(tasks)
	if err != nil {
		return nil, err
	}

	return filtered.(Tasks), nil
}

func (ds *Datastore) OneTask(uuid string) (*Task, error) {
	tasks, err := ds.AllTasks()
	if err != nil {
		return nil, err
	}

	filtered, err := qry.New().Eq("uuid", uuid).Query(tasks)
	if err != nil {
		return nil, err
	}

	if len(filtered.(Tasks)) > 0 {
		return filtered.(Tasks)[0], nil
	}

	return nil, nil
}

func (ds *Datastore) CreateTask(task *Task) error {
	return ds.Set(fmt.Sprintf("tasks/%s/%s", task.HostUUID, task.UUID), task)
}

func (ds *Datastore) UpdateTask(task *Task) error {
	now := time.Now()
	task.Updated = &now
	return ds.Set(fmt.Sprintf("tasks/%s/%s", task.HostUUID, task.UUID), task)
}

func (ds *Datastore) DeleteTask(uuid string) error {
	task, err := ds.OneTask(uuid)
	if err != nil {
		return err
	}

	if err := ds.Delete(fmt.Sprintf("tasks/%s", task.HostUUID, uuid)); err != nil {
		return err
	}
	return nil
}
