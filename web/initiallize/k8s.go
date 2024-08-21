package initiallize

import (
	"context"
	"fmt"
	"io/ioutil"
	"k8s-web/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func K8S() {
	kubeconfig := ".kube/config"
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	global.KubeConfigSet = clientset
}
func isInCluster() (isInCluster bool) {
	tokenFile := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	_, err := os.Stat(tokenFile)
	if err == nil {
		isInCluster = true
	}
	return
}
func K8SWithDiscovery() {
	if isInCluster() {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		global.KubeConfigSet = clientset
	} else {
		K8S()
	}
}
func K8SWithToken() {
	cAData, err := ioutil.ReadFile("k8s_use/identity/ca.crt")
	if err != nil {
		panic(err)
	}
	config := &rest.Config{
		Host:            "https://192.168.1.16:6443",
		BearerTokenFile: "k8s_use/identity/token",
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cAData,
		},
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	list, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}

}
