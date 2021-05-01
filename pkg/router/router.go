package router

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var AnnotaionZtServiceEnable string = "zerotier.takutakahashi.dev/enable"
var AnnotaionZtServiceSecret string = "zerotier.takutakahashi.dev/configuration-secret"
var AnnotaionZtServiceHostname string = "zerotier.takutakahashi.dev/hostname"

func BuildResources(svc corev1.Service) error {
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	// create cm
	cm, err := BuildConfig(svc)
	cc := clientset.CoreV1().ConfigMaps(svc.Namespace)
	_, err = cc.Get(cm.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = cc.Create(&cm)
		} else {
			return err
		}
	} else {
		_, err = cc.Update(&cm)
	}
	if err != nil {
		return err
	}

	// create dp
	dp := BuildDeployment(clientset, svc)
	dc := clientset.AppsV1().Deployments(svc.Namespace)
	_, err = dc.Get(dp.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = dc.Create(&dp)
		} else {
			return err
		}
	} else {
		_, err = dc.Update(&dp)
	}
	if err != nil {
		return err
	}

	return nil
}

func BuildConfig(svc corev1.Service) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-router-config", svc.Name),
			Namespace: svc.Namespace,
		},
		Data: map[string]string{},
	}
	for k, fp := range map[string]string{
		"haproxy.cfg": "./src/haproxy.cfg.tpl",
		"envoy.yaml":  "./src/envoy.yaml.tpl",
	} {
		tmpl, err := template.New(filepath.Base(fp)).ParseFiles(fp)
		if err != nil {
			return cm, err
		}
		result := bytes.Buffer{}
		err = tmpl.Execute(&result, svc)
		if err != nil {
			return cm, err
		}
		cm.Data[k] = result.String()
	}
	return cm, nil
}
func BuildDeployment(clientset *kubernetes.Clientset, svc corev1.Service) appsv1.Deployment {
	var replicas int32 = 1
	configSecretName, ok := svc.Annotations[AnnotaionZtServiceSecret]
	if !ok {
		configSecretName = "ztsvc-secret"
	}
	labels := svc.Labels
	if svc.Labels == nil {
		labels = map[string]string{}
	}
	priviledged := true
	labels["zerotier-svc"] = svc.Name
	dpname := fmt.Sprintf("%s-zt-router", svc.Name)
	dp := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        dpname,
			Namespace:   svc.Namespace,
			Labels:      labels,
			Annotations: svc.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: fmt.Sprintf("%s-config", dpname),
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-router-config", svc.Name),
									},
								},
							},
						},
					},
					InitContainers:     []corev1.Container{},
					ServiceAccountName: "ztsvc-node-daemon",
					Containers: []corev1.Container{
						{
							Name:  "proxy",
							Image: "haproxy:2.3",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      fmt.Sprintf("%s-config", dpname),
									ReadOnly:  true,
									MountPath: "/usr/local/etc/haproxy",
								},
							},
						},
						{
							Name:  "envoy",
							Image: "envoyproxy/envoy:v1.18-latest",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      fmt.Sprintf("%s-config", dpname),
									ReadOnly:  true,
									MountPath: "/etc/envoy",
								},
							},
						},
						{
							Name:  "zt",
							Image: "takutakahashi/zerotier-node-daemon:v0.3.1",
							SecurityContext: &corev1.SecurityContext{
								Privileged: &priviledged,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "NAMESPACE",
									Value: svc.Namespace,
								},
								{
									Name:  "DOMAIN",
									Value: svc.Annotations[AnnotaionZtServiceHostname],
								},
								{
									Name:  "NODE_NAME",
									Value: svc.Name,
								},
								{
									Name: "ZT_TOKEN",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: configSecretName,
											},
											Key: "ZT_TOKEN",
										},
									},
								},
								{
									Name: "NETWORK_ID",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: configSecretName,
											},
											Key: "NETWORK_ID",
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
