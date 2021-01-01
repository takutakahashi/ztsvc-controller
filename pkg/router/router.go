package router

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AnnotaionZtService string = "takutakahashi.dev/zerotier-network"
var AnnotaionZtToken string = "takutakahashi.dev/zerotier-token"

func BuildSecret(svc corev1.Service, secret corev1.Secret) corev1.Secret {
	return corev1.Secret{}
}
func BuildDeployment(svc corev1.Service) appsv1.Deployment {
	var replicas int32 = 1
	token, ok := svc.Annotations[AnnotaionZtService]
	if !ok {
		token = "token"
	}
	dp := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-zt-router", svc.Name),
			Namespace:   svc.Namespace,
			Labels:      svc.Labels,
			Annotations: svc.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: metav1.SetAsLabelSelector(svc.Labels),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{},
					Containers: []corev1.Container{
						{
							Name:  "zt",
							Image: "takutakahashi/zerotier-node-daemon",
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"},
									},
								},
								{
									Name: "ZT_TOKEN",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-zt-router-secret", svc.Name),
											},
											Key: token,
										},
									},
								},
								{
									Name: "NETWORK_ID",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-zt-router-secret", svc.Name),
											},
											Key: svc.Annotations[AnnotaionZtToken],
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return dp
}
