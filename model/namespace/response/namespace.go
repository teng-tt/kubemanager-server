package response

type Namespace struct {
	Name              string `json:"name"`
	CreationTimestamp int64  `json:"creationTimestamp"`
	Status            string `json:"status"`
}
