package channelling

import (
	"net/http"
	"regexp"
)

type Config struct {
	Title                           string                    // Title
	Ver                             string                    `json:"-"` // Version (not exported to Javascript)
	S                               string                    // Static URL prefix with version
	B                               string                    // Base URL
	Token                           string                    // Server token
	Renegotiation                   bool                      // Renegotiation flag
	StunURIs                        []string                  // STUN server URIs
	TurnURIs                        []string                  // TURN server URIs
	Tokens                          bool                      // True when we got a tokens file
	Version                         string                    // 服务器版本号
	UsersEnabled                    bool                      // 是否开启账户模式
	UsersAllowRegistration          bool                      // 是否开启用户注册
	UsersMode                       string                    // 账户模式
	DefaultRoomEnabled              bool                      // 是否开启默认房间 ("")
	Plugin                          string                    // 要加载的插件
	AuthorizeRoomCreation           bool                      // 创建房间是否要账户
	AuthorizeRoomJoin               bool                      // 加入房间是否要账户
	Modules                         []string                  // 开启的模块
	ModulesTable                    map[string]bool           `json:"-"` // 开启的模块的 映射表
	GlobalRoomID                    string                    `json:"-"` // 全局房间的 Id (not exported to Javascript)
	ContentSecurityPolicy           string                    `json:"-"` // HTML content security policy
	ContentSecurityPolicyReportOnly string                    `json:"-"` // HTML content security policy in report only mode
	RoomTypeDefault                 string                    `json:"-"` // 房间的默认类型
	RoomTypes                       map[*regexp.Regexp]string `json:"-"` // Map of regular expression -> room type
}

func (config *Config) WithModule(m string) bool {
	if val, ok := config.ModulesTable[m]; ok && val {
		return true
	}

	return false
}

func (config *Config) Get(request *http.Request) (int, interface{}, http.Header) {
	return 200, config, http.Header{"Content-Type": {"application/json; charset=utf-8"}}
}
