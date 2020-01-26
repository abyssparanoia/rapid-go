package gluefirestore

// DocRef ... firestore document reference
type DocRef struct {
	CollectionName string `json:"collection_name"`
	DocID          string `json:"doc_id"`
}
