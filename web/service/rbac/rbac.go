package rbac

import (
	"context"
	"k8s-web/global"
	rbac_req "k8s-web/model/rbac/request"
	rbac_res "k8s-web/model/rbac/response"
	"k8s-web/utils"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// @Author: morris
type RbacService struct {
}

func (RbacService) CreateOrUpdateRb(rbReq rbac_req.RoleBinding) error {
	ctx := context.TODO()
	//创建 cluster role binding
	if rbReq.Namespace == "" {
		rbK8sReq := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rbReq.Name,
				Namespace: rbReq.Namespace,
				Labels:    utils.ToMap(rbReq.Labels),
			},
			Subjects: func(saList []rbac_req.ServiceAccount) []rbacv1.Subject {
				subjects := make([]rbacv1.Subject, len(saList))
				for index, item := range saList {
					subjects[index] = rbacv1.Subject{
						Name:      item.Name,
						Kind:      "User",
						Namespace: item.Namespace,
					}
				}
				return subjects
			}(rbReq.Subjects),
			RoleRef: rbacv1.RoleRef{
				Name:     rbReq.RoleRef,
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
			},
		}
		clusterRbApi := global.KubeConfigSet.RbacV1().ClusterRoleBindings()
		if cluterRoleSrc, err := clusterRbApi.Get(ctx, rbReq.Name, metav1.GetOptions{}); err != nil {
			_, err := clusterRbApi.Create(ctx, rbK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			cluterRoleSrc.ObjectMeta.Labels = rbK8sReq.Labels
			cluterRoleSrc.Subjects = rbK8sReq.Subjects
			cluterRoleSrc.RoleRef = rbK8sReq.RoleRef
			_, err := clusterRbApi.Update(ctx, cluterRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		rbK8sReq := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rbReq.Name,
				Namespace: rbReq.Namespace,
				Labels:    utils.ToMap(rbReq.Labels),
			},
			Subjects: func(saList []rbac_req.ServiceAccount) []rbacv1.Subject {
				subjects := make([]rbacv1.Subject, len(saList))
				for index, item := range saList {
					subjects[index] = rbacv1.Subject{
						Name:      item.Name,
						Kind:      "User",
						Namespace: item.Namespace,
					}
				}
				return subjects
			}(rbReq.Subjects),
			RoleRef: rbacv1.RoleRef{
				Name:     rbReq.RoleRef,
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
			},
		}
		rbApi := global.KubeConfigSet.RbacV1().RoleBindings(rbK8sReq.Namespace)
		if rbK8sSrc, err := rbApi.Get(ctx, rbReq.Name, metav1.GetOptions{}); err != nil {
			_, err = rbApi.Create(ctx, rbK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			rbK8sSrc.ObjectMeta.Labels = rbK8sReq.Labels
			rbK8sSrc.Subjects = rbK8sReq.Subjects
			rbK8sSrc.RoleRef = rbK8sReq.RoleRef
			_, err = rbApi.Update(ctx, rbK8sSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (RbacService) DeleteRb(namespace, name string) error {
	ctx := context.TODO()
	if namespace != "" {
		err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else {
		err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// 查看RoleBinding详情
func (RbacService) GetRbDetail(namespace, name string) (rbac_req.RoleBinding, error) {
	ctx := context.TODO()
	rbReq := rbac_req.RoleBinding{}
	if namespace != "" {
		rbK8s, err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return rbReq, err
		}
		rbReq.Name = rbK8s.Name
		rbReq.Namespace = rbK8s.Namespace
		rbReq.Labels = utils.ToList(rbK8s.Labels)
		rbReq.RoleRef = rbK8s.RoleRef.Name
		rbReq.Subjects = func(subjects []rbacv1.Subject) []rbac_req.ServiceAccount {
			saList := make([]rbac_req.ServiceAccount, len(subjects))
			for i, subject := range subjects {
				saList[i] = rbac_req.ServiceAccount{
					Name:      subject.Name,
					Namespace: subject.Namespace,
				}
			}
			return saList
		}(rbK8s.Subjects)
	} else {
		rbK8s, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return rbReq, err
		}
		rbReq.Name = rbK8s.Name
		rbReq.Namespace = rbK8s.Namespace
		rbReq.Labels = utils.ToList(rbK8s.Labels)
		rbReq.RoleRef = rbK8s.RoleRef.Name
		rbReq.Subjects = func(subjects []rbacv1.Subject) []rbac_req.ServiceAccount {
			saList := make([]rbac_req.ServiceAccount, len(subjects))
			for i, subject := range subjects {
				saList[i] = rbac_req.ServiceAccount{
					Name:      subject.Name,
					Namespace: subject.Namespace,
				}
			}
			return saList
		}(rbK8s.Subjects)
	}
	return rbReq, nil
}

// 查看RoleBinding列表
func (RbacService) GetRbList(namespace, keyword string) ([]rbac_res.RoleBinding, error) {
	ctx := context.TODO()
	rbResList := make([]rbac_res.RoleBinding, 0)
	if namespace != "" {
		list, err := global.KubeConfigSet.RbacV1().RoleBindings(namespace).
			List(ctx, metav1.ListOptions{})
		if err != nil {
			return rbResList, err
		}
		for _, item := range list.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			rbResList = append(rbResList, rbac_res.RoleBinding{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	} else {
		list, err := global.KubeConfigSet.RbacV1().ClusterRoleBindings().
			List(ctx, metav1.ListOptions{})
		if err != nil {
			return rbResList, err
		}
		for _, item := range list.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			rbResList = append(rbResList, rbac_res.RoleBinding{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	}
	return rbResList, nil
}

func (RbacService) CreateOrUpdateRole(roleReq rbac_req.Role) error {
	ctx := context.TODO()
	if roleReq.Namespace == "" {
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		//创建 cluster role
		clusterRoleK8sReq := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleReq.Name,
				Namespace: roleReq.Namespace,
				Labels:    utils.ToMap(roleReq.Labels),
			},
			Rules: roleReq.Rules,
		}
		if clusterRoleSrc, err := clusterRoleApi.Get(ctx, roleReq.Name, metav1.GetOptions{}); err != nil {
			_, err = clusterRoleApi.Create(ctx, clusterRoleK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			clusterRoleSrc.ObjectMeta.Labels = clusterRoleK8sReq.Labels
			clusterRoleSrc.Rules = clusterRoleK8sReq.Rules
			_, err = clusterRoleApi.Update(ctx, clusterRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		// 创建 ns role
		nsRoleK8sReq := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleReq.Name,
				Namespace: roleReq.Namespace,
				Labels:    utils.ToMap(roleReq.Labels),
			},
			Rules: roleReq.Rules,
		}
		roleApi := global.KubeConfigSet.RbacV1().Roles(nsRoleK8sReq.Namespace)
		if nsRoleSrc, err := roleApi.Get(ctx, nsRoleK8sReq.Name, metav1.GetOptions{}); err != nil {
			_, err := roleApi.Create(ctx, nsRoleK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			nsRoleSrc.Labels = nsRoleK8sReq.Labels
			nsRoleSrc.Rules = nsRoleK8sReq.Rules
			_, err := roleApi.Update(ctx, nsRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (RbacService) GetRoleDetail(namespace, name string) (roleReq rbac_req.Role, err error) {
	ctx := context.TODO()
	if namespace != "" {
		roleApi := global.KubeConfigSet.RbacV1().Roles(namespace)
		roleK8s, errx := roleApi.Get(ctx, name, metav1.GetOptions{})
		if errx != nil {
			err = errx
			return
		}
		roleReq = rbac_req.Role{
			Name:      roleK8s.Name,
			Namespace: roleK8s.Namespace,
			Labels:    utils.ToList(roleK8s.Labels),
			Rules:     roleK8s.Rules,
		}
	} else {
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		roleK8s, errx := clusterRoleApi.Get(ctx, name, metav1.GetOptions{})
		if errx != nil {
			err = errx
			return
		}
		roleReq = rbac_req.Role{
			Name:      roleK8s.Name,
			Namespace: roleK8s.Namespace,
			Labels:    utils.ToList(roleK8s.Labels),
			Rules:     roleK8s.Rules,
		}
	}
	return
}

func (RbacService) GetRoleList(namespace, keyword string) ([]rbac_res.Role, error) {
	ctx := context.TODO()
	resRoleList := make([]rbac_res.Role, 0)
	if namespace != "" {
		roleApi := global.KubeConfigSet.RbacV1().Roles(namespace)
		roleK8sList, err := roleApi.List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleK8sList.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			resRoleList = append(resRoleList, rbac_res.Role{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	} else {
		clusterRoleApi := global.KubeConfigSet.RbacV1().ClusterRoles()
		roleK8sList, err := clusterRoleApi.List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleK8sList.Items {
			if !strings.Contains(item.Name, keyword) {
				continue
			}
			resRoleList = append(resRoleList, rbac_res.Role{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	}
	return resRoleList, nil
}
func (RbacService) DeleteRole(namespace, name string) error {
	ctx := context.TODO()
	if namespace != "" {
		err := global.KubeConfigSet.RbacV1().Roles(namespace).
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else {
		err := global.KubeConfigSet.RbacV1().ClusterRoles().
			Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (RbacService) CreateServiceAccount(saReq rbac_req.ServiceAccount) error {
	saK8s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saReq.Name,
			Namespace: saReq.Namespace,
			Labels:    utils.ToMap(saReq.Labels),
		},
	}
	_, err := global.KubeConfigSet.CoreV1().ServiceAccounts(saK8s.Namespace).
		Create(context.TODO(), saK8s, metav1.CreateOptions{})
	return err
}

func (RbacService) DeleteServiceAccount(namespace, name string) error {
	return global.KubeConfigSet.CoreV1().ServiceAccounts(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (RbacService) GetServiceAccountList(namespace, keyword string) ([]rbac_res.ServiceAccount, error) {
	list, err := global.KubeConfigSet.CoreV1().ServiceAccounts(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	resList := make([]rbac_res.ServiceAccount, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		resList = append(resList, rbac_res.ServiceAccount{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Unix(),
		})
	}
	return resList, err
}
