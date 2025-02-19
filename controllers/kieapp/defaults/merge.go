package defaults

import (
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/google/go-cmp/cmp"
	"github.com/imdario/mergo"
	oappsv1 "github.com/openshift/api/apps/v1"
	buildv1 "github.com/openshift/api/build/v1"
	oimagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/pkg/errors"
	api "github.com/spolti/kie-cloud-operator-new/api/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func merge(baseline api.Environment, overwrite api.Environment) (api.Environment, error) {
	var env api.Environment
	var err error
	env.Console = mergeCustomObject(baseline.Console, overwrite.Console)
	env.SmartRouter = mergeCustomObject(baseline.SmartRouter, overwrite.SmartRouter)
	env.Dashbuilder = mergeCustomObject(baseline.Dashbuilder, overwrite.Dashbuilder)
	env.Others, err = mergeCustomObjects(baseline.Others, overwrite.Others)
	if err != nil {
		return api.Environment{}, err
	}
	env.Servers, err = mergeCustomObjects(baseline.Servers, overwrite.Servers)
	if err != nil {
		return api.Environment{}, err
	}
	return env, nil
}

func mergeCustomObjects(baseline, overwrite []api.CustomObject) ([]api.CustomObject, error) {
	if len(overwrite) == 0 {
		return baseline, nil
	}
	if len(baseline) == 0 {
		return overwrite, nil
	}
	if len(baseline) != len(overwrite) {
		return nil, errors.New("incompatible objects with different array lengths cannot be merged")
	}
	var result []api.CustomObject
	for index := range baseline {
		mergedObject := mergeCustomObject(baseline[index], overwrite[index])
		result = append(result, mergedObject)
	}
	return result, nil
}

func mergeCustomObject(baseline api.CustomObject, overwrite api.CustomObject) api.CustomObject {
	var object api.CustomObject
	if overwrite.Omit {
		object.Omit = overwrite.Omit
	}
	object.PersistentVolumeClaims = mergePersistentVolumeClaims(baseline.PersistentVolumeClaims, overwrite.PersistentVolumeClaims)
	object.ServiceAccounts = mergeServiceAccounts(baseline.ServiceAccounts, overwrite.ServiceAccounts)
	object.Secrets = mergeSecrets(baseline.Secrets, overwrite.Secrets)
	object.Roles = mergeRoles(baseline.Roles, overwrite.Roles)
	object.RoleBindings = mergeRoleBindings(baseline.RoleBindings, overwrite.RoleBindings)
	object.DeploymentConfigs = mergeDeploymentConfigs(baseline.DeploymentConfigs, overwrite.DeploymentConfigs)
	object.StatefulSets = mergeStatefulSets(baseline.StatefulSets, overwrite.StatefulSets)
	object.ImageStreams = mergeImageStreams(baseline.ImageStreams, overwrite.ImageStreams)
	object.BuildConfigs = mergeBuildConfigs(baseline.BuildConfigs, overwrite.BuildConfigs)
	object.Services = mergeServices(baseline.Services, overwrite.Services)
	object.Routes = mergeRoutes(baseline.Routes, overwrite.Routes)
	object.ConfigMaps = mergeConfigMaps(baseline.ConfigMaps, overwrite.ConfigMaps)
	return object
}

func mergePersistentVolumeClaims(baseline []corev1.PersistentVolumeClaim, overwrite []corev1.PersistentVolumeClaim) []corev1.PersistentVolumeClaim {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getPersistentVolumeClaimReferenceSlice(baseline)
		overwriteRefs := getPersistentVolumeClaimReferenceSlice(overwrite)
		slice := make([]corev1.PersistentVolumeClaim, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getPersistentVolumeClaimReferenceSlice(objects []corev1.PersistentVolumeClaim) []api.OpenShiftObject {
	slice := make([]client.Object, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeServiceAccounts(baseline []corev1.ServiceAccount, overwrite []corev1.ServiceAccount) []corev1.ServiceAccount {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getServiceAccountReferenceSlice(baseline)
		overwriteRefs := getServiceAccountReferenceSlice(overwrite)
		slice := make([]corev1.ServiceAccount, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getServiceAccountReferenceSlice(objects []corev1.ServiceAccount) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeSecrets(baseline []corev1.Secret, overwrite []corev1.Secret) []corev1.Secret {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getSecretReferenceSlice(baseline)
		overwriteRefs := getSecretReferenceSlice(overwrite)
		slice := make([]corev1.Secret, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getSecretReferenceSlice(objects []corev1.Secret) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeRoles(baseline []rbacv1.Role, overwrite []rbacv1.Role) []rbacv1.Role {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getRoleReferenceSlice(baseline)
		overwriteRefs := getRoleReferenceSlice(overwrite)
		slice := make([]rbacv1.Role, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func mergeRoleBindings(baseline []rbacv1.RoleBinding, overwrite []rbacv1.RoleBinding) []rbacv1.RoleBinding {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getRoleBindingReferenceSlice(baseline)
		overwriteRefs := getRoleBindingReferenceSlice(overwrite)
		slice := make([]rbacv1.RoleBinding, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getRoleReferenceSlice(objects []rbacv1.Role) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func getRoleBindingReferenceSlice(objects []rbacv1.RoleBinding) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeDeploymentConfigs(baseline []oappsv1.DeploymentConfig, overwrite []oappsv1.DeploymentConfig) []oappsv1.DeploymentConfig {
	if len(overwrite) == 0 {
		return baseline
	}
	if len(baseline) == 0 {
		return overwrite
	}
	baselineRefs := getDeploymentConfigReferenceSlice(baseline)
	overwriteRefs := getDeploymentConfigReferenceSlice(overwrite)
	for overwriteIndex := range overwrite {
		overwriteItem := &overwrite[overwriteIndex]
		baselineIndex, _ := findOpenShiftObject(overwriteItem, baselineRefs)
		if baselineIndex >= 0 {
			baselineItem := baseline[baselineIndex]
			err := mergo.Merge(&overwriteItem.ObjectMeta, baselineItem.ObjectMeta)
			if err != nil {
				log.Error("Error merging interfaces. ", err)
				return nil
			}
			mergedSpec, err := mergeDCSpec(baselineItem.Spec, overwriteItem.Spec)
			if err != nil {
				log.Error("Error merging DeploymentConfig Specs. ", err)
				return nil
			}
			overwriteItem.Spec = mergedSpec
		}
	}
	slice := make([]oappsv1.DeploymentConfig, combinedSize(baselineRefs, overwriteRefs))
	err := mergeObjects(baselineRefs, overwriteRefs, slice)
	if err != nil {
		log.Error("Error merging objects. ", err)
		return nil
	}
	return slice

}

func mergeStatefulSets(baseline []appsv1.StatefulSet, overwrite []appsv1.StatefulSet) []appsv1.StatefulSet {
	if len(overwrite) == 0 {
		return baseline
	}
	if len(baseline) == 0 {
		return overwrite
	}
	baselineRefs := getStatefulSetReferenceSlice(baseline)
	overwriteRefs := getStatefulSetReferenceSlice(overwrite)
	for overwriteIndex := range overwrite {
		overwriteItem := &overwrite[overwriteIndex]
		baselineIndex, _ := findOpenShiftObject(overwriteItem, baselineRefs)
		if baselineIndex >= 0 {
			baselineItem := baseline[baselineIndex]
			err := mergo.Merge(&overwriteItem.ObjectMeta, baselineItem.ObjectMeta)
			if err != nil {
				log.Error("Error merging interfaces. ", err)
				return nil
			}
			mergedSpec, err := mergeStatefulSpec(baselineItem.Spec, overwriteItem.Spec)
			if err != nil {
				log.Error("Error merging DeploymentConfig Specs. ", err)
				return nil
			}
			overwriteItem.Spec = mergedSpec
		}
	}
	slice := make([]appsv1.StatefulSet, combinedSize(baselineRefs, overwriteRefs))
	err := mergeObjects(baselineRefs, overwriteRefs, slice)
	if err != nil {
		log.Error("Error merging objects. ", err)
		return nil
	}
	return slice

}

func mergeImageStreams(baseline, overwrite []oimagev1.ImageStream) []oimagev1.ImageStream {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getImageStreamReferenceSlice(baseline)
		overwriteRefs := getImageStreamReferenceSlice(overwrite)
		slice := make([]oimagev1.ImageStream, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getImageStreamReferenceSlice(objects []oimagev1.ImageStream) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeBuildConfigs(baseline, overwrite []buildv1.BuildConfig) []buildv1.BuildConfig {
	if len(overwrite) == 0 {
		return baseline
	}
	if len(baseline) == 0 {
		return overwrite
	}
	baselineRefs := getBuildConfigReferenceSlice(baseline)
	overwriteRefs := getBuildConfigReferenceSlice(overwrite)
	for overwriteIndex := range overwrite {
		overwriteItem := &overwrite[overwriteIndex]
		baselineIndex, _ := findOpenShiftObject(overwriteItem, baselineRefs)
		if baselineIndex >= 0 {
			baselineItem := baseline[baselineIndex]
			err := mergo.Merge(&overwriteItem.ObjectMeta, baselineItem.ObjectMeta)
			if err != nil {
				log.Error("Error merging interfaces. ", err)
				return nil
			}
			mergedSpec, err := mergeBuildSpec(baselineItem.Spec, overwriteItem.Spec)
			if err != nil {
				log.Error("Error merging BuildConfig Specs. ", err)
				return nil
			}
			overwriteItem.Spec = mergedSpec
		}
	}
	slice := make([]buildv1.BuildConfig, combinedSize(baselineRefs, overwriteRefs))
	err := mergeObjects(baselineRefs, overwriteRefs, slice)
	if err != nil {
		log.Error("Error merging objects. ", err)
		return nil
	}
	return slice
}

func mergeDCSpec(baseline oappsv1.DeploymentConfigSpec, overwrite oappsv1.DeploymentConfigSpec) (oappsv1.DeploymentConfigSpec, error) {
	mergedTemplate, err := mergeTemplate(baseline.Template, overwrite.Template)
	if err != nil {
		return oappsv1.DeploymentConfigSpec{}, err
	}
	overwrite.Template = mergedTemplate

	mergedTriggers, err := mergeTriggers(baseline.Triggers, overwrite.Triggers)
	if err != nil {
		return oappsv1.DeploymentConfigSpec{}, err
	}
	overwrite.Triggers = mergedTriggers

	err = mergo.Merge(&baseline, overwrite, mergo.WithOverride)
	if err != nil {
		return oappsv1.DeploymentConfigSpec{}, nil
	}
	return baseline, nil
}

func mergeStatefulSpec(baseline appsv1.StatefulSetSpec, overwrite appsv1.StatefulSetSpec) (appsv1.StatefulSetSpec, error) {
	mergedTemplate, err := mergeTemplate(&baseline.Template, &overwrite.Template)
	if err != nil {
		return appsv1.StatefulSetSpec{}, err
	}
	overwrite.Template = *mergedTemplate

	err = mergo.Merge(&baseline, overwrite, mergo.WithOverride)
	if err != nil {
		return appsv1.StatefulSetSpec{}, nil
	}
	return baseline, nil
}

func mergeBuildSpec(baseline buildv1.BuildConfigSpec, overwrite buildv1.BuildConfigSpec) (buildv1.BuildConfigSpec, error) {
	mergedTriggers, err := mergeBuildTriggers(baseline.Triggers, overwrite.Triggers)
	if err != nil {
		return buildv1.BuildConfigSpec{}, err
	}
	overwrite.Triggers = mergedTriggers
	overwrite.Strategy.SourceStrategy.Env = shared.EnvOverride(baseline.Strategy.SourceStrategy.Env, overwrite.Strategy.SourceStrategy.Env)

	err = mergo.Merge(&baseline, overwrite, mergo.WithOverride)
	if err != nil {
		return buildv1.BuildConfigSpec{}, nil
	}
	return baseline, nil
}

func mergeTemplate(baseline *corev1.PodTemplateSpec, overwrite *corev1.PodTemplateSpec) (*corev1.PodTemplateSpec, error) {
	if overwrite == nil {
		return baseline, nil
	}
	err := mergo.Merge(&overwrite.ObjectMeta, baseline.ObjectMeta)
	if err != nil {
		log.Error("Error merging interfaces. ", err)
		return nil, nil
	}
	mergedPodSpec, err := mergePodSpecs(baseline.Spec, overwrite.Spec)
	if err != nil {
		return nil, err
	}
	overwrite.Spec = mergedPodSpec

	err = mergo.Merge(baseline, *overwrite, mergo.WithOverride)
	if err != nil {
		return nil, err
	}
	return baseline, nil
}

func mergeTriggers(baseline oappsv1.DeploymentTriggerPolicies, overwrite oappsv1.DeploymentTriggerPolicies) (oappsv1.DeploymentTriggerPolicies, error) {
	var mergedTriggers []oappsv1.DeploymentTriggerPolicy
	for baselineIndex, baselineItem := range baseline {
		idx, found := findDeploymentTriggerPolicy(baselineItem, overwrite)
		if idx == -1 {
			log.Debugf("Not found, adding %v to slice\n", baselineItem)
		} else {
			log.Debugf("Will merge %v on top of %v\n", found, baselineItem)
			if baselineItem.ImageChangeParams != nil {
				if found.ImageChangeParams == nil {
					found.ImageChangeParams = baselineItem.ImageChangeParams
				}
			}
			err := mergo.Merge(&baseline[baselineIndex], found, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
		mergedTriggers = append(mergedTriggers, baseline[baselineIndex])
	}
	for overwriteIndex, overwriteItem := range overwrite {
		idx, _ := findDeploymentTriggerPolicy(overwriteItem, mergedTriggers)
		if idx == -1 {
			log.Debugf("Not found, appending %v to slice\n", overwriteItem)
			mergedTriggers = append(mergedTriggers, overwrite[overwriteIndex])
		}
	}
	return mergedTriggers, nil
}

func mergeBuildTriggers(baseline []buildv1.BuildTriggerPolicy, overwrite []buildv1.BuildTriggerPolicy) ([]buildv1.BuildTriggerPolicy, error) {
	var mergedTriggers []buildv1.BuildTriggerPolicy
	for baselineIndex, baselineItem := range baseline {
		idx, found := findBuildTriggerPolicy(baselineItem, overwrite)
		if idx == -1 {
			log.Debugf("Not found, adding %v to slice\n", baselineItem)
		} else {
			log.Debugf("Will merge %v on top of %v\n", found, baselineItem)
			err := mergo.Merge(&baseline[baselineIndex], found, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
		mergedTriggers = append(mergedTriggers, baseline[baselineIndex])
	}
	for overwriteIndex, overwriteItem := range overwrite {
		idx, _ := findBuildTriggerPolicy(overwriteItem, mergedTriggers)
		if idx == -1 {
			log.Debugf("Not found, appending %v to slice\n", overwriteItem)
			mergedTriggers = append(mergedTriggers, overwrite[overwriteIndex])
		}
	}
	return mergedTriggers, nil
}

// findDeploymentTriggerPolicy Finds a deploymentTrigger by Type. In case type == ImageChange
// the match will be returned if both are not empty
func findDeploymentTriggerPolicy(object oappsv1.DeploymentTriggerPolicy, slice []oappsv1.DeploymentTriggerPolicy) (int, oappsv1.DeploymentTriggerPolicy) {
	emptyImageChangeParams := &oappsv1.DeploymentTriggerImageChangeParams{}
	for index, candidate := range slice {
		if candidate.Type == object.Type {
			if object.Type == oappsv1.DeploymentTriggerOnImageChange {
				if !cmp.Equal(object.ImageChangeParams, emptyImageChangeParams) && !cmp.Equal(candidate.ImageChangeParams, emptyImageChangeParams) {
					return index, candidate
				}
			} else {
				return index, candidate
			}
		}
	}
	return -1, oappsv1.DeploymentTriggerPolicy{}
}

// findBuildTriggerPolicy Finds a buildTrigger by Type
func findBuildTriggerPolicy(object buildv1.BuildTriggerPolicy, slice []buildv1.BuildTriggerPolicy) (int, buildv1.BuildTriggerPolicy) {
	for index, candidate := range slice {
		if candidate.Type == object.Type {
			return index, candidate
		}
	}
	return -1, buildv1.BuildTriggerPolicy{}
}

func mergePodSpecs(baseline corev1.PodSpec, overwrite corev1.PodSpec) (corev1.PodSpec, error) {
	mergedContainers, err := mergeContainers(baseline.Containers, overwrite.Containers)
	if err != nil {
		return corev1.PodSpec{}, err
	}
	overwrite.Containers = mergedContainers

	mergedVolumes, err := mergeVolumes(baseline.Volumes, overwrite.Volumes)
	if err != nil {
		return corev1.PodSpec{}, err
	}
	overwrite.Volumes = mergedVolumes

	err = mergo.Merge(&baseline, overwrite, mergo.WithOverride)
	if err != nil {
		return corev1.PodSpec{}, err
	}
	return baseline, nil
}

func mergeContainers(baseline []corev1.Container, overwrite []corev1.Container) ([]corev1.Container, error) {
	if len(overwrite) == 0 {
		return baseline, nil
	} else if len(baseline) == 0 {
		return overwrite, nil
	} else if len(baseline) > 1 || len(overwrite) > 1 {
		err := errors.New("Merge algorithm does not yet support multiple containers within a deployment")
		return nil, err
	}
	if baseline[0].Env == nil {
		baseline[0].Env = make([]corev1.EnvVar, 0)
	}
	overwrite[0].Env = shared.EnvOverride(baseline[0].Env, overwrite[0].Env)
	mergedPorts, err := mergePorts(baseline[0].Ports, overwrite[0].Ports)
	if err != nil {
		return nil, err
	}
	overwrite[0].Ports = mergedPorts

	mergedVolumeMounts, err := mergeVolumeMounts(baseline[0].VolumeMounts, overwrite[0].VolumeMounts)
	if err != nil {
		return []corev1.Container{}, err
	}
	overwrite[0].VolumeMounts = mergedVolumeMounts

	err = mergo.Merge(&baseline[0], overwrite[0], mergo.WithOverride)
	if err != nil {
		return nil, err
	}
	return baseline, nil
}

func mergePorts(baseline []corev1.ContainerPort, overwrite []corev1.ContainerPort) ([]corev1.ContainerPort, error) {
	var slice []corev1.ContainerPort
	for index := range baseline {
		found := findContainerPort(baseline[index], overwrite)
		if found != (corev1.ContainerPort{}) {
			err := mergo.Merge(&baseline[index], found, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
		slice = append(slice, baseline[index])
	}
	for index := range overwrite {
		found := findContainerPort(overwrite[index], baseline)
		if found == (corev1.ContainerPort{}) {
			slice = append(slice, overwrite[index])
		}
	}
	return slice, nil
}

func findContainerPort(port corev1.ContainerPort, ports []corev1.ContainerPort) corev1.ContainerPort {
	for index := range ports {
		if port.Name == ports[index].Name {
			return ports[index]
		}
	}
	return corev1.ContainerPort{}
}

func mergeVolumes(baseline []corev1.Volume, overwrite []corev1.Volume) ([]corev1.Volume, error) {
	var mergedVolumes []corev1.Volume
	for baselineIndex, baselineItem := range baseline {
		idx, found := findVolume(baselineItem, overwrite)
		if idx == -1 {
			log.Debugf("Not found, adding %v to slice\n", baselineItem)
			mergedVolumes = append(mergedVolumes, baseline[baselineIndex])
		} else {
			log.Debugf("Will replace %v on top of %v\n", found, baselineItem)
			mergedVolumes = append(mergedVolumes, found)
		}
	}
	for overwriteIndex, overwriteItem := range overwrite {
		idx, _ := findVolume(overwriteItem, mergedVolumes)
		if idx == -1 {
			log.Debugf("Not found, appending %v to slice\n", overwriteItem)
			mergedVolumes = append(mergedVolumes, overwrite[overwriteIndex])
		}
	}
	return mergedVolumes, nil
}

func mergeVolumeMounts(baseline []corev1.VolumeMount, overwrite []corev1.VolumeMount) ([]corev1.VolumeMount, error) {
	var mergedVolumeMounts []corev1.VolumeMount
	for baselineIndex, baselineItem := range baseline {
		idx, found := findVolumeMount(baselineItem, overwrite)
		if idx == -1 {
			log.Debugf("Not found, adding %v to slice\n", baselineItem)
		} else {
			log.Debugf("Will merge %v on top of %v\n", found, baselineItem)
			err := mergo.Merge(&baseline[baselineIndex], found, mergo.WithOverride)
			if err != nil {
				return nil, err
			}
		}
		mergedVolumeMounts = append(mergedVolumeMounts, baseline[baselineIndex])
	}
	for overwriteIndex, overwriteItem := range overwrite {
		idx, _ := findVolumeMount(overwriteItem, mergedVolumeMounts)
		if idx == -1 {
			log.Debugf("Not found, appending %v to slice\n", overwriteItem)
			mergedVolumeMounts = append(mergedVolumeMounts, overwrite[overwriteIndex])
		}
	}
	return mergedVolumeMounts, nil
}

func findVolume(object corev1.Volume, slice []corev1.Volume) (int, corev1.Volume) {
	for index, candidate := range slice {
		if candidate.Name == object.Name {
			return index, candidate
		}
	}
	return -1, corev1.Volume{}
}

func findVolumeMount(object corev1.VolumeMount, slice []corev1.VolumeMount) (int, corev1.VolumeMount) {
	for index, candidate := range slice {
		if candidate.Name == object.Name {
			return index, candidate
		}
	}
	return -1, corev1.VolumeMount{}
}

func getDeploymentConfigReferenceSlice(objects []oappsv1.DeploymentConfig) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func getStatefulSetReferenceSlice(objects []appsv1.StatefulSet) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func getBuildConfigReferenceSlice(objects []buildv1.BuildConfig) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeServices(baseline []corev1.Service, overwrite []corev1.Service) []corev1.Service {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getServiceReferenceSlice(baseline)
		overwriteRefs := getServiceReferenceSlice(overwrite)
		for overwriteIndex := range overwrite {
			overwriteItem := &overwrite[overwriteIndex]
			baselineIndex, _ := findOpenShiftObject(overwriteItem, baselineRefs)
			if baselineIndex >= 0 {
				baselineItem := baseline[baselineIndex]
				err := mergo.Merge(&overwriteItem.ObjectMeta, baselineItem.ObjectMeta)
				if err != nil {
					log.Error("Error merging interfaces. ", err)
					return nil
				}
				overwriteItem.Spec.Ports = mergeServicePorts(baselineItem.Spec.Ports, overwriteItem.Spec.Ports)
			}
		}
		slice := make([]corev1.Service, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getServiceReferenceSlice(objects []corev1.Service) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeServicePorts(baseline []corev1.ServicePort, overwrite []corev1.ServicePort) []corev1.ServicePort {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		var mergedServicePorts []corev1.ServicePort
		for baselineIndex, baselinePort := range baseline {
			found, servicePort := findServicePort(baselinePort, overwrite)
			if found {
				mergedServicePorts = append(mergedServicePorts, servicePort)
			} else {
				mergedServicePorts = append(mergedServicePorts, baseline[baselineIndex])
			}
		}
		for overwriteIndex, overwritePort := range overwrite {
			found, _ := findServicePort(overwritePort, baseline)
			if !found {
				mergedServicePorts = append(mergedServicePorts, overwrite[overwriteIndex])
			}
		}
		return mergedServicePorts
	}
}

func findServicePort(port corev1.ServicePort, ports []corev1.ServicePort) (bool, corev1.ServicePort) {
	for index, candidate := range ports {
		if port.Name == candidate.Name {
			return true, ports[index]
		}
	}
	return false, corev1.ServicePort{}
}

func mergeRoutes(baseline []routev1.Route, overwrite []routev1.Route) []routev1.Route {
	if len(overwrite) == 0 {
		return baseline
	} else if len(baseline) == 0 {
		return overwrite
	} else {
		baselineRefs := getRouteReferenceSlice(baseline)
		overwriteRefs := getRouteReferenceSlice(overwrite)
		slice := make([]routev1.Route, combinedSize(baselineRefs, overwriteRefs))
		err := mergeObjects(baselineRefs, overwriteRefs, slice)
		if err != nil {
			log.Error("Error merging objects. ", err)
			return nil
		}
		return slice
	}
}

func getRouteReferenceSlice(objects []routev1.Route) []api.OpenShiftObject {
	slice := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		slice[index] = &objects[index]
	}
	return slice
}

func mergeConfigMaps(baseline []corev1.ConfigMap, overwrite []corev1.ConfigMap) []corev1.ConfigMap {
	if len(overwrite) == 0 {
		return baseline
	}
	if len(baseline) == 0 {
		return overwrite
	}
	baselineRefs := getConfigMapReferenceSlice(baseline)
	overwriteRefs := getConfigMapReferenceSlice(overwrite)
	for overwriteIndex := range overwrite {
		overwriteItem := &overwrite[overwriteIndex]
		baselineIndex, _ := findOpenShiftObject(overwriteItem, baselineRefs)
		if baselineIndex >= 0 {
			baselineItem := baseline[baselineIndex]
			err := mergo.Merge(&overwriteItem.ObjectMeta, baselineItem.ObjectMeta)
			if err != nil {
				log.Error("Error merging interfaces. ", err)
				return nil
			}
			err = mergo.Merge(&overwriteItem.Data, baselineItem.Data)
			if err != nil {
				log.Error("Error merging ConfigMap Data. ", err)
				return nil
			}
			err = mergo.Merge(&overwriteItem.BinaryData, baselineItem.BinaryData)
			if err != nil {
				log.Error("Error merging ConfigMap BinaryData. ", err)
				return nil
			}
		}
	}
	mergedConfigMaps := make([]corev1.ConfigMap, combinedSize(baselineRefs, overwriteRefs))
	err := mergeObjects(baselineRefs, overwriteRefs, mergedConfigMaps)
	if err != nil {
		log.Error("Error merging objects. ", err)
		return nil
	}
	return mergedConfigMaps
}

func getConfigMapReferenceSlice(objects []corev1.ConfigMap) []api.OpenShiftObject {
	references := make([]api.OpenShiftObject, len(objects))
	for index := range objects {
		references[index] = &objects[index]
	}
	return references
}

func combinedSize(baseline []api.OpenShiftObject, overwrite []api.OpenShiftObject) int {
	count := 0
	for _, object := range overwrite {
		_, found := findOpenShiftObject(object, baseline)
		if found == nil && object.GetAnnotations()["delete"] != "true" {
			//unique item with no counterpart in baseline, count it
			count++
		} else if found != nil && object.GetAnnotations()["delete"] == "true" {
			///Deletes the counterpart in baseline, deduct 1 since the counterpart is being counted below
			count--
		}
	}
	count += len(baseline)
	return count
}

func mergeObjects(baseline []api.OpenShiftObject, overwrite []api.OpenShiftObject, objectSlice interface{}) error {
	slice := reflect.ValueOf(objectSlice)
	sliceIndex := 0
	for _, object := range baseline {
		_, found := findOpenShiftObject(object, overwrite)
		if found == nil {
			slice.Index(sliceIndex).Set(reflect.ValueOf(object).Elem())
			sliceIndex++
			log.Debugf("Not found, added %s to beginning of slice\n", object)
		} else if found.GetAnnotations()["delete"] != "true" {
			err := mergo.Merge(object, found, mergo.WithOverride)
			if err != nil {
				return err
			}
			slice.Index(sliceIndex).Set(reflect.ValueOf(object).Elem())
			sliceIndex++
			if found.GetAnnotations() == nil {
				annotations := make(map[string]string)
				found.SetAnnotations(annotations)
			}
		}
	}
	for _, object := range overwrite {
		if object.GetAnnotations()["delete"] != "true" {
			_, found := findOpenShiftObject(object, baseline)
			if found == nil {
				slice.Index(sliceIndex).Set(reflect.ValueOf(object).Elem())
				sliceIndex++
			}
		}
	}
	return nil
}

func findOpenShiftObject(object api.OpenShiftObject, slice []api.OpenShiftObject) (int, api.OpenShiftObject) {
	for index, candidate := range slice {
		if candidate.GetName() == object.GetName() {
			return index, candidate
		}
	}
	return -1, nil
}
