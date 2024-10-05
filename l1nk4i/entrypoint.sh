#!/bin/bash

# 循环执行 myapp
while true; do
    ./myapp  # 运行应用

    # 检查上一个命令的退出状态
    if [ $? -eq 0 ]; then
        echo "myapp executed successfully."
        break  # 如果成功，退出循环
    else
        echo "myapp failed. Retrying in 2 seconds..."
        sleep 2  # 等待 2 秒后重试
    fi
done
