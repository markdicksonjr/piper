package piper

type Result struct {
	Err error
	Data interface{}
}

type MultiResult struct {
	Err error
	Results []*Result
}

type ResultHandler func(output MultiResult, isEndOfStream bool)

type Persister interface {
	PutRecords(records interface{}) error
}

type Filter interface {
	Filter(inbound interface{}) bool
}

type Mutator interface {
	Mutate(inbound []interface{}) ([]interface{}, error)
}

type Loader interface {
	LoadBatch(batchSize int) (data []interface{}, atEof bool, err error)
}

type ReadableRepository interface {
	Open() error

	// a simple "select * from collection where key IN (value)
	Query(collection string, key string, value string) (interface{}, error)

	// "false" for useJoin is a more performant alternative to a simple join, but will return fewer columns
	// e.g. "select * from Book where BookID IN (select BookID from BaseBook where BaseBookID = 1939);"
	// would be ReferenceQuery("Book", "BookID", "BaseBook", "BookID", "BaseBookID", "1939")
	ReferenceQuery(collection string, column string, subCollection string, subColumn string, subKey string, value string, useJoin bool) (interface{}, error)

	Close() error
}

type Pipe interface {
	Run(batchSize int, loader Loader, filter Filter, mutator Mutator, resultHandler ResultHandler) error
}
