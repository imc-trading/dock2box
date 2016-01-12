// +build linux

package lvm

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func Get() (LVM, error) {
	lvm := LVM{}

	pvs, err := GetPhysVols()
	if err != nil {
		return LVM{}, err
	}
	lvm.PhysVols = &pvs

	lvs, err := GetLogVols()
	if err != nil {
		return LVM{}, err
	}
	lvm.LogVols = &lvs

	vgs, err := GetVolGrps()
	if err != nil {
		return LVM{}, err
	}
	lvm.VolGrps = &vgs

	return lvm, nil
}

func GetPhysVols() ([]PhysVol, error) {
	pvs := []PhysVol{}

	_, err := exec.LookPath("pvs")
	if err != nil {
		return []PhysVol{}, errors.New("command doesn't exist: pvs")
	}

	o, err := exec.Command("pvs", "--units", "B").Output()
	if err != nil {
		return []PhysVol{}, err
	}

	for c, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if c < 1 || len(vals) < 1 {
			continue
		}

		pv := PhysVol{}

		pv.Name = vals[0]
		pv.VolGrp = vals[1]
		pv.Format = vals[2]
		pv.Attr = vals[3]

		pv.SizeGB, err = strconv.Atoi(strings.TrimRight(vals[4], "B"))
		if err != nil {
			return []PhysVol{}, err
		}
		pv.SizeGB = pv.SizeGB / 1024 / 1024 / 1024

		pv.FreeGB, err = strconv.Atoi(strings.TrimRight(vals[5], "B"))
		if err != nil {
			return []PhysVol{}, err
		}
		pv.FreeGB = pv.FreeGB / 1024 / 1024 / 1024

		pvs = append(pvs, pv)
	}

	return pvs, nil
}

func GetLogVols() ([]LogVol, error) {
	lvs := []LogVol{}

	_, err := exec.LookPath("lvs")
	if err != nil {
		return []LogVol{}, errors.New("command doesn't exist: lvs")
	}

	o, err := exec.Command("lvs", "--units", "B").Output()
	if err != nil {
		return []LogVol{}, err
	}

	for c, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if c < 1 || len(vals) < 1 {
			continue
		}

		lv := LogVol{}

		lv.Name = vals[0]
		lv.VolGrp = vals[1]
		lv.Attr = vals[2]

		lv.SizeGB, err = strconv.Atoi(strings.TrimRight(vals[3], "B"))
		if err != nil {
			return []LogVol{}, err
		}
		lv.SizeGB = lv.SizeGB / 1024 / 1024 / 1024

		lvs = append(lvs, lv)
	}

	return lvs, nil
}

func GetVolGrps() ([]VolGrp, error) {
	vgs := []VolGrp{}

	_, err := exec.LookPath("vgs")
	if err != nil {
		return []VolGrp{}, errors.New("command doesn't exist: vgs")
	}

	o, err := exec.Command("vgs", "--units", "B").Output()
	if err != nil {
		return []VolGrp{}, err
	}

	for c, line := range strings.Split(string(o), "\n") {
		vals := strings.Fields(line)
		if c < 1 || len(vals) < 1 {
			continue
		}

		vg := VolGrp{}

		vg.Name = vals[0]
		vg.Attr = vals[4]

		vg.SizeGB, err = strconv.Atoi(strings.TrimRight(vals[5], "B"))
		if err != nil {
			return []VolGrp{}, err
		}
		vg.SizeGB = vg.SizeGB / 1024 / 1024 / 1024

		vg.FreeGB, err = strconv.Atoi(strings.TrimRight(vals[6], "B"))
		if err != nil {
			return []VolGrp{}, err
		}
		vg.FreeGB = vg.FreeGB / 1024 / 1024 / 1024

		vgs = append(vgs, vg)
	}

	return vgs, nil
}
