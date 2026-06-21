#!/usr/bin/env bash
#
# TaskFlow end-to-end demo.
# Starts its own server, runs EVERY API operation, and prints each one's
# HTTP status + duration, then shows the server's own request log.
#
# Usage:   ./demo.sh         (from the taskflow/ directory)
#          make demo
#
set -uo pipefail

PORT="${PORT:-8090}"
BASE="http://localhost:${PORT}"
SECRET="demo-secret"
DB="$(mktemp -t taskflow_demo).db"
LOG="$(mktemp -t taskflow_demo).log"

# --- colors (fall back to plain if not a TTY) ---
if [ -t 1 ]; then
  BOLD=$'\033[1m'; DIM=$'\033[2m'; GREEN=$'\033[32m'; YELLOW=$'\033[33m'
  CYAN=$'\033[36m'; RED=$'\033[31m'; RESET=$'\033[0m'
else
  BOLD=""; DIM=""; GREEN=""; YELLOW=""; CYAN=""; RED=""; RESET=""
fi

cleanup() {
  if [ -n "${SRV_PID:-}" ]; then
    pkill -P "$SRV_PID" 2>/dev/null   # kill the compiled server that `go run` spawned
    kill "$SRV_PID" 2>/dev/null
    wait "$SRV_PID" 2>/dev/null       # reap quietly so the shell doesn't print "Terminated"
  fi
  rm -f "$DB" "$LOG"
}
trap cleanup EXIT

# --- start the server ---
echo "${BOLD}═══ TaskFlow API — live demo ═══${RESET}"
echo "Starting server on :${PORT} ..."
TASKFLOW_ADDR=":${PORT}" TASKFLOW_DB="$DB" TASKFLOW_JWT_SECRET="$SECRET" \
  go run . > "$LOG" 2>&1 &
SRV_PID=$!

# wait until it answers (or the process dies)
for _ in $(seq 1 100); do
  if curl -s -m1 "$BASE/health" >/dev/null 2>&1; then break; fi
  if ! kill -0 "$SRV_PID" 2>/dev/null; then
    echo "${RED}server failed to start:${RESET}"; cat "$LOG"; exit 1
  fi
  sleep 0.1
done
echo "${GREEN}ready${RESET}  ${DIM}(temp DB: $DB)${RESET}"
echo ""

TOKEN=""

# request: LABEL METHOD PATH [BODY] [auth] [show-body]
# prints a table row: label | method | status | duration(ms), and optionally the body.
request() {
  local label="$1" method="$2" path="$3" body="${4:-}" auth="${5:-}" show="${6:-}"
  local tmp; tmp="$(mktemp)"
  local args=(-s -o "$tmp" -w '%{http_code} %{time_total}' -X "$method")
  [ -n "$body" ] && args+=(--data "$body")
  [ "$auth" = "auth" ] && [ -n "$TOKEN" ] && args+=(-H "Authorization: Bearer $TOKEN")

  local meta code secs ms
  meta="$(curl "${args[@]}" "$BASE$path")"
  code="${meta%% *}"; secs="${meta##* }"
  ms="$(awk "BEGIN{printf \"%.1f\", ${secs}*1000}")"

  # color the status code
  local cc="$GREEN"
  case "$code" in 4*|5*) cc="$RED";; esac

  printf "  %-34s ${CYAN}%-6s${RESET} ${cc}%-3s${RESET}  ${YELLOW}%7s ms${RESET}\n" \
    "$label" "$method" "$code" "$ms"

  if [ "$show" = "show" ]; then
    sed 's/^/      /' "$tmp" | (python3 -m json.tool 2>/dev/null | sed 's/^/      /' || cat "$tmp")
  fi
  echo "$(cat "$tmp")" > /tmp/_tf_last_body
  rm -f "$tmp"
}

printf "  ${BOLD}%-34s %-6s %-3s  %10s${RESET}\n" "OPERATION" "METHOD" "ST" "DURATION"
echo   "  ────────────────────────────────────────────────────────────────"

# 1) Public endpoints
request "API index (/)"                 GET  "/"
request "Health check"                  GET  "/health"

# 2) Auth wall
request "List tasks WITHOUT token (401)" GET  "/tasks"

# 3) Register -> capture token
request "Register account"               POST "/auth/register" '{"email":"demo@taskflow.dev","password":"secret123"}'
TOKEN="$(sed -n 's/.*"token":"\([^"]*\)".*/\1/p' /tmp/_tf_last_body)"
echo "  ${DIM}captured JWT: ${TOKEN:0:24}...${RESET}"

# 4) Create tasks (the "changes")
echo ""
echo "  ${BOLD}Creating tasks:${RESET}"
request "Create task (priority=high)"    POST "/tasks" '{"title":"Finish the Go course","priority":"high"}'   auth show
request "Create task (default=medium)"   POST "/tasks" '{"title":"Push to GitHub"}'                            auth show
request "Create task (priority=low)"     POST "/tasks" '{"title":"Update CV","priority":"low"}'                auth show
request "Create invalid priority (400)"  POST "/tasks" '{"title":"bad","priority":"urgent"}'                   auth show

# 5) Read / filter
echo ""
echo "  ${BOLD}Reading & filtering:${RESET}"
request "List all tasks"                 GET  "/tasks"                 "" auth show
request "Filter ?priority=high"          GET  "/tasks?priority=high"   "" auth show
request "Get one task (/tasks/2)"        GET  "/tasks/2"               "" auth

# 6) Update
echo ""
echo "  ${BOLD}Updating & deleting:${RESET}"
request "Mark task 1 done"               PUT  "/tasks/1" '{"title":"Finish the Go course","priority":"high","done":true}' auth show
request "Filter ?done=true"              GET  "/tasks?done=true"       "" auth show
request "Delete task 3 (204)"            DELETE "/tasks/3"             "" auth
request "Get deleted task 3 (404)"       GET  "/tasks/3"               "" auth

# 7) Login
echo ""
echo "  ${BOLD}Login:${RESET}"
request "Login correct password"         POST "/auth/login" '{"email":"demo@taskflow.dev","password":"secret123"}'
request "Login WRONG password (401)"     POST "/auth/login" '{"email":"demo@taskflow.dev","password":"nope"}'

# --- show the server's own log (server-side durations from the Logging middleware) ---
echo ""
echo "  ${BOLD}Server log (Logging middleware — server-side durations):${RESET}"
sed 's/^/      /' "$LOG"

echo ""
echo "${GREEN}${BOLD}✓ demo complete — server stopped, temp DB removed.${RESET}"
rm -f /tmp/_tf_last_body
