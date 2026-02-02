# xboard2xui
xboard backend with xui

### xui 对接Xboard教程。

#### 下面以Reality为例，其他的我没试，理论上也可以。

<br>

---

#### 1.1.首先xui添加一个普通入站

按照常用的设置就行。下面“安全”选Reality，添加，正常生成配置信息，多余的留空即可。

![](img/addin.jpg?raw=true)

<br>

#### 1.2设置入站选项

![](img/xuiconfig.jpg?raw=true)

![](img/xuiconfig2.jpg?raw=true)

<br>

添加完后，会生成一个ID，这个是xui的NodeID，比如下面的NodeID就是"**1**"

![](img/nodeidxui.jpg?raw=true)

<br>

---

### 2.对接Xboard，

#### 2.1.第一步：
和普通的Xboard对接xrayr流程一样，生成一个密钥。这个就是ApiKey

然后Xboard添加协议，选"vless"，把xui的参数复制到Xboard里面。

![](img/xbconfig.jpg?raw=true)

注意，这里的参数要与之前设置的参数一致

### 节点地址填xui的地址，端口填xui入站的端口。
把`xui-address.vps.domain`换成你实际的IP/域名

<br>

#### 然后把xui生成的公钥，私钥，流控设置，参数。全部复制到Xboard里面。保持一样，这个是必须的。

shortID复制xui里面，复制第一个shortID到Xboard里面。这里是2e9f。参考你实际的shortID。

<details>
<summary>参考之前的</summary>
  
![](img/xuiconfig2.jpg?raw=true)

</details>

<br>

然后保存，记下NodeID。比如我这个是"11"

![](img/nodeidxb.jpg?raw=true)

<br>

---

#### 2.2下面填写配置

格式参考如下
```json
{
	"protocol": "vless",
	"flow-control": "xtls-rprx-vision",
	"xboard-skip-ssl-check": false,
	"xboard-config": {
		"ApiHost": "http://127.0.0.1:7001",
		"ApiKey": "OkkHRnmd9BKKfiIn6xwxyA",
		"NodeID": 11
	},
	"xui-skip-ssl-check": false,
	"xui-config": {
		"ApiHost": "http://127.0.0.1:2053",
		"user": "admin",
		"passwd": "admin",
		"NodeID": 1,
		"admin-path": "/xuipath"
	}
}
```

"ApiHost": "http://127.0.0.1:7001 这个是 xboard的地址，换成你实际地址。区分https/http

"ApiHost": "http://127.0.0.1:2053 这个是 xui的地址，换成你实际地址。区分https/http

两个NodeID填之前上面记下来的。

"flow-control"这个，填实际的流控就行。没有就填：""

"protocol"这个，填协议全称+小写，如："vless"。需与xboard协议定义一致

<br>

"user": "admin",
"passwd": "admin", 
这两个填xui密码。

"admin-path": "/xuipath"，这个填xui path，开头必须为/，结尾不能为/。 比如path为/xuipath/，填/xuipath就行。如果没有，就填"admin-path": ""

此外xui不能开启二步验证，对接api没有二步验证逻辑。

`"xui-skip-ssl-check": `如果证书与域名不匹配（如直接访问IP，而未使用证书域名），可信环境中可以跳过验证。

<br>

然后下载：[xboard2xui](https://github.com/oldfriendme/xboard2xui/releases)

把上面的config.json保存，然后启动对接：`xboard2xui config.json`

然后能看到xboard里面的小圆点，从红色变成黄色。说明节点上线了。

然后，可以打开xui，展开你的xui 对应的 NodeID的详细客户端信息，可以看到多个`xboard_`开头的用户，这就是xboard创建的用户，说明对接成功了。

![](img/xuiend.jpg?raw=true)

<br>

---

<br>

另外，路由以及出站，以及禁止，在xui里面配置，xboard的配置不起作用，会被xui本身覆盖

不过建议软件与xui安装在同一台VPS上，对接xui使用localhost+http协议，对接xboard使用https远程，就像对接xrayr一样的控制行为。

如果节点不通，可以对比下xui与xboard下发配置的差异，从而找出问题。

<br>

---

<br>

#### 常见问题

<details>
	
<summary>xui卡死</summary>

<br>

X-UI卡死问题早已有之，即使单独使用xui也有这种现象，表现症状：

- 1.xui虽然会显示用户在线与离线，但是xui本身的流量统计不再增加。（用户流量使用一直为0，或使用后不增加）
- 2.用户也正常连接，不超时（这个是xray内核负责的，面板卡死不影响内核）

<br>

当xui卡死后，即使有设备在使用，xboard用户也会显示黄色无人使用：![](img/xbbtn.jpg?raw=true)

<br>

#### 排查是否卡死：

打开X-UI日志

![](img/xuilog.jpg?raw=true)


查看是否出现出现大量的

![](img/xuiwarn.jpg?raw=true)


<br>

常见解决办法：

重启面板几乎可以解决99%的卡死问题：（注意：大多数情况只需要重启面板，并不需要重启xray与软件）

![](img/xuireboot.jpg?raw=true)

</details>
