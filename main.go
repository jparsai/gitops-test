package main

import (
	"context"
	"fmt"
	"os"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	fmt.Println("####### Running #######")
	ctx := context.Background()

	k8sClient, err := GetClient()
	if err != nil {
		fmt.Println("Failed to create client object", err)
	}

	if os.Getenv("CLUSTER_ROLE") != "" {
		fmt.Println("Creating Cluster Role")
		clusterRole := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-gitops-cluster-role",
			},
		}
		err = k8sClient.Create(ctx, clusterRole)
	} else {
		fmt.Println("Creating Role")
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-gitops-role",
				Namespace: os.Getenv("POD_NAMESPACE"),
			},
		}
		err = k8sClient.Create(ctx, role)
	}

	if err != nil {
		fmt.Println("Failed to create object", err)
	}
	time.Sleep(5 * time.Minute)
}

func GetClient() (client.Client, error) {
	overrides := clientcmd.ConfigOverrides{}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	err = rbacv1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}
