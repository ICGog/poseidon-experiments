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

package main

import (
	"flag"
	"github.com/ICGog/poseidon-experiments/pkg/benchmark"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	benchmarkName string
	kubeVersion   string
	kubeConfig    string
)

func createClientset() *kubernetes.Clientset {
	var config *rest.Config
	var err error
	if kubeVersion == "1.6" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		config, err = clientcmd.DefaultClientConfig.ClientConfig()
	}
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func init() {
	flag.StringVar(&benchmarkName, "benchmarkName", "coco", "The name of the benchmark to run. Options: coco")
	flag.StringVar(&kubeVersion, "kubeVersion", "1.5", "Kubernetes version. Options: 1.5 | 1.6")
	flag.StringVar(&kubeConfig, "kubeConfig", "/root/admin.conf", "Path to kubeConfig file")
	flag.Parse()
}

func main() {
	var runbenchmark benchmark.Benchmark
	switch benchmarkName {
	case "coco":
		runbenchmark = benchmark.CoCoBenchmark{
			Clientset: createClientset(),
		}
	default:
		glog.Fatalf("Unexpected benchmark name %s", benchmarkName)
	}
	runbenchmark.Setup()
	defer runbenchmark.Destroy()
	runbenchmark.Run()
}
