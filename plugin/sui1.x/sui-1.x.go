package main

import (
    "io"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"time"
	"sync"
	"strings"
	"errors"
	"crypto/tls"
	"encoding/base64"
)

type Root struct {
    Success bool   `json:"success"`
    Msg     string `json:"msg"`
    Obj     ObjItem `json:"obj"`
}

type ObjItem struct {
    Clients []Client `json:"clients"`
}

type Client struct {
    ID        int   `json:"id"`
    Enable    bool  `json:"enable"`
    Name      string `json:"name"`
    Inbounds  []int `json:"inbounds"`
    UUID      string   `json:"uuid"`
	Down      int64 `json:"down"`
    Up        int64 `json:"up"`
	Expiry    int64 `json:"expiry"`
    TotalUp   int64 `json:"totalUp"`
    TotalDown int64 `json:"totalDown"`
}

type UserPass struct {
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
    Name     string `json:"name,omitempty"`
}

type ShadowTLS struct {
    Name     string `json:"name"`
    Password string `json:"password"`
}

type Vmess struct {
    Name    string `json:"name"`
    UUID    string `json:"uuid"`
    AlterID int    `json:"alterId"`
}

type Vless struct {
    Name string `json:"name"`
    UUID string `json:"uuid"`
    Flow string `json:"flow"`
}

type Tuic struct {
    Name     string `json:"name"`
    UUID     string `json:"uuid"`
    Password string `json:"password"`
}

type Hysteria struct {
    Name    string `json:"name"`
    AuthStr string `json:"auth_str"`
}

type Config struct {
    Mixed         UserPass  `json:"mixed"`
    Socks         UserPass  `json:"socks"`
    HTTP          UserPass  `json:"http"`
    Shadowsocks   UserPass  `json:"shadowsocks"`
    Shadowsocks16 UserPass  `json:"shadowsocks16"`
    ShadowTLS     ShadowTLS `json:"shadowtls"`
    Vmess         Vmess     `json:"vmess"`
    Vless         Vless     `json:"vless"`
    Anytls        UserPass  `json:"anytls"`
    Trojan        UserPass  `json:"trojan"`
    Naive         UserPass  `json:"naive"`
    Hysteria      Hysteria  `json:"hysteria"`
    Tuic          Tuic      `json:"tuic"`
    Hysteria2     UserPass  `json:"hysteria2"`
}

type Payload struct {
    Enable     bool     `json:"enable"`
    Name       string   `json:"name"`
    Config     Config   `json:"config"`
    Inbounds   []int    `json:"inbounds"`
    Links      []string `json:"links"`
    Volume     int64    `json:"volume"`
    Expiry     int64    `json:"expiry"`
    Up         int64    `json:"up"`
    Down       int64    `json:"down"`
	Desc       string   `json:"desc"`
    Group      string   `json:"group"`
    Remark     string   `json:"remark"`
    DelayStart bool     `json:"delayStart"`
    AutoReset  bool     `json:"autoReset"`
    ResetDays  int      `json:"resetDays"`
    NextReset  int64    `json:"nextReset"`
    TotalUp    int64    `json:"totalUp"`
    TotalDown  int64    `json:"totalDown"`
    CreatedAt  int64    `json:"createdAt"`
    OnlineAt   int64    `json:"onlineAt"`
}

var Version = "sui-1.x"

var (
	APIToken string
	idfind sync.Map
)

func Login(user,passwd,api_token,urlhost,protocol string,skipSSL bool) (bool,error) {
	valid := map[string]struct{}{
		"vless": {},
		"vmess": {},
		"trojan": {},
		"shadowsocks": {},
		"hysteria": {},
		"hysteria2": {},
		"naive": {},
		"anytls": {},
		"shadowtls": {},
		"tuic": {},
	}
	
	if _, ok := valid[protocol]; !ok {
		return false,errors.New("ERR: Protocol not supported: "+protocol)
	}
	
    target := urlhost + "/apiv2/users"

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
	
	request.Header.Set("Token", api_token)
	
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

func AddData(class ClassData, urlhost, protocol string, group int, skipSSL bool) error {
	UserID:="xboard_sui_"+class.SubID+"_"+strconv.Itoa(group)
	userData := Config{
            Mixed: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Socks: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            HTTP: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Shadowsocks: UserPass{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            Shadowsocks16: UserPass{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            ShadowTLS: ShadowTLS{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            Vmess: Vmess{
                Name:    UserID,
                UUID:    class.UUID,
                AlterID: 0,
            },
            Vless: Vless{
                Name: UserID,
                UUID: class.UUID,
                Flow: class.FlowControl,
            },
            Anytls: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
            Trojan: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
            Naive: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Hysteria: Hysteria{
                Name:    UserID,
                AuthStr: class.UUID,
            },
            Tuic: Tuic{
                Name:     UserID,
                UUID:     class.UUID,
                Password: class.UUID,
            },
            Hysteria2: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
    }
    reqData := Payload{
		Enable: true,
        Name:   UserID,
        Config: userData,
        Inbounds: []int{group},
    }

    bodyData, err := json.Marshal(reqData)
    if err != nil {
        return err
    }
	
	form := url.Values{}
	form.Set("object", "clients")
	form.Set("action", "new")
	form.Set("data", string(bodyData))

    apiaddr := urlhost + "/apiv2/save"

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(form.Encode()))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Set("Token", APIToken)

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: skipSSL,
        },
        IdleConnTimeout: 55 * time.Second,
    }

    httpClient := &http.Client{
        Transport: transport,
        Timeout:   45 * time.Second,
    }

    resp, err := httpClient.Do(request)
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

    var jsonResp struct {
        Success bool   `json:"success"`
        Msg     string `json:"msg"`
    }

    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return err
    }

    if !jsonResp.Success {
        if jsonResp.Msg != "" {
            return errors.New("ERR: " + jsonResp.Msg)
        }
        return errors.New("ERR: request failed")
    }

    return nil
}

func DelData(class ClassData, urlhost, protocol string, group int, skipSSL bool) error {
	
	id,err:=getIDbyName(urlhost,class.SubID+"_"+strconv.Itoa(group),skipSSL)
	if err!=nil || id <0 {
		if err==nil {err=errors.New("id not found")}
		return err
	}
    form := url.Values{}
	form.Set("object", "clients")
	form.Set("action", "del")
	form.Set("data",strconv.Itoa(id))

    apiaddr := urlhost + "/apiv2/save"

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(form.Encode()))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Token", APIToken)

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: skipSSL,
        },
        IdleConnTimeout: 55 * time.Second,
    }

    httpClient := &http.Client{
        Transport: transport,
        Timeout:   45 * time.Second,
    }

    resp, err := httpClient.Do(request)
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
        return errors.New("xui resp body len=0")
    }

    var jsonResp struct {
        Success bool   `json:"success"`
        Msg     string `json:"msg"`
    }

    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return err
    }

    if !jsonResp.Success {
        if jsonResp.Msg != "" {
            return errors.New("ERR: " + jsonResp.Msg)
        }
        return errors.New("ERR: request failed")
    }

    return nil
}

func ListData(urlhost,protocol string,group int,skipSSL bool) ([]ClassData,error) {
    apiaddr:=urlhost+"/apiv2/load"
	
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
	
	request.Header.Set("Token", APIToken)

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
	for _, client := range jsonResp.Obj.Clients {
			if len(client.Inbounds)>0 && client.Inbounds[0]==group && strings.HasPrefix(client.Name,"xboard_sui_") {
				client.UUID=strings.TrimPrefix(strings.TrimSuffix(client.Name,"_"+strconv.Itoa(group)),"xboard_sui_")
				newclass=append(newclass,ClassData{
					UUID: client.UUID,
					Upflow: client.Up,
					Downflow: client.Down,
					GroupID: group,
					Email: client.Name,
					SubID: client.UUID,
					ExpiredAt: client.Expiry,
				})
			}
    }
	return newclass,nil
}

func ChangeData(class ClassData, urlhost, protocol string, group int, skipSSL bool) error {
    	UserID:="xboard_sui_"+class.SubID+"_"+strconv.Itoa(group)
	userData := Config{
            Mixed: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Socks: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            HTTP: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Shadowsocks: UserPass{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            Shadowsocks16: UserPass{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            ShadowTLS: ShadowTLS{
                Name:     UserID,
                Password: base64.StdEncoding.EncodeToString([]byte(class.UUID)),
            },
            Vmess: Vmess{
                Name:    UserID,
                UUID:    class.UUID,
                AlterID: 0,
            },
            Vless: Vless{
                Name: UserID,
                UUID: class.UUID,
                Flow: class.FlowControl,
            },
            Anytls: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
            Trojan: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
            Naive: UserPass{
                Username: UserID,
                Password: class.UUID,
            },
            Hysteria: Hysteria{
                Name:    UserID,
                AuthStr: class.UUID,
            },
            Tuic: Tuic{
                Name:     UserID,
                UUID:     class.UUID,
                Password: class.UUID,
            },
            Hysteria2: UserPass{
                Name:     UserID,
                Password: class.UUID,
            },
    }
    reqData := Payload{
		Enable: true,
        Name:   UserID,
        Config: userData,
        Inbounds: []int{group},
    }

    bodyData, err := json.Marshal(reqData)
    if err != nil {
        return err
    }
	
	form := url.Values{}
	form.Set("object", "clients")
	form.Set("action", "edit")
	form.Set("data", string(bodyData))

    apiaddr := urlhost + "/apiv2/save"

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(form.Encode()))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Set("Token", APIToken)

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: skipSSL,
        },
        IdleConnTimeout: 55 * time.Second,
    }

    httpClient := &http.Client{
        Transport: transport,
        Timeout:   45 * time.Second,
    }

    resp, err := httpClient.Do(request)
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

    var jsonResp struct {
        Success bool   `json:"success"`
        Msg     string `json:"msg"`
    }

    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return err
    }

    if !jsonResp.Success {
        if jsonResp.Msg != "" {
            return errors.New("ERR: " + jsonResp.Msg)
        }
        return errors.New("ERR: request failed")
    }

    return nil
}

func getIDbyName(urlhost,name string,skipSSL bool) (int,error) {
	val,ok:=idfind.Load(name)
	if ok {
		idfind.Delete(name)
		return val.(int),nil
	}
	apiaddr:=urlhost+"/apiv2/load"
	
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
        return -1,err
    }
	
	request.Header.Set("Token", APIToken)

    resp, err := client.Do(request)
    if err != nil {
        return -1,err
    }
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return -1,errors.New("Http StatusCode not resp OK");
	}
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return -1,err
    }
	var jsonResp Root
    if err := json.Unmarshal(body, &jsonResp); err != nil {
        return -1,err
    }
	if !jsonResp.Success {
		if jsonResp.Msg != "" {
			return -1,errors.New("ERR: " + jsonResp.Msg)
		} else {
			return -1,errors.New("ERR: xui return database is locked");
		}
    }
	
	for _, client := range jsonResp.Obj.Clients {
		if strings.Contains(client.Name, "_"+name) {
			return client.ID,nil
		}
    }
	return -1,nil
}