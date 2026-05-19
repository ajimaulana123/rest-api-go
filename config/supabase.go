package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

type Supabase struct {
	ProjectID  string
	URL        string
	AnonKey    string
	DBPassword string
	Region     string
	Pooler     string // aws-0 | aws-1 (project baru biasanya aws-1)
	DBMode     string // session (default) | transaction | direct
}

func loadSupabase() Supabase {
	region := strings.TrimSpace(os.Getenv("SUPABASE_REGION"))
	if region == "" {
		region = "ap-southeast-1"
	}

	mode := strings.TrimSpace(os.Getenv("SUPABASE_DB_MODE"))
	if mode == "" {
		mode = "session"
	}

	pooler := strings.TrimSpace(os.Getenv("SUPABASE_POOLER"))
	if pooler == "" {
		pooler = "aws-1"
	}

	sb := Supabase{
		ProjectID:  strings.TrimSpace(os.Getenv("SUPABASE_PROJECT_ID")),
		URL:        strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/"),
		AnonKey:    strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY")),
		DBPassword: os.Getenv("SUPABASE_DB_PASSWORD"),
		Region:     region,
		Pooler:     pooler,
		DBMode:     mode,
	}

	if sb.URL == "" || sb.AnonKey == "" {
		log.Fatal("SUPABASE_URL dan SUPABASE_ANON_KEY wajib diisi (Project Settings → API)")
	}

	return sb
}

func (s Supabase) ProjectRef() string {
	if s.ProjectID != "" {
		return s.ProjectID
	}

	u, err := url.Parse(s.URL)
	if err != nil {
		log.Fatalf("SUPABASE_URL tidak valid: %v", err)
	}

	host := u.Hostname()
	if strings.HasSuffix(host, ".supabase.co") {
		return strings.TrimSuffix(host, ".supabase.co")
	}

	log.Fatal("SUPABASE_PROJECT_ID wajib diisi jika URL bukan format *.supabase.co")
	return ""
}

func (s Supabase) DatabaseURL() string {
	if override := os.Getenv("DATABASE_URL"); override != "" {
		return override
	}

	if s.DBPassword == "" {
		log.Fatal("SUPABASE_DB_PASSWORD wajib diisi (password database, bukan anon key). Atau set DATABASE_URL langsung.")
	}

	ref := s.ProjectRef()
	poolerHost := fmt.Sprintf("%s-%s.pooler.supabase.com", s.Pooler, s.Region)

	switch s.DBMode {
	case "direct":
		user := url.UserPassword("postgres", s.DBPassword)
		return fmt.Sprintf("postgresql://%s@db.%s.supabase.co:5432/postgres", user.String(), ref)
	case "transaction":
		user := url.UserPassword("postgres."+ref, s.DBPassword)
		return fmt.Sprintf("postgresql://%s@%s:6543/postgres", user.String(), poolerHost)
	default: // session pooler — IPv4, cocok untuk Windows
		user := url.UserPassword("postgres."+ref, s.DBPassword)
		return fmt.Sprintf("postgresql://%s@%s:5432/postgres", user.String(), poolerHost)
	}
}
