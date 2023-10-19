#!/bin/bash
helm uninstall -n labels label-studio
kubectl -n labels delete pvc data-label-studio-postgresql-0
kubectl -n labels delete pvc redis-data-label-studio-redis-master-0
kubectl -n labels delete pvc label-studio-ls-pvc
kubectl -n labels delete pvc data-label-studio-postgresql-primary-0
kubectl -n labels delete pvc data-label-studio-postgresql-read-0
