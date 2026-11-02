package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// --- Models ---

type Profile struct {
	ID        int
	Name      string
	Tagline   string
	Bio       string
	Location  string
	Email     string
	Github    string
	Instagram string
	Linkedin  string
	Available bool
	UpdatedAt time.Time
}

type Skill struct {
	ID       int
	Name     string
	Category string
	Order    int
}

type Project struct {
	ID          int
	Title       string
	Description string
	TechStack   string
	GithubURL   string
	LiveURL     string
	Featured    bool
	Order       int
	CreatedAt   time.Time
}

// --- Public handlers ---

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "index.html", map[string]any{
		"Profile": getProfile(),
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	render(w, "about.html", map[string]any{
		"Profile": getProfile(),
		"Skills":  getSkills(),
	})
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	render(w, "projects.html", map[string]any{
		"Projects": getAllProjects(),
		"Profile":  getProfile(),
	})
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	render(w, "contact.html", map[string]any{
		"Profile": getProfile(),
	})
}

func handleCertificates(w http.ResponseWriter, r *http.Request) {
	render(w, "certificates.html", map[string]any{
		"Profile": getProfile(),
	})
}

// --- Admin handlers ---

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if isLoggedIn(r) {
			http.Redirect(w, r, "/"+cfg.AdminPath, http.StatusSeeOther)
			return
		}
		render(w, "admin_login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != cfg.Auth.Username || !checkPassword(cfg.Auth.PasswordHash, password) {
		render(w, "admin_login.html", map[string]any{
			"Error": "Username atau password salah.",
		})
		return
	}

	token, err := createSession(username)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(sessionDuration),
	})
	http.Redirect(w, r, "/"+cfg.AdminPath, http.StatusSeeOther)
}

func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		deleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/"+cfg.AdminPath+"/login", http.StatusSeeOther)
}

func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	projectCount := 0
	skillCount := 0
	db.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&projectCount)
	db.QueryRow(`SELECT COUNT(*) FROM skills`).Scan(&skillCount)

	render(w, "admin_dashboard.html", map[string]any{
		"ProjectCount": projectCount,
		"SkillCount":   skillCount,
		"AdminPath":    "/" + cfg.AdminPath,
	})
}

func handleAdminBio(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "admin_bio.html", map[string]any{
			"Profile":   getProfile(),
			"AdminPath": "/" + cfg.AdminPath,
		})
		return
	}

	_, err := db.Exec(`
		UPDATE profile SET
			name=$1, tagline=$2, bio=$3, location=$4,
			email=$5, github=$6, instagram=$7, linkedin=$8,
			available=$9, updated_at=NOW()
		WHERE id=1`,
		strings.TrimSpace(r.FormValue("name")),
		strings.TrimSpace(r.FormValue("tagline")),
		strings.TrimSpace(r.FormValue("bio")),
		strings.TrimSpace(r.FormValue("location")),
		strings.TrimSpace(r.FormValue("email")),
		strings.TrimSpace(r.FormValue("github")),
		strings.TrimSpace(r.FormValue("instagram")),
		strings.TrimSpace(r.FormValue("linkedin")),
		r.FormValue("available") == "on",
	)
	if err != nil {
		log.Printf("update profile: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+cfg.AdminPath+"/bio", http.StatusSeeOther)
}

func handleAdminSkills(w http.ResponseWriter, r *http.Request) {
	render(w, "admin_skills.html", map[string]any{
		"Skills":    getSkills(),
		"AdminPath": "/" + cfg.AdminPath,
	})
}

func handleAdminSkillCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	category := strings.TrimSpace(r.FormValue("category"))
	if name == "" {
		http.Redirect(w, r, "/"+cfg.AdminPath+"/skills", http.StatusSeeOther)
		return
	}
	db.Exec(`INSERT INTO skills (name, category) VALUES ($1, $2)`, name, category)
	http.Redirect(w, r, "/"+cfg.AdminPath+"/skills", http.StatusSeeOther)
}

func handleAdminSkillDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/"+cfg.AdminPath+"/skills/delete/")
	db.Exec(`DELETE FROM skills WHERE id=$1`, id)
	http.Redirect(w, r, "/"+cfg.AdminPath+"/skills", http.StatusSeeOther)
}

func handleAdminProjects(w http.ResponseWriter, r *http.Request) {
	render(w, "admin_projects.html", map[string]any{
		"Projects":  getAllProjects(),
		"AdminPath": "/" + cfg.AdminPath,
	})
}

func handleAdminProjectNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "admin_project_form.html", map[string]any{
			"Action":    "Tambah Project",
			"Project":   Project{},
			"AdminPath": "/" + cfg.AdminPath,
		})
		return
	}
	p := projectFromForm(r)
	_, err := db.Exec(`
		INSERT INTO projects (title, description, tech_stack, github_url, live_url, featured)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		p.Title, p.Description, p.TechStack, p.GithubURL, p.LiveURL, p.Featured,
	)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+cfg.AdminPath+"/projects", http.StatusSeeOther)
}

func handleAdminProjectEdit(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/"+cfg.AdminPath+"/projects/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		p := getProjectByID(id)
		if p == nil {
			http.NotFound(w, r)
			return
		}
		render(w, "admin_project_form.html", map[string]any{
			"Action":    "Edit Project",
			"Project":   p,
			"AdminPath": "/" + cfg.AdminPath,
		})
		return
	}

	p := projectFromForm(r)
	_, err = db.Exec(`
		UPDATE projects SET title=$1, description=$2, tech_stack=$3,
		github_url=$4, live_url=$5, featured=$6 WHERE id=$7`,
		p.Title, p.Description, p.TechStack, p.GithubURL, p.LiveURL, p.Featured, id,
	)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+cfg.AdminPath+"/projects", http.StatusSeeOther)
}

func handleAdminProjectDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/"+cfg.AdminPath+"/projects/delete/")
	db.Exec(`DELETE FROM projects WHERE id=$1`, idStr)
	http.Redirect(w, r, "/"+cfg.AdminPath+"/projects", http.StatusSeeOther)
}

// --- DB helpers ---

func getProfile() Profile {
	var p Profile
	err := db.QueryRow(`SELECT id,name,tagline,bio,location,email,github,instagram,linkedin,available,updated_at FROM profile WHERE id=1`).
		Scan(&p.ID, &p.Name, &p.Tagline, &p.Bio, &p.Location, &p.Email, &p.Github, &p.Instagram, &p.Linkedin, &p.Available, &p.UpdatedAt)
	if err != nil {
		log.Printf("getProfile: %v", err)
	}
	return p
}

func getSkills() []Skill {
	rows, err := db.Query(`SELECT id, name, category, order_index FROM skills ORDER BY category, order_index, name`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var list []Skill
	for rows.Next() {
		var s Skill
		rows.Scan(&s.ID, &s.Name, &s.Category, &s.Order)
		list = append(list, s)
	}
	return list
}

func getFeaturedProjects() []Project {
	rows, err := db.Query(`SELECT id,title,description,tech_stack,github_url,live_url,featured,order_index,created_at FROM projects WHERE featured=true ORDER BY order_index, created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanProjects(rows)
}

func getAllProjects() []Project {
	rows, err := db.Query(`SELECT id,title,description,tech_stack,github_url,live_url,featured,order_index,created_at FROM projects ORDER BY order_index, created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanProjects(rows)
}

func getProjectByID(id int) *Project {
	var p Project
	err := db.QueryRow(`SELECT id,title,description,tech_stack,github_url,live_url,featured,order_index,created_at FROM projects WHERE id=$1`, id).
		Scan(&p.ID, &p.Title, &p.Description, &p.TechStack, &p.GithubURL, &p.LiveURL, &p.Featured, &p.Order, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	return &p
}

func scanProjects(rows *sql.Rows) []Project {
	var list []Project
	for rows.Next() {
		var p Project
		rows.Scan(&p.ID, &p.Title, &p.Description, &p.TechStack, &p.GithubURL, &p.LiveURL, &p.Featured, &p.Order, &p.CreatedAt)
		list = append(list, p)
	}
	return list
}

func projectFromForm(r *http.Request) Project {
	return Project{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Description: strings.TrimSpace(r.FormValue("description")),
		TechStack:   strings.TrimSpace(r.FormValue("tech_stack")),
		GithubURL:   strings.TrimSpace(r.FormValue("github_url")),
		LiveURL:     strings.TrimSpace(r.FormValue("live_url")),
		Featured:    r.FormValue("featured") == "on",
	}
}
