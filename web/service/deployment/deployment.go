package deployment

import (
	"context"
	"k8s-web/global"
	deployment_req "k8s-web/model/deployment/request"
	deployment_res "k8s-web/model/deployment/response"
	"k8s-web/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type DeploymentService struct {
}

func (DeploymentService) CreateOrUpdateDeployment(deploymentReq deployment_req.Deployment) error {
	//转换为k8s结构
	podK8s := podConvert.PodReq2K8s(deploymentReq.Template)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentReq.Base.Name,
			Namespace: deploymentReq.Base.Namespace,
			Labels:    utils.ToMap(deploymentReq.Base.Labels),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentReq.Base.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: utils.ToMap(deploymentReq.Base.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}
	ctx := context.TODO()
	deploymentApi := global.KubeConfigSet.AppsV1().Deployments(deployment.Namespace)
	deploymentK8s, err := deploymentApi.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err == nil {
		deploymentK8s.Spec = deployment.Spec
		_, err = deploymentApi.Update(ctx, deploymentK8s, metav1.UpdateOptions{})
	} else {
		_, err = deploymentApi.Create(ctx, deployment, metav1.CreateOptions{})
	}
	return err
}
func (DeploymentService) DeleteDeployment(namespace, name string) error {
	return global.KubeConfigSet.AppsV1().Deployments(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}
func (DeploymentService) GetDeploymentDetail(namespace, name string) (deployment_req.Deployment, error) {
	var deloymentReq deployment_req.Deployment
	deploymentK8s, err := global.KubeConfigSet.AppsV1().Deployments(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return deloymentReq, err
	}
	podReq := podConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: deploymentK8s.Spec.Template.Labels,
		},
		Spec: deploymentK8s.Spec.Template.Spec,
	})
	deloymentReq = deployment_req.Deployment{
		Base: deployment_req.DeploymentBase{
			Name:      deploymentK8s.Name,
			Namespace: deploymentK8s.Namespace,
			Replicas:  *deploymentK8s.Spec.Replicas,
			Labels:    utils.ToList(deploymentK8s.Labels),
			Selector:  utils.ToList(deploymentK8s.Spec.Selector.MatchLabels),
		},
		Template: podReq,
	}
	return deloymentReq, err
}
func (DeploymentService) GetDeploymentList(namespace, keyword string) ([]deployment_res.Deployment, error) {
	deploymentResList := make([]deployment_res.Deployment, 0)
	list, err := global.KubeConfigSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return deploymentResList, err
	}
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		deploymentResList = append(deploymentResList, deployment_res.Deployment{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Unix(),
			Replicas:  *item.Spec.Replicas,
			Ready:     item.Status.Replicas,
			Available: item.Status.AvailableReplicas,
			UpToDate:  item.Status.UpdatedReplicas,
		})
	}
	return deploymentResList, err
}
