package server

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Health struct {
	mem  *mem.VirtualMemoryStat
	disk *disk.UsageStat
}

func newHealth(path string) (*Health, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	disk, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return &Health{mem: mem, disk: disk}, nil
}

// GetTotalMemory returns the total system memory in bytes.
func (h *Health) GetTotalMemory() uint64 {
	return h.mem.Total
}

// GetFreeMemory returns the available system memory in bytes.
func (h *Health) GetFreeMemory() uint64 {
	return h.mem.Available
}

func (h *Health) GetUsedDisk() uint64 {
	return h.disk.Used
}

func (h *Health) GetFreeDisk() uint64 {
	return h.disk.Free
}

func (h *Health) GetTotalDisk() uint64 {
	return h.disk.Total
}

func (h *Health) GetDiskPercent() float64 {
	return h.disk.UsedPercent
}
