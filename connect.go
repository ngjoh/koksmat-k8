package main

import (
	"os"
	"path/filepath"
	"runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ConnectToCluster connects to a Kubernetes cluster, either from inside a pod or locally.
func ConnectToCluster() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Check if we're inside a Kubernetes pod
	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// If not in a pod, use the local kubeconfig file
		home := getHomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// getHomeDir returns the home directory depending on the OS
func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE") // Windows
	}
	return os.Getenv("HOME") // Linux, macOS
}
