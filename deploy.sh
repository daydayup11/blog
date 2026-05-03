#!/bin/bash
# 一键部署脚本 — 个人博客（Docker 版）
# 用法:
#   ./deploy.sh           首次部署 / 更新
#   ./deploy.sh stop      停止服务
#   ./deploy.sh restart   重启服务
#   ./deploy.sh status    查看运行状态
#   ./deploy.sh logs      查看日志
#   ./deploy.sh update    拉取最新镜像并重启

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

info()    { echo -e "${GREEN}[✓]${NC} $1"; }
warn()    { echo -e "${YELLOW}[!]${NC} $1"; }
error()   { echo -e "${RED}[✗]${NC} $1"; exit 1; }
section() { echo -e "\n${CYAN}── $1${NC}"; }

# ── 命令处理 ────────────────────────────────────────────────────
case "${1:-}" in
  stop)
    section "停止服务"
    cd "$SCRIPT_DIR" && docker compose down
    info "服务已停止"
    exit 0
    ;;
  restart)
    section "重启服务"
    cd "$SCRIPT_DIR" && docker compose restart
    info "服务已重启"
    exit 0
    ;;
  status)
    cd "$SCRIPT_DIR" && docker compose ps
    exit 0
    ;;
  logs)
    cd "$SCRIPT_DIR" && docker compose logs -f --tail=100
    exit 0
    ;;
  update)
    section "更新到最新版本"
    cd "$SCRIPT_DIR"
    docker compose pull 2>/dev/null || warn "无法拉取远程镜像，使用本地版本"
    docker compose up -d --build
    info "更新完成"
    exit 0
    ;;
  ""|deploy)
    ;;
  *)
    echo "用法: $0 [deploy|stop|restart|status|logs|update]"
    exit 1
    ;;
esac

# ── 部署流程 ────────────────────────────────────────────────────
echo ""
echo "  ██████╗ ██╗      ██████╗  ██████╗ "
echo "  ██╔══██╗██║     ██╔═══██╗██╔════╝ "
echo "  ██████╔╝██║     ██║   ██║██║  ███╗"
echo "  ██╔══██╗██║     ██║   ██║██║   ██║"
echo "  ██████╔╝███████╗╚██████╔╝╚██████╔╝"
echo "  ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝ "
echo "  Docker 部署模式"
echo ""

# 1. 检查 Docker
section "检查环境"
if ! command -v docker &> /dev/null; then
  error "未找到 Docker，请先安装:\n  macOS: brew install --cask docker  或访问 https://docs.docker.com/get-docker/"
fi

if ! docker info &> /dev/null; then
  error "Docker 未运行，请先启动 Docker Desktop"
fi

DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null)
info "Docker 版本: $DOCKER_VERSION"

if ! docker compose version &> /dev/null; then
  error "需要 Docker Compose v2，请更新 Docker 到最新版本"
fi
info "Docker Compose: 已就绪"

# 2. 检查 .env 文件
section "检查配置"
if [ ! -f "$ENV_FILE" ]; then
  warn ".env 文件不存在，从模板创建..."
  cp "$SCRIPT_DIR/.env.example" "$ENV_FILE"
  echo ""
  echo -e "  ${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "  ${YELLOW}⚠ 请先编辑 .env 文件设置密码和密钥${NC}"
  echo ""
  echo "  必须修改以下两项："
  echo -e "    ${CYAN}ADMIN_PASS${NC}  = 后台登录密码"
  echo -e "    ${CYAN}JWT_SECRET${NC}  = 随机密钥（运行以下命令生成）"
  echo ""
  echo -e "    ${GREEN}openssl rand -hex 32${NC}"
  echo ""
  echo "  编辑完成后重新运行: ./deploy.sh"
  echo -e "  ${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo ""
  exit 1
fi

# 检查必填项是否已修改
source "$ENV_FILE" 2>/dev/null || true
if [ "${ADMIN_PASS:-}" = "your-strong-password-here" ] || [ -z "${ADMIN_PASS:-}" ]; then
  error "请在 .env 文件中设置 ADMIN_PASS（不能使用默认值）"
fi
if [ "${JWT_SECRET:-}" = "your-random-secret-here" ] || [ -z "${JWT_SECRET:-}" ]; then
  error "请在 .env 文件中设置 JWT_SECRET（可运行 openssl rand -hex 32 生成）"
fi
info ".env 配置检查通过"

# 3. 创建数据目录
section "初始化"
mkdir -p "$SCRIPT_DIR/data/uploads"
info "数据目录已就绪: ./data/"

# 4. 构建并启动
section "构建镜像"
cd "$SCRIPT_DIR"

# 尝试拉取远程镜像，失败则本地构建
if grep -q "DOCKER_IMAGE" "$ENV_FILE" 2>/dev/null && [ -n "${DOCKER_IMAGE:-}" ]; then
  info "拉取镜像: $DOCKER_IMAGE"
  docker pull "$DOCKER_IMAGE" 2>/dev/null && info "镜像拉取成功" || {
    warn "拉取失败，改用本地构建"
    docker compose build --no-cache
  }
else
  info "本地构建镜像（首次约需 3-5 分钟）..."
  docker compose build
fi

section "启动服务"
docker compose up -d

# 5. 等待健康检查
info "等待服务就绪..."
for i in $(seq 1 30); do
  sleep 1
  STATUS=$(docker compose ps --format json 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('Health',''))" 2>/dev/null || echo "")
  if curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT:-8080}/" 2>/dev/null | grep -q "200\|301"; then
    break
  fi
  if [ $i -eq 30 ]; then
    warn "服务启动超时，查看日志: ./deploy.sh logs"
  fi
done

# 6. 完成
PORT_DISPLAY="${PORT:-8080}"
echo ""
echo -e "  ${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "  ${GREEN}✓ 部署成功！${NC}"
echo ""
echo -e "  📖 前台:  ${CYAN}http://localhost:$PORT_DISPLAY${NC}"
echo -e "  🔧 后台:  ${CYAN}http://localhost:$PORT_DISPLAY/admin/index.html${NC}"
echo -e "  👤 账号:  ${YELLOW}${ADMIN_USER:-admin}${NC} / ${YELLOW}(见 .env 文件)${NC}"
echo ""
echo -e "  常用命令:"
echo -e "    ${GREEN}./deploy.sh logs${NC}     查看日志"
echo -e "    ${GREEN}./deploy.sh stop${NC}     停止服务"
echo -e "    ${GREEN}./deploy.sh update${NC}   更新到最新版"
echo -e "  ${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# macOS 自动打开浏览器
if command -v open &> /dev/null; then
  sleep 0.5
  open "http://localhost:$PORT_DISPLAY" 2>/dev/null || true
fi
