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
首先，我们会 import 一些标准库。和以前一样，我们需要 `controller-runtime` 库，client 包以及我们定义的有关 API 类型的软件包。
*/
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "tutorial.kubebuilder.io/project/api/v1"
)

/*
接下来，kubebuilder 为我们搭建了一个基本的 `reconciler` 结构体。
几乎每个 `reconciler` 都需要记录日志，并且需要能够获取对象，因此这个结构体是开箱即用的。
*/

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

/*
大多数 controllers 最终都会运行在 k8s 集群上，因此它们需要 RBAC 权限,
我们使用 controller-tools [RBAC markers](/reference/markers/rbac.md) 指定了这些权限。
这是运行所需的最小权限。
随着我们添加更多功能，我们将会重新定义这些权限。
*/

// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch

/*
`Reconcile` 方法对个某个单一的 object 执行 `reconciling` 动作，
我们的 [Request](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Request)只是一个 name，
但是 client 可以通过 name 信息从 cache 中获取到对应的 object。

我们返回的 result 为空，且 error 为 nil，
这表明 controller-runtime 已经成功 reconciled 了这个 object，无需进行任何重试，直到这个 object 被更改。

大多数 controller 都需要一个 logging handle 和一个 context，因此我们在这里进行了设置。

[context](https://golang.org/pkg/context/) 是用于允许取消请求或者用于跟踪之类的事情。
这是所有 client 方法的第一个参数。
`Background` context 只是一个 basic context，没有任何额外的数据或时间限制。

logging handle 用于记录日志, controller-runtime 通过 [logr](https://github.com/go-logr/logr) 库结构化日志记录。
稍后我们将看到，日志记录通过将 key-value 添加到静态消息中而起作用。
我们可以在 reconcile 方法的顶部提前分配一些 key-value ,以便查找在这个 reconciler 中所有的日志
*/
func (r *CronJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("cronjob", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

/*
最后，我们将此 reconciler 加到 manager，以便在启动 manager 时启动 reconciler。

现在，我们仅记录了 reconciler 在 `CronJob` 上的动作。稍后，我们将使用它来标记我们关心的 objects。
*/

func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Complete(r)
}
