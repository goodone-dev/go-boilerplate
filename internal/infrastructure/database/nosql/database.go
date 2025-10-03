package nosql

type Entity interface {
	CollectionName() string
	RepositoryName() string
}
