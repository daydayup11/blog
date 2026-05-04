#!/bin/bash
# 一键部署脚本 — 个人博客
# 用法:
#   ./deploy.sh           首次部署（交互式配置）
#   ./deploy.sh stop      停止服务
#   ./deploy.sh restart   重启服务
#   ./deploy.sh status    查看运行状态
#   ./deploy.sh logs      查看日志
#   ./deploy.sh update    更新到最新版本

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

info()    { echo -e "${GREEN}[✓]${NC} $1"; }
warn()    { echo -e "${YELLOW}[!]${NC} $1"; }
error()   { echo -e "${RED}[✗]${NC} $1"; exit 1; }
section() { echo -e "\n${CYAN}${BOLD}── $1${NC}"; }
ask()     { echo -e "${YELLOW}[?]${NC} $1"; }

banner() {
  echo ""
  echo "  ██████╗ ██╗      ██████╗  ██████╗ "
  echo "  ██╔══██╗██║     ██╔═══██╗██╔════╝ "
  echo "  ██████╔╝██║     ██║   ██║██║  ███╗"
  echo "  ██╔══██╗██║     ██║   ██║██║   ██║"
  echo "  ██████╔╝███████╗╚██████╔╝╚██████╔╝"
  echo "  ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝ "
  echo ""
}

# 兼容 docker compose（插件版）和 docker-compose（独立版）
if docker compose version &>/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose &>/dev/null; then
  DC="docker-compose"
else
  DC=""
fi

# ── 非首次部署的子命令 ───────────────────────────────────────────
case "${1:-}" in
  stop)
    [ -z "$DC" ] && error "未找到 Docker Compose"
    section "停止服务"
    cd "$SCRIPT_DIR" && $DC down
    info "服务已停止"
    exit 0
    ;;
  restart)
    [ -z "$DC" ] && error "未找到 Docker Compose"
    section "重启服务"
    cd "$SCRIPT_DIR" && $DC restart
    info "服务已重启"
    exit 0
    ;;
  status)
    [ -z "$DC" ] && error "未找到 Docker Compose"
    cd "$SCRIPT_DIR" && $DC ps
    exit 0
    ;;
  logs)
    [ -z "$DC" ] && error "未找到 Docker Compose"
    cd "$SCRIPT_DIR" && $DC logs -f --tail=100
    exit 0
    ;;
  update)
    [ -z "$DC" ] && error "未找到 Docker Compose"
    section "更新到最新版本"
    cd "$SCRIPT_DIR"
    $DC pull 2>/dev/null || warn "无法拉取远程镜像，改用本地重新构建"
    $DC up -d --build
    info "更新完成"
    exit 0
    ;;
  ""|deploy) ;;
  *)
    echo "用法: $0 [deploy|stop|restart|status|logs|update]"
    exit 1
    ;;
esac

# ══════════════════════════════════════════════════════════════════
# 首次部署流程
# ══════════════════════════════════════════════════════════════════
banner

# ── Step 1: 检查 Docker ──────────────────────────────────────────
section "检查环境"

if ! command -v docker &>/dev/null; then
  echo ""
  warn "未检测到 Docker，需要先安装"
  echo ""
  echo "  安装方式（选一个）："
  if [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "  ${CYAN}方式 A（推荐）${NC}：brew install --cask docker"
    echo -e "  ${CYAN}方式 B${NC}：https://docs.docker.com/desktop/install/mac-install/"
  else
    echo -e "  ${CYAN}Linux 一键安装${NC}：curl -fsSL https://get.docker.com | sh"
    echo -e "  ${CYAN}其他系统${NC}：https://docs.docker.com/get-docker/"
  fi
  echo ""
  read -r -p "  安装完成后按 Enter 继续，或 Ctrl+C 退出..." _
  echo ""
  command -v docker &>/dev/null || error "仍未检测到 Docker，请安装后重试"
fi

if ! docker info &>/dev/null 2>&1; then
  warn "Docker 已安装但未运行"
  if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "  请打开 Docker Desktop 应用，等待启动完成后按 Enter 继续"
    read -r -p "  " _
    sleep 3
  fi
  docker info &>/dev/null 2>&1 || error "Docker 仍未运行，请启动后重试"
fi
info "Docker $(docker version --format '{{.Server.Version}}' 2>/dev/null)"

if [ -z "$DC" ]; then
  warn "未找到 Docker Compose，尝试自动安装..."
  if [[ "$OSTYPE" == "darwin"* ]]; then
    brew install docker-compose 2>/dev/null && DC="docker-compose" || true
  fi
  [ -z "$DC" ] && error "请手动安装 Docker Compose：brew install docker-compose"
fi
info "Docker Compose $($DC version --short 2>/dev/null || echo 'OK')"

# ── Step 2: 配置向导 ─────────────────────────────────────────────
section "配置向导"

# 判断是否已有 .env（重新部署）
if [ -f "$ENV_FILE" ]; then
  source "$ENV_FILE" 2>/dev/null || true
  echo ""
  echo -e "  检测到已有配置文件，当前设置："
  echo -e "  ${CYAN}端口${NC}：${PORT:-8080}"
  echo -e "  ${CYAN}用户名${NC}：${ADMIN_USER:-admin}"
  echo -e "  ${CYAN}密码${NC}：已设置"
  echo ""
  read -r -p "  是否使用现有配置？[Y/n] " USE_EXISTING
  USE_EXISTING="${USE_EXISTING:-Y}"
  if [[ "$USE_EXISTING" =~ ^[Yy]$ ]]; then
    SKIP_CONFIG=true
  fi
fi

if [ "${SKIP_CONFIG:-false}" != "true" ]; then
  echo ""
  echo -e "  ${BOLD}接下来只需回答几个问题，全部有默认值，直接按 Enter 跳过${NC}"
  echo ""

  # 端口
  ask "服务端口 [默认: 8080]："
  read -r INPUT_PORT
  CONF_PORT="${INPUT_PORT:-8080}"

  # 用户名
  ask "后台管理用户名 [默认: admin]："
  read -r INPUT_USER
  CONF_USER="${INPUT_USER:-admin}"

  # 密码
  while true; do
    ask "后台管理密码 [直接 Enter 自动生成强密码]："
    read -r -s INPUT_PASS
    echo ""
    if [ -z "$INPUT_PASS" ]; then
      CONF_PASS=$(LC_ALL=C tr -dc 'A-Za-z0-9@#%' </dev/urandom | head -c 16 2>/dev/null || \
                  openssl rand -base64 12 | tr -dc 'A-Za-z0-9' | head -c 16)
      echo -e "  ${GREEN}已自动生成密码：${BOLD}${CONF_PASS}${NC}"
      echo -e "  ${YELLOW}请记下这个密码！${NC}"
      echo ""
      break
    else
      if [ ${#INPUT_PASS} -lt 6 ]; then
        warn "密码至少 6 位，请重新输入"
      else
        CONF_PASS="$INPUT_PASS"
        break
      fi
    fi
  done

  # JWT 密钥（完全自动，不问用户）
  CONF_JWT=$(openssl rand -hex 32 2>/dev/null || \
             LC_ALL=C tr -dc 'a-f0-9' </dev/urandom | head -c 64 2>/dev/null || \
             date +%s%N | sha256sum | head -c 64)

  # 写入 .env
  cat > "$ENV_FILE" <<EOF
# 个人博客配置文件
# 生成时间：$(date '+%Y-%m-%d %H:%M:%S')

PORT=$CONF_PORT
ADMIN_USER=$CONF_USER
ADMIN_PASS=$CONF_PASS
JWT_SECRET=$CONF_JWT
DB_PATH=/app/data/blog.db
FRONTEND_DIR=/app/frontend
EOF

  info "配置已保存到 .env"
fi

# 加载配置
source "$ENV_FILE"

# ── Step 3: 初始化数据目录 ───────────────────────────────────────
section "初始化"
mkdir -p "$SCRIPT_DIR/data/uploads"
info "数据目录就绪"

# ── Step 4: 构建并启动 ───────────────────────────────────────────
section "构建镜像"
cd "$SCRIPT_DIR"

echo "  首次构建约需 3-5 分钟，请耐心等待..."
echo ""

# 尝试拉取远程镜像，失败则本地构建
REMOTE_IMAGE="${DOCKER_IMAGE:-}"
if [ -n "$REMOTE_IMAGE" ]; then
  info "拉取镜像：$REMOTE_IMAGE"
  docker pull "$REMOTE_IMAGE" 2>/dev/null && info "镜像拉取成功" || {
    warn "拉取失败，改用本地构建"
    $DC build
  }
else
  $DC build
fi

section "启动服务"
$DC up -d

# ── Step 5: 等待就绪 ─────────────────────────────────────────────
info "等待服务启动..."
READY=false
for i in $(seq 1 30); do
  sleep 1
  if curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT:-8080}/" 2>/dev/null | grep -qE "200|301"; then
    READY=true
    break
  fi
done

# ── Step 6: 完成 ─────────────────────────────────────────────────
PORT_DISPLAY="${PORT:-8080}"
echo ""
echo -e "  ${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
if [ "$READY" = true ]; then
  echo -e "  ${GREEN}${BOLD}✓ 部署成功！${NC}"
else
  echo -e "  ${YELLOW}${BOLD}⚠ 服务可能仍在启动中${NC}"
fi
echo ""
echo -e "  📖 前台   ${CYAN}http://localhost:$PORT_DISPLAY${NC}"
echo -e "  🔧 后台   ${CYAN}http://localhost:$PORT_DISPLAY/admin${NC}"
echo -e "  👤 账号   ${BOLD}${ADMIN_USER:-admin}${NC} / ${BOLD}${ADMIN_PASS}${NC}"
echo ""
echo -e "  常用命令："
echo -e "    ${GREEN}./deploy.sh logs${NC}      查看运行日志"
echo -e "    ${GREEN}./deploy.sh stop${NC}      停止服务"
echo -e "    ${GREEN}./deploy.sh restart${NC}   重启服务"
echo -e "    ${GREEN}./deploy.sh update${NC}    更新到最新版本"
echo -e "  ${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# macOS 自动打开浏览器
if [ "$READY" = true ] && command -v open &>/dev/null; then
  sleep 0.5
  open "http://localhost:$PORT_DISPLAY" 2>/dev/null || true
fi
