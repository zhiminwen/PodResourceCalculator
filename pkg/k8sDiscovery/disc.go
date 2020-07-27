package k8sDiscovery

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func K8s() (kubernetes.Interface, *rest.Config, error) {
	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster == true {
		log.Infof("inside cluster, using in-cluster configuration")
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Errorf("Failed to get incluster config:%v", err)
			return nil, nil, err
		}
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Errorf("Failed to construct the clientSet:%v", err)
			return nil, nil, err
		}
		return clientSet, config, nil
	}

	log.Infof("outside of cluster")
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Errorf("Failed to build the config:%v", err)
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Failed to construct the clientSet:%v", err)
		return nil, nil, err
	}
	return clientSet, config, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

//for testing
func GetServerVersion(clientSet kubernetes.Interface) (string, error) {
	version, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		log.Errorf("Failed to get server version:%v", err)
		return "", err
	}
	return fmt.Sprintf("%s", version), nil
}
