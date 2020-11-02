package main

import (
	"context"
	"sync"

	"k8s.io/klog"

	"k8s.io/client-go/rest"

	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

var (
	_                      DeploymentClient = &deploymentClient{}
	DeploymentClientOnce   sync.Once
	GlobalDeploymentClient deploymentClient
)

type DeploymentClient interface {
	Create(name, image string, replicas int32) (*appsv1.Deployment, error)
	Update(name, image string, replicas int32) (*appsv1.Deployment, error)
	Delete(name string) error
	Get(name string) (*appsv1.Deployment, error)
}

type deploymentClient struct {
	client    typedappsv1.AppsV1Interface
	namespace string
}

func buildDeployment(namespace, name, image string, replicas int32) *appsv1.Deployment {
	labels := map[string]string{
		"app": name,
	}
	deployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
	return &deployment
}

func (c *deploymentClient) Create(name, image string, replicas int32) (*appsv1.Deployment, error) {
	klog.Infof("[CREATE] name: %s, image: %s, replicas: %d", name, image, replicas)
	deployment := buildDeployment(c.namespace, name, image, replicas)
	return c.client.Deployments(c.namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
}

func (c *deploymentClient) Update(name, image string, replicas int32) (*appsv1.Deployment, error) {
	klog.Infof("[UPDATE] name: %s, image: %s, replicas: %d", name, image, replicas)
	deployment := buildDeployment(c.namespace, name, image, replicas)
	return c.client.Deployments(c.namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
}

func (c *deploymentClient) Delete(name string) error {
	klog.Infof("[DELETE] name: %s", name)
	return c.client.Deployments(c.namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (c *deploymentClient) Get(name string) (*appsv1.Deployment, error) {
	return c.client.Deployments(c.namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func GetDeploymentClient(cfg *rest.Config, namespace string) DeploymentClient {
	DeploymentClientOnce.Do(func() {
		GlobalDeploymentClient = deploymentClient{
			client:    clientset.NewForConfigOrDie(cfg).AppsV1(),
			namespace: namespace,
		}
	})
	klog.Infof("[CLIENT] namespace: %s", namespace)
	return &GlobalDeploymentClient
}
