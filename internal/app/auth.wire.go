package app

import "pharmacy-modernization-project-model/internal/platform/auth"

func (a *App) wireAuth() error {
	if err := auth.NewBuilder().
		WithJWTConfig(
			a.Cfg.Auth.JWT.Secret,
			a.Cfg.Auth.JWT.Issuer,
			a.Cfg.Auth.JWT.Audience,
			a.Cfg.Auth.JWT.ClientIds,
			a.Cfg.Auth.JWT.Cookie.Name,
		).
		WithDevMode(a.Cfg.Auth.DevMode).
		WithEnvironment(a.Cfg.App.Env).
		WithLogger(a.Logger.Base).
		Build(); err != nil {
		return err
	}

	return nil
}
