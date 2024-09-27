package main

import (
	database "Nanouche/Controller"
	"crypto/sha256"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type PageHome struct {
	Posts           []database.Post
	IsConnecter     bool
	ConnectUserInfo string
	ConnectUserImg  string
	Email           string
	LoggedInUser    database.User
	Error           int
	Try             bool
}

type PagePost struct {
	OnePost         database.Post
	IdPost          string
	Responses       []database.Reponse
	ConnectUserName string
	ConnectUserImg  string
}

type PageNewPost struct {
	ConnectUserImg string
}

func initStruct() (PageHome, PagePost, PageNewPost) {
	var home PageHome
	home.Posts = database.GetAllPost()
	home.IsConnecter = false
	home.Error = 0
	home.Email = ""
	home.ConnectUserInfo = ""

	var post PagePost
	post.ConnectUserName = ""
	post.ConnectUserImg = ""

	var newpost PageNewPost

	return home, post, newpost
}

var tmplHome = template.Must(template.ParseFiles("./View/Html/home.html"))
var tmplProfil = template.Must(template.ParseFiles("./View/Html/profil.html"))
var tmplPost = template.Must(template.ParseFiles("./View/Html/post.html"))
var tmpl404 = template.Must(template.ParseFiles("./View/Html/404.html"))
var tmplNewPost = template.Must(template.ParseFiles("./View/Html/newpost.html"))
var HomeStruct, PostStruct, NewPostStruct = initStruct()

func main() {

	http.Handle("/View/css/", http.StripPrefix("/View/css/", http.FileServer(http.Dir("View/css"))))
	http.Handle("/View/images/", http.StripPrefix("/View/images/", http.FileServer(http.Dir("View/images"))))
	http.Handle("/View/js/", http.StripPrefix("/View/js/", http.FileServer(http.Dir("View/js"))))

	fmt.Printf("\n")
	fmt.Println("http://localhost:8080/")
	fmt.Printf("\n")

	http.HandleFunc("/404", pageNotFoundHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ho", homeTransiHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/updateProfile", updateProfileHandler)
	http.HandleFunc("/newpost", newPostHandler)
	http.HandleFunc("/newpostTransi", newPostTransiHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.ListenAndServe("localhost:8080", nil)
}

var filtres = "home"

func filter(HomeStruct *PageHome) {
	if filtres == "home" {
		HomeStruct.Posts = database.SelectByDescending("Date")
	} else if filtres == "animes" {
		HomeStruct.Posts = database.GetTagAnime()
	} else if filtres == "series" {
		HomeStruct.Posts = database.GetTagSerie()
	} else if filtres == "Plikes" {
		HomeStruct.Posts = database.SelectByDescending("nbrLikes")
	} else if filtres == "Pdislikes" {
		HomeStruct.Posts = database.SelectByDescending("NbrDislikes")
	} else if filtres == "Mlikes" {
		HomeStruct.Posts = database.SelectByAscending("nbrLikes")
	} else if filtres == "Mdislikes" {
		HomeStruct.Posts = database.SelectByAscending("NbrDislikes")
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/ho" && r.URL.Path != "/post" && r.URL.Path != "/newpost" && r.URL.Path != "/newpostTransi" && r.URL.Path != "/login" && r.URL.Path != "/logout" && r.URL.Path != "/profil" && r.URL.Path != "/updateProfile" {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	} else {
		err := tmplHome.Execute(w, HomeStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func homeTransiHandler(w http.ResponseWriter, r *http.Request) {
	buLikesDislikes := r.FormValue("bulike/dislike")
	if buLikesDislikes != "" {
		idBu := strings.Split(buLikesDislikes, ",")
		if idBu[1] == "like" {
			nbr := database.RecupNbr("nbrLikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrLikes", nbr, idBu[0])
		} else if idBu[1] == "dislike" {
			nbr := database.RecupNbr("nbrDislikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrDislikes", nbr, idBu[0])
		}
	}

	headerLinks := r.FormValue("link")
	if headerLinks != "" {
		filtres = headerLinks
	}

	BuMenuDeroulant := r.FormValue("BuMenuDeroulant")
	if BuMenuDeroulant != "" {
		filtres = BuMenuDeroulant
	}

	trucs := r.FormValue("creerPost")
	if trucs != "" {
		fmt.Println(trucs)
	}

	filter(&HomeStruct)
	http.Redirect(w, r, "/", http.StatusFound)

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	buLikesDislikes := r.FormValue("bulike/dislike")
	if buLikesDislikes != "" {
		idBu := strings.Split(buLikesDislikes, ",")
		if idBu[1] == "like" {
			nbr := database.RecupNbr("nbrLikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrLikes", nbr, idBu[0])
		} else if idBu[1] == "dislike" {
			nbr := database.RecupNbr("nbrDislikes", idBu[0])
			nbr += 1
			database.UpdateNbr("nbrDislikes", nbr, idBu[0])
		}
	}

	idPost := r.FormValue("buPost")
	if idPost == "retour" {
		PostStruct.IdPost = ""
		http.Redirect(w, r, "/ho", http.StatusFound)
	} else if idPost == "Envoyer" {
		rep := r.FormValue("response")
		database.DatabaseAndReponse([]string{PostStruct.IdPost, PostStruct.ConnectUserName, rep, transformDate(), PostStruct.ConnectUserImg})
	}
	if idPost != "Envoyer" && idPost != "" {
		PostStruct.IdPost = idPost
	}
	PostStruct.OnePost = database.GetOnePost(PostStruct.IdPost)
	PostStruct.Responses = database.GetResponses(PostStruct.IdPost)
	err := tmplPost.Execute(w, PostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	if !HomeStruct.IsConnecter {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := tmplNewPost.Execute(w, NewPostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newPostTransiHandler(w http.ResponseWriter, r *http.Request) {
	buton := r.FormValue("buCreerPost")
	if buton == "creer" {
		titre := r.FormValue("topic-name")
		tag := r.FormValue("category")
		description := r.FormValue("content")
		database.DatabaseAndPost([]string{HomeStruct.ConnectUserInfo, tag, titre, description, strconv.Itoa(0), strconv.Itoa(0), transformDate(), "117902422"})
	}
	filtres = "home"
	filter(&HomeStruct)
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	connectionEmail := r.FormValue("ConnectionEmail")
	connectionMdp := r.FormValue("ConnectionMdp")
	if connectionEmail != "" {
		user := database.GetUser(connectionEmail)
		if user.Password == HashPassword2(connectionMdp) && user.Email == connectionEmail {
			sessionID := uuid.New().String()
			expiration := time.Now().Add(1 * time.Hour)
			cookie := &http.Cookie{
				Name:     "Session",
				Value:    sessionID,
				Expires:  expiration,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			database.DatabaseAndSession([]string{connectionEmail, sessionID})
			HomeStruct.Email = user.Email
			HomeStruct.ConnectUserInfo = user.Username
			HomeStruct.ConnectUserImg = user.Img
			HomeStruct.IsConnecter = true
			HomeStruct.Error = 0

			PostStruct.ConnectUserName = user.Username
			PostStruct.ConnectUserImg = user.Img

			NewPostStruct.ConnectUserImg = user.Img
		} else {
			HomeStruct.Email = ""
			HomeStruct.ConnectUserInfo = ""
			HomeStruct.ConnectUserImg = ""
			HomeStruct.IsConnecter = false
			HomeStruct.Error = 2

			PostStruct.ConnectUserName = ""
			PostStruct.ConnectUserImg = ""

			NewPostStruct.ConnectUserImg = ""
		}
	}

	// Inscription ----------------------------------------------------------------
	inscriptionName := r.FormValue("InscriptionName")
	inscriptionEmail := r.FormValue("InscriptionEmail")
	inscriptionMdp := r.FormValue("InscriptionMdp")

	if inscriptionName != "" {
		if !database.GetEmail(inscriptionEmail) {
			img := renderImg()
			database.DatabaseAndUsers([]string{inscriptionEmail, inscriptionName, HashPassword2(inscriptionMdp), img})
			sessionID := uuid.New().String()
			expiration := time.Now().Add(1 * time.Hour)
			cookie := &http.Cookie{
				Name:     "Session",
				Value:    sessionID,
				Expires:  expiration,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			database.DatabaseAndSession([]string{inscriptionEmail, sessionID})
			HomeStruct.Email = inscriptionEmail
			HomeStruct.ConnectUserInfo = inscriptionName
			HomeStruct.ConnectUserImg = img
			HomeStruct.IsConnecter = true
			HomeStruct.Error = 0
			PostStruct.ConnectUserName = inscriptionName
			PostStruct.ConnectUserImg = img
			NewPostStruct.ConnectUserImg = img
		} else {
			HomeStruct.Error = 3
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("Session")
	if err == nil {
		database.DeleteSession(sessionCookie.Value)
	}
	HomeStruct.Email = ""
	HomeStruct.ConnectUserInfo = ""
	HomeStruct.IsConnecter = false
	PostStruct.ConnectUserName = ""
	PostStruct.ConnectUserImg = ""
	NewPostStruct.ConnectUserImg = ""

	cookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl404.Execute(w, HomeStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	if !HomeStruct.IsConnecter {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := tmplProfil.Execute(w, HomeStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !HomeStruct.IsConnecter {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les nouvelles informations du formulaire
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Mettre à jour les informations dans la base de données
	database.UpdateUserInfo(HomeStruct.Email, username, email, HashPassword2(password))
	HomeStruct.Email = email
	HomeStruct.ConnectUserInfo = username
	// Rediriger vers la page de profil mise à jour
	http.Redirect(w, r, "/", http.StatusFound)
}

// Utilitaires *************************************************************************************************************************

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func transformDate() string {
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02 15:04:05")
	return dateString
}

func renderImg() string {
	img := "1179024"
	nb := rand.Intn(15)
	boolean := true
	for boolean {
		if (nb != 0) && (nb != 1) {
			boolean = false
		} else {
			nb = rand.Intn(15)
		}
	}
	if nb < 10 {
		img += "0" + strconv.Itoa(nb)
	} else {
		img += strconv.Itoa(nb)
	}
	return img
}

func HashPassword2(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hashed := h.Sum(nil)
	return string(hashed)
}
