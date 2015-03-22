package heroku

import (
	"net/http"

	"github.com/bgentry/heroku-go"
	"github.com/remind101/empire/empire"
	"github.com/remind101/empire/empire/pkg/httpx"
	"golang.org/x/net/context"
)

type Domain heroku.Domain

func newDomain(d *empire.Domain) *Domain {
	return &Domain{
		Id:        d.ID,
		Hostname:  d.Hostname,
		CreatedAt: d.CreatedAt,
	}
}

type GetDomains struct {
	*empire.Empire
}

func (h *GetDomains) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a, err := findApp(ctx, h)
	if err != nil {
		return err
	}

	d, err := h.DomainsFindByApp(a)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	return Encode(w, d)
}

type PostDomainsForm struct {
	Hostname string `json:"hostname"`
}

type PostDomains struct {
	*empire.Empire
}

func (h *PostDomains) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a, err := findApp(ctx, h)
	if err != nil {
		return err
	}

	var form PostDomainsForm

	if err := Decode(r, &form); err != nil {
		return err
	}

	domain := &empire.Domain{
		AppName:  a.Name,
		Hostname: form.Hostname,
	}
	d, err := h.DomainsCreate(domain)
	if err != nil {
		return err
	}

	w.WriteHeader(201)
	return Encode(w, newDomain(d))
}

type DeleteDomain struct {
	*empire.Empire
}

func (h *DeleteDomain) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a, err := findApp(ctx, h)
	if err != nil {
		return err
	}

	vars := httpx.Vars(ctx)
	name := vars["hostname"]

	d, err := h.DomainsFindByHostname(name)
	if err != nil {
		return err
	}

	if d.AppName != a.Name {
		return ErrNotFound
	}

	if err = h.DomainsDestroy(d); err != nil {
		return err
	}

	return NoContent(w)
}
