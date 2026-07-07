package main

import (
    "encoding/json"
    "os"
	"sync"
	"errors"
)

type Request struct {
    ID     int             `json:"id"`
    Method string          `json:"method"`
    Data   json.RawMessage `json:"data"`
}

type Response struct {
    ID     int         `json:"id"`
    Result interface{} `json:"result"`
    Error  string      `json:"error"`
}

type ClassData struct {
    UUID      string `json:"uuid"`
	FlowControl string `json:"flowcontrol"`
    Upflow    int64  `json:"upflow"`
    Downflow  int64  `json:"downflow"`
	GroupID   int  `json:"groupid"`
    Email     string `json:"email"`
    SubID     string `json:"subid"`
    ExpiredAt int64  `json:"expired_at"`
}

var lock sync.Mutex

func main() {
	if len(os.Args) <=1 {
		return
	}
	if os.Args[1]!="--rpc=json" {
		os.Stderr.Write([]byte("ERR: rpc call type not found\n"))
		return
	}
    decoder := json.NewDecoder(os.Stdin)
    encoder := json.NewEncoder(os.Stdout)

    for {
		var req Request
        if err := decoder.Decode(&req); err != nil {
			os.Stderr.Write([]byte("plugin exit err:"+err.Error()+"\n"))
            return
        }

        switch req.Method {

        case "Login","AddData","DelData","ListData","ChangeData":
		req_data:=req.Data
		var resp Response
		resp.ID = req.ID
		go func(){
            r, err := rpcCall(req.Method,req_data)
            if err != nil {
                resp.Error = err.Error()
            } else {
                resp.Result = r
            }
			lock.Lock()
			encoder.Encode(resp)
			lock.Unlock()
		}()
        case "PING":
		    var resp Response
			resp.ID = req.ID
            resp.Result = "PONG"
			lock.Lock()
			encoder.Encode(resp)
			lock.Unlock()
		default:
		    var resp Response
			resp.ID = req.ID
            resp.Error = "unknown method"
			lock.Lock()
			encoder.Encode(resp)
			lock.Unlock()
        }
    }
}

func rpcCall(rpcMethod string,jdata json.RawMessage) (interface{},error) {
switch rpcMethod {
    case "Login":
	type input_args struct {
		User string `json:"user"`
		Passwd string `json:"passwd"`
		Api_token string `json:"api_token"`
		Urlhost string `json:"urlhost"`
		Protocol string `json:"protocol"`
		SkipSSL bool `json:"skipSSL"`
	}
	
	var args input_args
	err:=json.Unmarshal(jdata, &args)
	if err!=nil {
		return false,err
	}
	return Login(args.User,args.Passwd,args.Api_token,args.Urlhost,args.Protocol,args.SkipSSL)

	case "AddData":
	type input_args struct {
		Class ClassData `json:"class"`
		Protocol string `json:"protocol"`
		Group int `json:"group"`
		Urlhost string `json:"urlhost"`
		SkipSSL bool `json:"skipSSL"`
	}
	
	var args input_args
	err:=json.Unmarshal(jdata, &args)
	if err!=nil {
		return nil,err
	}
	err=AddData(args.Class,args.Urlhost,args.Protocol,args.Group,args.SkipSSL)
	return nil,err
	
	case "DelData":
	type input_args struct {
		Class ClassData `json:"class"`
		Protocol string `json:"protocol"`
		Group int `json:"group"`
		Urlhost string `json:"urlhost"`
		SkipSSL bool `json:"skipSSL"`
	}
	
	var args input_args
	err:=json.Unmarshal(jdata, &args)
	if err!=nil {
		return nil,err
	}
	err=DelData(args.Class,args.Urlhost,args.Protocol,args.Group,args.SkipSSL)
	return nil,err
	
	case "ListData":
	type input_args struct {
		Protocol string `json:"protocol"`
		Group int `json:"group"`
		Urlhost string `json:"urlhost"`
		SkipSSL bool `json:"skipSSL"`
	}
	
	var args input_args
	err:=json.Unmarshal(jdata, &args)
	if err!=nil {
		return nil,err
	}
	class,err:=ListData(args.Urlhost,args.Protocol,args.Group,args.SkipSSL)
	if err!=nil {
		return nil,err
	}
	b,err:=json.Marshal(class)
	if err!=nil {
		return nil,err
	}
	return json.RawMessage(b),nil
	
	case "ChangeData":
	type input_args struct {
		Class ClassData `json:"class"`
		Protocol string `json:"protocol"`
		Group int `json:"group"`
		Urlhost string `json:"urlhost"`
		SkipSSL bool `json:"skipSSL"`
	}
	
	var args input_args
	err:=json.Unmarshal(jdata, &args)
	if err!=nil {
		return nil,err
	}
	err=ChangeData(args.Class,args.Urlhost,args.Protocol,args.Group,args.SkipSSL)
	return nil,err
	
default:
	return nil,errors.New("rpc method not found")
}
}