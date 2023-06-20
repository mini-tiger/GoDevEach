/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GuestbookSpec defines the desired state of Guestbook
// 期望状态
type GuestbookSpec struct {
	// 业务服务对应的镜像，包括名称:tag
	Image string `json:"image"`
	// service占用的宿主机端口，外部请求通过此端口访问pod的服务
	PushPort *int32 `json:"pushPort"`
	Port     *int32 `json:"port,omitempty"`
	// 单个pod的QPS上限
	SinglePodQPS *int32 `json:"singlePodQPS"`
	// 当前整个业务的总QPS需求量
	TotalQPS *int32 `json:"totalQPS"`
}

// GuestbookStatus defines the observed state of Guestbook
// 实际状态，该数据结构中的值都是业务代码计算出来的
type GuestbookStatus struct {
	// 当前kubernetes中实际发生的总QPS
	RealQPS *int32 `json:"realQPS,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=gb
//+kubebuilder:printcolumn:name="SinglePodQPS",type="integer",JSONPath=".spec.singlePodQPS",description="single pod qps field"
//+kubebuilder:printcolumn:name="TotalQPS",type="integer",JSONPath=".spec.totalQPS",description="TotalQPS field"
//+kubebuilder:printcolumn:name="Port",type="integer",JSONPath=".spec.port",description="port field"
//+kubebuilder:printcolumn:name="PushPort",type="integer",JSONPath=".spec.pushPort",description="pushport field"

// Guestbook is the Schema for the guestbooks API
type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}

func (in *Guestbook) String() string {
	var realQPS string

	if nil == in.Status.RealQPS {
		realQPS = "nil"
	} else {
		realQPS = strconv.Itoa(int(*(in.Status.RealQPS)))
	}

	return fmt.Sprintf("Image [%s], Port [%d],PushPort [%d], SinglePodQPS [%d], TotalQPS [%d], RealQPS [%s]",
		in.Spec.Image,
		*(in.Spec.Port),
		*(in.Spec.PushPort),
		*(in.Spec.SinglePodQPS),
		*(in.Spec.TotalQPS),
		realQPS)
}

//+kubebuilder:object:root=true

// GuestbookList contains a list of Guestbook
type GuestbookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Guestbook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Guestbook{}, &GuestbookList{})
}
