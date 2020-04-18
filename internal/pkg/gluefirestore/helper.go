package gluefirestore

import (
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/abyssparanoia/rapid-go/internal/pkg/util"
)

// GenerateDocumentRef ... generate document
func GenerateDocumentRef(fCli *firestore.Client, docRefs []*DocRef) *firestore.DocumentRef {
	var dst *firestore.DocumentRef
	for i, docRef := range docRefs {
		if i == 0 {
			dst = fCli.Collection(docRef.CollectionName).Doc(docRef.DocID)
		} else {
			dst = dst.Collection(docRef.CollectionName).Doc(docRef.DocID)
		}
	}
	return dst
}

// Get ... get a document
func Get(ctx context.Context, docRef *firestore.DocumentRef, dst interface{}) (bool, error) {
	dsnp, err := docRef.Get(ctx)
	if dsnp != nil && !dsnp.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	return true, nil
}

// GetMulti ... get multi documents
func GetMulti(ctx context.Context, fCli *firestore.Client, docRefs []*firestore.DocumentRef, dsts interface{}) error {
	dsnps, err := fCli.GetAll(ctx, docRefs)
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for _, dsnp := range dsnps {
		if !dsnp.Exists() {
			continue
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// GetByQuery ... query a document
func GetByQuery(ctx context.Context, query firestore.Query, dst interface{}) (bool, error) {
	it := query.Documents(ctx)
	defer it.Stop()
	dsnp, err := it.Next()
	if err == iterator.Done {
		return false, nil
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	return true, nil
}

// ListByQuery ... query multi documents
func ListByQuery(ctx context.Context, query firestore.Query, dsts interface{}) error {
	it := query.Documents(ctx)
	defer it.Stop()
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for {
		dsnp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// ListByQueryCursor ... query multi documents with cursor
func ListByQueryCursor(ctx context.Context, query firestore.Query, limit int, cursor *firestore.DocumentSnapshot, dsts interface{}) (*firestore.DocumentSnapshot, error) {
	if cursor != nil {
		query = query.StartAfter(cursor)
	}
	it := query.Limit(limit).Documents(ctx)
	defer it.Stop()
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	var lastDsnp *firestore.DocumentSnapshot
	for {
		dsnp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(v)
		if err != nil {
			return nil, err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		rv.Set(reflect.Append(rv, rrv))
		lastDsnp = dsnp
	}
	if rv.Len() == limit {
		return lastDsnp, nil
	}
	return nil, nil
}

// TxGet ... get a single in transaction
func TxGet(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef, dst interface{}) (bool, error) {
	dsnp, err := tx.Get(docRef)
	if dsnp != nil && !dsnp.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	return true, nil
}

// TxGetMulti ... get multi documents in transaction
func TxGetMulti(ctx context.Context, tx *firestore.Transaction, docRefs []*firestore.DocumentRef, dsts interface{}) error {
	dsnps, err := tx.GetAll(docRefs)
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for _, dsnp := range dsnps {
		if !dsnp.Exists() {
			continue
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// TxGetByQuery ... query a document in transaction
func TxGetByQuery(ctx context.Context, tx *firestore.Transaction, query firestore.Query, dst interface{}) (bool, error) {
	it := tx.Documents(query)
	defer it.Stop()
	dsnp, err := it.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = dsnp.DataTo(dst)
	if err != nil {
		return false, err
	}
	setDocByDst(dst, dsnp.Ref)
	return true, nil
}

// TxListByQuery ... query multi documents in transaction
func TxListByQuery(ctx context.Context, tx *firestore.Transaction, query firestore.Query, dsts interface{}) error {
	it := tx.Documents(query)
	defer it.Stop()
	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	for {
		dsnp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		v := reflect.New(rrt).Interface()
		err = dsnp.DataTo(&v)
		if err != nil {
			return err
		}
		rrv := reflect.ValueOf(v)
		setDocByDsts(rrv, rrt, dsnp.Ref)
		rv.Set(reflect.Append(rv, rrv))
	}
	return nil
}

// Create ... create a document
func Create(ctx context.Context, colRef *firestore.CollectionRef, src interface{}) error {
	ref, _, err := colRef.Add(ctx, src)
	if err != nil {
		return err
	}
	setDocByDst(src, ref)
	return nil
}

// BtCreate ... create a document with batch
func BtCreate(ctx context.Context, bt *firestore.WriteBatch, colRef *firestore.CollectionRef, src interface{}) {
	id := util.StrUniqueID()
	ref := colRef.Doc(id)
	bt.Create(ref, src)
	setDocByDst(src, ref)
}

// TxCreate ... create a document in transaction
func TxCreate(ctx context.Context, tx *firestore.Transaction, colRef *firestore.CollectionRef, src interface{}) error {
	id := util.StrUniqueID()
	ref := colRef.Doc(id)
	err := tx.Create(ref, src)
	if err != nil {
		return err
	}
	setDocByDst(src, ref)
	return nil
}

// Update ... update a document
func Update(ctx context.Context, docRef *firestore.DocumentRef, kv map[string]interface{}) error {
	srcs := []firestore.Update{}
	for k, v := range kv {
		src := firestore.Update{Path: k, Value: v}
		srcs = append(srcs, src)
	}
	_, err := docRef.Update(ctx, srcs)
	if err != nil {
		return err
	}
	return nil
}

// BtUpdate ... update a document with batch
func BtUpdate(ctx context.Context, bt *firestore.WriteBatch, docRef *firestore.DocumentRef, kv map[string]interface{}) {
	srcs := []firestore.Update{}
	for k, v := range kv {
		src := firestore.Update{Path: k, Value: v}
		srcs = append(srcs, src)
	}
	_ = bt.Update(docRef, srcs)
}

// TxUpdate ... update a document in transaction
func TxUpdate(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef, kv map[string]interface{}) error {
	srcs := []firestore.Update{}
	for k, v := range kv {
		src := firestore.Update{Path: k, Value: v}
		srcs = append(srcs, src)
	}
	err := tx.Update(docRef, srcs)
	if err != nil {
		return err
	}
	return nil
}

// Set ... set a document
func Set(ctx context.Context, docRef *firestore.DocumentRef, src interface{}) error {
	_, err := docRef.Set(ctx, src)
	if err != nil {
		return err
	}
	setDocByDst(src, docRef)
	return nil
}

// BtSet ... set a document with batch
func BtSet(ctx context.Context, bt *firestore.WriteBatch, docRef *firestore.DocumentRef, src interface{}) {
	_ = bt.Set(docRef, src)
	setDocByDst(src, docRef)
}

// TxSet ... set a document in transaction
func TxSet(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef, src interface{}) error {
	err := tx.Set(docRef, src)
	if err != nil {
		return err
	}
	setDocByDst(src, docRef)
	return nil
}

// Delete ... delete a document
func Delete(ctx context.Context, docRef *firestore.DocumentRef) error {
	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// BtDelete ... delete a document with batch
func BtDelete(ctx context.Context, bt *firestore.WriteBatch, docRef *firestore.DocumentRef) {
	_ = bt.Delete(docRef)
}

// TxDelete ... delete a document in transaction
func TxDelete(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef) error {
	err := tx.Delete(docRef)
	if err != nil {
		return err
	}
	return nil
}
