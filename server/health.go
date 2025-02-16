package server

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type health struct {
	mem  *mem.VirtualMemoryStat
	disk *disk.UsageStat
}

func newHealth(path string) (*health, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	disk, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return &health{mem: mem, disk: disk}, nil
}

// GetTotalMemory returns the total system memory in bytes.
func (h *health) GetTotalMemory() uint64 {
	return h.mem.Total
}

// GetFreeMemory returns the available system memory in bytes.
func (h *health) GetFreeMemory() uint64 {
	return h.mem.Available
}

func (h *health) GetUsedDisk() uint64 {
	return h.disk.Used
}

func (h *health) GetFreeDisk() uint64 {
	return h.disk.Free
}

func (h *health) GetTotalDisk() uint64 {
	return h.disk.Total
}

func (h *health) GetDiskPercent() float64 {
	return h.disk.UsedPercent
}
