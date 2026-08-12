package metadata

type MetadataCandidate struct {
	Source   string `json:"source"`
	SourceID string `json:"sourceId"`

	Title    string `json:"title"`
	AltTitle string `json:"altTitle"`
	Released string `json:"released"`

	Developers []string `json:"developers"`

	CoverURL string `json:"coverUrl"`
}
