package router

import (
	"net/http"
	"strings"

	"github.com/sqs/mux"
)

// RepoSpecPathPattern is the path pattern for repository URIs.
//
// TODO(sqs): match the "R$rid" format too.
var RepoSpecPathPattern = `{RepoSpec:(?:[^/.@][^/@]*/)+(?:[^/.@][^/@]*)}`

// RepoPathPattern is the path pattern for repository URIs with optional revisions.
var RepoPathPattern = RepoSpecPathPattern + `{Rev:(?:@[\w-.]+)?}`

// FixRepoVars is a mux.PostMatchFunc that cleans and normalizes the repository URI.
func FixRepoVars(req *http.Request, match *mux.RouteMatch, r *mux.Route) {
	if rev, present := match.Vars["Rev"]; present {
		if rev == "" {
			delete(match.Vars, "Rev")
		} else {
			match.Vars["Rev"] = strings.TrimPrefix(rev, "@")
		}
	}
}

// PrepareRepoRouteVars is a mux.BuildVarsFunc that converts from a cleaned
// and normalized repository URI to a repository component in the route.
func PrepareRepoRouteVars(vars map[string]string) map[string]string {
	if rev, present := vars["Rev"]; !present {
		vars["Rev"] = ""
	} else if rev != "" {
		vars["Rev"] = "@" + rev
	}
	return vars
}
