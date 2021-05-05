package recipes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"gopkg.in/yaml.v2"
)

type RecipeFileFetcherImpl struct {
	HTTPGetFunc  func(string) (*http.Response, error)
	readFileFunc func(string) ([]byte, error)
}

func NewRecipeFileFetcher() RecipeFileFetcher {
	f := RecipeFileFetcherImpl{}
	f.HTTPGetFunc = defaultHTTPGetFunc
	f.readFileFunc = defaultReadFileFunc
	return &f
}

func defaultHTTPGetFunc(recipeURL string) (*http.Response, error) {
	return http.Get(recipeURL)
}

func defaultReadFileFunc(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (f *RecipeFileFetcherImpl) FetchRecipeFile(recipeURL *url.URL) (*types.OpenInstallationRecipe, error) {
	response, err := f.HTTPGetFunc(recipeURL.String())
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("received non-2xx status code %d when retrieving recipe", response.StatusCode)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return NewRecipeFile(string(body))
}

func (f *RecipeFileFetcherImpl) LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error) {
	out, err := f.readFileFunc(filename)
	if err != nil {
		return nil, err
	}

	return NewRecipeFile(string(out))
}

func NewRecipeFile(recipeFileString string) (*types.OpenInstallationRecipe, error) {
	var f types.OpenInstallationRecipe
	err := yaml.Unmarshal([]byte(recipeFileString), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
