package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

// GetTotalMemory returns the total system memory in bytes.
func GetTotalMemory() (uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get total memory: %w", err)
	}
	return v.Total, nil
}

// GetFreeMemory returns the available system memory in bytes.
func GetFreeMemory() (uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get available memory: %w", err)
	}
	return v.Available, nil
}
