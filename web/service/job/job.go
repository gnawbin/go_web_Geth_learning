package job

import (
	"context"
	"errors"
	"fmt"
	"k8s-web/global"
	"k8s-web/model/base"
	job_req "k8s-web/model/job/request"
	job_res "k8s-web/model/job/response"
	"k8s-web/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/pointer"
	"strings"
	"time"
)

// @Author: morris
type JobSetService struct {
}

func (JobSetService) CreateOrUpdateJob(jobReq job_req.Job) error {
	//转换为k8s结构
	podK8s := podConvert.PodReq2K8s(jobReq.Template)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobReq.Base.Name,
			Namespace: jobReq.Base.Namespace,
			Labels:    utils.ToMap(jobReq.Base.Labels),
		},
		Spec: batchv1.JobSpec{
			//Selector: &metav1.LabelSelector{
			//	MatchLabels: utils.ToMap(jobReq.Base.Selector),
			//},
			ActiveDeadlineSeconds: pointer.Int64(60),
			Completions:           &jobReq.Base.Completions,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	jobApi := global.KubeConfigSet.BatchV1().Jobs(job.Namespace)
	jobK8s, err := jobApi.Get(ctx, job.Name, metav1.GetOptions{})
	if err == nil {
		//参数校验
		jobCopy := *job
		vJobName := job.Name + "-validate"
		jobCopy.Name = vJobName
		_, err = jobApi.Create(ctx, &jobCopy, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			return err
		}
		//开启监听
		var labelSelector []string
		for k, v := range jobK8s.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		var podLabelSelector []string
		for k, v := range jobK8s.Spec.Template.Labels {
			podLabelSelector = append(podLabelSelector, fmt.Sprintf("%s=%s", k, v))
		}
		watcher, errin := jobApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if errin != nil {
			return errin
		}
		notify := make(chan error)
		go func(thisJob *batchv1.Job, notify chan error) {
			//监听删除信号后，创建
			for {
				select {
				case e := <-watcher.ResultChan():
					switch e.Type {
					case watch.Deleted:
						//删除关联Pod
						if list, errx := global.KubeConfigSet.CoreV1().Pods(jobK8s.Namespace).
							List(ctx, metav1.ListOptions{
								LabelSelector: strings.Join(podLabelSelector, ","),
							}); errx == nil {
							for _, item := range list.Items {
								//delete pod
								background := metav1.DeletePropagationBackground
								err = global.KubeConfigSet.CoreV1().Pods(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{
									GracePeriodSeconds: pointer.Int64(0),
									PropagationPolicy:  &background,
								})
							}
						}
						_, errin = jobApi.Create(ctx, thisJob, metav1.CreateOptions{})
						notify <- errin
						return
					}
				case <-time.After(5 * time.Second):
					notify <- errors.New("更新Job超时")
					return
					//fmt.Println("timeout")
				}
			}
		}(job, notify)
		//删除
		background := metav1.DeletePropagationForeground
		err = jobApi.Delete(ctx, job.Name, metav1.DeleteOptions{
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
		_, err = jobApi.Create(ctx, job, metav1.CreateOptions{})
	}
	return err
}
func (JobSetService) DeleteJob(namespace, name string) error {
	jobApi := global.KubeConfigSet.BatchV1().Jobs(namespace)
	ctx := context.TODO()
	jobK8s, err := jobApi.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//开启监听
	var labelSelector []string
	for k, v := range jobK8s.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%s=%s", k, v))
	}
	watcher, err := jobApi.Watch(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})
	if err != nil {
		return err
	}
	var podLabelSelector []string
	for k, v := range jobK8s.Spec.Template.Labels {
		podLabelSelector = append(podLabelSelector, fmt.Sprintf("%s=%s", k, v))
	}
	err = jobApi.Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	notify := make(chan error)
	go func(thisJob *batchv1.Job, notify chan error) {
		//监听删除信号后，创建
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Deleted:
					//删除关联Pod
					if list, errx := global.KubeConfigSet.CoreV1().Pods(jobK8s.Namespace).
						List(ctx, metav1.ListOptions{
							LabelSelector: strings.Join(podLabelSelector, ","),
						}); errx == nil {
						//清理job关联的Pod
						for _, item := range list.Items {
							//delete pod
							background := metav1.DeletePropagationBackground
							err = global.KubeConfigSet.CoreV1().Pods(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{
								GracePeriodSeconds: pointer.Int64(0),
								PropagationPolicy:  &background,
							})
						}
					}
					notify <- nil
					return
				}
			case <-time.After(5 * time.Second):
				notify <- errors.New("删除Job超时")
				return
			}
		}
	}(jobK8s, notify)
	select {
	case errx := <-notify:
		if errx != nil {
			return errx
		}
	}
	return nil
}
func (JobSetService) GetJobDetail(namespace, name string) (job_req.Job, error) {
	var jobReq job_req.Job
	jobK8s, err := global.KubeConfigSet.BatchV1().Jobs(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return jobReq, err
	}
	podReq := podConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: jobK8s.Spec.Template.Labels,
		},
		Spec: jobK8s.Spec.Template.Spec,
	})
	newPodLabels := make([]base.ListMapItem, 0)
	for _, label := range podReq.Base.Labels {
		if strings.Contains(label.Key, "controller-uid") {
			continue
		}
		if strings.Contains(label.Key, "job-name") {
			continue
		}
		newPodLabels = append(newPodLabels, label)
	}
	podReq.Base.Labels = newPodLabels
	jobReq = job_req.Job{
		Base: job_req.JobBase{
			Name:      jobK8s.Name,
			Namespace: jobK8s.Namespace,
			Labels:    utils.ToList(jobK8s.Labels),
			//Selector:    utils.ToList(jobK8s.Spec.Selector.MatchLabels),
			Completions: *jobK8s.Spec.Completions,
		},
		Template: podReq,
	}
	return jobReq, err
}
func (JobSetService) GetJobList(namespace, keyword string) ([]job_res.Job, error) {
	jobResList := make([]job_res.Job, 0)
	list, err := global.KubeConfigSet.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return jobResList, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		//item.Status.
		jobResList = append(jobResList, job_res.Job{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Completions: *item.Spec.Completions,
			Age:         item.CreationTimestamp.Unix(),
			Succeeded:   item.Status.Succeeded,
		})
		if item.Status.CompletionTime != nil &&
			item.Status.StartTime != nil {
			jobResList[len(jobResList)-1].Duration =
				item.Status.CompletionTime.Unix() - item.Status.StartTime.Unix()
		}
	}
	return jobResList, err
}
