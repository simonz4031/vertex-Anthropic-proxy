package translation

var ModelNameMapping = map[string]string{
    "claude-3.5-sonnet-20240620": "claude-3-5-sonnet",
    "claude-3-5-sonnet-20240620": "claude-3-5-sonnet",
    "claude-3.5-sonnet":          "claude-3-5-sonnet",
    "claude-3-5-sonnet":          "claude-3-5-sonnet",
    // Add more mappings as needed
}

func NormalizeModelName(modelName string) string {
    if normalizedName, exists := ModelNameMapping[modelName]; exists {
        return normalizedName
    }
    return modelName // Return the original name if no mapping exists
}