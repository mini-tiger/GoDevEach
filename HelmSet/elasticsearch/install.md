```shell

docker pull docker.io/bitnami/bitnami-shell:10-debian-10-r403
docker tag docker.io/bitnami/bitnami-shell:10-debian-10-r403 harbor.dev.21vianet.com/taojun/bitnami-shell:10-debian-10-r403

docker pull bitnami/kibana:7.17.9-debian-11-r26 && docker pull bitnami/elasticsearch:7.17.9-debian-11-r30

docker tag bitnami/kibana:7.17.9-debian-11-r26 harbor.dev.21vianet.com/taojun/kibana:7.17.9-debian-11-r26
docker push harbor.dev.21vianet.com/taojun/kibana:7.17.9-debian-11-r26
docker tag bitnami/elasticsearch:7.17.9-debian-11-r30 harbor.dev.21vianet.com/taojun/elasticsearch:7.17.9-debian-11-r30
docker push harbor.dev.21vianet.com/taojun/elasticsearch:7.17.9-debian-11-r30
```