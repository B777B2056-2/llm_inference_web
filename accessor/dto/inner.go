package dto

type BatchInferenceKafkaPayload struct {
	TraceId            string          `json:"trace_id"`
	BatchInferenceId   string          `json:"batch_inference_id"`
	BatchInferenceName string          `json:"batch_inference_name"`
	UserId             int             `json:"user_id"`
	SamplingParams     InferenceParams `json:"sampling_params"`
	Prompts            []string        `json:"prompts"`
}
