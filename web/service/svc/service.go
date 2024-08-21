package svc

import (
	"context"
	"k8s-web/global"
	svc_req "k8s-web/model/svc/request"
	svc_res "k8s-web/model/svc/response"
	"k8s-web/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

// @Author: morris
type SvcService struct {
}

func (SvcService) GetSvcList(namespace string, keyword string) ([]svc_res.Service, error) {
	list, err := global.KubeConfigSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	serviceResList := make([]svc_res.Service, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		serviceResList = append(serviceResList, svc_res.Service{
			Name:       item.Name,
			Namespace:  item.Namespace,
			Type:       item.Spec.Type,
			ClusterIP:  item.Spec.ClusterIP,
			ExternalIP: item.Spec.ExternalIPs,
			Age:        item.CreationTimestamp.Unix(),
		})
	}
	return serviceResList, nil
}

func (SvcService) GetSvcDetail(namespace string, name string) (svcReq svc_req.Service, err error) {
	serviceK8s, err := global.KubeConfigSet.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return
	}
	servicePorts := make([]svc_req.ServicePort, 0)
	for _, port := range serviceK8s.Spec.Ports {
		servicePorts = append(servicePorts, svc_req.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: port.TargetPort.IntVal,
			NodePort:   port.NodePort,
		})
	}
	svcReq = svc_req.Service{
		Name:      serviceK8s.Name,
		Namespace: serviceK8s.Namespace,
		Labels:    utils.ToList(serviceK8s.Labels),
		Type:      serviceK8s.Spec.Type,
		Selector:  utils.ToList(serviceK8s.Spec.Selector),
		Ports:     servicePorts,
	}
	return
}

func (SvcService) DeleteSvc(namespace string, name string) error {
	return global.KubeConfigSet.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (SvcService) CreateOrUpdateSvc(serviceReq svc_req.Service) error {
	servicePorts := make([]corev1.ServicePort, 0)
	for _, port := range serviceReq.Ports {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name: port.Name,
			Port: port.Port,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: port.TargetPort,
			},
			NodePort: port.NodePort,
		})
	}
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceReq.Name,
			Namespace: serviceReq.Namespace,
			Labels:    utils.ToMap(serviceReq.Labels),
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceReq.Type,
			Selector: utils.ToMap(serviceReq.Selector),
			Ports:    servicePorts,
		},
	}
	serviceApi := global.KubeConfigSet.CoreV1().Services(service.Namespace)
	ctx := context.TODO()
	serviceK8s, err := serviceApi.Get(ctx, service.Name, metav1.GetOptions{})
	if err == nil {
		serviceK8s.Spec = service.Spec
		_, err = serviceApi.Update(ctx, serviceK8s, metav1.UpdateOptions{})
	} else {
		_, err = serviceApi.Create(ctx, &service, metav1.CreateOptions{})
	}
	return err
}
