# minio 不支持
### 1. 清除配置
```shell
  extraEnvironmentVars:
    MINIO_STORAGE_ACCESS_KEY: "KKnTKEF2PA7qOVse"
    MINIO_STORAGE_SECRET_KEY: "uGQSnAdQ6xfu1GMqyBg5jS4cOS86WlNs"
    MINIO_STORAGE_BUCKET_NAME: "neolink"
    MINIO_STORAGE_ENDPOINT: "http://minio:9000"
    MINIO_SKIP: "False"
```

# minio 支持
### 1. 先安装minio

### 2. 配置
```shell
  extraEnvironmentVars:
    MINIO_STORAGE_ACCESS_KEY: "KKnTKEF2PA7qOVse"
    MINIO_STORAGE_SECRET_KEY: "uGQSnAdQ6xfu1GMqyBg5jS4cOS86WlNs"
    MINIO_STORAGE_BUCKET_NAME: "neolink"
    MINIO_STORAGE_ENDPOINT: "http://minio:9000"
    MINIO_SKIP: "False"
```