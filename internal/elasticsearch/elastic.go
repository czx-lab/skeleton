package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	optins []elastic.ClientOptionFunc
	conn   *elastic.Client
}

type Document struct {
	Index string
	Id    string
	Body  any
}

type Index struct {
	Name    string
	Mapping string
}

var once sync.Once
var instance *Elasticsearch

func NewElastic(optins ...elastic.ClientOptionFunc) (*Elasticsearch, error) {
	client, err := elastic.NewClient(optins...)
	if err != nil {
		return nil, err
	}
	once.Do(func() {
		instance = &Elasticsearch{
			optins: optins,
			conn:   client,
		}
	})
	return instance, nil
}

func (e *Elasticsearch) CreateIndex(index Index) (*elastic.IndicesCreateResult, error) {
	exist, err := e.HasIndex(index.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("index already exists")
	}
	return e.conn.CreateIndex(index.Name).BodyString(index.Mapping).Do(context.Background())
}

func (e *Elasticsearch) DelIndex(name string) (*elastic.IndicesDeleteResponse, error) {
	exist, err := e.HasIndex(name)
	if err != nil {
		return nil, err
	}
	if exist {
		return e.conn.DeleteIndex(name).Do(context.Background())
	}
	return nil, errors.New("index does not exist")
}

func (e *Elasticsearch) HasIndex(name string) (bool, error) {
	return e.conn.IndexExists(name).Do(context.Background())
}

func (e *Elasticsearch) CreateDoc(doc Document) (*elastic.IndexResponse, error) {
	exist, err := e.HasIndex(doc.Index)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("index does not exist")
	}
	return e.conn.Index().Index(doc.Index).Id(doc.Id).BodyJson(doc.Body).Do(context.Background())
}

type createRes struct {
	index string
	num   int
	err   error
}

func (e *Elasticsearch) BatchCreateDoc(docs []Document) (map[string]int, error) {
	if len(docs) < 1 {
		return nil, errors.New("the added data cannot be empty")
	}
	documents := make(map[string][]Document)
	for _, doc := range docs {
		documents[doc.Index] = append(documents[doc.Index], doc)
	}
	result := make(map[string]int)
	wg := &sync.WaitGroup{}
	res := make(chan createRes, len(documents))
	for index, doc := range documents {
		wg.Add(1)
		go e.createDoc(wg, index, doc, res)
	}
	wg.Wait()
	close(res)
	for v := range res {
		result[v.index] = v.num
		if v.err != nil {
			return result, v.err
		}
	}
	return result, nil
}

func (e *Elasticsearch) createDoc(wg *sync.WaitGroup, index string, doc []Document, res chan createRes) {
	defer wg.Done()
	bulk := e.conn.Bulk().Index(index)
	requests := make([]elastic.BulkableRequest, len(doc))
	for _, item := range doc {
		requests = append(requests, elastic.NewBulkIndexRequest().Id(item.Id).Doc(item.Body))
	}
	_, err := bulk.Add(requests...).Do(context.Background())
	if err != nil {
		res <- createRes{
			index: index,
			num:   0,
			err:   err,
		}
		return
	}
	num := bulk.NumberOfActions()
	if num < 1 {
		res <- createRes{
			index: index,
			num:   num,
			err:   fmt.Errorf("the data added to the index %s cannot be empty", index),
		}
		return
	}
	res <- createRes{
		index: index,
		num:   num,
	}
}

type UpdateDoc struct {
	Document
	Query  elastic.Query
	Script *elastic.Script
}

func (e *Elasticsearch) BatchUpdateDoc(update Document) (*elastic.UpdateResponse, error) {
	return e.conn.Update().Index(update.Index).Id(update.Id).Doc(update.Body).Do(context.Background())
}

func (e *Elasticsearch) UpdateDocById(update UpdateDoc) (*elastic.UpdateResponse, error) {
	return e.conn.Update().Index(update.Index).Id(update.Id).Script(update.Script).Do(context.Background())
}

func (e *Elasticsearch) UpdateDoc(update UpdateDoc) (*elastic.BulkIndexByScrollResponse, error) {
	return e.conn.UpdateByQuery(update.Document.Index).Query(update.Query).Script(update.Script).ProceedOnVersionConflict().Do(context.Background())
}

type SearchQuery struct {
	Index     []string
	Query     elastic.Query
	From      int
	Size      int
	Pretty    bool
	Sort      bool
	SortField string
}

type BatchQuery struct {
	Index string
	Id    string
}

func (e *Elasticsearch) BatchQueryDoc(search SearchQuery) (*elastic.SearchResult, error) {
	if len(search.Index) < 1 {
		return nil, errors.New("index cannot be empty")
	}
	query := e.conn.Search().Index(search.Index...)
	if search.Query != nil {
		query.Query(search.Query)
	}
	if search.SortField != "" {
		query.Sort(search.SortField, search.Sort)
	}
	if search.From > 0 {
		query.From(search.From)
	}
	if search.Size < 1 {
		query.Size(20)
	}
	return query.Pretty(search.Pretty).Do(context.Background())
}

func (e *Elasticsearch) BatchQueryDocByIds(querys ...BatchQuery) ([][]byte, error) {
	if len(querys) == 0 {
		return nil, errors.New("index id is empty")
	}
	items := make([]*elastic.MultiGetItem, len(querys))
	for _, item := range querys {
		items = append(items, elastic.NewMultiGetItem().Index(item.Index).Id(item.Id))
	}
	result, err := e.conn.MultiGet().Add(items...).Do(context.Background())
	if err != nil {
		return nil, err
	}
	docs := make([][]byte, len(result.Docs))
	for _, doc := range result.Docs {
		tmp, _ := doc.Source.MarshalJSON()
		docs = append(docs, tmp)
	}
	return docs, nil
}

func (e *Elasticsearch) QueryDocById(index, id string) ([]byte, error) {
	result, err := e.conn.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, nil
	}
	data, err := result.Source.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return data, nil
}

type QueryDel struct {
	Index string
	Query elastic.Query
}

func (e *Elasticsearch) DelDocById(query BatchQuery) (*elastic.DeleteResponse, error) {
	return e.conn.Delete().Index(query.Index).Id(query.Id).Do(context.Background())
}

func (e *Elasticsearch) DelQueryDoc(query QueryDel) (*elastic.BulkIndexByScrollResponse, error) {
	return e.conn.DeleteByQuery(query.Index).Query(query.Query).ProceedOnVersionConflict().Do(context.Background())
}
