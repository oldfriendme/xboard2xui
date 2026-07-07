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
	TgID      int `json:"tgId"`
	SubID     string `json:"subId"`
	Comment   string `json:"comment"`
	Reset     int    `json:"reset"`
	Security  string `json:"security"`
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

var Version = "3xui-3.1.x"

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

func AddData(class ClassData, urlhost, protocol string, group int, skipSSL bool) error {
    type AddClientRequest struct {
        Client     Client `json:"client"`
        InboundIDs []int  `json:"inboundIds"`
    }

    email := ""
    client := Client{
        Flow:       class.FlowControl,
        SubID:      class.SubID+"_"+strconv.Itoa(group),
        TotalGB:    0,
        ExpiryTime: class.ExpiredAt,
        Reset:      0,
        LimitIP:    0,
        Enable:     true,
        Comment:    "",
    }

    switch protocol {
    case "vless", "vmess":
        email = "xboard_" + class.UUID + strconv.Itoa(group)
        client.ID = class.UUID
        client.Email = email
        client.Security = "auto"

    case "trojan":
        email = "xboard_tj_" + class.UUID + strconv.Itoa(group)
        client.Passwd = class.UUID
        client.Email = email

    case "shadowsocks":
        email = "xboard_ss_" + class.UUID + strconv.Itoa(group)
        client.Passwd = class.UUID
        client.Email = email
        client.Method = "aes-128-gcm"

    case "hysteria", "hysteria2":
        email = "xboard_hy2_" + class.UUID + strconv.Itoa(group)
        client.Auth = class.UUID
        client.Email = email

    default:
        return errors.New("unsupported protocol")
    }

    reqData := AddClientRequest{
        Client:     client,
        InboundIDs: []int{group},
    }

    bodyData, err := json.Marshal(reqData)
    if err != nil {
        return err
    }

    apiaddr := urlhost + "/panel/api/clients/add"

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(string(bodyData)))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", "Bearer "+APIToken)

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
    var email string

    switch protocol {
    case "vless", "vmess":
        email = "xboard_" + class.UUID + strconv.Itoa(group)

    case "trojan":
        email = "xboard_tj_" + class.UUID + strconv.Itoa(group)

    case "shadowsocks":
        email = "xboard_ss_" + class.UUID + strconv.Itoa(group)

    case "hysteria", "hysteria2":
        email = "xboard_hy2_" + class.UUID + strconv.Itoa(group)

    default:
        return errors.New("unsupported protocol")
    }

    apiaddr := urlhost + "/panel/api/clients/del/" + url.PathEscape(email)

    request, err := http.NewRequest("POST", apiaddr, nil)
    if err != nil {
        return err
    }

    request.Header.Set("Authorization", "Bearer "+APIToken)

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
				SubID: strings.TrimSuffix(client.SubID,"_"+strconv.Itoa(group)),
				ExpiredAt: client.ExpiryTime,
			})
           }
			break
        }
    }
	return newclass,nil
}

func ChangeData(class ClassData, urlhost, protocol string, group int, skipSSL bool) error {
    var email string

    client := Client{
        Flow:       class.FlowControl,
        SubID:      class.SubID+"_"+strconv.Itoa(group),
        TotalGB:    0,
        ExpiryTime: class.ExpiredAt,
        Reset:      0,
        LimitIP:    0,
        Enable:     true,
        Comment:    "",
    }

    switch protocol {
    case "vless", "vmess":
        email = "xboard_" + class.UUID + strconv.Itoa(group)
        client.Email = email
        client.ID = class.UUID
        client.Security = "auto"

    case "trojan":
        email = "xboard_tj_" + class.UUID + strconv.Itoa(group)
        client.Email = email
        client.Passwd = class.UUID

    case "shadowsocks":
        email = "xboard_ss_" + class.UUID + strconv.Itoa(group)
        client.Email = email
        client.Passwd = class.UUID
        client.Method = "aes-128-gcm"

    case "hysteria", "hysteria2":
        email = "xboard_hy2_" + class.UUID + strconv.Itoa(group)
        client.Email = email
        client.Auth = class.UUID

    default:
        return errors.New("unsupported protocol")
    }

    bodyData, err := json.Marshal(client)
    if err != nil {
        return err
    }

    apiaddr := urlhost + "/panel/api/clients/update/" + url.PathEscape(email)

    request, err := http.NewRequest("POST", apiaddr, strings.NewReader(string(bodyData)))
    if err != nil {
        return err
    }

    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", "Bearer "+APIToken)

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