package main

import (
    "io"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"time"
	"strings"
	"errors"
	"crypto/tls"
)

type Client struct {
	ID        string `json:"id"`
	Flow      string `json:"flow"`
	Email     string `json:"email"`
	Passwd    string `json:"password"`
	Auth      string `json:"auth"`
	Method    string `json:"method"`
	LimitIP   int    `json:"limitIp"`
	TotalGB   int64  `json:"totalGB"`
	ExpiryTime int64 `json:"expiryTime"`
	Enable    bool   `json:"enable"`
	TgID      string `json:"tgId"`
	SubID     string `json:"subId"`
	Comment   string `json:"comment"`
	Reset     int    `json:"reset"`
}

type Payload struct {
	Clients []Client `json:"clients"`
}

type ClientStat struct {
	ID         int    `json:"id"`
	InboundID  int    `json:"inboundId"`
	Enable     bool   `json:"enable"`
	Email      string `json:"email"`
	UUID       string `json:"uuid"`
	SubID      string `json:"subId"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	AllTime    int64  `json:"allTime"`
	ExpiryTime int64  `json:"expiryTime"`
	Total      int64  `json:"total"`
	Reset      int    `json:"reset"`
	LastOnline int64  `json:"lastOnline"`
}

type ObjItem struct {
	ID          int          `json:"id"`
	ClientStats []ClientStat `json:"clientStats"`
}

type Root struct {
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
	Obj     []ObjItem `json:"obj"`
}

var Version = "3xui-3.0.x"

var (
	APIToken string
)

func Login(user,passwd,api_token,urlhost,protocol string,skipSSL bool) (bool,error) {
	valid := map[string]struct{}{
		"vless": {},
		"vmess": {},
		"trojan": {},
		"shadowsocks": {},
		"hysteria": {},
	}
	
	if _, ok := valid[protocol]; !ok {
		return false,errors.New("ERR: Protocol not supported: "+protocol)
	}
	
    target := urlhost + "/panel/api/server/status"

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
		IdleConnTimeout:     50 * time.Second,
	}
	
	client := &http.Client{
			Transport: transport,
			Timeout: time.Second*40,
	}
	
	request, err := http.NewRequest("GET", target, nil)
	
	if err != nil {
		return false,err;
	}
	
	request.Header.Set("Authorization", "Bearer "+api_token)
	
	resp, err := client.Do(request)
	if err != nil {
		return false,err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return false,errors.New("connect xui err, Urlpath or api-token err")
	}
	
	type Response struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}
	
	var jsonResp Response
	
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
    if err != nil {
        return false,err
    }
	
	if !jsonResp.Success{
		return false,errors.New("ERR: " + jsonResp.Msg)
	}
	APIToken=api_token
	return false,nil //loop,err
}

func AddData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error {
	var userData string
    if protocol == "vless" {
		p := Payload{
		Clients: []Client{
			{
				ID:         class.UUID,
				Flow:       class.FlowControl,
				Email:      "xboard_"+class.UUID+strconv.Itoa(group),
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       "",
				SubID:      class.SubID,
				Comment:    "",
				Reset:      0,
			},
		},
	}
	out, err := json.Marshal(p)
	if err != nil {
		return err
	}
	userData=string(out)
	} else if protocol == "vmess" {
		p := Payload{
		Clients: []Client{
			{
				ID:         class.UUID,
				Flow:       class.FlowControl,
				Email:      "xboard__"+class.UUID+strconv.Itoa(group),
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       "",
				SubID:      class.SubID,
				Comment:    "",
				Reset:      0,
			},
		},
	}
	out, err := json.Marshal(p)
	if err != nil {
		return err
	}
	userData=string(out)
	} else if protocol == "trojan" {
		p := Payload{
		Clients: []Client{
			{
				Passwd:     class.UUID,
				Flow:       "",
				Email:      "xboard_tj_"+class.UUID+strconv.Itoa(group),
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       "",
				SubID:      class.SubID,
				Comment:    "",
				Reset:      0,
			},
		},
	}
	out, err := json.Marshal(p)
	if err != nil {
		return err
	}
	userData=string(out)
	} else if protocol == "shadowsocks" {
		p := Payload{
		Clients: []Client{
			{
				Passwd:     class.UUID,
				Method:     "",
				Email:      "xboard_ss_"+class.UUID+strconv.Itoa(group),
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       "",
				SubID:      class.SubID,
				Comment:    "",
				Reset:      0,
			},
		},
	}
	out, err := json.Marshal(p)
	if err != nil {
		return err
	}
	userData=string(out)
	} else if protocol == "hysteria" {
		p := Payload{
		Clients: []Client{
			{
				Auth:       class.UUID,
				Email:      "xboard_hy2_"+class.UUID+strconv.Itoa(group),
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       "",
				SubID:      class.SubID,
				Comment:    "",
				Reset:      0,
			},
		},
	}
	out, err := json.Marshal(p)
	if err != nil {
		return err
	}
	userData=string(out)
	}
	
	form := url.Values{}
	form.Set("id", strconv.Itoa(group))
	form.Set("settings", userData)
	apiaddr:=urlhost+"/panel/api/inbounds/addClient"
	request, err := http.NewRequest("POST", apiaddr, strings.NewReader(form.Encode()))
	if err != nil {
        return err;
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+APIToken)
	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
		IdleConnTimeout:     55 * time.Second,
	}
		
	client := &http.Client{
			Transport: transport,
			Timeout: time.Second*45,
	}

    resp, err := client.Do(request)
    if err != nil {
        return err
    }
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return errors.New("Http StatusCode not resp OK");
	}
	
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
	if len(body) == 0 {
		return errors.New("xui resp body len=0");
	}
	var jsonResp Root
    if err := json.Unmarshal(body, &jsonResp); err != nil {
		return err
    }
	if !jsonResp.Success {
		if jsonResp.Msg != "" {
			return errors.New("ERR: " + jsonResp.Msg)
		} else {
			return errors.New("ERR: xui return database is locked");
		}
    }
	return nil
}

func DelData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error {
	userinfo:=""
    if protocol == "shadowsocks" {
		userinfo="xboard_ss_"+class.UUID+strconv.Itoa(group)
	}else{
		userinfo=class.UUID
	}
	
	apiaddr:=urlhost+"/panel/api/inbounds/"+strconv.Itoa(group)+"/delClient/"+userinfo 

	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
		IdleConnTimeout:     55 * time.Second,
	}
		
	client := &http.Client{
			Transport: transport,
			Timeout: time.Second*45,
	}
		
    request, err := http.NewRequest("POST", apiaddr, nil)
    if err != nil {
        return err;
    }
	
	request.Header.Set("Authorization", "Bearer "+APIToken)
	
    resp, err := client.Do(request)
    if err != nil {
        return err
    }
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return errors.New("del user fail, StatusCode!=200")
	}
	
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
	
	if len(body) == 0 {
		return errors.New("xui resp body len=0");
	}
	
	var jsonResp Root
    if err := json.Unmarshal(body, &jsonResp); err != nil {
		return err
    }
	if !jsonResp.Success {
		if jsonResp.Msg != "" {
			return errors.New("ERR: " + jsonResp.Msg)
		} else {
			return errors.New("ERR: xui return database is locked");
		}
    }
	return nil
}


func ListData(urlhost,protocol string,group int,skipSSL bool) ([]ClassData,error) {
    apiaddr:=urlhost+"/panel/api/inbounds/list"
	
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
		IdleConnTimeout:     56 * time.Second,
	}
		
	client := &http.Client{
			Transport: transport,
			Timeout: time.Second*40,
	}
		
    request, err := http.NewRequest("GET", apiaddr, nil)
    if err != nil {
        return nil,err
    }
	
	request.Header.Set("Authorization", "Bearer "+APIToken)

    resp, err := client.Do(request)
    if err != nil {
        return nil,err
    }
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil,errors.New("Http StatusCode not resp OK");
	}
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil,err
    }
	var jsonResp Root
    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return nil,err
    }
	if !jsonResp.Success {
		if jsonResp.Msg != "" {
			return nil,errors.New("ERR: " + jsonResp.Msg)
		} else {
			return nil,errors.New("ERR: xui return database is locked");
		}
    }
	
	var newclass []ClassData
	for _, obj := range jsonResp.Obj {
        if obj.ID == group {
           for _, client := range obj.ClientStats {
			newclass=append(newclass,ClassData{
				UUID: client.UUID,
				Upflow: client.Up,
				Downflow: client.Down,
				GroupID: obj.ID,
				Email: client.Email,
				SubID: client.SubID,
				ExpiredAt: client.ExpiryTime,
			})
           }
			break
        }
    }
	return newclass,nil
}

func ChangeData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error {
    var (
        userData string
        clientID string
    )

    if protocol == "vless" {
        p := Payload{
            Clients: []Client{
                {
                    ID:         class.UUID,
                    Flow:       class.FlowControl,
                    Email:      "xboard_" + class.UUID + strconv.Itoa(group),
                    LimitIP:    0,
                    TotalGB:    0,
                    ExpiryTime: class.ExpiredAt,
                    Enable:     true,
                    TgID:       "",
                    SubID:      class.SubID,
                    Comment:    "",
                    Reset:      0,
                },
            },
        }

        out, err := json.Marshal(p)
        if err != nil {
            return err
        }

        userData = string(out)
        clientID = class.UUID

    } else if protocol == "vmess" {
        p := Payload{
            Clients: []Client{
                {
                    ID:         class.UUID,
                    Flow:       class.FlowControl,
                    Email:      "xboard__" + class.UUID + strconv.Itoa(group),
                    LimitIP:    0,
                    TotalGB:    0,
                    ExpiryTime: class.ExpiredAt,
                    Enable:     true,
                    TgID:       "",
                    SubID:      class.SubID,
                    Comment:    "",
                    Reset:      0,
                },
            },
        }

        out, err := json.Marshal(p)
        if err != nil {
            return err
        }

        userData = string(out)
        clientID = class.UUID

    } else if protocol == "trojan" {

        p := Payload{
            Clients: []Client{
                {
                    Passwd:     class.UUID,
                    Flow:       "",
                    Email:      "xboard_tj_" + class.UUID + strconv.Itoa(group),
                    LimitIP:    0,
                    TotalGB:    0,
                    ExpiryTime: class.ExpiredAt,
                    Enable:     true,
                    TgID:       "",
                    SubID:      class.SubID,
                    Comment:    "",
                    Reset:      0,
                },
            },
        }

        out, err := json.Marshal(p)
        if err != nil {
            return err
        }

        userData = string(out)
        clientID = class.UUID

    } else if protocol == "shadowsocks" {

        p := Payload{
            Clients: []Client{
                {
                    Passwd:     class.UUID,
                    Method:     "",
                    Email:      "xboard_ss_" + class.UUID + strconv.Itoa(group),
                    LimitIP:    0,
                    TotalGB:    0,
                    ExpiryTime: class.ExpiredAt,
                    Enable:     true,
                    TgID:       "",
                    SubID:      class.SubID,
                    Comment:    "",
                    Reset:      0,
                },
            },
        }

        out, err := json.Marshal(p)
        if err != nil {
            return err
        }

        userData = string(out)
        clientID = class.UUID

    } else if protocol == "hysteria" {

        p := Payload{
            Clients: []Client{
                {
                    Auth:       class.UUID,
                    Email:      "xboard_hy2_" + class.UUID + strconv.Itoa(group),
                    LimitIP:    0,
                    TotalGB:    0,
                    ExpiryTime: class.ExpiredAt,
                    Enable:     true,
                    TgID:       "",
                    SubID:      class.SubID,
                    Comment:    "",
                    Reset:      0,
                },
            },
        }

        out, err := json.Marshal(p)
        if err != nil {
            return err
        }

        userData = string(out)
        clientID = class.UUID

    } else {
        return errors.New("unsupported protocol")
    }

    form := url.Values{}
    form.Set("id", strconv.Itoa(group))
    form.Set("settings", userData)

    apiaddr := urlhost + "/panel/api/inbounds/updateClient/" + clientID

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(form.Encode()))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Set("Authorization", "Bearer "+APIToken)

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: skipSSL,
        },
        IdleConnTimeout: 55 * time.Second,
    }

    client := &http.Client{
        Transport: transport,
        Timeout:   45 * time.Second,
    }

    resp, err := client.Do(request)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return errors.New("Http StatusCode not resp OK")
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    if len(body) == 0 {
        return errors.New("xui resp body len=0")
    }

    var jsonResp Root

    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return err
    }

    if !jsonResp.Success {
        if jsonResp.Msg != "" {
            return errors.New("ERR: " + jsonResp.Msg)
        }
        return errors.New("ERR: xui return database is locked")
    }

    return nil
}