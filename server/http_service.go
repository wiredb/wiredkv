package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/auula/wiredkv/clog"
	"github.com/auula/wiredkv/types"
	"github.com/auula/wiredkv/utils"
	"github.com/auula/wiredkv/vfs"
	"github.com/gorilla/mux"
)

const version = "wiredb/0.1.1"

var (
	root         *mux.Router
	authPassword string
	allowIpList  []string
)

// http://192.168.101.225:2668/{types}/{key}
// POST 创建 http://192.168.101.225:2668/zset/user-01-score
// PUT  更新 http://192.168.101.225:2668/zset/user-01-score
// GET  获取 http://192.168.101.225:2668/table/user-01-shop-cart

func init() {
	root = mux.NewRouter()
	root.HandleFunc("/", statusController)
	root.HandleFunc("/tables/{key}", tablesController)
	root.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		okResponse(w, http.StatusNotFound, nil, "404 Not Found - Oops!")
	})
	root.Use(authMiddleware)
}

type ResponseBody struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
}

type SystemInfo struct {
	KeyCount    int    `json:"key_count"`
	Version     string `json:"version"`
	GCStatus    int8   `json:"gc_status"`
	MemoryFree  string `json:"memory_free"`
	MemoryTotal string `json:"memory_total"`
}

func statusController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		okResponse(w, http.StatusMethodNotAllowed, nil, "HTTP Protocol Method Not Allowed!")
		return
	}

	info := SystemInfo{
		Version:  version,
		KeyCount: storage.KeysCount(),
		GCStatus: storage.RegionGCStatus(),
	}

	freemem, err := utils.GetFreeMemory()
	if err == nil {
		info.MemoryFree = fmt.Sprintf("%.2fGB", float64(freemem)/1024/1024/1024)
	} else {
		okResponse(w, http.StatusInternalServerError, nil, err.Error())
		clog.Errorf("HTTP server status controller: %s", err)
		return
	}

	totalmem, err := utils.GetTotalMemory()
	if err == nil {
		info.MemoryTotal = fmt.Sprintf("%.2fGB", float64(totalmem)/1024/1024/1024)
	} else {
		okResponse(w, http.StatusInternalServerError, nil, err.Error())
		clog.Errorf("HTTP server status controller: %s", err)
		return
	}

	okResponse(w, http.StatusOK, info, "")
}

func tablesController(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)["key"]
	if !ok || key == "" {
		okResponse(w, http.StatusOK, nil, "missing required parameter key!")
		return
	}

	switch r.Method {
	case http.MethodGet:
		seg, err := storage.FetchSegment(key)
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		table, err := seg.ToTables()
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		okResponse(w, http.StatusOK, table, "")
	case http.MethodPut:
		var tables types.Tables
		err := json.NewDecoder(r.Body).Decode(&tables)
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		seg, err := vfs.NewSegment(key, tables, tables.TTL)
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		err = storage.PutSegment(key, seg)
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		okResponse(w, http.StatusOK, nil, "request processed successfully!")
	case http.MethodPost:
		okResponse(w, http.StatusOK, nil, "request processed successfully!")
	case http.MethodDelete:
		err := storage.DeleteSegment(key)
		if err != nil {
			okResponse(w, http.StatusInternalServerError, nil, err.Error())
			clog.Errorf("HTTP server tables controller: %s", err)
			return
		}
		okResponse(w, http.StatusOK, nil, "delete data successfully!")
	default:
		okResponse(w, http.StatusMethodNotAllowed, nil, "HTTP Protocol Method Not Allowed!")
	}

}

func okResponse(w http.ResponseWriter, code int, result interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", version)
	w.WriteHeader(code)

	resp := ResponseBody{
		Code:    code,
		Result:  result,
		Message: message,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.Error(err)
	}
}

// 中间件函数，进行 Basic Auth 鉴权
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 使用标准 Auth 头
		token := r.Header.Get("Auth")
		clog.Debugf("HTTP request header authorization: %v", r.Header)

		// 获取客户端 IP 地址
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		// 检查 IP 白名单
		isAllowedIP := false
		if len(allowIpList) > 0 {
			for _, allowedIP := range allowIpList {
				if strings.Split(ip, ":")[0] == allowedIP {
					isAllowedIP = true
					break
				}
			}
		}
		if !isAllowedIP {
			clog.Warnf("Unauthorized IP address: %s", ip)
			okResponse(w, http.StatusUnauthorized, nil, fmt.Sprintf("Your IP %s is not allowed!", ip))
			return
		}

		if token != authPassword {
			clog.Warnf("Unauthorized access attempt from client %s", ip)
			okResponse(w, http.StatusUnauthorized, nil, "Access not authorised!")
			return
		}

		clog.Infof("Client %s connection successfully", ip)
		next.ServeHTTP(w, r)

	})
}
