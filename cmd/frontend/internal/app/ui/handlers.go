package ui

import (
	"github.com/prince1809/sourcegraph/cmd/frontend/envvar"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/app/jscontext"
	"github.com/prince1809/sourcegraph/pkg/actor"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"html/template"
	"net/http"
	"net/url"
)

type InjectedHTML struct {
	HeadTop    template.HTML
	HeadBottom template.HTML
	BodyTop    template.HTML
	BodyBottom template.HTML
}

type Metadata struct {
	// Title is the title of the page for Twitter cards, OpenGraph, etc.
	// etc. "Open in Sourcegraph"
	Title string

	// Description is the description of the page for Twitter cards, OpenGraph,
	// etc. e.g. "View this link in Sourcegraph Editor".
	Description string
}

type Common struct {
	Injected InjectedHTML
	Metadata *Metadata
	Context  jscontext.JSContext
	AssetURL string
	Title    string
	Error    *pageError
}

// newCommon builds a *Common data structure, returning an error if one occurs.
// In the event of the repository having been renamed, the request is handled
// by newCommon and nil, nil is returned. Basic usage looks like:
//
//　common, err := newCommon(w, r, serveError)
// if err != nil {
//     return err
// }
// if common == nil {
//     return nil // request was handled
// }
//
// In the case of a repository that is cloning, a Common data structure is
// returned but it has an incomplete RevSpec.
func newCommon(w http.ResponseWriter, r *http.Request, title string, serveError func(w http.ResponseWriter, r *http.Request, err error, statusCode int)) (*Common, error) {
	common := &Common{
		Injected: InjectedHTML{
			HeadTop:    template.HTML(conf.Get().Critical.HtmlHeadTop),
			HeadBottom: template.HTML(conf.Get().Critical.HtmlHeadBottom),
			BodyTop:    template.HTML(conf.Get().Critical.HtmlBodyTop),
			BodyBottom: template.HTML(conf.Get().Critical.HtmlBodyBottom),
		},
	}
	return common, nil
}

func serveHome(w http.ResponseWriter, r *http.Request) error {
	common, err := newCommon(w, r, "sourcegraph", serveError)
	if err != nil {
		return err
	}
	if common == nil {
		return nil // request was handled
	}

	if envvar.SourcegraphDotComMode() && !actor.FromContext(r.Context()).IsAuthenticated() {
		// The user is not signed in and tried to access Sourcegraph.com. redirects to
		// about.sourcegraph.com so they see general info page.
		http.Redirect(w, r, (&url.URL{Scheme: aboutRedirectScheme, Host: aboutRedirectHost}).String(), http.StatusTemporaryRedirect)
		return nil
	}

	// On non-Sourcegraph.com instances, there is no separate homepage, so redirect to /search.
	r.URL.Path = "/search"
	http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
	return nil
}
