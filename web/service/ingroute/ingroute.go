package ingroute

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-web/global"
	ingroute_k8s "k8s-web/model/ingroute/k8s"
	ingroute_req "k8s-web/model/ingroute/request"
	ingroute_res "k8s-web/model/ingroute/response"
	"k8s-web/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type IngRouteService struct {
}

func (IngRouteService) GetIngRouteMiddlewareList(namespace string) (mwList []string, err error) {
	//查询middleware 列表
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/middlewares", namespace)
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	mwList = make([]string, 0)
	var middlewareList ingroute_k8s.MiddlewareList
	err = json.Unmarshal(raw, &middlewareList)
	if err != nil {
		return
	}
	for _, item := range middlewareList.Items {
		mwList = append(mwList, item.Metadata.Name)
	}
	return
}

func (IngRouteService) DeleteIngRoute(namespace, name string) error {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes/%s", namespace, name)
	_, err := global.KubeConfigSet.RESTClient().Delete().AbsPath(url).DoRaw(context.TODO())
	return err
}

func (IngRouteService) CreateOrUpdateIngRoute(ingressRouteReq ingroute_req.IngressRoute) (err error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", ingressRouteReq.Namespace)
	//convert to k8s structure
	ingressRoute := ingroute_k8s.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefik.io/v1alpha1",
			Kind:       "IngressRoute",
		},
		Metadata: metav1.ObjectMeta{
			Name:      ingressRouteReq.Name,
			Namespace: ingressRouteReq.Namespace,
			Labels:    utils.ToMap(ingressRouteReq.Labels),
		},
		Spec: ingressRouteReq.IngressRouteSpec,
	}
	//已存在则更新
	result, err := json.Marshal(ingressRoute)
	if err != nil {
		return
	}
	//查询是否存在
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).Name(ingressRouteReq.Name).DoRaw(context.TODO())
	if err == nil {
		//修改row
		var ingressRouteK8s ingroute_k8s.IngressRoute
		err = json.Unmarshal(raw, &ingressRouteK8s)
		if err != nil {
			return
		}
		//update
		ingressRouteK8s.Spec = ingressRoute.Spec
		resultx, errMar := json.Marshal(ingressRouteK8s)
		if errMar != nil {
			return errMar
		}
		_, err = global.KubeConfigSet.RESTClient().Put().
			Name(ingressRouteK8s.Metadata.Name).
			AbsPath(url).
			Body(resultx).DoRaw(context.TODO())
	} else {
		//create
		_, err = global.KubeConfigSet.RESTClient().Post().AbsPath(url).Body(result).DoRaw(context.TODO())
	}
	return
}

func (IngRouteService) GetIngRouteList(namespace, keyword string) (ingRouteResList []ingroute_res.IngressRoute, err error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", namespace)
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	if err != nil {
		return
	}
	//tmpMap := make(map[string]any)
	var ingRouteList ingroute_k8s.IngressRouteList
	err = json.Unmarshal(raw, &ingRouteList)
	if err != nil {
		return
	}
	ingRouteResList = make([]ingroute_res.IngressRoute, 0)
	for _, item := range ingRouteList.Items {
		if !strings.Contains(item.Metadata.Name, keyword) {
			continue
		}
		ingRouteResList = append(ingRouteResList, ingroute_res.IngressRoute{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Age:       item.Metadata.CreationTimestamp.Unix(),
		})
	}
	return
}
func (IngRouteService) GetIngRouteDetail(namespace, name string) (ingRouteReq *ingroute_req.IngressRoute, err error) {
	url := fmt.Sprintf("/apis/traefik.io/v1alpha1/namespaces/%s/ingressroutes", namespace)
	url = url + "/" + name
	raw, err := global.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(context.TODO())
	if err != nil {
		return
	}
	var ingRoute ingroute_k8s.IngressRoute
	err = json.Unmarshal(raw, &ingRoute)
	if err != nil {
		return
	}
	ingRouteReq = &ingroute_req.IngressRoute{
		Name:             ingRoute.Metadata.Name,
		Namespace:        ingRoute.Metadata.Namespace,
		Labels:           utils.ToList(ingRoute.Metadata.Labels),
		IngressRouteSpec: ingRoute.Spec,
	}
	return
}
