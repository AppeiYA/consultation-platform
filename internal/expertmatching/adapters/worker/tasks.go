package worker

const TypeStartMatching = "expertmatching:start"

type StartMatchingPayload struct {
	RunID string `json:"run_id"`
}
