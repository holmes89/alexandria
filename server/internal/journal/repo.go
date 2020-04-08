package journal

type Repository interface {
	FindAllEntries() ([]Entry, error)
	CreateEntry(Entry) (Entry, error)
}
