/*
Package inceis "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package inceis

import (
	"github.com/MerinEREN/iiPackages/api/account"
	"github.com/MerinEREN/iiPackages/api/contents"
	"github.com/MerinEREN/iiPackages/api/demand"
	"github.com/MerinEREN/iiPackages/api/demands"
	"github.com/MerinEREN/iiPackages/api/languages"
	"github.com/MerinEREN/iiPackages/api/offers"
	"github.com/MerinEREN/iiPackages/api/page"
	"github.com/MerinEREN/iiPackages/api/pages"
	"github.com/MerinEREN/iiPackages/api/role"
	"github.com/MerinEREN/iiPackages/api/roleType"
	"github.com/MerinEREN/iiPackages/api/roleTypes"
	"github.com/MerinEREN/iiPackages/api/roles"
	"github.com/MerinEREN/iiPackages/api/servicePacks"
	"github.com/MerinEREN/iiPackages/api/settingsAccount"
	"github.com/MerinEREN/iiPackages/api/signin"
	"github.com/MerinEREN/iiPackages/api/signout"
	"github.com/MerinEREN/iiPackages/api/tag"
	"github.com/MerinEREN/iiPackages/api/tags"
	"github.com/MerinEREN/iiPackages/api/timeline"
	"github.com/MerinEREN/iiPackages/api/user"
	"github.com/MerinEREN/iiPackages/api/userRoles"
	"github.com/MerinEREN/iiPackages/api/userTags"
	"github.com/MerinEREN/iiPackages/api/users"
	"github.com/MerinEREN/iiPackages/session"
	"strings"
	// "github.com/MerinEREN/iiPackages/cookie"
	"github.com/MerinEREN/iiPackages/page/template"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
	// "regexp"
	"time"
)

var _ memcache.Item // For debugging, delete when done.

var (
// CHANGE THE REGEXP BELOW !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// validPath = regexp.MustCompile("^/|[/A-Za-z0-9]$")
)

func init() {
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/",
		http.TimeoutHandler(http.HandlerFunc(makeHandlerFunc(signin.Handler)),
			1000*time.Millisecond,
			"This is http.TimeoutHandler(handler, time.Duration, message) "+
				"message bitch =)"))
	// http.HandleFunc("/", makeHandlerFunc(signin.Handler))
	http.HandleFunc("/contents", makeHandlerFunc(contents.Handler))
	http.HandleFunc("/users/", makeHandlerFunc(user.Handler))
	http.HandleFunc("/userTags/", makeHandlerFunc(userTags.Handler))
	http.HandleFunc("/userRoles/", makeHandlerFunc(userRoles.Handler))
	http.HandleFunc("/languages", makeHandlerFunc(languages.Handler))
	http.HandleFunc("/signout", makeHandlerFunc(signout.Handler))
	http.HandleFunc("/accounts/", makeHandlerFunc(account.Handler))
	http.HandleFunc("/demands", makeHandlerFunc(demands.Handler))
	http.HandleFunc("/demands/", makeHandlerFunc(demand.Handler))
	http.HandleFunc("/offers", makeHandlerFunc(offers.Handler))
	http.HandleFunc("/servicePacks", makeHandlerFunc(servicePacks.Handler))
	http.HandleFunc("/timeline", makeHandlerFunc(timeline.Handler))
	http.HandleFunc("/pages", makeHandlerFunc(pages.Handler))
	http.HandleFunc("/pages/", makeHandlerFunc(page.Handler))
	http.HandleFunc("/tags", makeHandlerFunc(tags.Handler))
	http.HandleFunc("/tags/", makeHandlerFunc(tag.Handler))
	http.HandleFunc("/roles", makeHandlerFunc(roles.Handler))
	http.HandleFunc("/roles/", makeHandlerFunc(role.Handler))
	http.HandleFunc("/roleTypes", makeHandlerFunc(roleTypes.Handler))
	http.HandleFunc("/roleTypes/", makeHandlerFunc(roleType.Handler))
	http.HandleFunc("/settingsAccount", makeHandlerFunc(settingsAccount.Handler))
	http.HandleFunc("/users", makeHandlerFunc(users.Handler))
	// http.HandleFunc("/accounts", makeHandlerFunc(accounts.Handler))
	// http.HandleFunc("/signUp", makeHandlerFunc(signUpHandler))
	// http.HandleFunc("/logIn", makeHandlerFunc(logInHandler))
	// http.HandleFunc("/accounts", makeHandlerFunc(accountsHandler))
	/* if http.PostForm("/logIn", data); err != nil {
		http.Err(w, "Internal server error while login",
			http.StatusBadRequest)
	} */
	fs := http.FileServer(http.Dir("../iiClient/public"))
	// http.Handle("/css/", fs)
	http.Handle("/img/", fs)
	http.Handle("/js/", fs)
	/* log.Printf("About to listen on 10443. " +
	"Go to https://192.168.1.100:10443/ " +
	"or https://localhost:10443/") */
	// Redirecting to a port or a domain etc.
	// go http.ListenAndServe(":8080",
	// http.RedirectHandler("https://192.168.1.100:10443", 301))
	// err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	// ListenAndServe and ListenAndServeTLS always returns a non-nil error !!!
	// log.Fatal(err)
}

/* func logInHandler(w http.ResponseWriter, r *http.Request, s string) {
	p, err := content.Get(r, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "POST" {
		key := "email"
		email := r.PostFormValue(key)
		key = "password"
		password := r.PostFormValue(key)
		acc, err := account.VerifyUser(c, email, password)
		switch err {
		case account.EmailNotExist:
			fmt.Fprintln(w, err)
		case account.ExistingEmail:
			for _, u := range acc.Users {
				if u.Email == email {
					// ALLWAYS CREATE COOKIE BEFORE EXECUTING TEMPLATE
					cookie.Set(w, r, "session", u.UUID)
				}
			}
			// NEWER EXECUTE TEMPLATE OR WRITE ANYTHING TO THE BODY BEFORE
			// REDIRECT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			http.Redirect(w, r, "/accounts/"+acc.Name, http.StatusSeeOther)
		case account.InvalidPassword:
			fmt.Fprintln(w, err)
		default:
			// Status code could be wrong
			http.Error(w, err.Error(), http.StatusNotImplemented)
			log.Fatalln(err)
		}
	}
	template.RenderLogIn(w, p)
} */

type handlerFuncWithSessionParam func(s *session.Session)

// CHANGE THE SESSION THING !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func makeHandlerFunc(fn handlerFuncWithSessionParam) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/* m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Printf("Invalid Path: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		} */
		/* for _, val := range m {
			fmt.Println(val)
		}*/
		// CHANGE CONTENT AND TEMPLATE THINGS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		s := new(session.Session)
		s.Init(w, r)
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			// Authenticate the client
			// Check should be in here to be able to make content data requests
			// and redirect content page request.
			if s.U == nil && r.URL.Path != "/" {
				log.Println("REDIRECTTT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			log.Println("Getting template !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			template.RenderIndex(w)
			// } else if strings.Contains(r.Header.Get("Accept"), "application/json") {
		} else {
			log.Println("Getting data !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			// SECURITY CHECK
			// add other landing page request calls later.
			if s.U != nil || (r.URL.Path == "/contents" || r.URL.Path == "/") {
				fn(s)
			}
		}
	}
}

// HEADER ALWAYS SHOULD BE SET BEFORE ANYTHING WRITE A PAGE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!
// w.Header().Set("Content-Type", "text/html"; charset=utf-8")
//fmt.Fprintln(w, things...) // Writes to the body
