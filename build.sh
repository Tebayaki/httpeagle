#!/bin/bash

# 镜像名称和标签
IMAGE_NAME="httpeagle"
IMAGE_TAG="latest"
CONTAINER_NAME="httpeagle-server"

# 构建镜像
echo "🔨 Building Docker image..."
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

# 检查是否构建成功
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi