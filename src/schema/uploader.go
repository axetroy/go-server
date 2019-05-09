package schema

type FileResponse struct {
	Hash     string `json:"hash"`
	Filename string `json:"filename"`
	Origin   string `json:"origin"`
	Size     int64  `json:"size"`
}
