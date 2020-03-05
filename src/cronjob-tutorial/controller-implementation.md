# 实现一个 controller

我们的 CronJob controller 的基本逻辑是：

1. 按名称加载 CronJob

2. 列出所有 active jobs，并更新状态

3. 根据历史记录清理 old jobs

4. 检查 Job 是否已被 suspended（如果被 suspended，请不要执行任何操作）

5. 获取到下一次要 schedule 的 Job

6. 运行新的 Job, 确定新 Job 没有超过 deadline 时间，且不会被我们 concurrency 规则 block

7. 如果 Job 正在运行或者它应该下次运行，请重新排队

{{#literatego ./testdata/project/controllers/cronjob_controller.go}}

很明显，我们现在有了一个可以运行 controller ，让我们针对集群进行测试，然后，如果我们没有任何问题，把它部署到集群中！
让我们针对集群进行测试，然后，如果我们没有任何问题，请对其进行部署！
