# xboard2xui plugin

### xboard2xui plugin 插件说明

plugin用于告诉xboard2xui如何操作xui，详细可查看个目录设计

插件功能v0.9以上可用

<br>

### 插件开发指南

插件需要定义4个函数，以告诉程序如何操作xui

分别是 "Login","AddData","DelData","ChangeData","ListData"

配置json示例如下：
```json
	"xui-config": {
		"ApiHost": "http://127.0.0.1:2053",
		"user": "admin",
		"passwd": "admin",
		"api-token":"MAZAAAAAAAAAAAAAAAAAAA",
		"NodeID": 1,
		"admin-path": "",
		"plugin":"./xboard2xui-3xui2.9.x-linux-amd64",
		"plugin_watchdog":false
	}
```

`"plugin":"${plugin_file}"` //插件路径，可使用绝对路径或相对路径

`"plugin_watchdog":false` //是否启动插件的watchdog功能，默认不开启

<br>

#### 注意:

`xboard2xui-3xui2.9.x-*` 此类开头为插件程序，不能单独运行，必须被`"plugin":"${plugin_file}"`加载运行


---

<br>

## 插件函数详情

### Login 函数

```go
func Login(user,passwd,api_token,urlhost,protocol string,skipSSL bool) (bool,error)
```

告诉程序如何登录，程序可选择API-Token还是user:passwd获取session。

登录流程结束后，程序需要自行定义逻辑以持久化API-Token或者session


#### 传入参数

`user,passwd,api_token` 此3个参数为config.json定义的参数直接传递，即：

```json
"user": "admin",
"passwd": "admin",
"api-token":"MAZAAAAAAAAAAAAAAAAAAA",
```

插件可以选择使用user,passwd还是api_token

<br>

#### 返回结果

函数返回 loop,err。第一个loop告诉程序是否定时循环执行Login，第二个表明是否发生错误

<br>

### AddData 函数

```go
func AddData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error
```

定义程序如何添加数据

<br>

### DelData 函数

```go
func DelData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error
```

定义程序如何删除数据

<br>

### ChangeData 函数

```go
ChangeData(class ClassData,urlhost,protocol string,group int,skipSSL bool) error
```

定义程序如何修改数据

<br>

### ListData 函数

```go
func ListData(urlhost,protocol string,group int,skipSSL bool) ([]ClassData,error)
```

定义程序获取数据

<br>

---

## 插件其他事项

- 插件需要chmod+x权限才能被加载进内存。

- 插件设计可自行使用变量保存session数据，在插件运行期间，数据会在内存中持久化。

- 插件不止可以操作xui，如果进行合适的定义，甚至还能操作s-ui，h-ui。甚至xray原版，singbox裸核
