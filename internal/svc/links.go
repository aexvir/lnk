package svc

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aexvir/lnk/internal/logging"
	"github.com/aexvir/lnk/internal/storage"
	"github.com/aexvir/lnk/internal/translation"
	"github.com/aexvir/lnk/proto"
)

type LinkStore interface {
	CreateLink(target string, slug *string) (string, error)
	GetLink(slug string) (*storage.Link, error)
	DeleteLink(slug string) error
	AllLinks() []*storage.Link

	GetTarget(slug string) (string, error)
	RegisterHit(slug string)
}

type LinksService struct {
	proto.UnimplementedLinksServer

	store LinkStore
	log   *logging.Logger
}

func NewLinksService(store LinkStore) LinksService {
	log := logging.NewLogger("lnk.links")
	return LinksService{
		store: store,
		log:   log,
	}
}

func (lgs *LinksService) ListLinks(ctx context.Context, _ *emptypb.Empty) (*proto.LinkList, error) {
	lgs.log.Write("ListLinks", "_")
	links := lgs.store.AllLinks()

	var list proto.LinkList
	for _, link := range links {
		list.Links = append(list.Links, translation.DbLinkToProto(link))
	}

	return &list, nil
}

func (lgs *LinksService) CreateLink(ctx context.Context, req *proto.CreateLinkReq) (*proto.LinkId, error) {
	lgs.log.Write("CreateLink", req.String())

	link, err := lgs.store.CreateLink(req.Target, req.Slug)
	if err != nil {
		return nil, fmt.Errorf("error creating link: %w", err)
	}

	return &proto.LinkId{
		Slug: link,
	}, nil
}

func (lgs *LinksService) GetLink(ctx context.Context, req *proto.LinkId) (*proto.LinkDetails, error) {
	lgs.log.Write("GetLink", "slug: %s", req.Slug)

	link, err := lgs.store.GetLink(req.Slug)
	if err != nil {
		return nil, fmt.Errorf("error getting link: %w", err)
	}

	return translation.DbLinkToProto(link), nil
}

func (lgs *LinksService) DeleteLink(ctx context.Context, req *proto.LinkId) (*emptypb.Empty, error) {
	lgs.log.Write("DeleteLink", "slug: %s", req.Slug)

	return nil, lgs.store.DeleteLink(req.Slug)
}

// LinkRedirectHandler fetches a shortened link by slug and if there is a link
// with that slug, it redirects to that link's target url.
func LinkRedirectHandler(store LinkStore) http.HandlerFunc {
	log := logging.NewLogger("lnk.redirect")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respond(w, http.StatusMethodNotAllowed, "only get requests allowed")
			return
		}

		slug := path.Base(r.URL.Path)
		log.Write("visit", "slug: %s", slug)

		target, err := store.GetTarget(slug)
		if err != nil {
			respond(w, http.StatusNotFound, err.Error())
			return
		}

		store.RegisterHit(slug)
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	}
}

func respond(w http.ResponseWriter, status int, msg string, args ...any) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(msg, args...)))
}
