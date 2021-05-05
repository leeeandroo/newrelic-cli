package types

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	InfraAgentRecipeName = "infrastructure-agent-installer"
	LoggingRecipeName    = "logs-integration"
)

var (
	RecipeVariables = map[string]string{}
)

func (r *OpenInstallationRecipe) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var recipe map[string]interface{}
	err := unmarshal(&recipe)
	if err != nil {
		return err
	}

	r.Dependencies = interfaceSliceToStringSlice(recipe["dependencies"].([]interface{}))
	r.Description = toStringByFieldName("description", recipe)
	r.DisplayName = toStringByFieldName("displayName", recipe)
	r.File = toStringByFieldName("file", recipe)
	r.ID = toStringByFieldName("id", recipe)
	r.InputVars = expandInputVars(recipe)

	installAsString, err := expandInstalllMapToString(recipe)
	if err != nil {
		return err
	}
	r.Install = installAsString

	r.InstallTargets = expandInstallTargets(recipe)
	r.Keywords = interfaceSliceToStringSlice(recipe["keywords"].([]interface{}))
	r.LogMatch = expandLogMatch(recipe)
	r.Name = toStringByFieldName("name", recipe)
	r.PostInstall = expandPostInstall(recipe)
	r.PreInstall = expandPreInstall(recipe)
	r.ProcessMatch = interfaceSliceToStringSlice(recipe["processMatch"].([]interface{}))

	// Quickstarts are not quite ready in the API yet.
	// The Nerdgraph type are incorrect and will get getting
	// updated when the Quickstarts feature is worked on.
	// r.Quickstarts = expandQuickStarts(recipe)

	r.Repository = toStringByFieldName("repository", recipe)

	if v, ok := recipe["stability"]; ok {
		r.Stability = OpenInstallationStability(v.(string))
	}

	r.SuccessLinkConfig = expandSuccessLinkConfig(recipe)

	if v, ok := recipe["validationNrql"]; ok {
		r.ValidationNRQL = NRQL(v.(string))
	}

	return err
}

func expandSuccessLinkConfig(recipe map[string]interface{}) OpenInstallationSuccessLinkConfig {
	v, ok := recipe["successLinkConfig"]
	if !ok {
		return OpenInstallationSuccessLinkConfig{}
	}

	dataIn := v.(map[interface{}]interface{})
	reData := map[string]interface{}{}
	for k, v := range dataIn {
		reData[k.(string)] = v
	}

	dataOut := OpenInstallationSuccessLinkConfig{
		Filter: toStringByFieldName("filter", reData),
	}

	if v, ok := reData["type"]; ok {
		dataOut.Type = OpenInstallationSuccessLinkType(v.(string))
	}

	return dataOut
}

func expandInstallTargets(recipe map[string]interface{}) []OpenInstallationRecipeInstallTarget {
	v, ok := recipe["installTargets"]
	if !ok {
		return []OpenInstallationRecipeInstallTarget{}
	}

	dataIn := v.([]interface{})
	dataOut := make([]OpenInstallationRecipeInstallTarget, len(dataIn))
	dataz := make([]map[string]interface{}, len(dataIn))
	for i, vv := range dataIn {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			dataz[i] = varr
		}
	}

	for i, v := range dataz {
		vOut := OpenInstallationRecipeInstallTarget{
			KernelArch:      toStringByFieldName("kernelArch", v),
			KernelVersion:   toStringByFieldName("kernelVersion", v),
			PlatformVersion: toStringByFieldName("platformVersion", v),
		}

		if v, ok := v["os"]; ok {
			vOut.Os = OpenInstallationOperatingSystem(v.(string))
		}

		if v, ok := v["platform"]; ok {
			vOut.Platform = OpenInstallationPlatform(v.(string))
		}

		if v, ok := v["platformFamily"]; ok {
			vOut.PlatformFamily = OpenInstallationPlatformFamily(v.(string))
		}

		if v, ok := v["type"]; ok {
			vOut.Type = OpenInstallationTargetType(v.(string))
		}

		dataOut[i] = vOut
	}

	return dataOut
}

func expandPreInstall(recipe map[string]interface{}) OpenInstallationPreInstallConfiguration {
	v, ok := recipe["preInstall"]
	if !ok {
		return OpenInstallationPreInstallConfiguration{}
	}

	vv := v.(map[interface{}]interface{})
	infoOut := map[string]interface{}{}
	for k, v := range vv {
		infoOut[k.(string)] = v
	}

	return OpenInstallationPreInstallConfiguration{
		Info:   toStringByFieldName("info", infoOut),
		Prompt: toStringByFieldName("prompt", infoOut),
	}
}

func expandPostInstall(recipe map[string]interface{}) OpenInstallationPostInstallConfiguration {
	v, ok := recipe["postInstall"]
	if !ok {
		return OpenInstallationPostInstallConfiguration{}
	}

	vv := v.(map[interface{}]interface{})
	infoOut := map[string]interface{}{}
	for k, v := range vv {
		infoOut[k.(string)] = v
	}

	return OpenInstallationPostInstallConfiguration{
		Info: toStringByFieldName("info", infoOut),
	}
}

func expandInputVars(recipe map[string]interface{}) []OpenInstallationRecipeInputVariable {
	v, ok := recipe["inputVars"]
	if !ok {
		return []OpenInstallationRecipeInputVariable{}
	}

	vars := v.([]interface{})
	varsOut := make([]OpenInstallationRecipeInputVariable, len(vars))

	varz := make([]map[string]interface{}, len(vars))
	for i, vv := range vars {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			varz[i] = varr
		}
	}

	for i, v := range varz {
		vOut := OpenInstallationRecipeInputVariable{
			Default: toStringByFieldName("default", v),
			Name:    toStringByFieldName("name", v),
			Prompt:  toStringByFieldName("prompt", v),
			Secret:  toBoolByFieldName("secret", v),
		}

		varsOut[i] = vOut
	}

	return varsOut
}

func expandLogMatch(recipe map[string]interface{}) []OpenInstallationLogMatch {
	v, ok := recipe["logMatch"]
	if !ok {
		return []OpenInstallationLogMatch{}
	}

	dataIn := v.([]interface{})
	dataOut := make([]OpenInstallationLogMatch, len(dataIn))
	dataz := make([]map[string]interface{}, len(dataIn))
	for i, vv := range dataIn {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			dataz[i] = varr
		}
	}

	for i, v := range dataz {
		vOut := OpenInstallationLogMatch{
			Attributes: expandLogAttributes(v),
			File:       toStringByFieldName("file", v),
			Name:       toStringByFieldName("name", v),
			Pattern:    toStringByFieldName("pattern", v),
			Systemd:    toStringByFieldName("systemd", v),
		}

		dataOut[i] = vOut
	}

	return dataOut
}

func expandLogAttributes(data map[string]interface{}) OpenInstallationAttributes {
	attributesOut := OpenInstallationAttributes{}

	v, ok := data["attributes"]
	if !ok {
		return attributesOut
	}

	attrs := v.(map[interface{}]interface{})
	attrsOut := map[string]string{}
	for k, v := range attrs {
		attrsOut[k.(string)] = v.(string)
	}

	if v, ok := attrsOut["logtype"]; ok {
		attributesOut.Logtype = v
	}

	return attributesOut
}

func toBoolByFieldName(fieldName string, data map[string]interface{}) bool {
	if v, ok := data[fieldName]; ok {
		return v.(bool)
	}

	return false
}

func toStringByFieldName(fieldName string, data map[string]interface{}) string {
	if v, ok := data[fieldName]; ok {
		return v.(string)
	}

	return ""
}

func expandInstalllMapToString(recipeIn map[string]interface{}) (string, error) {
	installIn, ok := recipeIn["install"]
	if !ok {
		return "", fmt.Errorf("error unmarshaling installation recipe: field 'install' is empty or undefined")
	}

	installOut := map[string]interface{}{}
	installMap := installIn.(map[interface{}]interface{})
	for k, v := range installMap {
		installOut[k.(string)] = v
	}

	installAsString, err := yaml.Marshal(installOut)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling recipe.install to string: %s", err)
	}

	return string(installAsString), nil
}

func interfaceSliceToStringSlice(slice []interface{}) []string {
	out := make([]string, len(slice))

	for i, v := range slice {
		out[i] = v.(string)
	}

	return out
}

func (r *OpenInstallationRecipe) PostInstallMessage() string {
	if r.PostInstall.Info != "" {
		return r.PostInstall.Info
	}

	return ""
}

func (r *OpenInstallationRecipe) PreInstallMessage() string {
	if r.PreInstall.Info != "" {
		return r.PreInstall.Info
	}

	return ""
}

type RecipeVars map[string]string

// AddVar is responsible for including a new variable on the recipe Vars
// struct, which is used by go-task executor.
func (r *OpenInstallationRecipe) AddVar(key string, value interface{}) {
	// if len(r.Vars) == 0 {
	// 	r.Vars = make(map[string]interface{})
	// }

	// r.Vars[key] = value
}

func (r *OpenInstallationRecipe) SetRecipeVar(key string, value string) {
	RecipeVariables[key] = value
}

func (r *OpenInstallationRecipe) IsApm() bool {
	return r.HasKeyword("apm")
}

func (r *OpenInstallationRecipe) HasHostTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.HOST)
}

func (r *OpenInstallationRecipe) HasApplicationTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.APPLICATION)
}

func (r *OpenInstallationRecipe) HasKeyword(keyword string) bool {
	if len(r.Keywords) == 0 {
		return false
	}

	for _, single := range r.Keywords {
		if strings.EqualFold(single, keyword) {
			return true
		}
	}

	return false
}

func (r *OpenInstallationRecipe) HasTargetType(t OpenInstallationTargetType) bool {
	if len(r.InstallTargets) == 0 {
		return false
	}

	for _, target := range r.InstallTargets {
		if target.Type == t {
			return true
		}
	}

	return false
}
