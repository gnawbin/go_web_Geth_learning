package ingress

import (
	"context"
	"k8s-web/global"
	ingress_req "k8s-web/model/ingress/request"
	ingress_res "k8s-web/model/ingress/response"
	"k8s-web/utils"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type IngressService struct {
}

func (IngressService) GetIngressDetail(namespace string, name string) (ingresReq ingress_req.Ingress, err error) {
	//详情转换
	ingressK8s, err := global.KubeConfigSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return
	}
	rules := make([]ingress_req.IngressRule, 0)
	for _, rule := range ingressK8s.Spec.Rules {
		rules = append(rules, ingress_req.IngressRule{
			Host:  rule.Host,
			Value: rule.IngressRuleValue,
		})
	}
	ingresReq = ingress_req.Ingress{
		Name:      ingressK8s.Name,
		Namespace: ingressK8s.Namespace,
		Labels:    utils.ToList(ingressK8s.Labels),
		Rules:     rules,
	}
	return
}

func (IngressService) GetIngressList(namespace string, keyword string) ([]ingress_res.Ingress, error) {
	list, err := global.KubeConfigSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	ingressResList := make([]ingress_res.Ingress, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		hosts := make([]string, 0)
		//ports := make([]string, 0)
		for _, rule := range item.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}
		ingressResList = append(ingressResList, ingress_res.Ingress{
			Name:      item.Name,
			Namespace: item.Namespace,
			Hosts:     strings.Join(hosts, ","),
			Age:       item.CreationTimestamp.Unix(),
		})
	}
	return ingressResList, nil
}

func (IngressService) DeleteIngress(namespace string, name string) error {
	return global.KubeConfigSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (IngressService) CreateOrUpdateIngress(ingressReq ingress_req.Ingress) error {
	ingressRules := make([]networkingv1.IngressRule, 0)
	for _, rule := range ingressReq.Rules {
		ingressRules = append(ingressRules, networkingv1.IngressRule{
			Host:             rule.Host,
			IngressRuleValue: rule.Value,
		})
	}
	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressReq.Name,
			Namespace: ingressReq.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: ingressRules,
		},
	}
	ingressesApi := global.KubeConfigSet.NetworkingV1().Ingresses(ingress.Namespace)
	ctx := context.TODO()
	ingressK8s, err := ingressesApi.Get(ctx, ingress.Name, metav1.GetOptions{})
	if err == nil {
		ingressK8s.Spec = ingress.Spec
		_, err = ingressesApi.Update(ctx, ingressK8s, metav1.UpdateOptions{})
	} else {
		_, err = ingressesApi.Create(ctx, &ingress, metav1.CreateOptions{})
	}
	return err
}
