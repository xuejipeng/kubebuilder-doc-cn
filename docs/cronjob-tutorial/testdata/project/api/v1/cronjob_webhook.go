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
	"github.com/robfig/cron"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// +kubebuilder:docs-gen:collapse=Go imports

/*
下一步，为我们的 webhooks 创建一个 logger。
*/

var cronjoblog = logf.Log.WithName("cronjob-resource")

/*
然后, 我们通过 manager 构建 webhook。
*/

func (r *CronJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

/*
注意，下面我们使用 kubebuilder 的 marker 语句生成 webhook manifests(配置 webhook 的 yaml 文件)，
这个 marker 会生成一个 mutating Webhook 的 manifests。

关于每个参数的解释都可以在[这里](../reference/markers/webhook.md)找到。
*/

// +kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io

/*
我们使用 `webhook.Defaulter` 接口将 webhook 的默认值设置为我们的 CRD，
这将自动生成一个 Webhook(defaulting webhooks)，并调用它的 Default 方法。

`Default` 方法用来改变接受到的内容，并设置默认值。
*/

var _ webhook.Defaulter = &CronJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CronJob) Default() {
	cronjoblog.Info("default", "name", r.Name)

	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

/*
这个 marker 负责生成一个 validating webhook manifest。
*/

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,versions=v1,name=vcronjob.kb.io

/*
我们可以通过声明式验证(declarative validation)来验证我们的 CRD，
通常情况下，声明式验证就足够了，但是有时更高级的用例需要进行复杂的验证。

例如，我们将在后面看到我们使验证(declarative validation)来验证 cron schedule(是否是 * * * * * 这种格式) 的格式是否正确，
而不是编写复杂的正则表达式来验证。

如果实现了 `webhook.Validator` 接口，将会自动生成一个 Webhook(validating webhooks) 来调用我们验证方法。

`ValidateCreate`、`ValidateUpdate` 和 `ValidateDelete` 方法分别在创建，更新和删除 resrouces 时验证它们收到的信息。
我们将 `ValidateCreate` 与 `ValidateUpdate` 方法分开，因为某些字段可能是固定不变的，
他们只能在 `ValidateCreate` 方法中被调用，这样会提高一些安全性，
`ValidateDelete` 和 `ValidateUpdate` 方法也被分开，以便在进行删除操作时进行单独的验证。
但是在这里，我们在 `ValidateDelete` 中什么也没有做，
只是对 `ValidateCreate` 和 `ValidateUpdate` 使用同一个方法进行了验证，
因为我们不需要在删除时验证任何内容。
*/

var _ webhook.Validator = &CronJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateCreate() error {
	cronjoblog.Info("validate create", "name", r.Name)

	return r.validateCronJob()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateUpdate(old runtime.Object) error {
	cronjoblog.Info("validate update", "name", r.Name)

	return r.validateCronJob()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateDelete() error {
	cronjoblog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

/*
验证 CronJob 的 name 和 spec 字段
*/

func (r *CronJob) validateCronJob() error {
	var allErrs field.ErrorList
	if err := r.validateCronJobName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateCronJobSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		r.Name, allErrs)
}

/*
一些字段是用 OpenAPI schema 方式进行验证，
你可以在 [设计API](./api-design.md) 中找到有关 kubebuilder 的验证 markers(注释前缀为`//+kubebuilder:validation`)。
你也可以通过运行 `controller-gen crd -w` 或 在[这里](/reference/markers/crd-validation.md)
找到所有关于使用 markers 验证的格式信息。
*/

func (r *CronJob) validateCronJobSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateScheduleFormat(
		r.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))
}

/*
我们在这里验证 [cron](https://en.wikipedia.org/wiki/Cron) schedule 的格式
*/

func validateScheduleFormat(schedule string, fldPath *field.Path) *field.Error {
	if _, err := cron.ParseStandard(schedule); err != nil {
		return field.Invalid(fldPath, schedule, err.Error())
	}
	return nil
}

/*
验证字符串字段的长度可以以声明方式完成。

但是 `ObjectMeta.Name` 字段是在 `apimachinery` 库下的 package 中定义的，
因此我们无法以声明性的方式对其进行验证。
*/

func (r *CronJob) validateCronJobName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The cronjob controller appends
		// a 11-character suffix to the cronjob (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore cronjob
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 52 characters")
	}
	return nil
}

// +kubebuilder:docs-gen:collapse=Validate object name
