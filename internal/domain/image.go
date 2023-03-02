package domain

// Image type should contain:
// ID which is auto-generated in the database as an int64 type
// Path to the saved image
// ContentType for proper image retrieval
type Image struct {
	ID          int64
	Path        string
	ContentType string
}
