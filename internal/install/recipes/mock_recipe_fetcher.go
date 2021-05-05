package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type MockRecipeFetcher struct {
	FetchRecipeErr                error
	FetchRecipesErr               error
	FetchRecommendationsErr       error
	FetchRecipeCallCount          int
	FetchRecipesCallCount         int
	FetchRecommendationsCallCount int
	FetchRecipeVals               []types.OpenInstallationRecipe
	FetchRecipeVal                *types.OpenInstallationRecipe
	FetchRecipesVal               []types.OpenInstallationRecipe
	FetchRecommendationsVal       []types.OpenInstallationRecipe
	FetchRecipeNameCount          map[string]int
}

func NewMockRecipeFetcher() *MockRecipeFetcher {
	f := MockRecipeFetcher{}
	f.FetchRecipesVal = []types.OpenInstallationRecipe{}
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeNameCount = make(map[string]int)
	return &f
}

func (f *MockRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.OpenInstallationRecipe, error) {
	f.FetchRecipeCallCount++
	f.FetchRecipeNameCount[friendlyName]++

	if len(f.FetchRecipeVals) > 0 {
		i := utils.MinOf(f.FetchRecipeCallCount, len(f.FetchRecipeVals)) - 1
		return &f.FetchRecipeVals[i], f.FetchRecipesErr
	}

	return f.FetchRecipeVal, f.FetchRecipeErr
}

func (f *MockRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	f.FetchRecipesCallCount++
	return f.FetchRecipesVal, f.FetchRecipesErr
}

func (f *MockRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	f.FetchRecommendationsCallCount++
	return f.FetchRecommendationsVal, f.FetchRecommendationsErr
}
