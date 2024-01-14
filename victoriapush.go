/*
victoriapush pushes metrics to VictoriaMetrics using the prometheus exposition format

send metrics with custom labels in a map

assign global labels
*/
package victoriapush

import "context"

// add URL and globalLabels to context and return
func InitVictoria(URL string, globalLabels map[string]string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "URL", URL)
	ctx = context.WithValue(ctx, "globalLabels", globalLabels)

	return ctx
}

// Replace globalLabels map with a new one
func ReplaceGlobalLabels(ctx context.Context, globalLabels map[string]string) context.Context {
	ctx = context.WithValue(ctx, "globalLabels", globalLabels)
	return ctx
}

// Add a map of labels to globalLabels
func AddGlobalLabels(ctx context.Context, globalLabels map[string]string) context.Context {

	//get globalLabels from context, iterate over new map and add each to globalLabels from ctx
	currGlobalLabels := ctx.Value("globalLabels").(map[string]string)
	for label, v := range globalLabels {
		currGlobalLabels[label] = v
	}

	//add updated globalLabels back to ctx
	ctx = context.WithValue(ctx, "globalLabels", currGlobalLabels)
	return ctx
}

// Remove the list of labels from the globalLabels map IF it matches
func RemGlobalLabels(ctx context.Context, labelsForRem []string) context.Context {

	//delete all labels in labelsForRem from globalLabels
	currGlobalLabels := ctx.Value("globalLabels").(map[string]string)
	for _, label := range labelsForRem {
		delete(currGlobalLabels, label)
	}

	//add updated globalLabels back to ctx
	ctx = context.WithValue(ctx, "globalLabels", currGlobalLabels)
	return ctx
}
