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
我们从一些 import 开始，正如你看到的，我们使用的 package 比脚手架帮我们生成的多，在使用它们时，我们会逐一讨论。
*/
package controllers

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-logr/logr"
	"github.com/robfig/cron"
	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batch "tutorial.kubebuilder.io/project/api/v1"
)

/*
接下来，我们需要一个 Clock 字段，它帮我们在测试中伪装计时。
*/

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Clock
}

/*
我们模拟时钟，以便在测试时更容易跳转，
当前时间只是调用了 `time.Now` 函数。
*/
type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

// clock knows how to get the current time.
// It can be used to fake out timing for testing.
type Clock interface {
	Now() time.Time
}

// +kubebuilder:docs-gen:collapse=Clock

/*
请注意，我们需要更多的 RBAC 权限 -- 由于我们现在正在创建和管理 Job，因此我们需要这些权限，
这意味着要添加更多 [markers](/reference/markers/rbac.md)，所以我们增加了最下面两行。
*/

// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

/*
现在， 我们到了 controller 的核心部分 -- 实现 reconciler 的逻辑
*/
var (
	scheduledTimeAnnotation = "batch.tutorial.kubebuilder.io/scheduled-at"
)

func (r *CronJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("cronjob", req.NamespacedName)

	/*
			### 1: 按 namespace 加载 CronJob

			我们使用 client 获取 CronJob。
		    所有的 client 方法都将 context（用来取消请求）作为其第一个参数，
		    并将所讨论的 object 作为其最后一个参数。 Get 方法有点特殊，
		    因为它使用 [`NamespacedName`](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/client#ObjectKey)
		    作为中间参数（大多数没有中间参数，如下所示）。

		    最后，许多 client 方法也采用可变参数选项(也就是 "...")。
	*/
	var cronJob batch.CronJob
	if err := r.Get(ctx, req.NamespacedName, &cronJob); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	/*
			### 2: 列出所有 active jobs，并更新状态
			要完全更新我们的状态，我们需要列出此 namespace 中属于此 CronJob 的所有 Job。
		    与 Get 方法类似，我们可以使用 List 方法列出 Job。
		    注意，我们使用可变参数选项来设置 `client.InNamespace` 和 `client.MatchingFields`。
	*/
	var childJobs kbatch.JobList
	if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}

	/*
			当得到所有的 Job 后，我们把 Job 的状态分为 active、successful和 failed,
		    并跟踪他们最近的运行情况，以便将其记录在 status 中。
		    请记住，status 应该可以从整体的状态重新构造，
		    因此从 root object 的状态读取信息通常不是一个好主意。
		    相反，您应该在每次运行时重新构建它。
		    这就是我们在这里要做的。

			我们可以使用 status conditions 来检查作业是“完成”、成功或失败。
			我们将把这种逻辑放在帮助程序中，以使我们的代码更整洁。
	*/

	// find the active list of jobs
	var activeJobs []*kbatch.Job
	var successfulJobs []*kbatch.Job
	var failedJobs []*kbatch.Job
	var mostRecentTime *time.Time // find the last run so we can update the status

	/*
			如果一项工作的 "succeeded" 或 "failed" 的 condition 标记为 "true"，我们认为该工作 "finished"。
		    `Status.conditions` 使我们可以向 objects 添加可扩展的状态信息，
			其他人和 controller 可以通过检查这些状态信息以确定 Job 完成和健康状况。
	*/
	isJobFinished := func(job *kbatch.Job) (bool, kbatch.JobConditionType) {
		for _, c := range job.Status.Conditions {
			if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
				return true, c.Type
			}
		}

		return false, ""
	}
	// +kubebuilder:docs-gen:collapse=isJobFinished

	/*
		我们将使用匿名方法从创建 Job 时添加的 annotation 中获取到 Job 计划执行的时间。
	*/
	getScheduledTimeForJob := func(job *kbatch.Job) (*time.Time, error) {
		timeRaw := job.Annotations[scheduledTimeAnnotation]
		if len(timeRaw) == 0 {
			return nil, nil
		}

		timeParsed, err := time.Parse(time.RFC3339, timeRaw)
		if err != nil {
			return nil, err
		}
		return &timeParsed, nil
	}
	// +kubebuilder:docs-gen:collapse=getScheduledTimeForJob

	for i, job := range childJobs.Items {
		_, finishedType := isJobFinished(&job)
		switch finishedType {
		case "": // ongoing
			activeJobs = append(activeJobs, &childJobs.Items[i])
		case kbatch.JobFailed:
			failedJobs = append(failedJobs, &childJobs.Items[i])
		case kbatch.JobComplete:
			successfulJobs = append(successfulJobs, &childJobs.Items[i])
		}

		// We'll store the launch time in an annotation, so we'll reconstitute that from
		// the active jobs themselves.
		scheduledTimeForJob, err := getScheduledTimeForJob(&job)
		if err != nil {
			log.Error(err, "unable to parse schedule time for child job", "job", &job)
			continue
		}
		if scheduledTimeForJob != nil {
			if mostRecentTime == nil {
				mostRecentTime = scheduledTimeForJob
			} else if mostRecentTime.Before(*scheduledTimeForJob) {
				mostRecentTime = scheduledTimeForJob
			}
		}
	}

	if mostRecentTime != nil {
		cronJob.Status.LastScheduleTime = &metav1.Time{Time: *mostRecentTime}
	} else {
		cronJob.Status.LastScheduleTime = nil
	}
	cronJob.Status.Active = nil
	for _, activeJob := range activeJobs {
		jobRef, err := ref.GetReference(r.Scheme, activeJob)
		if err != nil {
			log.Error(err, "unable to make reference to active job", "job", activeJob)
			continue
		}
		cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)
	}

	/*  在这里，我们将以略高的日志记录级别记录观察到的 Job 数量，以进行调试。
	请注意，我们是使用固定消息在键值对中附加额外的信息，而不是使用字符串格式。
	这样可以更轻松地过滤和查询日志行。
	*/
	log.V(1).Info("job count", "active jobs", len(activeJobs), "successful jobs", len(successfulJobs), "failed jobs", len(failedJobs))

	/*
		我们根据我们得到的时间更新 CRD 的状态。
		和以前一样，我们使用 client。
		为了更新 subresource 的状态，我们将使用 client 的 `Status().Update` 方法

		subresource status 会忽略对 spec 的更改,
		因此与其他任何 update 发生冲突的可能性较小, 并且它可以具有单独的权限
		状态子资源会忽略对规范的更改，因此与其他任何更新发生冲突的可能性较小，并且可以具有单独的权限。
	*/
	if err := r.Status().Update(ctx, &cronJob); err != nil {
		log.Error(err, "unable to update CronJob status")
		return ctrl.Result{}, err
	}

	/*
		一旦我们正确 update 了我们的 status，我们可以确保整体的状态是符合我们在 spec 中指定的。

		### 3: 根据历史记录清理过期 jobs

		首先，我们会尝试清理过期的 jobs，以免留下太多闲杂事。
	*/

	// NB: deleting these is "best effort" -- if we fail on a particular one,
	// we won't requeue just to finish the deleting.
	if cronJob.Spec.FailedJobsHistoryLimit != nil {
		sort.Slice(failedJobs, func(i, j int) bool {
			if failedJobs[i].Status.StartTime == nil {
				return failedJobs[j].Status.StartTime != nil
			}
			return failedJobs[i].Status.StartTime.Before(failedJobs[j].Status.StartTime)
		})
		for i, job := range failedJobs {
			if int32(i) >= int32(len(failedJobs))-*cronJob.Spec.FailedJobsHistoryLimit {
				break
			}
			if err := r.Delete(ctx, job, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete old failed job", "job", job)
			} else {
				log.V(0).Info("deleted old failed job", "job", job)
			}
		}
	}

	if cronJob.Spec.SuccessfulJobsHistoryLimit != nil {
		sort.Slice(successfulJobs, func(i, j int) bool {
			if successfulJobs[i].Status.StartTime == nil {
				return successfulJobs[j].Status.StartTime != nil
			}
			return successfulJobs[i].Status.StartTime.Before(successfulJobs[j].Status.StartTime)
		})
		for i, job := range successfulJobs {
			if int32(i) >= int32(len(successfulJobs))-*cronJob.Spec.SuccessfulJobsHistoryLimit {
				break
			}
			if err := r.Delete(ctx, job, client.PropagationPolicy(metav1.DeletePropagationBackground)); (err) != nil {
				log.Error(err, "unable to delete old successful job", "job", job)
			} else {
				log.V(0).Info("deleted old successful job", "job", job)
			}
		}
	}

	/* ### 4: 检查 Job 是否已被 suspended

	如果这个 object 已被 suspended，且我们不想运行任何其他 Jobs，我们会立即 return。
	这对调试 Job 的问题非常有用。
	*/

	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		log.V(1).Info("cronjob suspended, skipping")
		return ctrl.Result{}, nil
	}

	/*
		### 5: 获取到下一次要 schedule 的 Job

		如果 Job 没有被暂停，则需要计算下一次 schedule 的 Job，以及是否有尚未处理的 Job。
	*/

	/*
			我们将使用有用的 cron 库来计算下一个 scheduled 时间。
			我们将从上次 Job 开始的时间计算下一次运行的时间，如果找不到上次运行的时间，就创建一个 CronJob。

			如果错过的 Job 数量太多，并且我们没有设置任何 deadlines，我们将释放这个 Job，
		    以免造成 controller 重启。

			否则，我们将只返回错过的 Job（我们将使用最后一个运行的 Job）和下一次要运行的 Job，
			以便让我们知道何时该重新进行 reconcile。
	*/
	getNextSchedule := func(cronJob *batch.CronJob, now time.Time) (lastMissed time.Time, next time.Time, err error) {
		sched, err := cron.ParseStandard(cronJob.Spec.Schedule)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("Unparseable schedule %q: %v", cronJob.Spec.Schedule, err)
		}

		// for optimization purposes, cheat a bit and start from our last observed run time
		// we could reconstitute this here, but there's not much point, since we've
		// just updated it.
		var earliestTime time.Time
		if cronJob.Status.LastScheduleTime != nil {
			earliestTime = cronJob.Status.LastScheduleTime.Time
		} else {
			earliestTime = cronJob.ObjectMeta.CreationTimestamp.Time
		}
		if cronJob.Spec.StartingDeadlineSeconds != nil {
			// controller is not going to schedule anything below this point
			schedulingDeadline := now.Add(-time.Second * time.Duration(*cronJob.Spec.StartingDeadlineSeconds))

			if schedulingDeadline.After(earliestTime) {
				earliestTime = schedulingDeadline
			}
		}
		if earliestTime.After(now) {
			return time.Time{}, sched.Next(now), nil
		}

		starts := 0
		for t := sched.Next(earliestTime); !t.After(now); t = sched.Next(t) {
			lastMissed = t
			// An object might miss several starts. For example, if
			// controller gets wedged on Friday at 5:01pm when everyone has
			// gone home, and someone comes in on Tuesday AM and discovers
			// the problem and restarts the controller, then all the hourly
			// jobs, more than 80 of them for one hourly scheduledJob, should
			// all start running with no further intervention (if the scheduledJob
			// allows concurrency and late starts).
			//
			// However, if there is a bug somewhere, or incorrect clock
			// on controller's server or apiservers (for setting creationTimestamp)
			// then there could be so many missed start times (it could be off
			// by decades or more), that it would eat up all the CPU and memory
			// of this controller. In that case, we want to not try to list
			// all the missed start times.
			starts++
			if starts > 100 {
				// We can't get the most recent times so just return an empty slice
				return time.Time{}, time.Time{}, fmt.Errorf("Too many missed start times (> 100). Set or decrease .spec.startingDeadlineSeconds or check clock skew.")
			}
		}
		return lastMissed, sched.Next(now), nil
	}
	// +kubebuilder:docs-gen:collapse=getNextSchedule

	// figure out the next times that we need to create
	// jobs at (or anything we missed).
	missedRun, nextRun, err := getNextSchedule(&cronJob, r.Now())
	if err != nil {
		log.Error(err, "unable to figure out CronJob schedule")
		// we don't really care about requeuing until we get an update that
		// fixes the schedule, so don't return an error
		return ctrl.Result{}, nil
	}

	/*
		我们将准备最终的 request 以重新排队直到下一个工作，然后确定 Job 是否真的需要运行。
		在下一个 Job 来到之前，我们会把最终的 request 放到队列中，以判断 Job 是否真的需要运行。
	*/
	scheduledResult := ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())} // save this so we can re-use it elsewhere
	log = log.WithValues("now", r.Now(), "next run", nextRun)

	/*
		### 6: 运行新的 Job, 确定新 Job 没有超过 deadline 时间，且不会被我们 concurrency 规则 block

		如果我们错过了一个 Job 的运行时间点，但是 Job 还在 deadline 时间内，我们需要再次运行这个 Job
	*/
	if missedRun.IsZero() {
		log.V(1).Info("no upcoming scheduled times, sleeping until next")
		return scheduledResult, nil
	}

	// make sure we're not too late to start the run
	log = log.WithValues("current run", missedRun)
	tooLate := false
	if cronJob.Spec.StartingDeadlineSeconds != nil {
		tooLate = missedRun.Add(time.Duration(*cronJob.Spec.StartingDeadlineSeconds) * time.Second).Before(r.Now())
	}
	if tooLate {
		log.V(1).Info("missed starting deadline for last run, sleeping till next")
		// TODO(directxman12): events
		return scheduledResult, nil
	}

	/*
		如果我们不得不运行一个 Job，那么需要等现有 Job 完成之后，替换现有 Job 或添加一个新 Job。
		如果我们的信息由于 cache 的延迟而过时，那么我们将在获得 up-to-date 信息时重新排队。
	*/
	// figure out how to run this job -- concurrency policy might forbid us from running
	// multiple at the same time...
	if cronJob.Spec.ConcurrencyPolicy == batch.ForbidConcurrent && len(activeJobs) > 0 {
		log.V(1).Info("concurrency policy blocks concurrent runs, skipping", "num active", len(activeJobs))
		return scheduledResult, nil
	}

	// ...or instruct us to replace existing ones...
	if cronJob.Spec.ConcurrencyPolicy == batch.ReplaceConcurrent {
		for _, activeJob := range activeJobs {
			// we don't care if the job was already deleted
			if err := r.Delete(ctx, activeJob, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete active job", "job", activeJob)
				return ctrl.Result{}, err
			}
		}
	}

	/*
		一旦弄清楚如何处理现有 Job，我们便会真正创建所需的工作。
	*/

	/*
		我们需要基于 CronJob 的 template 构建 Job, 我们将从 template 复制 spec，然后复制一些基本的 字段。

		然后，我们将设置 "scheduled time" annotation，
		以便我们可以在每次 reconcile 时重新构造我们的 `LastScheduleTime` 字段。

		最后，我们需要设置所有者参考。 这使 Kubernetes 垃圾收集器在删除 CronJob 时清理 Job，
		并允许控制器运行时找出给定作业被更改（added，deleted，completes）时需要协调哪个 cronjob。
	*/
	constructJobForCronJob := func(cronJob *batch.CronJob, scheduledTime time.Time) (*kbatch.Job, error) {
		// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
		name := fmt.Sprintf("%s-%d", cronJob.Name, scheduledTime.Unix())

		job := &kbatch.Job{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
				Name:        name,
				Namespace:   cronJob.Namespace,
			},
			Spec: *cronJob.Spec.JobTemplate.Spec.DeepCopy(),
		}
		for k, v := range cronJob.Spec.JobTemplate.Annotations {
			job.Annotations[k] = v
		}
		job.Annotations[scheduledTimeAnnotation] = scheduledTime.Format(time.RFC3339)
		for k, v := range cronJob.Spec.JobTemplate.Labels {
			job.Labels[k] = v
		}
		if err := ctrl.SetControllerReference(cronJob, job, r.Scheme); err != nil {
			return nil, err
		}

		return job, nil
	}
	// +kubebuilder:docs-gen:collapse=constructJobForCronJob

	// actually make the job...
	job, err := constructJobForCronJob(&cronJob, missedRun)
	if err != nil {
		log.Error(err, "unable to construct job from template")
		// don't bother requeuing until we get a change to the spec
		return scheduledResult, nil
	}

	// ...and create it on the cluster
	if err := r.Create(ctx, job); err != nil {
		log.Error(err, "unable to create Job for CronJob", "job", job)
		return ctrl.Result{}, err
	}

	log.V(1).Info("created Job for CronJob run", "job", job)

	/*
		### 7: 如果 Job 正在运行或者它应该下次运行，请重新排队

		最后，我们将返回上面准备的结果，这表示我们想在下次运行需要重新排队时使用。
		这被视为 maximum deadline -- 如果两者之间有其他变化，例如我们的 Job 开始或完成，
		或者我们得到一个修改信息，我们需要尽快 reconcile。
	*/
	// we'll requeue once we see the running job, and update our status
	return scheduledResult, nil
}

/*
### Setup

最后，我们将更新我们的设置。
为了使我们的 reconciler 可以按其所有者快速查找 Jobs，我们需要一个 index。
我们声明一个 index key，以后可以与 client 一起使用它作为伪字段名称，然后描述如何从 Job 对象中提取 index key。
索引器将自动为我们处理 namespace，因此，如果 Job 具有 CronJob 所有者，我们只需提取所有者名称。

此外，我们会通知 manager 该 controller 拥有一些 Jobs，
以便在作业发生更改，删除等情况时，它将自动在基础 CronJob 上调用 Reconcile。
*/
var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = batch.GroupVersion.String()
)

func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// set up a real clock, since we're not in a test
	if r.Clock == nil {
		r.Clock = realClock{}
	}

	if err := mgr.GetFieldIndexer().IndexField(&kbatch.Job{}, jobOwnerKey, func(rawObj runtime.Object) []string {
		// grab the job object, extract the owner...
		job := rawObj.(*kbatch.Job)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		// ...make sure it's a CronJob...
		if owner.APIVersion != apiGVStr || owner.Kind != "CronJob" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batch.CronJob{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
