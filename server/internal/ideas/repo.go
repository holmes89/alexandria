package ideas

type Repository interface {
	GetIdeas() ([]Idea, error)
	GetIdeaByID(id string) (Idea, error)
	CreateIdea(idea Idea) (Idea, error)
	AddIdeaResource(resource IdeaResource) error
	RemoveIdeaResource(resource IdeaResource) error
}
