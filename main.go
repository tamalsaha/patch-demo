package main

import (
	"encoding/json"
	"fmt"

	"gomodules.xyz/pointer"
	"gomodules.xyz/x/crypto/rand"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func main() {
	cur := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("stash"),
			Namespace: "default",
			Labels: map[string]string{
				"app": "patch-demo",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: pointer.Int32P(1),
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "patch-demo",
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "busybox",
							Image:           "busybox",
							ImagePullPolicy: core.PullIfNotPresent,
							Command: []string{
								"sleep",
								"3600",
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      TestSourceDataVolumeName,
									MountPath: TestSourceDataMountPath,
								},
							},
						},
					},
					Volumes: []core.Volume{
						{
							Name: TestSourceDataVolumeName,
							VolumeSource: core.VolumeSource{
								GitRepo: &core.GitRepoVolumeSource{
									Repository: "https://github.com/appscode/stash-data.git",
								},
							},
						},
					},
				},
			},
		},
	}

	curJson, err := json.Marshal(cur)
	if err != nil {
		panic(err)
	}

	mod := cur.DeepCopy()
	mod.Spec.Template.Spec.Containers[0].VolumeMounts = nil

	modJson, err := json.Marshal(mod)
	if err != nil {
		panic(err)
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, apps.StatefulSet{})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(patch))
}
