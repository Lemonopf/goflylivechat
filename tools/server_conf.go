package tools

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const ServerConfPath = "config/server.json"

type ServerConf struct {
	JwtSecret        string
	SessionSecret    string
	Sub2apiJwtSecret string
}

var (
	serverConf     *ServerConf
	serverConfOnce sync.Once
)

// GetServerConf 获取服务端密钥配置，优先级：环境变量 > config/server.json > 自动生成并落盘
func GetServerConf() *ServerConf {
	serverConfOnce.Do(func() {
		serverConf = loadServerConf()
	})
	return serverConf
}

func loadServerConf() *ServerConf {
	conf := &ServerConf{}
	if info, err := ioutil.ReadFile(ServerConfPath); err == nil {
		if err := json.Unmarshal(info, conf); err != nil {
			log.Printf("解析 %s 失败: %v，将使用随机密钥\n", ServerConfPath, err)
		}
	}
	changed := false
	if os.Getenv("GOFLY_JWT_SECRET") != "" {
		conf.JwtSecret = os.Getenv("GOFLY_JWT_SECRET")
	} else if conf.JwtSecret == "" {
		conf.JwtSecret = randHex(32)
		changed = true
	}
	if os.Getenv("GOFLY_SESSION_SECRET") != "" {
		conf.SessionSecret = os.Getenv("GOFLY_SESSION_SECRET")
	} else if conf.SessionSecret == "" {
		conf.SessionSecret = randHex(32)
		changed = true
	}
	// sub2api 的 jwt.secret，用于离线校验嵌入 sub2api 时传入的用户 token，不自动生成
	if os.Getenv("GOFLY_SUB2API_JWT_SECRET") != "" {
		conf.Sub2apiJwtSecret = os.Getenv("GOFLY_SUB2API_JWT_SECRET")
	}
	if changed {
		saveServerConf(conf)
	}
	return conf
}

func saveServerConf(conf *ServerConf) {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Printf("序列化 server 配置失败: %v\n", err)
		return
	}
	if err := ioutil.WriteFile(ServerConfPath, data, 0600); err != nil {
		log.Printf("写入 %s 失败: %v（密钥仅在本次进程内有效）\n", ServerConfPath, err)
	}
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
