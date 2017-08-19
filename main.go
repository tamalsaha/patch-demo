package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/appscode/go/crypto/rand"
	gt "github.com/appscode/go/types"
	"github.com/appscode/log"
	"github.com/mattbaird/jsonpatch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube/config"))
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}
	kubeClient := clientset.NewForConfigOrDie(config)

	ko := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("stash"),
			Namespace: "default",
			Labels: map[string]string{
				"app": "patch-demo",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: gt.Int32P(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "patch-demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "busybox",
							Image:           "busybox",
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Command: []string{
								"sleep",
								"3600",
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      TestSourceDataVolumeName,
									MountPath: TestSourceDataMountPath,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: TestSourceDataVolumeName,
							VolumeSource: apiv1.VolumeSource{
								GitRepo: &apiv1.GitRepoVolumeSource{
									Repository: "https://github.com/appscode/stash-data.git",
								},
							},
						},
					},
				},
			},
		},
	}
	ko, err = kubeClient.AppsV1beta1().Deployments(ko.Namespace).Create(ko)
	if err != nil {
		log.Fatalln(err)
	}

	oJson, err := json.Marshal(ko)
	if err != nil {
		log.Fatalln(err)
	}

	// ----------------------------------------------------------------------------------------------------------------

	if ko.Annotations == nil {
		ko.Annotations = map[string]string{}
	}
	ko.Annotations["example.com"] = "123"
	ko.Spec.Replicas = gt.Int32P(2)
	ko.Spec.Template.Spec.Containers = append(ko.Spec.Template.Spec.Containers, apiv1.Container{
		Name:            "bnew",
		Image:           "busybox",
		ImagePullPolicy: apiv1.PullIfNotPresent,
		Command: []string{
			"sleep",
			"3600",
		},
		VolumeMounts: []apiv1.VolumeMount{
			{
				Name:      TestSourceDataVolumeName,
				MountPath: TestSourceDataMountPath,
			},
		},
	})
	mJson, err := json.Marshal(ko)
	if err != nil {
		log.Fatalln(err)
	}

	patch, err := jsonpatch.CreatePatch(oJson, mJson)
	if err != nil {
		log.Fatalln(err)
	}
	pb, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(pb))

	final, err := kubeClient.AppsV1beta1().Deployments(ko.Namespace).Patch(ko.Name, types.JSONPatchType, pb)
	if err != nil {
		log.Fatalln(err)
	}

	fb, err := json.MarshalIndent(final, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(fb))
}
