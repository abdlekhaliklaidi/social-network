package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"forum/auth"
	"forum/database"
	"forum/handlers"

	"github.com/rs/cors"
)

func main() {
	//  CORS  localhost:3000
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.DB.Close()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	// fileServer := http.FileServer(http.Dir("./social-network/components"))
	// http.Handle("/social-network/components/", http.StripPrefix("/social-network/components", fileServer))

	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/show_posts", handlers.ShowPosts)
	http.HandleFunc("/post_submit", handlers.PostSubmit)
	http.HandleFunc("/comment_submit", handlers.CommentSubmit)
	http.HandleFunc("/interact", handlers.HandleInteract)
	http.HandleFunc("/get_categories", handlers.GetCategories)
	http.HandleFunc("/Connections", handlers.Connections)

	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/check-session", auth.CheckSessionHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/register", auth.RegisterHandler)

	//  CORS
	handler := c.Handler(http.DefaultServeMux)

	log.Println("Server started on :4848")
	fmt.Println("http://localhost:4848/")
	err := http.ListenAndServe(":4848", handler)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	time.Sleep(10 * time.Second)
}

// 1. فهم Next.js:
// قبل أن تبدأ في الانتقال، تأكد من فهمك لكيفية عمل Next.js. هذا يشمل:

// التوجيه الديناميكي (Dynamic Routing)
// التقديم من جهة الخادم (SSR) والتقديم الثابت (SSG)
// إمكانية التفاعل مع APIs من خلال API Routes
// ستفد من static generation أو server-side rendering حسب حاجتك، حيث أن Next.js يدعم هذا مباشرة.
// قم بترتيب الهيكلية الخاصة بالمحتوى باستخدام getStaticProps أو getServerSideProps لتحميل البيانات من API أو قاعدة بياناتك.
