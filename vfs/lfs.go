package vfs

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/auula/vasedb/utils"
)

var (
	once              sync.Once
	indexShard        = 5
	instance          *LogStructuredFS
	fsPerm            = fs.FileMode(0755)
	dataFileExtension = ".vsdb"
	indexFullFileName = "idx.vsi"
	dataFileMetadata  = []byte{0xDB, 0x0, 0x0, 0x1}
)

const RWCA = os.O_RDWR | os.O_CREATE | os.O_APPEND

type Options struct {
	Path   string
	FsPerm os.FileMode
}

// INode represents a file system node with metadata.
type INode struct {
	RegionID    uint16    // Unique identifier for the INode
	Offset      uint32    // Offset within the file
	CreatedTime time.Time // Creation time of the INode
	EexpireTime time.Time // Expiration time of the INode
}

type indexMap struct {
	mu    sync.RWMutex      // 每个分片使用独立的锁
	index map[uint64]*INode // 存储映射
}

// LogStructuredFS represents the virtual file storage system.
type LogStructuredFS struct {
	mu        sync.Mutex
	offset    uint64
	regionID  uint64
	directory string
	indexs    []*indexMap         // Index mapping for INode references
	active    *os.File            // Currently active file for writing
	regions   map[uint64]*os.File // Archived files keyed by unique file ID
}

// 根据某种哈希函数（如简单的模运算）来选择分片
func (lfs *LogStructuredFS) getShardIndex(key uint64) *indexMap {
	return lfs.indexs[key%uint64(indexShard)]
}

// 使用 `getShardIndex` 获取分片，并加锁进行操作
func (lfs *LogStructuredFS) AddINode(key uint64, inode *INode) {
	shard := lfs.getShardIndex(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.index[key] = inode
}

func (lfs *LogStructuredFS) GetINode(key uint64) (*INode, bool) {
	shard := lfs.getShardIndex(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	inode, exists := shard.index[key]
	return inode, exists
}

func (lfs *LogStructuredFS) BatchINodes(inodes ...*INode) {

}

func HashSum64(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

func (lfs *LogStructuredFS) ChangeReigons() error {
	lfs.mu.Lock()
	defer lfs.mu.Unlock()
	err := lfs.active.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync active file: %w", err)
	}

	err = lfs.active.Close()
	if err != nil {
		return fmt.Errorf("failed to close active file: %w", err)
	}

	return nil
}

func (lfs *LogStructuredFS) createActiveReigon() error {
	lfs.mu.Lock()
	defer lfs.mu.Unlock()
	lfs.regionID += 1
	fileName, err := newDataFileName(lfs.regionID)
	if err != nil {
		return fmt.Errorf("failed to new data file name: %w", err)
	}

	activeFile, err := os.OpenFile(filepath.Join(lfs.directory, fileName), RWCA, fsPerm)
	if err != nil {
		return fmt.Errorf("failed to create active file: %w", err)
	}

	n, err := activeFile.Write(dataFileMetadata)
	if err != nil {
		return fmt.Errorf("failed to write data file metadata: %w", err)
	}

	if n != len(dataFileMetadata) {
		return errors.New("failed to new file metadata write")
	}

	lfs.active = activeFile
	lfs.offset = uint64(len(dataFileMetadata))

	return nil
}

func (lfs *LogStructuredFS) recoverRegions() error {
	lfs.mu.Lock()
	defer lfs.mu.Unlock()
	files, err := os.ReadDir(lfs.directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) <= 0 {
		lfs.regionID = 1
	} else {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), dataFileExtension) {
				regions, err := os.Open(filepath.Join(lfs.directory, file.Name()))
				if err != nil {
					return fmt.Errorf("failed to open data file: %w", err)
				}

				regionID, err := dataFileNameToUint64(file.Name())
				if err != nil {
					return fmt.Errorf("failed to get regions id: %w", err)
				}
				lfs.regions[regionID] = regions
			}
		}

		var regionIds []uint64
		for v := range lfs.regions {
			regionIds = append(regionIds, v)
		}
		// 对 regionIds 切片从小到大排序
		sort.Slice(regionIds, func(i, j int) bool {
			return regionIds[i] < regionIds[j]
		})
		// 找到最新数据文件的版本
		lfs.regionID = regionIds[len(regionIds)-1]
	}

	return nil
}

func OpenFS(opt *Options) (*LogStructuredFS, error) {
	var top_err error
	once.Do(func() {
		if instance != nil {
			return
		}

		err := checkFileSystem(opt.Path)
		if err != nil {
			top_err = err
			return
		}

		fsPerm = opt.FsPerm
		instance = &LogStructuredFS{
			indexs:    make([]*indexMap, indexShard),
			regions:   make(map[uint64]*os.File, 10),
			offset:    uint64(len(dataFileMetadata)),
			regionID:  0,
			directory: opt.Path,
		}

		for i := 0; i < indexShard; i++ {
			instance.indexs[i] = &indexMap{
				mu:    sync.RWMutex{},
				index: make(map[uint64]*INode),
			}
		}

		// 先对已有的数据文件执行恢复操作，并且初始化内存中的数据版本号
		err = instance.recoverRegions()
		if err != nil {
			top_err = fmt.Errorf("failed to recover data regions: %w", err)
			return
		}

		err = instance.createActiveReigon()
		if err != nil {
			top_err = fmt.Errorf("failed to create active regions: %w", err)
			return
		}

	})

	if top_err != nil {
		return nil, fmt.Errorf("failed to open file system: %w", top_err)
	}

	// 单例子模式，但是挡不住其他包通过 new(LogStructuredFS) 也能创建一个实例，那这样根本不起作用了
	return instance, nil
}

func (lfs *LogStructuredFS) CloseFS() error {
	lfs.mu.Lock()
	defer lfs.mu.Unlock()
	for _, file := range lfs.regions {
		err := utils.CloseFile(file)
		if err != nil {
			return fmt.Errorf("failed to close region file: %w", err)
		}
	}

	err := utils.CloseFile(lfs.active)
	if err != nil {
		return fmt.Errorf("failed to close active region: %w", err)
	}

	// 如果有 index 文件的快照，就从 index 文件快照进行恢复，如果没有就全局扫描
	return nil
}

func validateFileHeader(file *os.File) error {
	var fileHeader [4]byte
	n, err := file.Read(fileHeader[:])
	if err != nil {
		return err
	}

	if n != len(dataFileMetadata) {
		return errors.New("file is too short to contain valid signature")
	}

	if !bytes.Equal(fileHeader[:], dataFileMetadata[:]) {
		return fmt.Errorf("unsupported data file version: %v", file.Name())
	}

	return nil
}

func checkFileSystem(path string) error {
	if !utils.IsExist(path) {
		err := os.MkdirAll(path, fsPerm)
		if err != nil {
			return err
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(files) > 0 {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), dataFileExtension) {
				if len(file.Name()) == 8 && strings.HasPrefix(file.Name(), "000") {
					file, err := os.Open(filepath.Join(path, file.Name()))
					if err != nil {
						return fmt.Errorf("failed to check data file: %w", err)
					}
					defer file.Close()

					err = validateFileHeader(file)
					if err != nil {
						return fmt.Errorf("failed to validated data file header: %w", err)
					}
				}
			}

			if !file.IsDir() && file.Name() == indexFullFileName {
				file, err := os.Open(indexFullFileName)
				if err != nil {
					return fmt.Errorf("failed to check index file: %w", err)
				}
				defer file.Close()

				err = validateFileHeader(file)
				if err != nil {
					return fmt.Errorf("failed to validated index file header: %w", err)
				}
			}
		}
	}

	return nil
}

func newDataFileName(regionID uint64) (string, error) {
	fileName := fmt.Sprintf("%08d%s", regionID, dataFileExtension)
	if len(fileName) == 8 && strings.HasPrefix(fileName, "000") {
		return fileName, nil
	}
	return "", fmt.Errorf("new region id %d cannot be converted to a valid file name", regionID)
}

// fileNameToUint16 将文件名（如 0000001.vsdb）中的数字部分转换为 uint16
func dataFileNameToUint64(fileName string) (uint64, error) {
	parts := strings.Split(fileName, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid file name format: %s", fileName)
	}

	// 转换为 uint16
	number, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number from file name: %w", err)
	}

	return uint64(number), nil
}

// Uint16ToFileName 将 uint16 转换为文件名格式（如 1 转为 0000001.vsdb）
func uint64ToDataFileName(number uint64) string {
	return fmt.Sprintf("%08d%s", number, dataFileExtension)
}
