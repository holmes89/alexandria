package links

import (
	"alexandria/internal/tags"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Repository interface {
	FindAllLinks() ([]Link, error)
	FindLinkByID(id string) (Link, error)
	CreateLink(Link) (Link, error)
}

type Service interface {
	FindAll() ([]Link, error)
	FindByID(id string) (Link, error)
	Create(Link) (Link, error)
	AddTag(id string, tag string) error
	RemoveTag(id string, tag string) error
}

type service struct {
	repo     Repository
	tagsRepo tags.Repository
}

func NewService(repo Repository, tagsRepo tags.Repository) Service {
	return &service{
		repo:     repo,
		tagsRepo: tagsRepo,
	}
}

func (s *service) FindAll() ([]Link, error) {
	return s.repo.FindAllLinks()
}

func (s *service) Create(entity Link) (linkEntity Link, err error) {

	link := entity.Link
	if link == "" {
		return linkEntity, errors.New("link cannot be empty")
	}

	// Request the HTML page.
	res, err := http.Get(link)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch website")
		return linkEntity, errors.New("unable to fetch website")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.WithError(err).Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return linkEntity, errors.New("unable to fetch website")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	title := doc.Find("title").First().Text()
	icon, ok := doc.Find("link[rel='apple-touch-icon']").Attr("href")
	if !ok || len(icon) >= 2048 {
		logrus.Warn("unable to find apple-touch-icon looking for other icon")
		icon, ok = doc.Find("link[rel='icon']").Attr("href")
		if !ok || len(icon) >= 2048  {
			logrus.Warn("unable to find icon looking for other icon")
			icon, ok = doc.Find("link[rel='shortcut icon']").Attr("href")
			if !ok {
				icon = ""
			}
		}
	}

	linkEntity = Link{
		Link:        link,
		DisplayName: title,
		IconPath:    parseIcon(link, icon),
	}

	return s.repo.CreateLink(linkEntity)
}
func (s *service) FindByID(id string) (Link, error) {
	return s.repo.FindLinkByID(id)
}
func (s *service) AddTag(id string, tag string) error {
	return s.tagsRepo.AddResourceTag(id, tags.LinksResource, tag)
}

func (s *service) RemoveTag(id string, tag string) error {
	return s.tagsRepo.RemoveResourceTag(id, tag)
}

func parseIcon(website, path string) string {
	if path == "" || strings.Contains(path, "//") {
		return path
	}

	u, err := url.Parse(website)
	if err != nil {
		logrus.WithError(err).Error("unable to parse site")
		return ""
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, path)
}
