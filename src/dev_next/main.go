package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
	"os"

    _ "github.com/mattn/go-sqlite3"
)

var (
    XboardDBPath string
    XuiDBPath string
    SyncInterval = 90 * time.Second
	xuiID int
	xboardID int
	xuiFlow string
)

type XboardUser struct {
    ID            int
    Email         string
    UUID          string
    U             int64
    D             int64
    TransferEnable int64
    GroupID       int
    PlanID        int
}

type XuiClient struct {
    ID          string `json:"id"`
    Security    string `json:"security"`
    Password    string `json:"password"`
    Flow        string `json:"flow"`
    Email       string `json:"email"`
    LimitIp     int    `json:"limitIp"`
    TotalGB     int64  `json:"totalGB"`
    ExpiryTime  int64  `json:"expiryTime"`
    Enable      bool   `json:"enable"`
    TgId        int    `json:"tgId"`
    SubId       string `json:"subId"`
    Comment     string `json:"comment"`
    Reset       int    `json:"reset"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}

type XuiSettings struct {
    Clients    []XuiClient `json:"clients"`
    Decryption string      `json:"decryption"`
    Encryption string      `json:"encryption"`
}

type XuiInbound struct {
    ID       int
    Port     int
    Protocol string
    Settings string
    Remark   string
}

type XuiClientTraffic struct {
    InboundID int
    Email     string
    Up        int64
    Down      int64
    Enable    bool
}

func main() {
	argc:=len(os.Args)
	if argc <= 1 {
		fmt.Println("xboard2xui [config.json] [logfile]")
		return
	}
	
	setUP:=readConf(os.Args[1])
	XboardDBPath=setUP.XboardConfig.Database
	XuiDBPath=setUP.XuiConfig.Database
	xuiID=setUP.XuiConfig.NodeID
	xboardID=setUP.XuiConfig.NodeID
	xuiFlow=setUP.FlowControl
    log.Println("INFO: Init xboard server info, Please wait a few minutes.")

    xboardDB, err := sql.Open("sqlite3", XboardDBPath)
    if err != nil {
        fmt.Println("Err: xboard init fail",err)
    }
    defer xboardDB.Close()

    xuiDB, err := sql.Open("sqlite3", XuiDBPath)
    if err != nil {
        fmt.Println("Err: fail to connect xui",err)
    }
    defer xuiDB.Close()

    if err := syncUsers(xboardDB, xuiDB); err != nil {
       fmt.Println("Err: Read xboard user info fail:", err)
    }

    if err := syncTraffic(xboardDB, xuiDB); err != nil {
        log.Println("ERR: xui return database is locked", err)
    }

	log.Println("INFO: Complete Init xboard server info.")
	
	if argc >2 {
	l,err :=os.OpenFile(os.Args[2],os.O_WRONLY|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(l)
	}
	
    ticker := time.NewTicker(SyncInterval)
    defer ticker.Stop()

    for range ticker.C {
		err := syncUsers(xboardDB, xuiDB);
        if err != nil {
            log.Println("upload Traffic err:",err)
        }
		err = syncTraffic(xboardDB, xuiDB);
        if err != nil {
            log.Println("upload Traffic err:",err)
        }
    }
}