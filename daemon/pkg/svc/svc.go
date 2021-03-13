package svc

import (
	"context"
	"fmt"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/zerotier"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Ensure(node zerotier.Node, namespace string) error {
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	extSvcName := fmt.Sprintf("%s-zt-router-ext", node.Name)
	sc := clientset.CoreV1().Services(namespace)
	s := &v1.Service{}
	s, err = sc.Get(context.TODO(), extSvcName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		s = &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      extSvcName,
				Namespace: namespace,
				Annotations: map[string]string{
					"external-dns.alpha.kubernetes.io/hostname": node.Domain,
				},
			},
			Spec: v1.ServiceSpec{
				Type:         v1.ServiceTypeExternalName,
				ExternalName: node.IP,
			},
		}
		_, err = sc.Create(context.TODO(), s, metav1.CreateOptions{})
		return err
	} else if err != nil {
		return err
	}
	_, err = sc.Update(context.TODO(), s, metav1.UpdateOptions{})
	return err
}
