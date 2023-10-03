package routes

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/salberternst/fsn_landing_page/pkg/middleware"
	"github.com/salberternst/fsn_landing_page/pkg/models"
)

func AddRoutes(r *gin.Engine) {
	api := r.Group("/")

	// only use the middlewares in the api group
	api.Use(middleware.TokenMiddleware())
	addUserRoutes(api)
	addHealthRoutes(api)

	remote, err := url.Parse("http://thingsboard.192-168-178-60.nip.io")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		credentials, err := models.GetThingsboardToken(req.Header.Get("X-Access-Token"))
		if err != nil {
			req.Context().Done()
			req.Close = true
			return
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials.Token))
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	r.Any("/api/*proxyPath", gin.WrapH(proxy))
}
