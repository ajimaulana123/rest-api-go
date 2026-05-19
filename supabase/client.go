package supabase

import (
	"be/config"

	supago "github.com/supabase-community/supabase-go"
)

func NewClient(cfg config.Supabase) (*supago.Client, error) {
	return supago.NewClient(cfg.URL, cfg.AnonKey, &supago.ClientOptions{})
}
