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
	fmt.Println("Process started...")
	ctx := context.Background()

	k8sClient, err := GetClient()
	if err != nil {
		fmt.Println("Failed to create client object", err)
	}

	if os.Getenv("CREATE_CLUSTER_ROLE") == "true" {

		fmt.Println("Creating ClusterRole.")
		clusterRole := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-gitops-cluster-role",
			},
		}

		if err = k8sClient.Create(ctx, clusterRole); err != nil {
			fmt.Println("Failed to create ClusterRole", err)
			return
		}

		fmt.Println("Succesfully created the ClusterRole: ", clusterRole.Name)

		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-gitops-cluster-role-binding",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      "argocd-argocd-application-controller",
					Namespace: os.Getenv("POD_NAMESPACE"),
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "ClusterRole",
				Name:     "test-gitops-cluster-role",
			},
		}

		if err = k8sClient.Create(ctx, clusterRoleBinding); err != nil {
			fmt.Println("Failed to create ClusterRoleBinding", err)
			return
		}

		fmt.Println("Succesfully created the ClusterRoleBinding: ", clusterRoleBinding.Name)

	} else {

		fmt.Println("Creating Role")

		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-gitops-role",
				Namespace: os.Getenv("POD_NAMESPACE"),
			},
		}

		if err = k8sClient.Create(ctx, role); err != nil {
			fmt.Println("Failed to create Role", err)
			return
		}

		fmt.Println("Succesfully created the Role: ", role.Name, " in namespace: ", role.Namespace)
	}

	time.Sleep(10 * time.Minute)
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
