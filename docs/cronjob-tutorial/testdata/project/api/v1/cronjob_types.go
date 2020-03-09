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

package v1

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:docs-gen:collapse=Imports

/*
首先，让我们看一下 `spec`，如我们之前说的，`spec` 表明**期望的状态**，所以有关 controller 的任何 "inputs" 都在这里声明。

一个基本的 CronJob 需要以下功能：

- 一个 schedule(调度器) -- (CronJob 中的 "Cron")
- 用来运行 Job 的模板 -- (CronJob 中的 "Job")

为了让用户体现更好，我们也需要一些额外的功能：

- 一个截止时间（`StartingDeadlineSeconds`）, 如果错过了这个截止时间，Job 将会等到下一个调度时间点再被调度。
- 如果多个 Job 同时启动要怎么做（`ConcurrencyPolicy`）（等待？停掉最老的一个？还是同时运行？）
- 一个暂停(`Suspend`)功能，以防止 Job 在运行过程中出现什么错误。
- 限制历史 Job 的数量（`SuccessfulJobsHistoryLimit`）

请记住，因为 job 不会读取自己的状态，所以我们需要用一些方式去跟踪一个 CronJob 是否已经运行过 Job，
我们可以用至少一个 old job 来完成这个功能。

我们会用几个标记 (`//+comment`) 去定义一些额外的数据，这些标记在
[controller-tools](https://github.com/kubernetes-sigs/controller-tools)
生成 CRD manifest 时会被使用。
稍后我们还将看到，controller-tools 还会使用 GoDoc 来为字段生成描述信息。
*/

// CronJobSpec defines the desired state of CronJob
type CronJobSpec struct {
	// +kubebuilder:validation:MinLength=0

	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule"`

	// +kubebuilder:validation:Minimum=0

	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// Specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	// +optional
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// Specifies the job that will be created when executing a CronJob.
	JobTemplate batchv1beta1.JobTemplateSpec `json:"jobTemplate"`

	// +kubebuilder:validation:Minimum=0

	// The number of successful finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`

	// +kubebuilder:validation:Minimum=0

	// The number of failed finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

/*
我们自定义了一个类型（`ConcurrencyPolicy`）来保存我们的并发策略。
它实际上只是一个字符串，但是这个类型名称提供了额外的说明，
我们还将验证附加到类型上，而不是字段上，从而使验证更易于重用。
`// +kubebuilder:validation:Enum=Allow;Forbid;Replace` 表示这个 `ConcurrencyPolicy` 只接受 `Allow`、`Forbid`、`Replace` 这三个值。
*/

// ConcurrencyPolicy describes how the job will be handled.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is AllowConcurrent.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows CronJobs to run concurrently.
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent forbids concurrent runs, skipping next run if previous
	// hasn't finished yet.
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent cancels currently running job and replaces it with a new one.
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

/*
下面，让我们设计 `status`,它用来存储观察到的状态。它包含我们希望用户或其他 controllers 能够获取的所有信息

我们将保留一份正在运行的 Job 列表，以及最后一次成功运行 Job 的时间。
注意，如上所述，我们使用 metav1.Time 而不是 time.Time 来获得稳定的序列化。
*/

// CronJobStatus defines the observed state of CronJob
type CronJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// A list of pointers to currently running jobs.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

/*
<span id="installation"></span>
最后我们只剩下了 `CronJob` 结构体，如之前说的，我们不需要修改该结构体的任何内容，
但是我们希望操作改 Kind 像操作 Kubernetes 内置资源一样，所以我们需要增加一个 mark `+kubebuilder:subresource:status`,
来声明一个status subresource, 关于 subresource 的更多信息可以参考
[k8s 文档](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#subresources)
*/

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CronJob is the Schema for the cronjobs API
type CronJob struct {
	/*
	 */
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

func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}

// +kubebuilder:docs-gen:collapse=Root Object Definitions
