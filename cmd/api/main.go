package main

import (
	"log"

	"github.com/patriciabonaldy/bequest_challenge/cmd/api/bootstrap"
)

// @title API document title
// @version version(1.0)
// @description Description of specifications
// @Precautions when using termsOfService specifications

// @contact.name API supporter
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name license(Mandatory)
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
