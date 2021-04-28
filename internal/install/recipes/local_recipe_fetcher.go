package recipes

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type LocalRecipeFetcher struct {
	Path string
}

func (f *LocalRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.OpenInstallationRecipe, error) {
	recipes, err := f.FetchRecommendations(ctx, manifest)
	if err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		if recipe.Name == friendlyName {
			return &recipe, nil
		}
	}

	return nil, fmt.Errorf("%s: %w", friendlyName, ErrRecipeNotFound)
}

func (f *LocalRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	recipes, err := f.FetchRecipes(ctx, manifest)
	if err != nil {
		return nil, err
	}

	return manifest.ConstrainRecipes(recipes), nil
}

func (f *LocalRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	var recipes []types.OpenInstallationRecipe
	var err error

	if f.Path == "" {
		return nil, fmt.Errorf("unable to load recipes from empty path spec")
	}

	recipes, err = loadRecipesFromDir(ctx, f.Path)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func loadRecipesFromDir(ctx context.Context, path string) ([]types.OpenInstallationRecipe, error) {
	recipePaths := []string{}

	log.WithFields(log.Fields{
		"path": path,
	}).Debug("loading recipes")

	err := filepath.Walk(
		path,
		func(path string, info os.FileInfo, err error) error {
			ext := filepath.Ext(path)

			if ext == ".yml" || ext == ".yaml" {
				recipePaths = append(recipePaths, path)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	recipes := []types.OpenInstallationRecipe{}

	for _, path := range recipePaths {
		var r types.OpenInstallationRecipe

		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Error(err)
			continue
		}

		err = yaml.Unmarshal(content, &r)
		if err != nil {
			log.Error(err)
			continue
		}

		recipes = append(recipes, r)
	}

	return recipes, nil
}
