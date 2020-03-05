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
首先，我们有一些 `package-level` 的标记，`+kubebuilder:object:generate=true` 表示该程序包中有 Kubernetes 对象，
`+groupName=batch.tutorial.kubebuilder.io` 表示该程序包的 Api group 是 `batch.tutorial.kubebuilder.io`。
`object` 生成器使用前者，而 CRD 生成器使用后者, 来为由此包创建的 CRD 生成正确的元数据。
*/

// Package v1 contains API Schema definitions for the batch v1 API group
// +kubebuilder:object:generate=true
// +groupName=batch.tutorial.kubebuilder.io
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

/*
然后，我们定义一些全局变量帮助我们建立 Scheme。
由于我们需要在 controller 中使用此程序包中的所有 types，
因此需要一种方便的方法（或约定）可以将所有 types 添加到其他 `Scheme` 中。
`SchemeBuilder` 让我们做这件事情变的容易。
*/
var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "batch.tutorial.kubebuilder.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
