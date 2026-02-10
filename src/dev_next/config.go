package main

import (
    "encoding/json"
	"fmt"
    "os"
)

type Config struct {
    Protocol            string        `json:"protocol"`
    FlowControl         string        `json:"flow-control"`
    XboardSkipSSLCheck  bool          `json:"xboard-skip-ssl-check"`
    XboardConfig        XboardConfig  `json:"xboard-config"`
    XuiSkipSSLCheck     bool          `json:"xui-skip-ssl-check"`
    XuiConfig           XuiConfig     `json:"xui-config"`
}

type XboardConfig struct {
    ApiHost string `json:"ApiHost"`
    ApiKey  string `json:"ApiKey"`
    NodeID  int    `json:"NodeID"`
	Database string `json:"database"`
}

type XuiConfig struct {
    ApiHost   string `json:"ApiHost"`
    User      string `json:"user"`
    Passwd    string `json:"passwd"`
    NodeID    int    `json:"NodeID"`
    AdminPath string `json:"admin-path"`
	Database string `json:"database"`
}

func readConf(conffile string) Config {
    data, err := os.ReadFile(conffile)
    if err != nil {
        fmt.Println("ERR: Read config err:",err)
		os.Exit(-1)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        fmt.Println("ERR: Load config json err:",err)
		os.Exit(-1)
    }
	
	mp:=cfg.Protocol
	
	if !(mp == "vless" || mp == "vmess" || mp == "trojan" || mp == "shadowsocks" ) {
		fmt.Println("ERR: Protocol not supported ",mp)
		os.Exit(-1)
	}
	
    return cfg
}