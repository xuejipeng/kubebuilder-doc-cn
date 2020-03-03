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

我们的 main 包从一些基本的 import 开始。 值得关注的库：

- 比较核心的库： [controller-runtime](https://godoc.org/sigs.k8s.io/controller-runtime)
- controller-runtime 的默认日志库：Zap（稍后会详细介绍）

*/

package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

/*

每组控制器都需要一个 [*Scheme*](https://book.kubebuilder.io/cronjob-tutorial/gvks.html#err-but-whats-that-scheme-thing)，
Scheme 会提供 Kinds 与 Go types 之间的映射关系。
在编写 API 定义时，我们将进一步讨论 Kinds。所以，现在你只需要记住这一点。

*/
var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	// +kubebuilder:scaffold:scheme
}

/*
这是，我们的 main.go 功能相当简单：

- 为 metrics 绑定一些基本的 flags。

- 我们实例化一个 manager，用于跟踪我们运行的所有 controllers，并设置 shared caches 和可以连接到 API server 的 k8s clients 实例，并将 Scheme 配置传入 manager。

- 我们运行我们的 manager, 而 manager 又运行所有的 controllers 和 webhook。manager 会一直处于运行状态，直到收到正常关闭信号为止。
这样，当我们的 operator 运行在 Kubernetes 上时，我们可以通过优雅的方式终止这个 Pod。

虽然目前我们还没有任何东西可以运行，但是请记住使用了 `+kubebuilder:scaffold:builder` 注释的地方 -- 事情很快就会变得有趣。
*/

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme, MetricsBindAddress: metricsAddr})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	/*
	   请注意：下面代码如果指定了 Namespace 字段, controllers 仍然可以 watch cluster-scoped 的资源(例如Node), 对于 namespace 级别的资源，缓存将仅保存指定 namespace 中的对象。
	*/
	mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		Namespace:          namespace,
		MetricsBindAddress: metricsAddr,
	})

	/*
	   上面的示例将 operator 应用范围局限到单个 namespace。在这种情况下，建议将默认的 ClusterRole 和 ClusterRoleBinding 分别替换​​为 Role 和 RoleBinding 来将授权限制于此名称空间。
	   有关更多信息，请参见如何使用 [RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)。

	   除此之外，你可以使用 MultiNamespacedCacheBuilder 来监视一组特定的 namespaces：
	*/

	var namespaces []string // List of Namespaces

	mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		NewCache:           cache.MultiNamespacedCacheBuilder(namespaces),
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})

	/*
		有关更多信息，请参见 [MultiNamespacedCacheBuilder](https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/cache#MultiNamespacedCacheBuilder)
	*/

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
