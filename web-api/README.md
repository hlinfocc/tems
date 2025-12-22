tems web-api 后端服务及接口

## 运行项目

1. 确保已安装 Go 1.24 或更高版本
2. 克隆项目仓库：
   ```
   cd web-api
   ```
3. 安装依赖：
   ```
   go mod tidy
   ```
4. 配置环境变量（可选）：

5. 运行项目：
   ```
   go run main.go
   ```

## 配置环境变量

在项目根目录下创建 `.env` 文件，配置环境变量。内容如下：

```
# 数据库连接信息
# 数据库主机地址
DB_HOST=127.0.0.1
# 数据库端口
DB_PORT=5432
# 数据库用户名
DB_USER=postgres
# 数据库密码
DB_PASSWD=qq123456
# 数据库名称
DB_NAME=temsdb

# 服务配置
# 服务监听地址
SERVER_HOST=0.0.0.0
# 服务监听端口
SERVER_PORT=55555
# 图片预览服务URL前缀
SERVER_IMG_PREVIEW_URL=http://127.0.0.1:55555
# 小程序是否开放注册(默认为true)
OPEN_REGISTER=true

# 小程序配置
# 小程序AppID
WECHART_APPID=wxxxxxx
# 小程序AppSecret
WECHART_SECRET=5f56666666766fffgggg6196
```


## 项目目录结构

```
├── web-api/            # 后端API服务
│   ├── config/         # 配置文件
│   ├── controllers/    # 控制器层
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── routes/         # 路由定义
│   ├── upload/         # 上传文件处理
│   ├── utils/          # 工具函数
│   └── main.go         # 主程序入口
```
