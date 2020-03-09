/*

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
// +kubebuilder:docs-gen:collapse=Apache License

/*
最开始我们导入了 `meta/v1` 库，该库通常不会直接使用，但是它包含了 Kubernetes 所有 kind 的 metadata 结构体类型。
*/
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
下一步，我们为我们的 Kind 定义 Spec(`CronJobSpec`) 和 Status(`CronJobStatus`) 的类型。
Kubernetes 的功能是将期望状态(`Spec`)与集群实际状态(其他对象的`Status`)和外部状态进行协调。
然后记录下观察到的状态（`Status`）。
因此，每个具有功能的对象都包含 `Spec` 和 `Status` 。
但是 ConfigMap 之类的一些类型不遵循此模式，因为它们不会维护一种状态，但是大多数类型都需要。
*/
// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CronJobSpec defines the desired state of CronJob
type CronJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// CronJobStatus defines the observed state of CronJob
type CronJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

/*<span id="moment"></span>*/
/*
然后，我们定义了 Kinds 对应的结构体类型，`CronJob` 和 `CronJobList`。
`CronJob` 是我们的 root type，用来描述 `CronJob Kind`。和所有 Kubernetes 对象一样，
它包含 `TypeMeta` (用来定义 API version 和 Kind) 和 `ObjectMeta` (用来定义 name、namespace 和 labels等一些字段)

`CronJobList` 包含了一个 `CronJob` 的切片，它是用来批量操作 `Kind` 的，比如 LIST 操作

通常，我们不会修改它们 -- 所有的修改都是在 Spec 和 Status 上进行的。

`+kubebuilder:object:root` 注释称为标记(marker)。
稍后我们还会看到更多它们，它们提供了一些元数据，
来告诉 [controller-tools](https://github.com/kubernetes-sigs/controller-tools)(我们的代码和 YAML 生成器) 一些额外的信息。
这个注释告诉 `object` 这是一种 `root type Kind`。
然后，`object` 生成器会为我们生成 [runtime.Object](https://godoc.org/k8s.io/apimachinery/pkg/runtime#Object) 接口的实现，
这是所有 Kinds 必须实现的接口。
*/

// +kubebuilder:object:root=true

// CronJob is the Schema for the cronjobs API
type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronJobSpec   `json:"spec,omitempty"`
	Status CronJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CronJobList contains a list of CronJob
type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronJob `json:"items"`
}

/*
最后，我们将 Kinds 注册到 API group 中。这使我们可以将此 API group 中的 Kind 添加到任何 [Scheme](https://godoc.org/k8s.io/apimachinery/pkg/runtime#Scheme) 中。
*/
func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}
