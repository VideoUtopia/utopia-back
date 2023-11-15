#!/bin/bash

# 检查8080端口的占用情况
port=8080
pid=$(lsof -t -i:$port)

if [ -z "$pid" ]; then
  echo "端口 $port 没有被占用。"
else
  echo "端口 $port 被进程 $pid 占用。"

  # 终止占用8080端口的进程
  echo "终止进程 $pid ..."
  kill $pid

fi

echo "开始启动Go程序..."
nohup ./server > output.log 2>&1 &
sleep 2 # 等待一段时间以确保程序已经启动
pid=$!
if [ -n "$(ps -p $pid -o pid=)" ]; then
  echo "Go程序已成功启动。"
  exit 0
else
  echo "Go程序启动失败，请检查错误信息。"
  exit 1
fi