/*
 Copyright 2018 Padduck, LLC
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/web/api"
	"github.com/pufferpanel/pufferpanel/web/auth"
	"github.com/pufferpanel/pufferpanel/web/daemon"
	"github.com/pufferpanel/pufferpanel/web/handlers"
	"github.com/pufferpanel/pufferpanel/web/oauth2"
	"strings"
)

const ClientPath = "client/dist"
const IndexFile = ClientPath + "/index.html"

var noHandle404 = []string{"/api/", "/oauth2/", "/daemon/"}

func RegisterRoutes(e *gin.Engine) {
	api.RegisterRoutes(e.Group("/api", handlers.HasOAuth2Token))
	oauth2.RegisterRoutes(e.Group("/oauth2"))
	auth.RegisterRoutes(e.Group("/auth"))
	daemon.RegisterRoutes(e.Group("/daemon", handlers.HasOAuth2Token))

	e.Static("/css", ClientPath+"/css")
	e.Static("/fonts", ClientPath+"/fonts")
	e.Static("/img", ClientPath+"/img")
	e.Static("/js", ClientPath+"/js")
	e.StaticFile("/favicon.png", ClientPath+"/favicon.png")
	e.StaticFile("/favicon.ico", ClientPath+"/favicon.ico")
	//e.StaticFile("/", IndexFile)
	e.NoRoute(handlers.AuthMiddleware, handle404)
}

func handle404(c *gin.Context) {
	for _, v := range noHandle404 {
		if strings.HasPrefix(c.Request.URL.Path, v) {
			c.AbortWithStatus(404)
			return
		}
	}

	if strings.HasSuffix(c.Request.URL.Path, ".js") {
		c.Writer.Header().Set("Content-Type", "application/js")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	if strings.HasSuffix(c.Request.URL.Path, ".css") {
		c.Writer.Header().Set("Content-Type", "application/css")
		c.File(ClientPath + c.Request.URL.Path)
		return
	}

	c.File(IndexFile)
}