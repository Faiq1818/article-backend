package article

import (
	"fmt"
	"net/http"
)

func (h *AuthHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	fmt.Print(queryParams.Get("page"))

}
