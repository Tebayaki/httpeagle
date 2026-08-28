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
    echo ""
    echo "📦 Image size:"
    docker images ${IMAGE_NAME}:${IMAGE_TAG} --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
    echo ""
    echo "🚀 Run container:"
    echo "docker run -d \\"
    echo "  --name ${CONTAINER_NAME} \\"
    echo "  -p 41596:41596 \\"
    echo "  -v /path/to/your/images:/images \\"
    echo "  -v /path/to/your/certs:/certs \\"
    echo "  ${IMAGE_NAME}:${IMAGE_TAG}"
else
    echo "❌ Build failed!"
    exit 1
fi