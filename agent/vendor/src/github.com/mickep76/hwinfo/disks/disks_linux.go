// +build linux

package disks

import (
	"github.com/mickep76/hwinfo/common"
	"path/filepath"
	"strconv"
)

// Get system disks.
func Get() ([]Disk, error) {
	da := []Disk{}

	files, err := filepath.Glob("/sys/class/block/*")
	if err != nil {
		return []Disk{}, err
	}

	for _, path := range files {
		o, err := common.LoadFiles([]string{
			filepath.Join(path, "dev"),
			filepath.Join(path, "size"),
		})
		if err != nil {
			return []Disk{}, err
		}

		d := Disk{}

		d.Name = filepath.Base(path)
		d.Device = o["dev"]

		d.SizeGB, err = strconv.Atoi(o["size"])
		if err != nil {
			return []Disk{}, err
		}
		d.SizeGB = d.SizeGB * 512 / 1024 / 1024 / 1024

		da = append(da, d)
	}

	return da, nil
}
