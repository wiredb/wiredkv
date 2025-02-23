package server

import (
	"fmt"
	"net/http"

	"github.com/auula/wiredkv/types"
	"github.com/auula/wiredkv/utils"
	"github.com/auula/wiredkv/vfs"
	"github.com/gin-gonic/gin"
)

var storage *vfs.LogStructuredFS

// 返回一个处理函数的工厂函数
func GetListController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func PutListController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func DeleteListController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func GetTableController(ctx *gin.Context) {
	_, seg, err := storage.FetchSegment(ctx.Param("key"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "key data not found.",
		})
		return
	}

	tables, err := seg.ToTable()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tables)
}

func PutTableController(ctx *gin.Context) {
	key := ctx.Param("key")

	var tables types.Table
	err := ctx.ShouldBindJSON(&tables)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	seg, err := vfs.NewSegment(key, tables, tables.TTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = storage.PutSegment(key, seg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "request processed succeed.",
	})
}

func DeleteTableController(ctx *gin.Context) {
	key := ctx.Param("key")

	err := storage.DeleteSegment(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "delete data succeed.",
	})
}

func GetZsetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func PutZsetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func DeleteZsetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func GetTextController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func PutTextController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func DeleteTextController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func GetNumberController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func PutNumberController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func DeleteNumberController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func GetSetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func PutSetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func DeleteSetController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "List of users",
	})
}

func GetHealthController(ctx *gin.Context) {
	health, err := newHealth(storage.GetDirectory())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, SystemInfo{
		Version:     version,
		GCState:     storage.GCState(),
		KeyCount:    storage.KeysCount(),
		DiskFree:    fmt.Sprintf("%.2fGB", utils.BytesToGB(health.GetFreeDisk())),
		DiskUsed:    fmt.Sprintf("%.2fGB", utils.BytesToGB(health.GetUsedDisk())),
		DiskTotal:   fmt.Sprintf("%.2fGB", utils.BytesToGB(health.GetTotalDisk())),
		MemoryFree:  fmt.Sprintf("%.2fGB", utils.BytesToGB(health.GetFreeMemory())),
		MemoryTotal: fmt.Sprintf("%.2fGB", utils.BytesToGB(health.GetTotalMemory())),
		DiskPercent: fmt.Sprintf("%.2f%%", health.GetDiskPercent()),
	})
}

func Error404Handler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "Oops! 404 Not Found!",
	})
}
