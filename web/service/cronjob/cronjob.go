package cronjob

import (
	"context"
	"errors"
	"fmt"
	"k8s-web/global"
	cronjob_req "k8s-web/model/cronjob/request"
	cronjob_res "k8s-web/model/cronjob/response"
	"k8s-web/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
	"time"
)

// @Author: morris
type CronJobSetService struct {
}

func (CronJobSetService) CreateOrUpdateCronJob(cronJobReq cronjob_req.CronJob) error {
	//转换为k8s结构
	podK8s := podConvert.PodReq2K8s(cronJobReq.Template)
	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronJobReq.Base.Name,
			Namespace: cronJobReq.Base.Namespace,
			Labels:    utils.ToMap(cronJobReq.Base.Labels),
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   cronJobReq.Base.Schedule,
			Suspend:                    &cronJobReq.Base.Suspend,
			ConcurrencyPolicy:          cronJobReq.Base.ConcurrencyPolicy,
			SuccessfulJobsHistoryLimit: &cronJobReq.Base.SuccessfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     &cronJobReq.Base.FailedJobsHistoryLimit,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					BackoffLimit: &cronJobReq.Base.JobBase.BackoffLimit,
					Completions:  &cronJobReq.Base.JobBase.Completions,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: podK8s.ObjectMeta,
						Spec:       podK8s.Spec,
					},
				},
			},
		},
	}
	ctx := context.TODO()
	cronjobApi := global.KubeConfigSet.BatchV1().CronJobs(cronJob.Namespace)
	cronJobK8s, err := cronjobApi.Get(ctx, cronJob.Name, metav1.GetOptions{})
	if err == nil {
		//参数校验
		cronJobCopy := *cronJob
		vCronJobName := cronJobCopy.Name + "-validate"
		cronJobCopy.Name = vCronJobName
		_, err = cronjobApi.Create(ctx, &cronJobCopy, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			return err
		}
		//开启监听
		var labelSelector []string
		for k, v := range cronJobK8s.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		var podLabelSelector []string
		for k, v := range cronJobK8s.Spec.JobTemplate.Spec.Template.Labels {
			podLabelSelector = append(podLabelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		//监听cronjob删除状态
		watcher, errin := cronjobApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if errin != nil {
			return errin
		}
		notify := make(chan error)
		go func(thisCronJob *batchv1.CronJob, notify chan error) {
			//监听删除信号后，创建
			for {
				select {
				case e := <-watcher.ResultChan():
					switch e.Type {
					case watch.Deleted:
						//删除关联Job
						_, errin = cronjobApi.Create(ctx, thisCronJob, metav1.CreateOptions{})
						notify <- errin
						return
					}
				case <-time.After(5 * time.Second):
					notify <- errors.New("更新Job超时")
					return
					//fmt.Println("timeout")
				}
			}
		}(cronJob, notify)
		//删除
		background := metav1.DeletePropagationForeground
		err = cronjobApi.Delete(ctx, cronJob.Name, metav1.DeleteOptions{
			PropagationPolicy: &background,
		})
		if err != nil {
			return err
		}
		//监听删除后重新创建的信息
		select {
		case errx := <-notify:
			if errx != nil {
				return errx
			}
		}
	} else {
		_, err = cronjobApi.Create(ctx, cronJob, metav1.CreateOptions{})
	}
	return err
}
func (CronJobSetService) DeleteCronJob(namespace, name string) error {
	cronJobApi := global.KubeConfigSet.BatchV1().CronJobs(namespace)
	ctx := context.TODO()
	cronJobK8s, err := cronJobApi.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//开启监听
	var labelSelector []string
	for k, v := range cronJobK8s.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
	}
	watcher, err := cronJobApi.Watch(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})
	if err != nil {
		return err
	}
	err = cronJobApi.Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	notify := make(chan error)
	go func(thisJob *batchv1.CronJob, notify chan error) {
		//监听删除信号后，创建
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Deleted:
					notify <- nil
					return
				}
			case <-time.After(5 * time.Second):
				notify <- errors.New("删除CronJob超时")
				return
			}
		}
	}(cronJobK8s, notify)
	select {
	case errx := <-notify:
		if errx != nil {
			return errx
		}
	}
	return nil
}
func (CronJobSetService) GetCronJobDetail(namespace, name string) (cronjob_req.CronJob, error) {
	var cronJobReq cronjob_req.CronJob
	cronJobK8s, err := global.KubeConfigSet.BatchV1().CronJobs(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return cronJobReq, err
	}
	podReq := podConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: cronJobK8s.Spec.JobTemplate.Spec.Template.Labels,
		},
		Spec: cronJobK8s.Spec.JobTemplate.Spec.Template.Spec,
	})
	cronJobReq = cronjob_req.CronJob{
		Base: cronjob_req.CronJobBase{
			Name:                       cronJobK8s.Name,
			Namespace:                  cronJobK8s.Namespace,
			Labels:                     utils.ToList(cronJobK8s.Labels),
			Schedule:                   cronJobK8s.Spec.Schedule,
			ConcurrencyPolicy:          cronJobK8s.Spec.ConcurrencyPolicy,
			Suspend:                    *cronJobK8s.Spec.Suspend,
			SuccessfulJobsHistoryLimit: *cronJobK8s.Spec.SuccessfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     *cronJobK8s.Spec.SuccessfulJobsHistoryLimit,
			JobBase: cronjob_req.JobBase{
				BackoffLimit: *cronJobK8s.Spec.JobTemplate.Spec.BackoffLimit,
				Completions:  *cronJobK8s.Spec.JobTemplate.Spec.Completions,
			},
		},
		Template: podReq,
	}
	return cronJobReq, err
}
func (CronJobSetService) GetCronJobList(namespace, keyword string) ([]cronjob_res.CronJob, error) {
	cronJobResList := make([]cronjob_res.CronJob, 0)
	list, err := global.KubeConfigSet.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return cronJobResList, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		cronJob := cronjob_res.CronJob{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Unix(),
			Suspend:   item.Spec.Suspend,
			Schedule:  item.Spec.Schedule,
			Active:    len(item.Status.Active),
		}
		if item.Status.LastScheduleTime != nil {
			cronJob.LastSchedule = item.Status.LastScheduleTime.Unix()
		}
		//item.Status.
		cronJobResList = append(cronJobResList, cronJob)
	}
	return cronJobResList, err
}
