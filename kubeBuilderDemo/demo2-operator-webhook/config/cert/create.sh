# 生成 CA 证书
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=host.docker.internal" -days 10000 -out ca.crt

# 签发本地证书
openssl genrsa -out tls.key 2048
openssl req -new -SHA256 -newkey rsa:2048 -nodes -keyout tls.key -out tls.csr -subj "/C=CN/ST=Shanghai/L=Shanghai/O=/OU=/CN=host.docker.internal"
openssl req -new -key tls.key -out tls.csr -config csr.conf
openssl x509 -req -in tls.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out tls.crt -days 10000 -extensions v3_ext -extfile csr.conf