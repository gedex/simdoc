package middleware

// @todo uncommented this until middleware can captures parsed URL Params.
// https://github.com/zenazn/goji/issues/76
import (
	"net/http"
	"strconv"

	"github.com/goji/context"
	"github.com/zenazn/goji/web"

	"github.com/gedex/simdoc/pkg/datastore"
)

// DocumentToContextInjector injects document information into the context.
func DocumentToContextInjector(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		docIdStr := c.URLParams["docId"]

		docId, _ := strconv.ParseInt(docIdStr, 10, 64)

		if docId > 0 {
			checkDoc(c, w, docId)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func checkDoc(c *web.C, w http.ResponseWriter, docId int64) {
	ctx := context.FromC(*c)
	user := ToUser(c)

	// @todo check document participants?

	doc, err := datastore.GetDocumentById(ctx, docId)
	switch {
	case err != nil && user == nil:
		w.WriteHeader(http.StatusUnauthorized)
		return
	case err != nil && user != nil:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	DocToC(c, doc)
}
