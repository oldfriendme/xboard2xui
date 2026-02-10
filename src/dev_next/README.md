# dev next

dev nextbeta v1.x计划，添加添加同步数据库测试，解决网络拥塞以及延迟问题导致节点掉线以及同步性能问题。

```json
{
	"protocol": "vless",
	"flow-control": "xtls-rprx-vision",
	"xboard-skip-ssl-check": false,
	"xboard-config": {
		"ApiHost": "http://127.0.0.1:7001",
		"ApiKey": "OkkHRnmd9BKKfiIn6xwxyA",
		"NodeID": 11
		"database": "/mnt/xboard.sqlite" //直接定义数据库，与上述字段互斥
	},
	"xui-skip-ssl-check": false,
	"xui-config": {
		"ApiHost": "http://127.0.0.1:2053",
		"user": "admin",
		"passwd": "admin",
		"NodeID": 1,
		"admin-path": "/xuipath"
		"database": "/mnt/xui.db" //直接定义数据库，与上述字段互斥
	}
}
```