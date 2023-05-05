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
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	types "github.com/openfaas/faas-provider/types"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// initialReplicasCount how many replicas to start of creating for a function
const initialReplicasCount = 1

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler(functionNamespace string, factory k8s.FunctionFactory) http.HandlerFunc {
	secrets := k8s.NewSecretsClient(factory.Client)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := types.FunctionDeployment{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to unmarshal request: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			wrappedErr := fmt.Errorf("validation failed: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		namespace := functionNamespace
		if len(request.Namespace) > 0 {
			namespace = request.Namespace
		}

		existingSecrets, err := secrets.GetSecrets(namespace, request.Secrets)
		if err != nil {
			wrappedErr := fmt.Errorf("unable to fetch secrets: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		deploymentSpec, specErr := makeDeploymentSpec(request, existingSecrets, factory)

		var profileList []k8s.Profile
		if request.Annotations != nil {
			profileNamespace := factory.Config.ProfilesNamespace
			profileList, err = factory.GetProfiles(ctx, profileNamespace, *request.Annotations)
			if err != nil {
				wrappedErr := fmt.Errorf("failed create Deployment spec: %s", err.Error())
				log.Println(wrappedErr)
				http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
				return
			}
		}
		for _, profile := range profileList {
			factory.ApplyProfile(profile, deploymentSpec)
		}

		if specErr != nil {
			wrappedErr := fmt.Errorf("failed create Deployment spec: %s", specErr.Error())
			log.Println(wrappedErr)
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		deploy := factory.Client.AppsV1().Deployments(namespace)

		_, err = deploy.Create(context.TODO(), deploymentSpec, metav1.CreateOptions{})
		if err != nil {
			wrappedErr := fmt.Errorf("unable create Deployment: %s", err.Error())
			log.Println(wrappedErr)
			http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Deployment created: %s.%s\n", request.Service, namespace)

		service := factory.Client.CoreV1().Services(namespace)
		serviceSpec := makeServiceSpec(request, factory)
		_, err = service.Create(context.TODO(), serviceSpec, metav1.CreateOptions{})

		if err != nil {
			wrappedErr := fmt.Errorf("failed create Service: %s", err.Error())
			log.Println(wrappedErr)
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Service created: %s.%s\n", request.Service, namespace)

		w.WriteHeader(http.StatusAccepted)
	}
}

func makeDeploymentSpec(request types.FunctionDeployment, existingSecrets map[string]*apiv1.Secret, factory k8s.FunctionFactory) (*appsv1.Deployment, error) {
	envVars := buildEnvVars(&request)

	initialReplicas := int32p(initialReplicasCount)
	labels := map[string]string{
		"faas_function": request.Service,
	}

	if request.Labels != nil {
		if min := getMinReplicaCount(*request.Labels); min != nil {
			initialReplicas = min
		}
		for k, v := range *request.Labels {
			labels[k] = v
		}
	}

	nodeSelector := createSelector(request.Constraints)

	resources, resourceErr := createResources(request)

	if resourceErr != nil {
		return nil, resourceErr
	}

	var imagePullPolicy apiv1.PullPolicy
	switch factory.Config.ImagePullPolicy {
	case "Never":
		imagePullPolicy = apiv1.PullNever
	case "IfNotPresent":
		imagePullPolicy = apiv1.PullIfNotPresent
	default:
		imagePullPolicy = apiv1.PullAlways
	}

	annotations := buildAnnotations(request)

	probes, err := factory.MakeProbes(request)
	if err != nil {
		return nil, err
	}

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
			Labels: map[string]string{
				"faas_function": request.Service,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"faas_function": request.Service,
				},
			},
			Replicas: initialReplicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32p(10),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        request.Service,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: apiv1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []apiv1.Container{
						{
							Name:  request.Service,
							Image: request.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: factory.Config.RuntimeHTTPPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env:             envVars,
							Resources:       *resources,
							ImagePullPolicy: imagePullPolicy,
							LivenessProbe:   probes.Liveness,
							ReadinessProbe:  probes.Readiness,
							SecurityContext: &corev1.SecurityContext{
								ReadOnlyRootFilesystem: &request.ReadOnlyRootFilesystem,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					DNSPolicy:     corev1.DNSClusterFirst,
				},
			},
		},
	}
	//hwh
	deploymentSpec.Spec.Template.Spec.NodeName = "kube-node-7"

	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	"root", "123456", "192.168.1.141", "tianxuan")
	//log.Println(dsn)
	//
	//db, err := gorm.Open("mysql", dsn)
	//if err != nil {
	//	log.Printf("连接到数据库失败%v", err)
	//} else {
	//	log.Printf("连接到数据库成功%s", "tianxuan")
	//	defer db.Close()
	//	var nodes []Node
	//	db.Find(&nodes)
	//	log.Println(nodes[0])
	//	log.Println("数据库操作结束")
	//}

	//hwh
	factory.ConfigureReadOnlyRootFilesystem(request, deploymentSpec)
	factory.ConfigureContainerUserID(deploymentSpec)

	if err := factory.ConfigureSecrets(request, deploymentSpec, existingSecrets); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}

//hwh
func MakeDeployHandler3(functionNamespace string, factory k8s.FunctionFactory) http.HandlerFunc {
	secrets := k8s.NewSecretsClient(factory.Client)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := types.FunctionDeployment{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to unmarshal request: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			wrappedErr := fmt.Errorf("validation failed: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		namespace := functionNamespace
		if len(request.Namespace) > 0 {
			namespace = request.Namespace
		}

		existingSecrets, err := secrets.GetSecrets(namespace, request.Secrets)
		if err != nil {
			wrappedErr := fmt.Errorf("unable to fetch secrets: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		nodemap := make(map[string]int32)
		nodemap["kube-node-7"] = 2
		nodemap["kube-node-5"] = 6
		for s, i := range nodemap {
			log.Println("进入nodemap，开始扩容")
			deploymentSpec, specErr := makeDeploymentSpec3(request, existingSecrets, factory, i, s)

			var profileList []k8s.Profile
			if request.Annotations != nil {
				profileNamespace := factory.Config.ProfilesNamespace
				profileList, err = factory.GetProfiles(ctx, profileNamespace, *request.Annotations)
				if err != nil {
					wrappedErr := fmt.Errorf("failed create Deployment spec: %s", err.Error())
					log.Println(wrappedErr)
					http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
					return
				}
			}
			for _, profile := range profileList {
				factory.ApplyProfile(profile, deploymentSpec)
			}

			if specErr != nil {
				wrappedErr := fmt.Errorf("failed create Deployment spec: %s", specErr.Error())
				log.Println(wrappedErr)
				http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
				return
			}

			deploy := factory.Client.AppsV1().Deployments(namespace)

			_, err = deploy.Create(context.TODO(), deploymentSpec, metav1.CreateOptions{})
			if err != nil {
				wrappedErr := fmt.Errorf("unable create Deployment: %s", err.Error())
				log.Println(wrappedErr)
				http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("Deployment created: %s.%s\n", request.Service, namespace)

			service := factory.Client.CoreV1().Services(namespace)
			serviceSpec := makeServiceSpec(request, factory)
			_, err = service.Create(context.TODO(), serviceSpec, metav1.CreateOptions{})

			if err != nil {
				wrappedErr := fmt.Errorf("failed create Service: %s", err.Error())
				log.Println(wrappedErr)
				http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
				return
			}

		}

		log.Printf("Service created: %s.%s\n", request.Service, namespace)

		w.WriteHeader(http.StatusAccepted)
	}
}

func makeDeploymentSpec3(request types.FunctionDeployment, existingSecrets map[string]*apiv1.Secret, factory k8s.FunctionFactory, replices int32, nodename string) (*appsv1.Deployment, error) {
	envVars := buildEnvVars(&request)

	initialReplicas := int32p(initialReplicasCount)
	labels := map[string]string{
		"faas_function": request.Service,
	}

	if request.Labels != nil {
		if min := getMinReplicaCount(*request.Labels); min != nil {
			initialReplicas = min
		}
		for k, v := range *request.Labels {
			labels[k] = v
		}
	}

	nodeSelector := createSelector(request.Constraints)

	resources, resourceErr := createResources(request)

	if resourceErr != nil {
		return nil, resourceErr
	}

	var imagePullPolicy apiv1.PullPolicy
	switch factory.Config.ImagePullPolicy {
	case "Never":
		imagePullPolicy = apiv1.PullNever
	case "IfNotPresent":
		imagePullPolicy = apiv1.PullIfNotPresent
	default:
		imagePullPolicy = apiv1.PullAlways
	}

	annotations := buildAnnotations(request)

	probes, err := factory.MakeProbes(request)
	if err != nil {
		return nil, err
	}

	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: annotations,
			Labels: map[string]string{
				"faas_function": request.Service,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"faas_function": request.Service,
				},
			},
			Replicas: initialReplicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
			RevisionHistoryLimit: int32p(10),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        request.Service,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: apiv1.PodSpec{
					NodeSelector: nodeSelector,
					Containers: []apiv1.Container{
						{
							Name:  request.Service,
							Image: request.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: factory.Config.RuntimeHTTPPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env:             envVars,
							Resources:       *resources,
							ImagePullPolicy: imagePullPolicy,
							LivenessProbe:   probes.Liveness,
							ReadinessProbe:  probes.Readiness,
							SecurityContext: &corev1.SecurityContext{
								ReadOnlyRootFilesystem: &request.ReadOnlyRootFilesystem,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					DNSPolicy:     corev1.DNSClusterFirst,
				},
			},
		},
	}

	//hwh 数据库操作，先不执行
	if false {
		deploymentSpec.Spec.Template.Spec.NodeName = "kube-node-7"

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			"root", "123456", "192.168.1.141", "tianxuan")
		log.Println(dsn)

		db, err := gorm.Open("mysql", dsn)
		if err != nil {
			log.Printf("连接到数据库失败%v", err)
		} else {
			log.Printf("连接到数据库成功%s", "tianxuan")
			defer db.Close()
			var nodes []Node
			db.Find(&nodes)
			log.Println(nodes[0])
			log.Println("数据库操作结束")
		}
	}

	deploymentSpec.Spec.Template.Spec.NodeName = nodename
	deploymentSpec.Spec.Replicas = int32p(replices)

	//hwh
	factory.ConfigureReadOnlyRootFilesystem(request, deploymentSpec)
	factory.ConfigureContainerUserID(deploymentSpec)

	if err := factory.ConfigureSecrets(request, deploymentSpec, existingSecrets); err != nil {
		return nil, err
	}

	return deploymentSpec, nil
}

//指定节点创建pod
func MakeDeployHandler5(functionNamespace string, factory k8s.FunctionFactory) http.HandlerFunc {
	//secrets := k8s.NewSecretsClient(factory.Client)

	return func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := types.FunctionDeployment{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to unmarshal request: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			wrappedErr := fmt.Errorf("validation failed: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		namespace := functionNamespace
		if len(request.Namespace) > 0 {
			namespace = request.Namespace
		}
		//
		//existingSecrets, err := secrets.GetSecrets(namespace, request.Secrets)
		//if err != nil {
		//	wrappedErr := fmt.Errorf("unable to fetch secrets: %s", err.Error())
		//	http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
		//	return
		//}
		//
		//deploymentSpec, specErr := makeDeploymentSpec(request, existingSecrets, factory)
		//
		//var profileList []k8s.Profile
		//if request.Annotations != nil {
		//	profileNamespace := factory.Config.ProfilesNamespace
		//	profileList, err = factory.GetProfiles(ctx, profileNamespace, *request.Annotations)
		//	if err != nil {
		//		wrappedErr := fmt.Errorf("failed create Deployment spec: %s", err.Error())
		//		log.Println(wrappedErr)
		//		http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
		//		return
		//	}
		//}
		//for _, profile := range profileList {
		//	factory.ApplyProfile(profile, deploymentSpec)
		//}
		//
		//if specErr != nil {
		//	wrappedErr := fmt.Errorf("failed create Deployment spec: %s", specErr.Error())
		//	log.Println(wrappedErr)
		//	http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
		//	return
		//}
		//
		//deploy := factory.Client.AppsV1().Deployments(namespace)
		//
		//_, err = deploy.Create(context.TODO(), deploymentSpec, metav1.CreateOptions{})
		//if err != nil {
		//	wrappedErr := fmt.Errorf("unable create Deployment: %s", err.Error())
		//	log.Println(wrappedErr)
		//	http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//log.Printf("Deployment created: %s.%s\n", request.Service, namespace)
		//hwh 数据库操作，先不执行

		//计算放置节点
		//var deploynode string
		deploynode := "kube-node-7"
		var nowfc Fc
		for _, fc := range Fcs {
			if fc.Name == "hello-python" {
				nowfc = fc
				break
			}
		}

		//Node稳定度排序
		sort.Slice(Nodes, func(i, j int) bool {
			if Std(Nodes[i].CpuFloat)< Std(Nodes[j].CpuFloat){
				return false
			}else {
				return true
			}
		})

		if nowfc.Delay <= 1 {
			for _, node := range Nodes {
				if node.CpuFloat[0] > nowfc.Cpu {
					deploynode = node.Name
					fc_lo := FcLocate{
						FcName: nowfc.Name,
						Locate: deploynode,
						Number: 1,
					}
					Db.Create(&fc_lo)
					node.CpuFloat[0] -= nowfc.Cpu
					break
				}
			}
		} else {
		OuterLoop:
			for _, node := range Nodes {
				for i := 0; i <= int(nowfc.Delay); i++ {
					if node.CpuFloat[i] < nowfc.Cpu {
						break
					}
					if i+1>int(nowfc.Delay){
						//找到
						deploynode = node.Name
						fc_lo := FcLocate{
							FcName: nowfc.Name,
							Locate: deploynode,
							Number: 1,
						}
						Db.Create(&fc_lo)

						for j := 0; j <= i; j++ {
							node.CpuFloat[j] -= nowfc.Cpu
						}

						break OuterLoop
					}
				}
			}
		}
		var fc Fc
		Db.Model(&fc).Where("name = ?", "hello-python").Update("replicas", 1,"nodename", deploynode)

		DeployPod("hello-python", deploynode, factory)
		//DeployPod("hello-python","kube-node-7",factory)

		service := factory.Client.CoreV1().Services(namespace)
		serviceSpec := makeServiceSpec(request, factory)
		_, err = service.Create(context.TODO(), serviceSpec, metav1.CreateOptions{})

		if err != nil {
			wrappedErr := fmt.Errorf("failed create Service: %s", err.Error())
			log.Println(wrappedErr)
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Service created: %s.%s\n", request.Service, namespace)

		w.WriteHeader(http.StatusAccepted)
	}
}

//hwh

func makeServiceSpec(request types.FunctionDeployment, factory k8s.FunctionFactory) *corev1.Service {

	serviceSpec := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        request.Service,
			Annotations: buildAnnotations(request),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"faas_function": request.Service,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     factory.Config.RuntimeHTTPPort,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: factory.Config.RuntimeHTTPPort,
					},
				},
			},
		},
	}

	return serviceSpec
}

func buildAnnotations(request types.FunctionDeployment) map[string]string {
	var annotations map[string]string
	if request.Annotations != nil {
		annotations = *request.Annotations
	} else {
		annotations = map[string]string{}
	}

	if _, ok := annotations["prometheus.io.scrape"]; !ok {
		annotations["prometheus.io.scrape"] = "false"
	}
	return annotations
}

func buildEnvVars(request *types.FunctionDeployment) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}

	if len(request.EnvProcess) > 0 {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k8s.EnvProcessName,
			Value: request.EnvProcess,
		})
	}

	for k, v := range request.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	sort.SliceStable(envVars, func(i, j int) bool {
		return strings.Compare(envVars[i].Name, envVars[j].Name) == -1
	})

	return envVars
}

func int32p(i int32) *int32 {
	return &i
}

func createSelector(constraints []string) map[string]string {
	selector := make(map[string]string)

	if len(constraints) > 0 {
		for _, constraint := range constraints {
			parts := strings.Split(constraint, "=")

			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}
	}

	return selector
}

func createResources(request types.FunctionDeployment) (*apiv1.ResourceRequirements, error) {
	resources := &apiv1.ResourceRequirements{
		Limits:   apiv1.ResourceList{},
		Requests: apiv1.ResourceList{},
	}

	// Set Memory limits
	if request.Limits != nil && len(request.Limits.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.Memory)
		if err != nil {
			return resources, err
		}
		resources.Limits[apiv1.ResourceMemory] = qty
	}

	if request.Requests != nil && len(request.Requests.Memory) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.Memory)
		if err != nil {
			return resources, err
		}
		resources.Requests[apiv1.ResourceMemory] = qty
	}

	// Set CPU limits
	if request.Limits != nil && len(request.Limits.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Limits.CPU)
		if err != nil {
			return resources, err
		}
		resources.Limits[apiv1.ResourceCPU] = qty
	}

	if request.Requests != nil && len(request.Requests.CPU) > 0 {
		qty, err := resource.ParseQuantity(request.Requests.CPU)
		if err != nil {
			return resources, err
		}
		resources.Requests[apiv1.ResourceCPU] = qty
	}

	return resources, nil
}

func getMinReplicaCount(labels map[string]string) *int32 {
	if value, exists := labels["com.openfaas.scale.min"]; exists {
		minReplicas, err := strconv.Atoi(value)
		if err == nil && minReplicas > 0 {
			return int32p(int32(minReplicas))
		}

		log.Println(err)
	}

	return nil
}
