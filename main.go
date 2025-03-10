package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/auth"
	"social-network/database"
	"sync"
	"time"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.DB.Close()

	fileServer := http.FileServer(http.Dir("./app"))
	http.Handle("/app/", http.StripPrefix("/app", fileServer))

	// http.HandleFunc("/", handlers.HomePage)
	// http.HandleFunc("/show_posts", handlers.ShowPosts)
	// http.HandleFunc("/post_submit", handlers.PostSubmit)
	// http.HandleFunc("/comment_submit", handlers.CommentSubmit)
	// http.HandleFunc("/interact", handlers.HandleInteract)
	// http.HandleFunc("/get_categories", handlers.GetCategories)
	// http.HandleFunc("/Connections", handlers.Connections)

	http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/check-session", auth.CheckSessionHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/register", auth.RegisterHandler)

	log.Println("Server started on :4848")
	fmt.Println("http://localhost:4848/")
	err := http.ListenAndServe(":4848", nil)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	time.Sleep(10 * time.Second)
}
