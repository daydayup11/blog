# 05 部署指南

## 本地开发

### 前提条件

- Go 1.22 及以上（安装：https://go.dev/dl/）
- Git

### 启动步骤

```bash
# 1. 克隆项目
git clone <repo-url>
cd blog

# 2. 启动后端（会自动创建数据库）
cd backend
go run ./cmd/server
```

访问 http://localhost:8080

后台：http://localhost:8080/admin/index.html（账号 admin / admin123）

### 前端开发模式（无需后端）

如果只改前端，可以不启动 Go 服务，用 Python 起一个静态服务器：

```bash
cd frontend
python3 -m http.server 3000
```

然后访问 http://localhost:3000/blog.html?mock=1（加 `?mock=1` 使用模拟数据）

---

## Docker 部署

### 构建镜像

```bash
cd backend
docker build -t blog-api:latest .
```

> 注意：Dockerfile 使用了 `CGO_ENABLED=1`（SQLite 需要），构建时间稍长（约 2-5 分钟）。

### 运行容器

```bash
docker run -d \
  --name blog \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/../frontend:/app/frontend \
  -e JWT_SECRET=your-secret-key \
  -e ADMIN_PASS=your-password \
  -e FRONTEND_DIR=/app/frontend \
  blog-api:latest
```

**挂载说明：**
- `-v $(pwd)/data:/app/data`：持久化 SQLite 数据库文件
- `-v $(pwd)/../frontend:/app/frontend`：挂载前端文件

### 环境变量说明

| 变量 | 必改 | 说明 |
|---|---|---|
| `JWT_SECRET` | ✅ | JWT 签名密钥，随机字符串即可 |
| `ADMIN_PASS` | ✅ | 后台密码 |
| `ADMIN_USER` | 可选 | 后台用户名，默认 admin |
| `PORT` | 可选 | 监听端口，默认 8080 |
| `DB_PATH` | 可选 | 数据库路径，默认 ./data/blog.db |
| `FRONTEND_DIR` | 可选 | 前端目录，默认 ../frontend |

---

## 内网穿透（临时分享）

适合临时展示给他人，无需服务器。使用 Cloudflare Tunnel（免费）：

```bash
# 安装（macOS）
brew install cloudflared

# 启动隧道（确保后端已在 8080 运行）
cloudflared tunnel --url http://localhost:8080
```

输出中会有一个随机域名，如 `https://xxx.trycloudflare.com`，把这个链接发给对方即可访问。

**注意：**
- 每次启动都会生成新的随机域名
- 进程退出隧道即断开
- 如需固定域名，参考下方"固定域名"章节

---

## 生产服务器部署

### 推荐方案：VPS + Nginx 反向代理

适合有自己域名和服务器的情况。

**1. 服务器上拉代码**

```bash
git clone <repo-url>
cd blog
```

**2. 编译后端**

```bash
cd backend
go build -o ../server ./cmd/server
```

**3. 用 systemd 管理进程（Linux）**

创建 `/etc/systemd/system/blog.service`：

```ini
[Unit]
Description=Blog Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/blog/backend
ExecStart=/path/to/blog/server
Environment=JWT_SECRET=your-secret
Environment=ADMIN_PASS=your-password
Environment=DB_PATH=/path/to/blog/data/blog.db
Environment=FRONTEND_DIR=/path/to/blog/frontend
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable blog
sudo systemctl start blog
```

**4. Nginx 配置**

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate     /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**5. 申请 HTTPS 证书**

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 自动申请并配置
sudo certbot --nginx -d yourdomain.com
```

---

## 固定域名内网穿透

如果想每次都使用同一个域名（需要 Cloudflare 账号 + 自己的域名）：

```bash
# 1. 登录 Cloudflare 账号
cloudflared tunnel login

# 2. 创建命名隧道
cloudflared tunnel create blog

# 3. 创建配置文件 ~/.cloudflared/config.yml
tunnel: <tunnel-id>
credentials-file: ~/.cloudflared/<tunnel-id>.json

ingress:
  - hostname: blog.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404

# 4. 在 Cloudflare DNS 添加 CNAME 记录
cloudflared tunnel route dns blog blog.yourdomain.com

# 5. 启动
cloudflared tunnel run blog
```

之后 `blog.yourdomain.com` 就是固定的公网地址。

---

## 数据备份

数据库是单个文件，备份很简单：

```bash
# 备份
cp data/blog.db data/blog.db.backup

# 或定时备份（crontab）
0 2 * * * cp /path/to/blog/data/blog.db /backup/blog-$(date +%Y%m%d).db
```

上传的图片在 `data/uploads/` 目录，也需要一并备份。
