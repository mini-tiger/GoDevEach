https://github.com/GoogleContainerTools/kaniko/releases
1. images
```shell
docker pull gcr.io/kaniko-project/executor:v1.9.1-debug
docker tag gcr.io/kaniko-project/executor:v1.9.1-debug harbor.dev.21vianet.com/taojun/kaniko-project-executor:v1.9.1-debug
```

2. help
```shell
 docker run -it --rm harbor.dev.21vianet.com/taojun/kaniko-project-executor:v1.9.1-debug --help
```

3. build package in docker

```shell
# .docker/config.json 导入容器内 上传镜像不用传入密钥
# 将Dockerfile 和相关文件导入容器 工作目录{/workspaceDemo}
# --context dockerfile 上下文,支持多种[dir,stdin...]
docker run --rm -entrypoint="" \
 -v "$HOME"/.docker/config.json:/kaniko/.docker/config.json \
 -v "$PWD"/workspaceDemo:/workspaceDemo  \
 harbor.dev.21vianet.com/taojun/kaniko-project-executor:v1.9.1-debug \
    --skip-tls-verify \
    --context=/workspaceDemo \
 --dockerfile=/workspaceDemo/Dockerfile \
    --destination=harbor.dev.21vianet.com/taojun/ubuntu20.04-ssh:latest \
    -v debug
```
4. build package in kubernetes
   https://www.orchome.com/10644
   https://github.com/GoogleContainerTools/kaniko
   https://www.xiexianbin.cn/docker/images/kaniko/index.html