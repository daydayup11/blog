#!/bin/bash
# 一键启动脚本 — 个人博客
# 用法: ./start.sh [stop|restart|status|logs]

set -e

# ── 配置（可按需修改）───────────────────────────────────────────
PORT=${PORT:-8080}
ADMIN_USER=${ADMIN_USER:-admin}
ADMIN_PASS=${ADMIN_PASS:-admin123}
JWT_SECRET=${JWT_SECRET:-$(LC_ALL=C tr -dc 'a-zA-Z0-9' < /dev/urandom | head -c 32 2>/dev/null || echo "dev-secret-please-change")}
DB_PATH=${DB_PATH:-./data/blog.db}
# ────────────────────────────────────────────────────────────────

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"
PID_FILE="$SCRIPT_DIR/.blog.pid"
LOG_FILE="$SCRIPT_DIR/.blog.log"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

info()    { echo -e "${GREEN}[✓]${NC} $1"; }
warn()    { echo -e "${YELLOW}[!]${NC} $1"; }
error()   { echo -e "${RED}[✗]${NC} $1"; exit 1; }
section() { echo -e "\n${CYAN}── $1${NC}"; }

# ── 辅助函数 ────────────────────────────────────────────────────

is_running() {
  [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null
}

check_port() {
  lsof -ti:"$PORT" > /dev/null 2>&1
}

stop_server() {
  if [ -f "$PID_FILE" ]; then
    local pid; pid=$(cat "$PID_FILE")
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid"
      sleep 1
      info "服务已停止 (PID $pid)"
    fi
    rm -f "$PID_FILE"
  else
    # Try to kill whatever is on the port
    lsof -ti:"$PORT" | xargs kill -9 2>/dev/null || true
    info "端口 $PORT 已清理"
  fi
}

# ── 命令处理 ────────────────────────────────────────────────────

case "${1:-}" in
  stop)
    section "停止服务"
    stop_server
    exit 0
    ;;
  status)
    if is_running; then
      info "运行中 (PID $(cat "$PID_FILE")) — http://localhost:$PORT"
    else
      warn "未运行"
    fi
    exit 0
    ;;
  logs)
    [ -f "$LOG_FILE" ] && tail -f "$LOG_FILE" || warn "日志文件不存在"
    exit 0
    ;;
  restart)
    stop_server
    ;;
  ""|start)
    ;;
  *)
    echo "用法: $0 [start|stop|restart|status|logs]"
    exit 1
    ;;
esac

# ── 启动流程 ────────────────────────────────────────────────────

echo ""
echo "  ██████╗ ██╗      ██████╗  ██████╗ "
echo "  ██╔══██╗██║     ██╔═══██╗██╔════╝ "
echo "  ██████╔╝██║     ██║   ██║██║  ███╗"
echo "  ██╔══██╗██║     ██║   ██║██║   ██║"
echo "  ██████╔╝███████╗╚██████╔╝╚██████╔╝"
echo "  ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝ "
echo ""

# 1. 检查 Go
section "检查环境"
if ! command -v go &> /dev/null; then
  error "未找到 Go，请先安装: https://go.dev/dl/\n  macOS: brew install go"
fi
GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | head -1)
info "Go 版本: $GO_VERSION"

# 2. 检查端口
if check_port; then
  warn "端口 $PORT 已被占用，尝试释放..."
  lsof -ti:"$PORT" | xargs kill -9 2>/dev/null || true
  sleep 1
  if check_port; then
    error "端口 $PORT 仍被占用，请手动处理: lsof -i:$PORT"
  fi
fi

# 3. 编译
section "编译后端"
cd "$BACKEND_DIR"

if [ ! -f "go.mod" ]; then
  error "找不到 backend/go.mod，请确认项目结构"
fi

info "下载依赖..."
go mod download 2>/dev/null && info "依赖已就绪" || warn "依赖下载有警告，继续..."

info "编译中..."
go build -o "$SCRIPT_DIR/.blog_server" ./cmd/server 2>&1 | grep -v "^$" || true

if [ ! -f "$SCRIPT_DIR/.blog_server" ]; then
  error "编译失败，查看错误信息"
fi
info "编译成功"

# 4. 创建数据目录
section "初始化"
mkdir -p "$SCRIPT_DIR/data/uploads"
info "数据目录: $SCRIPT_DIR/data/"

# 5. 启动服务
section "启动服务"

export PORT="$PORT"
export DB_PATH="$DB_PATH"
export JWT_SECRET="$JWT_SECRET"
export ADMIN_USER="$ADMIN_USER"
export ADMIN_PASS="$ADMIN_PASS"
export FRONTEND_DIR="$FRONTEND_DIR"

cd "$SCRIPT_DIR"
nohup ./.blog_server > "$LOG_FILE" 2>&1 &
echo $! > "$PID_FILE"
SERVER_PID=$!

# 等待服务就绪（最多 10 秒）
info "等待服务就绪..."
for i in $(seq 1 20); do
  sleep 0.5
  if curl -s -o /dev/null -w "%{http_code}" "http://localhost:$PORT/" 2>/dev/null | grep -q "200\|301"; then
    break
  fi
  if ! kill -0 "$SERVER_PID" 2>/dev/null; then
    error "服务启动失败，查看日志: cat $LOG_FILE"
  fi
done

if ! kill -0 "$SERVER_PID" 2>/dev/null; then
  error "服务进程已退出，查看日志: cat $LOG_FILE"
fi

# 6. 完成
echo ""
echo -e "  ${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "  ${GREEN}✓ 启动成功！${NC}"
echo ""
echo -e "  📖 前台:  ${CYAN}http://localhost:$PORT${NC}"
echo -e "  🔧 后台:  ${CYAN}http://localhost:$PORT/admin/index.html${NC}"
echo -e "  👤 账号:  ${YELLOW}$ADMIN_USER${NC} / ${YELLOW}$ADMIN_PASS${NC}"
echo -e "  📋 日志:  ./start.sh logs"
echo -e "  🛑 停止:  ./start.sh stop"
echo -e "  ${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 自动打开浏览器（macOS）
if command -v open &> /dev/null; then
  sleep 0.5
  open "http://localhost:$PORT" 2>/dev/null || true
fi
