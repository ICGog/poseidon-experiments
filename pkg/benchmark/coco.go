// Poseidon
// Copyright (c) The Poseidon Authors.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT
// LIMITATION ANY IMPLIED WARRANTIES OR CONDITIONS OF TITLE, FITNESS FOR
// A PARTICULAR PURPOSE, MERCHANTABLITY OR NON-INFRINGEMENT.
//
// See the Apache Version 2.0 License for specific language governing
// permissions and limitations under the License.

package benchmark

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"math/rand"
	"time"
)

const BENCHMARK_NAMESPACE = "benchmark"

func (this CoCoBenchmark) Run() {
	glog.Info("Running CocoBenchmark...")
	// 8 x 8
	this.createJob(fmt.Sprintf("test-cpuspin-%d", rand.Uint32()),
		"icgog/cpuspin:dev", "Turtle", "900m", "128Mi", 8,
		[]string{"/bin/sh", "-c", "/tmp/cpu_spin 60"})
	// 8 x 8
	// this.createJob(fmt.Sprintf("test-memstream-1m-%d", rand.Uint32()),
	// 	"icgog/memstream:dev", "Rabbit", "900m", "128Mi", 8,
	// 	[]string{"/bin/sh", "-c", "/tmp/mem_stream 1048576"})
	// 8 x 8
	// this.createJob(fmt.Sprintf("test-memstream-50m-%d", rand.Uint32()),
	// 	"icgog/memstream:dev", "Devil", "900m", "128Mi", 8,
	// 	[]string{"/bin/sh", "-c", "/tmp/mem_stream 52428800"})
	// 2 x 6
	// this.createJob(fmt.Sprintf("test-iostream-write-%d", rand.Uint32()),
	// 	"icgog/iostream:dev", "Devil", "400m", "4096Mi", 6,
	// 	[]string{"/bin/sh", "-c", "/usr/bin/fio /tmp/fio-seqwrite.fio"})
	// 2 x 6
	// this.createJob(fmt.Sprintf("test-iostream-read-%d", rand.Uint32()),
	// 	"icgog/iostream:dev", "Devil", "400m", "4096Mi", 6,
	// 	[]string{"/bin/sh", "-c", "/usr/bin/fio /tmp/fio-seqread.fio"})
	//	this.createPod(fmt.Sprintf("test-nginx-pod-%d", rand.Uint32()), "nginx:latest", "Sheep")
}

func (this CoCoBenchmark) createPod(name, image, tasktype, requestCPU, requestMem string, cmd []string) {
	annots := make(map[string]string)
	annots["scheduler.alpha.kubernetes.io/name"] = "poseidon"
	labels := make(map[string]string)
	labels["scheduler"] = "poseidon"
	labels["taskType"] = tasktype
	_, err := this.Clientset.Pods(BENCHMARK_NAMESPACE).Create(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annots,
			Labels:      labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            name,
				Image:           image,
				ImagePullPolicy: "IfNotPresent",
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceName(v1.ResourceCPU):    resource.MustParse(requestCPU),
						v1.ResourceName(v1.ResourceMemory): resource.MustParse(requestMem),
					},
				},
				Command: cmd,
			}},
			RestartPolicy: "Never",
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(80 * time.Second))
}

func (this CoCoBenchmark) createJob(name, image, tasktype, requestCPU, requestMem string, numtasks int32, cmd []string) {
	annots := make(map[string]string)
	annots["scheduler.alpha.kubernetes.io/name"] = "poseidon"
	labels := make(map[string]string)
	labels["scheduler"] = "poseidon"
	labels["taskType"] = tasktype
	_, err := this.Clientset.Batch().Jobs(BENCHMARK_NAMESPACE).Create(&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annots,
			Labels:      labels,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &numtasks,
			Completions: &numtasks,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: annots,
					Labels:      labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: "IfNotPresent",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceName(v1.ResourceCPU):    resource.MustParse(requestCPU),
									v1.ResourceName(v1.ResourceMemory): resource.MustParse(requestMem),
								}},
							Command: cmd,
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(80 * time.Second))
}

func (this CoCoBenchmark) Setup() {
	_, err := this.Clientset.Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: BENCHMARK_NAMESPACE},
	})
	if err != nil {
		panic(err)
	}
}

func (this CoCoBenchmark) Destroy() {
	// Delete namespace
	err := this.Clientset.Namespaces().Delete(BENCHMARK_NAMESPACE, &metav1.DeleteOptions{})
	// Delete all pods
	err = this.Clientset.Pods(BENCHMARK_NAMESPACE).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
}
