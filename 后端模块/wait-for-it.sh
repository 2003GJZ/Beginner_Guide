#!/bin/bash
# 等待其他服务准备就绪的脚本

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until nc -z "$host" "$port"; do
  >&2 echo "等待 $host:$port 准备就绪..."
  sleep 1
done

>&2 echo "$host:$port 已准备就绪"

exec $cmd
