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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	"math/rand"
	"time"
)

const BENCHMARK_NAMESPACE = "benchmark"

func (this CoCoBenchmark) Run() {
	glog.Info("Running CocoBenchmark...")
	this.createPod()
}

func (this CoCoBenchmark) createPod() {
	name := fmt.Sprintf("test-nginx-pod-%d", rand.Uint32())
	annots := make(map[string]string)
	annots["scheduler.alpha.kubernetes.io/name"] = "poseidon"
	labels := make(map[string]string)
	labels["scheduler"] = "poseidon"
	_, err := this.Clientset.Pods(BENCHMARK_NAMESPACE).Create(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annots,
			Labels:      labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            fmt.Sprintf("container-%s", name),
				Image:           "nginx:latest",
				ImagePullPolicy: "IfNotPresent",
			}}},
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(20 * time.Second))
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
