#!/bin/bash
helm -n labels uninstall minio
kubectl -n labels delete pvc data-minio-0
kubectl -n labels delete pvc data-minio-1
kubectl -n labels delete pvc data-minio-2
kubectl -n labels delete pvc data-minio-3
