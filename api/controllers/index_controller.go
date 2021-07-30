package controllers

import (
	"html/template"
	"net/http"
	"sort"

	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

// Heartbeat godoc
// @Summary Get heartbeat of the system
// @Description Get heartbeat of the system. This API can be used for Service discovery and to get the Health Status.
// @Tags Heartbeat
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /heartbeat [get]
func (server *Server) Heartbeat(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome User Manegement API. I am live!!")
}

func (server *Server) IndexPage(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["github"] = "Github"
	m["google"] = "Google"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	var indexTemplate = `{{range $key,$value:=.Providers}}
    	<p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
		{{end}}`

	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, providerIndex)
}
