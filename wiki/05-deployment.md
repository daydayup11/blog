# 05 部署指南

## 两种部署方式对比

| | Docker 部署（推荐） | 直接运行（开发者）|
|---|---|---|
| 适合 | 所有人，无需懂 Go | 开发者，需要 Go 环境 |
| 依赖 | 只需 Docker | Go 1.22+ |
| 启动命令 | `./deploy.sh` | `./start.sh` |
| 数据位置 | `./data/`（挂载到容器）| `./data/` |
| 更新方式 | `./deploy.sh update` | `git pull && ./start.sh restart` |

---

## 方式一：Docker 部署（推荐）

### 前提

安装 Docker：
- macOS：`brew install --cask docker` 或 https://docs.docker.com/get-docker/
- Linux：`curl -fsSL https://get.docker.com | sh`
- Windows：https://docs.docker.com/desktop/install/windows-install/

### 步骤

**1. 克隆项目**

```bash
git clone <repo-url>
cd blog
```

**2. 配置环境变量**

```bash
cp .env.example .env
```

用编辑器打开 `.env`，**必须修改**以下两项：

```bash
# 后台登录密码（改成你自己的）
ADMIN_PASS=your-strong-password

# JWT 密钥（随机字符串，运行下面命令生成）
# openssl rand -hex 32
JWT_SECRET=粘贴上面命令的输出
```

**3. 一键启动**

```bash
./deploy.sh
```

首次运行会构建 Docker 镜像，约需 3-5 分钟。启动后自动打开浏览器。

访问地址：http://localhost:8080
后台地址：http://localhost:8080/admin/index.html

### 常用命令

```bash
./deploy.sh              # 启动 / 重新部署
./deploy.sh stop         # 停止服务
./deploy.sh restart      # 重启服务
./deploy.sh status       # 查看运行状态
./deploy.sh logs         # 实时查看日志（Ctrl+C 退出）
./deploy.sh update       # 拉取最新版本并重启
```

### 修改端口

在 `.env` 里改：

```bash
PORT=9000
```

然后 `./deploy.sh restart`。

---

## 方式二：直接运行（开发者）

### 前提

安装 Go 1.22+：
- macOS：`brew install go`
- 其他：https://go.dev/dl/

### 步骤

```bash
git clone <repo-url>
cd blog
./start.sh
```

脚本会自动编译后端并启动，首次约需 1-2 分钟。

### 常用命令

```bash
./start.sh               # 启动
./start.sh stop          # 停止
./start.sh restart       # 重启
./start.sh status        # 查看运行状态（显示 PID）
./start.sh logs          # 实时查看日志
```

---

## 服务器长期部署

### 推荐方案：VPS + Docker + Nginx 反向代理

适合有自己域名和服务器（阿里云/腾讯云/Vultr 等）的情况。

**1. 服务器上安装 Docker**

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

**2. 克隆项目并配置**

```bash
git clone <repo-url>
cd blog
cp .env.example .env
# 编辑 .env 设置密码和密钥
./deploy.sh
```

**3. 配置 Nginx 反向代理**

安装 Nginx：

```bash
sudo apt install nginx
```

创建配置文件 `/etc/nginx/sites-available/blog`：

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

    # 上传文件大小限制（默认 1MB，按需调整）
    client_max_body_size 20M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

**4. 申请 HTTPS 证书（免费）**

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

证书自动续期，无需手动操作。

---

## 版本更新

### 使用 Docker Hub 镜像更新（发布版）

项目有版本发布时，在 `.env` 里指定版本：

```bash
# .env
DOCKER_IMAGE=daiyutong/blog:v1.1.0
```

然后更新：

```bash
./deploy.sh update
```

### 从源码更新（开发者）

```bash
git pull
./deploy.sh        # Docker 版重新构建
# 或
./start.sh restart # 直接运行版重新编译
```

---

## 发布新版本（项目维护者）

项目通过 Git tag 触发自动构建，推送 tag 后 GitHub Actions 会：
1. 构建 `linux/amd64` 和 `linux/arm64` 双架构镜像
2. 推送到 Docker Hub（`:v1.0.0` + `:latest`）
3. 在 GitHub 创建 Release 并附上部署说明

```bash
# 发布 v1.0.0
git tag v1.0.0
git push origin v1.0.0
```

**前提**：需要在 GitHub 仓库 Settings → Secrets → Actions 中配置：
- `DOCKERHUB_USERNAME`：Docker Hub 用户名
- `DOCKERHUB_TOKEN`：Docker Hub Access Token（在 Docker Hub → Account Settings → Security 生成）

---

## 数据备份

数据库是单个文件，备份非常简单：

```bash
# 手动备份
cp data/blog.db data/blog.db.backup

# 定时备份（Linux crontab）
# 每天凌晨 2 点备份，保留最近 7 天
0 2 * * * cd /path/to/blog && cp data/blog.db data/blog-$(date +\%Y\%m\%d).db
0 2 * * * find /path/to/blog/data -name 'blog-*.db' -mtime +7 -delete
```

上传的图片在 `data/uploads/` 目录，与数据库一并备份：

```bash
tar -czf backup-$(date +%Y%m%d).tar.gz data/
```

---

## 内网穿透（临时分享）

适合临时演示，无需服务器：

```bash
# 安装（macOS）
brew install cloudflared

# 确保服务在运行
./start.sh  # 或 ./deploy.sh

# 启动隧道
cloudflared tunnel --url http://localhost:8080
```

输出中的 `https://xxx.trycloudflare.com` 即为公网地址，发给对方即可访问。

> 注意：每次启动域名随机，进程退出即失效。固定域名需注册 Cloudflare 账号，参考 Cloudflare Tunnel 文档。

---

## 常见问题

**Q：端口被占用怎么办？**
```bash
# 查看占用端口的进程
lsof -i:8080
# 或修改 .env 里的 PORT 换一个端口
```

**Q：忘记后台密码怎么办？**

修改 `.env` 里的 `ADMIN_PASS`，然后重启：
```bash
./deploy.sh restart  # 或 ./start.sh restart
```

**Q：Docker 构建很慢？**

首次构建需要下载 Go 编译器和依赖，约 3-5 分钟属于正常。后续更新会利用缓存，速度快很多。

**Q：数据在哪里，怎么迁移到新服务器？**

所有数据在 `data/` 目录：
```bash
# 打包
tar -czf blog-data.tar.gz data/
# 传到新服务器后解压即可
tar -xzf blog-data.tar.gz
```
