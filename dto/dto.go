package dto

type BackendErrorResponse struct {
	Description *string   `json:"description,omitempty"`
	ErrorCode   *int      `json:"errorCode,omitempty"`
	Meta        *MetaData `json:"meta,omitempty"`
}

type MetaData struct {
	Path      *string `json:"path,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}
