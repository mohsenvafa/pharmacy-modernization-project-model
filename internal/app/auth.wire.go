package app

import "pharmacy-modernization-project-model/internal/platform/auth"

func (a *App) wireAuth() error {
	builder := auth.NewBuilder().
		WithJWTConfig(
			a.Cfg.Auth.JWT.Issuer,
			a.Cfg.Auth.JWT.Audience,
			a.Cfg.Auth.JWT.ClientIds,
			a.Cfg.Auth.JWT.Cookie.Name,
		).
		WithDevMode(a.Cfg.Auth.DevMode).
		WithEnvironment(a.Cfg.App.Env).
		WithLogger(a.Logger.Base)

	// Add JWKS configuration
	builder = builder.WithJWKSConfig(
		a.Cfg.Auth.JWT.JWKSURL,
		a.Cfg.Auth.JWT.JWKSCache,
		a.Cfg.Auth.JWT.SigningMethods,
	)

	if err := builder.Build(); err != nil {
		return err
	}

	return nil
}
