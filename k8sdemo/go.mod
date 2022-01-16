module k8sdemo

go 1.16

require (
	k8s.io/api v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.20.0
)
// Kubernetes version 匹配 client-go version
// https://github.com/kubernetes/client-go
