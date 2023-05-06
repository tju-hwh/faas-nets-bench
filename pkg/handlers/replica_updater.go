// Copyright (c) Alex Ellis 2017. All rights reserved.
// Copyright 2020 OpenFaaS Author(s)
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openfaas/faas-netes/pkg/k8s"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/openfaas/faas-provider/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MakeReplicaUpdater updates desired count of replicas
func MakeReplicaUpdater(defaultNamespace string, clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas")

		vars := mux.Vars(r)

		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		req := types.ScaleServiceRequest{}

		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		options := metav1.GetOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
		}

		deployment, err := clientset.AppsV1().Deployments(lookupNamespace).Get(context.TODO(), functionName, options)
		//hwh 扩容指定节点
		deployment.Spec.Template.Spec.NodeName = "kube-node-7"
		//hwh
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to lookup function deployment " + functionName))
			log.Println(err)
			return
		}

		oldReplicas := *deployment.Spec.Replicas
		replicas := int32(req.Replicas)

		log.Printf("Set replicas - %s %s, %d/%d\n", functionName, lookupNamespace, replicas, oldReplicas)

		deployment.Spec.Replicas = &replicas

		_, err = clientset.AppsV1().Deployments(lookupNamespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		//
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to update function deployment " + functionName))
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// 改动3
func MakeReplicaUpdater3(defaultNamespace string, clientset *kubernetes.Clientset, factory k8s.FunctionFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas")

		vars := mux.Vars(r)

		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		req := types.ScaleServiceRequest{}

		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		options := metav1.GetOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
		}

		deployment, err := clientset.AppsV1().Deployments(lookupNamespace).Get(context.TODO(), functionName, options)
		//hwh 扩容指定节点
		deployment.Spec.Template.Spec.NodeName = "kube-node-5"
		//hwh
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to lookup function deployment " + functionName))
			log.Println(err)
			return
		}

		oldReplicas := *deployment.Spec.Replicas
		replicas := int32(req.Replicas)

		log.Printf("Set replicas - %s %s, %d/%d\n", functionName, lookupNamespace, replicas, oldReplicas)

		deployment.Spec.Replicas = &replicas

		//_, err = clientset.AppsV1().Deployments(lookupNamespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		//魔改,失败
		if replicas == 10 {
			log.Println("开始魔改：")
			deployment.Spec.Template.Spec.NodeName = "kube-node-5"
			deployment.Spec.Replicas = int32p(3)

			_, err = factory.Client.AppsV1().Deployments(lookupNamespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
			if err != nil {
				log.Printf("create deployment err:%e", err)
			}

			deployment.Spec.Template.Spec.NodeName = "kube-node-7"
			deployment.Spec.Replicas = int32p(7)
			_, err = factory.Client.AppsV1().Deployments(lookupNamespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
			if err != nil {
				log.Printf("create deployment err:%e", err)
			}

		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to update function deployment " + functionName))
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// 改动4 指定节点扩容
func MakeReplicaUpdater4(defaultNamespace string, clientset *kubernetes.Clientset, factory k8s.FunctionFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas 指定节点扩容")

		vars := mux.Vars(r)

		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")

		lookupNamespace := defaultNamespace

		if len(namespace) > 0 {
			lookupNamespace = namespace
		}

		req := types.ScaleServiceRequest{}

		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		options := metav1.GetOptions{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
		}
		//旧的deployment
		deployment, err := clientset.AppsV1().Deployments(lookupNamespace).Get(context.TODO(), functionName, options)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to lookup function deployment " + functionName))
			log.Println(err)
			return
		}
		oldReplicas := *deployment.Spec.Replicas
		replicas := int32(req.Replicas)
		log.Printf("Set replicas - %s %s, %d/%d\n", functionName, lookupNamespace, replicas, oldReplicas)
		//更新旧的deployment副本数
		//deployment.Spec.Replicas = &replicas
		//更新deployment
		//_, err = clientset.AppsV1().Deployments(lookupNamespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		////
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Unable to update function deployment " + functionName))
		//	log.Println(err)
		//	return
		//}
		//w.WriteHeader(http.StatusAccepted)
		//

		//pod := &corev1.Pod{
		//
		//	TypeMeta: metav1.TypeMeta{
		//		Kind:       "Pod",
		//		APIVersion: "v1",
		//	},
		//	ObjectMeta: metav1.ObjectMeta{
		//		Name:        functionName + "-pod-" + RandomText(8),
		//		Annotations: deployment.Spec.Template.ObjectMeta.Annotations,
		//		Labels:      deployment.Spec.Template.ObjectMeta.Labels,
		//	},
		//	Spec: corev1.PodSpec{
		//		NodeName:   "deployPosition",
		//		Containers: deployment.Spec.Template.Spec.Containers,
		//	},
		//}
		// 实例化一个Pod数据结构
		log.Println(functionName, " 222")

		annotations := make(map[string]string)
		annotations["prometheus.io.scrape"] = "false"
		labels := map[string]string{
			"faas_function": functionName,
		}
		pod := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:        functionName + "-pod-" + RandomText(8),
				Annotations: annotations,
				Labels:      labels,
			},
			Spec: corev1.PodSpec{
				NodeName: "deployPosition",
				Containers: []corev1.Container{
					{
						Name:  functionName,
						Image: "1007287849/hello-python: latest",
						//Resources:       deployment.Spec.Template.Spec.Containers[0].Resources,
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{
								Name:          "http",
								Protocol:      corev1.ProtocolSCTP,
								ContainerPort: 8080,
							},
						},
					},
				},
			},
		}

		nodemap := make(map[string]int)
		//nodemap["kube-node-7"] = 2
		nodemap["kube-node-5"] = 1
		log.Println(namespace, "  111")
		log.Println(pod)

		for s, i := range nodemap {
			log.Println("进入nodemap，开始扩容")
			pod.Spec.NodeName = s
			for j := 0; j < i; j++ {
				result, err := factory.Client.CoreV1().Pods("openfaas-fn-hwh").Create(context.TODO(), pod, metav1.CreateOptions{})

				if err != nil {
					log.Printf("create pod err: %e", err)
					panic(err.Error())
				} else {
					log.Printf("Create pod %s success \n", result.Name)
				}
			}
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

// 改动5 指定节点扩容
func MakeReplicaUpdater5(defaultNamespace string, clientset *kubernetes.Clientset, factory k8s.FunctionFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas")
		vars := mux.Vars(r)
		functionName := vars["name"]
		q := r.URL.Query()
		namespace := q.Get("namespace")
		lookupNamespace := defaultNamespace
		if len(namespace) > 0 {
			lookupNamespace = namespace
		}
		req := types.ScaleServiceRequest{}
		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}
		//options := metav1.GetOptions{
		//	TypeMeta: metav1.TypeMeta{
		//		Kind:       "Deployment",
		//		APIVersion: "apps/v1",
		//	},
		//}
		//clientset.AppsV1().ReplicaSets(lookupNamespace).Get(context.TODO(),functionName,options)
		//deployment, err := clientset.AppsV1().Deployments(lookupNamespace).Get(context.TODO(), functionName, options)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Unable to lookup function deployment " + functionName))
		//	log.Println(err)
		//	return
		//}
		//oldReplicas := *deployment.Spec.Replicas
		oldReplicas := int32(1)
		replicas := int32(req.Replicas)
		log.Printf("Set replicas - %s %s, %d/%d\n", functionName, lookupNamespace, replicas, oldReplicas)

		//增加指定节点扩容逻辑
		difference := replicas - oldReplicas
		log.Println("Update replicas 指定节点扩容")
		for i := 0; i < int(difference); i++ {
			log.Println("第 ", i, "次扩容")
			DeployPod("hello-python", "kube-node-7", factory)
		}

		var nodes []Node
		Db.Find(&nodes)
		log.Println(nodes[0])
		log.Println("数据库操作ok")
		//结束指定节点扩容逻辑

		//deployment.Spec.Replicas = &replicas
		//_, err = clientset.AppsV1().Deployments(lookupNamespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		////
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Unable to update function deployment " + functionName))
		//	log.Println(err)
		//	return
		//}

		w.WriteHeader(http.StatusAccepted)

	}
}

//hwh

func RandomText(len int) (randText string) {
	//set random seed
	rand.Seed(time.Now().UnixNano())
	var captcha string
	for i := 0; i < len; i++ {
		//generate number form 0 to 9
		num := rand.Intn(10)
		//change number to string
		captcha += strconv.Itoa(num)
	}
	return captcha
}

//
// deployPod 指定pod放置位置的放置
func DeployPod(functionName string, deployPosition string, factory k8s.FunctionFactory) int {

	//获取需要放置的function所需的资源量
	cpuNum, memoryNum := 100, 40
	imageName := "1007287849/hello-python:latest"
	//将int拼接为字符串 str格式
	cpuLimit := fmt.Sprintf("%dm", cpuNum)
	memoryLimit := fmt.Sprintf("%dMi", memoryNum)

	labels := map[string]string{
		"faas_function":     functionName,
		"pod-template-hash": "7dd4f496f8",
	}
	//将资源量转化为对应数据类型
	resources, _ := createMyResources(cpuLimit, memoryLimit)

	annotations := make(map[string]string)
	annotations["prometheus.io.scrape"] = "false"
	// 得到deployment的客户端
	// podClient := factory.
	// 	CoreV1().
	// 	Pods("openfaas-fn-zmy")

	// 实例化一个Pod数据结构
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        functionName + "-pod-" + RandomText(8),
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: corev1.PodSpec{
			NodeName: deployPosition,
			Containers: []corev1.Container{
				{
					Name:            functionName,
					Image:           imageName,
					Resources:       *resources,
					ImagePullPolicy: "IfNotPresent",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 8080,
						},
					},
				},
			},
		},
	}

	//直接创建
	result, err := factory.Client.CoreV1().Pods("openfaas-fn-hwh").Create(context.TODO(), pod, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
		log.Printf("create pod error: %e", err.Error())
	}

	log.Printf("Create pod %s \n", result.Name)
	// resultSp, err := podClient.Get(context.TODO(), result.Name, metav1.GetOptions{})
	// if err!=nil {
	// 	panic(err.Error())
	// }
	return 1
}

func createMyResources(cpuLimit string, memoryLimit string) (*corev1.ResourceRequirements, error) {
	//新建一个resources资源变量
	resources := &corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}
	//
	cpuRequest := "10m"
	memoryRequest := "10Mi"

	// Set Memory limits
	if len(memoryLimit) > 0 {
		qty, err := resource.ParseQuantity(memoryLimit)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceMemory] = qty
	}

	if len(memoryRequest) > 0 {
		qty, err := resource.ParseQuantity(memoryRequest)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceMemory] = qty
	}

	// Set CPU limits
	if len(cpuLimit) > 0 {
		qty, err := resource.ParseQuantity(cpuLimit)
		if err != nil {
			return resources, err
		}
		resources.Limits[corev1.ResourceCPU] = qty
	}

	if len(cpuRequest) > 0 {
		qty, err := resource.ParseQuantity(cpuRequest)
		if err != nil {
			return resources, err
		}
		resources.Requests[corev1.ResourceCPU] = qty
	}

	return resources, nil
}
