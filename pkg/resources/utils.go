//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package resources

import (
	"os"
	"strconv"
	"strings"

	operatorv1alpha1 "github.com/ibm/ibm-metering-sender-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const CommonServicesProductName = "IBM Cloud Platform Common Services"
const CommonServicesProductID = "068a62892a1e4db39641342e592daa25"
const CommonServicesProductVersion = "3.4.0"
const MeteringComponentName = "meteringsvc"
const MeteringReleaseName = "metering"
const SenderDeploymentName = "metering-sender"
const MeteringDependencies = "ibm-common-services.auth-idp, mongodb, cert-manager"

var DefaultStatusForCR = []string{"none"}
var DefaultMode int32 = 420

var log = logf.Log.WithName("resource_utils")

func BuildMongoDBEnvVars(mongoDB operatorv1alpha1.MeteringSpecMongoDB) []corev1.EnvVar {
	mongoDBEnvVars := []corev1.EnvVar{
		{
			Name:  "HC_MONGO_HOST",
			Value: mongoDB.Host,
		},
		{
			Name:  "HC_MONGO_PORT",
			Value: strconv.Itoa(mongoDB.Port),
		},
		{
			Name: "HC_MONGO_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mongoDB.UsernameSecret,
					},
					Key:      mongoDB.UsernameKey,
					Optional: &TrueVar,
				},
			},
		},
		{
			Name: "HC_MONGO_PASS",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: mongoDB.PasswordSecret,
					},
					Key:      mongoDB.PasswordKey,
					Optional: &TrueVar,
				},
			},
		},
		{
			Name:  "HC_MONGO_ISSSL",
			Value: "true",
		},
		{
			Name:  "HC_MONGO_SSL_CA",
			Value: "/certs/mongodb-ca/tls.crt",
		},
		{
			Name:  "HC_MONGO_SSL_CERT",
			Value: "/certs/mongodb-client/tls.crt",
		},
		{
			Name:  "HC_MONGO_SSL_KEY",
			Value: "/certs/mongodb-client/tls.key",
		},
	}
	return mongoDBEnvVars
}

func BuildCommonClusterEnvVars(instanceNamespace, instanceIAMnamespace string) []corev1.EnvVar {
	reqLogger := log.WithValues("func", "BuildCommonClusterEnvVars")

	var iamNamespace string
	if instanceIAMnamespace != "" {
		reqLogger.Info("IAMnamespace=" + instanceIAMnamespace)
		iamNamespace = instanceIAMnamespace
	} else {
		reqLogger.Info("IAMnamespace is blank, use instance=" + instanceNamespace)
		iamNamespace = instanceNamespace
	}

	clusterEnvVars := []corev1.EnvVar{
		{
			Name:  "IAM_NAMESPACE",
			Value: iamNamespace,
		},
	}
	return clusterEnvVars
}

func BuildSenderClusterEnvVars(instanceNamespace, instanceClusterNamespace,
	instanceClusterName, hubKubeConfigSecret string) []corev1.EnvVar {

	reqLogger := log.WithValues("func", "BuildSenderClusterEnvVars")

	var clusterName string
	if instanceClusterName != "" {
		clusterName = instanceClusterName
	} else {
		clusterName = DefaultClusterName
	}

	var clusterNamespace string
	if instanceClusterNamespace != "" {
		reqLogger.Info("clusterNamespace=" + instanceClusterNamespace)
		clusterNamespace = instanceClusterNamespace
	} else {
		reqLogger.Info("clusterNamespace is blank, use instance=" + instanceNamespace)
		clusterNamespace = instanceNamespace
	}

	clusterEnvVars := []corev1.EnvVar{
		{
			Name:  "HC_CLUSTER_NAME",
			Value: clusterName,
		},
		{
			Name:  "HC_CLUSTER_NAMESPACE",
			Value: clusterNamespace,
		},
		{
			Name:  "HC_HUB_CONFIG",
			Value: hubKubeConfigSecret,
		},
	}

	return clusterEnvVars
}

// returns the labels associated with the resource being created
func LabelsForMetadata(deploymentName string) map[string]string {
	return map[string]string{"app.kubernetes.io/name": deploymentName, "app.kubernetes.io/component": MeteringComponentName,
		"app.kubernetes.io/managed-by": "operator", "app.kubernetes.io/instance": MeteringReleaseName, "release": MeteringReleaseName}
}

// returns the labels for selecting the resources belonging to the given metering CR name
func LabelsForSelector(deploymentName string, crType string, crName string) map[string]string {
	return map[string]string{"app": deploymentName, "component": MeteringComponentName, crType: crName}
}

// returns the labels associated with the Pod being created
func LabelsForPodMetadata(deploymentName string, crType string, crName string) map[string]string {
	podLabels := LabelsForMetadata(deploymentName)
	selectorLabels := LabelsForSelector(deploymentName, crType, crName)
	for key, value := range selectorLabels {
		podLabels[key] = value
	}
	return podLabels
}

// GetPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []corev1.Pod) []string {
	reqLogger := log.WithValues("func", "GetPodNames")
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
		reqLogger.Info("pod name=" + pod.Name)
	}
	return podNames
}

// GetImageID returns the ID of an operand image, either <imageName>@<SHA> or <imageName>:<tag>
func GetImageID(instanceImageRegistry, instanceImageTagPostfix, defaultImageRegistry,
	imageName, envVarName, defaultImageTag string) string {
	reqLogger := log.WithValues("func", "GetImageID")

	// determine if the image registry has been overridden by the CR
	var imageRegistry, imageID string
	if instanceImageRegistry == "" {
		imageRegistry = defaultImageRegistry
		reqLogger.Info("use default imageRegistry=" + imageRegistry)
	} else {
		imageRegistry = instanceImageRegistry
		reqLogger.Info("use instance imageRegistry=" + imageRegistry)
	}

	// determine if an image SHA or tag has been set in an env var.
	// if not, use the default tag (mainly used during development).
	imageTagOrSHA := os.Getenv(envVarName)
	if len(imageTagOrSHA) > 0 {
		// use the value from the env var to build the image ID.
		// a SHA value looks like "sha256:nnnn".
		// a tag value looks like "3.5.0".
		if strings.HasPrefix(imageTagOrSHA, "sha256:") {
			// use the SHA value
			imageID = imageRegistry + "/" + imageName + "@" + imageTagOrSHA
		} else {
			// use the tag value
			imageID = imageRegistry + "/" + imageName + ":" + imageTagOrSHA + instanceImageTagPostfix
		}
	} else {
		// use the default tag to build the image ID
		imageID = imageRegistry + "/" + imageName + ":" + defaultImageTag + instanceImageTagPostfix
	}

	return imageID
}

// checkerCommand is the command to be executed by the secret-check container.
// mongoDB contains the password names from the CR.
// additionalInfo contains info about additional secrets to check.
func BuildSecretCheckContainer(deploymentName, imageName, checkerCommand string,
	mongoDB operatorv1alpha1.MeteringSpecMongoDB, additionalInfo *SecretCheckData) corev1.Container {

	containerName := deploymentName + "-secret-check"
	nameList := mongoDB.UsernameSecret + " " + mongoDB.PasswordSecret + " " +
		mongoDB.ClusterCertsSecret + " " + mongoDB.ClientCertsSecret
	usernameSecretDir := "muser-" + mongoDB.UsernameSecret
	passwordSecretDir := "mpass-" + mongoDB.PasswordSecret
	dirList := usernameSecretDir + " " + passwordSecretDir + " " +
		mongoDB.ClusterCertsSecret + " " + mongoDB.ClientCertsSecret
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "mongodb-ca-cert",
			MountPath: "/sec/" + mongoDB.ClusterCertsSecret,
		},
		{
			Name:      "mongodb-client-cert",
			MountPath: "/sec/" + mongoDB.ClientCertsSecret,
		},
		{
			Name:      usernameSecretDir,
			MountPath: "/sec/" + usernameSecretDir,
		},
		{
			Name:      passwordSecretDir,
			MountPath: "/sec/" + passwordSecretDir,
		},
	}
	if additionalInfo != nil {
		nameList += " "
		nameList += additionalInfo.Names
		dirList += " "
		dirList += additionalInfo.Dirs
		volumeMounts = append(volumeMounts, additionalInfo.VolumeMounts...)
	}

	var secretCheckContainer = corev1.Container{
		Image:           imageName,
		Name:            containerName,
		ImagePullPolicy: corev1.PullAlways,
		Command: []string{
			"sh",
			"-c",
			checkerCommand,
		},
		Env: []corev1.EnvVar{
			{
				Name:  "SECRET_LIST",
				Value: nameList,
			},
			{
				Name:  "SECRET_DIR_LIST",
				Value: dirList,
			},
		},
		VolumeMounts:    volumeMounts,
		Resources:       commonInitResources,
		SecurityContext: &commonSecurityContext,
	}
	return secretCheckContainer
}

func BuildInitContainer(deploymentName, imageName string, envVars []corev1.EnvVar) corev1.Container {
	containerName := deploymentName + "-init"
	var initContainer = corev1.Container{
		Image:           imageName,
		Name:            containerName,
		ImagePullPolicy: corev1.PullAlways,
		Command: []string{
			"node",
			"/datamanager/lib/metering_init.js",
			"verifyOnlyMongo",
		},
		// CommonEnvVars and mongoDBEnvVars will be added by the controller
		Env:             envVars,
		VolumeMounts:    commonInitVolumeMounts,
		Resources:       commonInitResources,
		SecurityContext: &commonSecurityContext,
	}
	return initContainer
}

//AnnotationsForPod returns the annotations associated with the pod being created
func AnnotationsForPod() map[string]string {
	return map[string]string{"productName": CommonServicesProductName, "productID": CommonServicesProductID,
		"productVersion": CommonServicesProductVersion, "productMetric": "FREE", "clusterhealth.ibm.com/dependencies": MeteringDependencies}
}

// GetServiceAccountName returns the service account name or default if it is not set in the environment
func GetServiceAccountName() string {

	sa := "default"

	envSa := os.Getenv("SA_NAME")
	if len(envSa) > 0 {
		sa = envSa
	}
	return sa
}

// set loglevelType to "log4js" when building volumes for metering-mcmui.
// set loglevelType to "loglevel" when building volumes for any other component.
func BuildCommonVolumes(mongoDB operatorv1alpha1.MeteringSpecMongoDB, loglevelPrefix, loglevelType string) []corev1.Volume {
	// example for metering-ui
	//   Name: loglevel
	//     Key: metering-ui-loglevel.json
	//     Path: loglevel.json
	// example for metering-mcmui
	//   Name: log4js
	//     Key: metering-mcmui-log4js.json
	//     Path: log4js.json
	loglevelKey := loglevelPrefix + "-" + loglevelType + ".json"
	loglevelPath := loglevelType + ".json"

	// CS??? removed icp-serviceid-apikey-secret
	commonVolumes := []corev1.Volume{
		{
			Name: "mongodb-ca-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  mongoDB.ClusterCertsSecret,
					DefaultMode: &DefaultMode,
					Optional:    &TrueVar,
				},
			},
		},
		{
			Name: "mongodb-client-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  mongoDB.ClientCertsSecret,
					DefaultMode: &DefaultMode,
					Optional:    &TrueVar,
				},
			},
		},
		{
			Name: "muser-icp-mongodb-admin",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  mongoDB.UsernameSecret,
					DefaultMode: &DefaultMode,
					Optional:    &TrueVar,
				},
			},
		},
		{
			Name: "mpass-icp-mongodb-admin",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  mongoDB.PasswordSecret,
					DefaultMode: &DefaultMode,
					Optional:    &TrueVar,
				},
			},
		},
		{
			Name: loglevelType,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "metering-logging-configuration",
					},
					Items: []corev1.KeyToPath{
						{
							Key:  loglevelKey,
							Path: loglevelPath,
						},
					},
					DefaultMode: &DefaultMode,
					Optional:    &TrueVar,
				},
			},
		},
	}
	return commonVolumes
}
