package geminiranker

import "google.golang.org/genai"

func rankingConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"candidates": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"consultant_id": map[string]any{
								"type": "string",
							},
							"score": map[string]any{
								"type": "number",
							},
							"reasons": map[string]any{
								"type": "array",
								"items": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"factor": map[string]any{
											"type": "string",
										},
										"detail": map[string]any{
											"type": "string",
										},
									},
									"required": []string{
										"factor",
										"detail",
									},
								},
							},
						},
						"required": []string{
							"consultant_id",
							"score",
							"reasons",
						},
					},
				},
			},
			"required": []string{
				"candidates",
			},
		},
	}
}