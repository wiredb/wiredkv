package utils

// BytesToGB converts a given size in bytes to gigabytes (GB).
func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
