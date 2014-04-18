package yunyun-utils

/*
 * define variables
 */
var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)
var templates = template.Must(template.ParseFiles("login.html", "account.html", "register.html"))
var validPath = regexp.MustCompile("^/(login|account|register)/?([a-zA-Z0-9]*)$")

/* 
 * webservice functions
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
func accountHandler(w http.ResponseWriter, r *http.Request, title string) {
  user := r.FormValue("user")
  pass := r.FormValue("pass")
  //fmt.Println("user: '" + user + "'")
  //fmt.Println("pass: '" + pass + "'")
  err := loginUser(user, pass)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	  http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	renderTemplate(w, "account", &Page{User: user, Body: []byte("")})
}
func registerHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{User: title}
	}
	renderTemplate(w, "register", p)
}
func loginHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{User: title}
	}
	renderTemplate(w, "login", p)
}