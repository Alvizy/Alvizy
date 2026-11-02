package main

import (
	"log"
	"net/http"
)

func main() {
	MustLoadConfig("config.json")
	MustInitDB(cfg.Database.DSN)

	if err := RunMigrations(); err != nil {
		log.Fatalf("migrasi gagal: %v", err)
	}

	loadTemplates()

	mux := http.NewServeMux()
	admin := "/" + cfg.AdminPath

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/about", handleAbout)
	mux.HandleFunc("/projects", handleProjects)
	mux.HandleFunc("/contact", handleContact)
	mux.HandleFunc("/certificates", handleCertificates)

	// Admin auth (tidak dilindungi middleware)
	mux.HandleFunc(admin+"/login", handleAdminLogin)
	mux.HandleFunc(admin+"/logout", handleAdminLogout)

	// Admin protected
	adminMux := http.NewServeMux()
	adminMux.HandleFunc(admin, handleAdminDashboard)
	adminMux.HandleFunc(admin+"/bio", handleAdminBio)
	adminMux.HandleFunc(admin+"/skills", handleAdminSkills)
	adminMux.HandleFunc(admin+"/skills/create", handleAdminSkillCreate)
	adminMux.HandleFunc(admin+"/skills/delete/", handleAdminSkillDelete)
	adminMux.HandleFunc(admin+"/projects", handleAdminProjects)
	adminMux.HandleFunc(admin+"/projects/new", handleAdminProjectNew)
	adminMux.HandleFunc(admin+"/projects/edit/", handleAdminProjectEdit)
	adminMux.HandleFunc(admin+"/projects/delete/", handleAdminProjectDelete)

	mux.Handle(admin+"/", AdminMiddleware(adminMux))

	log.Printf("server jalan di http://localhost%s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
