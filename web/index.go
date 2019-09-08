package web

import (
	"html/template"
	"net/http"
)

var indexPage = template.Must(template.New("").Parse(`
<!doctype html>
<html lang="en">

<head>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>

<body>
    <div class="container">
        <h1>Search products</h1>
        <div class="alert alert-info" role="alert">
            Hint: you can use any username/password to access the search
        </div>
        <form method="{{.Method}}" action="{{.Action}}">
            <input type="hidden" name="pretty" value="yes" />
            <div class="form-group">
                <label for="q">Search for<sup>*</sup></label>
                <input id="q" type="text" name="q" class="form-control" required/>
            </div>
            <div class="form-group">
                <label for="size">return</label>
                <input id="size" type="number" class="form-control" name="size" />
                <small class="form-text text-muted">results per page</small>
            </div>
            <div class="form-group">
                <label for="from">starting from</label>
                <input id="from" type="number" class="form-control" name="from" />
            </div>
            <div class="form-group">
                <label for="sort">sorted by</label>
                <input id="sort" type="text" name="sort" class="form-control" placeholder="i.e. price:desc" />
            </div>
            <div class="form-group">
                <label for="filter">and filtered by</label>
                <input id="filter" type="text" name="filter" class="form-control" placeholder="use Lucene syntax" />
            </div>
            <button type="reset" class="btn">Reset</button>
            <button type="submit" class="btn btn-primary">Go!</button>
        </form>
    </div>
</body>

</html>
`))

// IndexHandler serves a static HTML page to execute search query
func IndexHandler(method, action string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		indexPage.Execute(w, struct {
			Method, Action string
		}{method, action})
	})
}
