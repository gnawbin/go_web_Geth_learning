package secret

import (
	"k8s-web/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
import secretreq "k8s-web/model/secret/request"

// @Author: morris
type Req2K8s struct {
}

func (Req2K8s) SecretReq2K8sConvert(secret secretreq.Secret) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
			Labels:    utils.ToMap(secret.Labels),
		},
		Type:       secret.Type,
		StringData: utils.ToMap(secret.Data),
	}
}
