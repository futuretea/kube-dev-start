package main

import (
	"flag"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func wait(client DeploymentClient, name string, isDelete bool) {
	for {
		// GET
		deployment, err := client.Get(name)
		if err != nil {
			if errors.IsNotFound(err) {
				break
			}
			klog.Errorf("[GET] err: %s", err)
		}
		klog.Infof("[GET] name: %s, available replicas: %d", name, deployment.Status.AvailableReplicas)
		if isDelete {
			if deployment.Status.AvailableReplicas == 0 {
				break
			}
		} else {
			if deployment.Status.AvailableReplicas == *deployment.Spec.Replicas {
				break
			}
		}
	}
}

func CURD(client DeploymentClient, name, image string) error {
	var replicas int32 = 1

	// CREATE
	_, err := client.Create(name, image, replicas)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			klog.Info("[CREATE] is already exists")
		} else {
			klog.Errorf("[CREATE] err: %s", err)
			return err
		}
	}
	wait(client, name, false)

	// UPDATE
	_, err = client.Update(name, image, replicas+1)
	if err != nil {
		klog.Errorf("[UPDATE] err: %s", err)
		return err
	}
	wait(client, name, false)

	// DELETE
	err = client.Delete(name)
	if err != nil {
		klog.Errorf("[DELETE] err: %s", err)
		return err
	}
	wait(client, name, true)
	return nil
}

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		klog.Errorf("[CONFIG] err: %s", err)
	}
	kubeConfigFile := homedir + "/.kube/config"
	kubeConfig := flag.String("kubeconfig", kubeConfigFile, "Path to a kube config.")
	namespace := flag.String("namespace", "default", "namespace")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Fatalln(err)
	}
	client := GetDeploymentClient(config, *namespace)
	err = CURD(client, "kube-dev-start", "nginx")
	if err != nil {
		klog.Errorf("[CURD] err: %s", err)
	}
}
