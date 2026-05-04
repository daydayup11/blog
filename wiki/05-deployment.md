# 05 部署指南

## 两种部署方式

| | Docker 部署（推荐） | 直接运行（开发者）|
|---|---|---|
| 适合 | 所有人，无需懂 Go | 开发者，需要 Go 环境 |
| 依赖 | Docker | Go 1.22+ |
| 启动 | `./deploy.sh` | `./start.sh` |
| 配置 | 交互式引导，自动生成密钥 | 环境变量或默认值 |

---

## 方式一：Docker 部署（推荐）

### 安装 Docker

- **macOS**：`brew install --cask docker` 或 https://docs.docker.com/desktop/install/mac-install/
- **Linux**：`curl -fsSL https://get.docker.com | sh`
- **Windows**：https://docs.docker.com/desktop/install/windows-install/

安装 Docker Compose（如果未附带）：

```bash
brew install docker-compose   # macOS
```

### 一键部署

```bash
git clone <repo-url>
cd blog
./deploy.sh
```

脚本会全程引导，**无需手动编辑任何文件**：

1. 自动检测 Docker 是否安装和运行，未就绪时给出提示
2. 询问端口、用户名、密码（全部有默认值，直接回车跳过）
3. 密码留空时自动生成 16 位强密码并显示，**请记下**
4. JWT 密钥完全自动生成，无需关心
5. 构建镜像、启动服务，完成后显示访问地址

```
[?] 服务端口 [默认: 8080]：          ← 直接回车
[?] 后台管理用户名 [默认: admin]：    ← 直接回车
[?] 后台管理密码 [直接 Enter 自动生成强密码]：  ← 直接回车
  已自动生成密码：Kuqk4x6KpIEf84k4   ← 记下这个密码
```

部署完成后自动打开浏览器（macOS）。

### 常用命令

```bash
./deploy.sh              # 首次部署 / 重新部署
./deploy.sh stop         # 停止服务
./deploy.sh restart      # 重启服务
./deploy.sh logs         # 实时查看日志（Ctrl+C 退出）
./deploy.sh status       # 查看运行状态
./deploy.sh update       # 拉取最新版本并重启
```

### 重新部署

再次运行 `./deploy.sh` 时，脚本会检测到已有配置并询问是否复用，无需重新填写。

### 修改配置

直接编辑 `.env` 文件，然后重启：

```bash
nano .env
./deploy.sh restart
```

---

## 方式二：直接运行（开发者）

**依赖：** Go 1.22+（[下载](https://go.dev/dl/)，macOS：`brew install go`）

```bash
git clone <repo-url>
cd blog
./start.sh
```

脚本自动编译后端并后台运行，首次约需 1-2 分钟。

### 常用命令

```bash
./start.sh               # 启动（自动编译）
./start.sh stop          # 停止
./start.sh restart       # 重启
./start.sh logs          # 实时查看日志
./start.sh status        # 查看运行状态（显示 PID）
```

---

## 访问地址

| 地址 | 说明 |
|---|---|
| `http://localhost:8080` | 前台网站 |
| `http://localhost:8080/admin` | 后台管理 |

---

## 服务器长期部署

### VPS + Nginx + HTTPS

适合有自己域名和服务器的情况。

**1. 服务器上部署**

```bash
git clone <repo-url>
cd blog
./deploy.sh
```

**2. 安装 Nginx**

```bash
# Ubuntu/Debian
sudo apt install nginx
```

**3. 配置反向代理**

创建 `/etc/nginx/sites-available/blog`：

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

    client_max_body_size 20M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

**4. 申请免费 HTTPS 证书**

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

证书自动续期，无需手动操作。

---

## 内网穿透（临时分享）

无需服务器，适合临时演示给他人：

```bash
# 安装（macOS）
brew install cloudflared

# 启动隧道（确保服务已运行）
cloudflared tunnel --url http://localhost:8080
```

输出的 `https://xxx.trycloudflare.com` 即为公网地址，发给对方即可访问。

> 注意：每次域名随机，进程退出即失效。

---

## 版本更新

### 更新已有部署

```bash
./deploy.sh update
```

脚本会拉取最新镜像（或重新构建）并重启服务，数据不受影响。

### 发布新版本（维护者）

通过 Git tag 触发 GitHub Actions 自动构建并推送到 Docker Hub：

```bash
git tag v1.0.0
git push origin v1.0.0
```

需要在 GitHub 仓库 Settings → Secrets → Actions 中配置：
- `DOCKERHUB_USERNAME`：Docker Hub 用户名
- `DOCKERHUB_TOKEN`：Docker Hub Access Token（在 Docker Hub → Account Settings → Security 生成）

---

## 数据备份

所有数据在 `data/` 目录（SQLite 数据库 + 上传图片），备份简单：

```bash
# 手动备份
tar -czf backup-$(date +%Y%m%d).tar.gz data/

# 定时备份（Linux crontab，每天凌晨 2 点，保留 7 天）
0 2 * * * cd /path/to/blog && tar -czf data/backup-$(date +\%Y\%m\%d).tar.gz data/blog.db data/uploads/
0 2 * * * find /path/to/blog/data -name 'backup-*.tar.gz' -mtime +7 -delete
```

迁移到新服务器：

```bash
# 旧服务器打包
tar -czf blog-data.tar.gz data/

# 新服务器解压后重新部署
tar -xzf blog-data.tar.gz
./deploy.sh
```

---

## 常见问题

**Q：端口被占用？**
```bash
lsof -i:8080          # 查看占用进程
# 或在 .env 里改 PORT=9000，然后 ./deploy.sh restart
```

**Q：忘记后台密码？**

编辑 `.env` 改 `ADMIN_PASS`，然后重启：
```bash
./deploy.sh restart
```

**Q：Docker 构建很慢？**

首次构建需下载 Go 编译环境和依赖，约 3-5 分钟正常。后续更新会利用缓存，速度快很多。

**Q：数据会丢失吗？**

数据存在宿主机 `data/` 目录并挂载到容器，更新或重启不影响数据。
