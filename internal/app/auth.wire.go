package app

import "pharmacy-modernization-project-model/internal/platform/auth"

func (a *App) wireAuth() error {
	builder := auth.NewBuilder().
		WithJWTConfig(a.Cfg.Auth.JWT.Cookie.Name).
		WithDevMode(a.Cfg.Auth.DevMode).
		WithEnvironment(a.Cfg.App.Env).
		WithLogger(a.Logger.Base)

	// Convert string token types config to TokenType map
	tokenTypesConfig := make(map[auth.TokenType]auth.TokenTypeConfig)
	for tokenTypeStr, config := range a.Cfg.Auth.JWT.TokenTypesConfig {
		tokenTypesConfig[auth.TokenType(tokenTypeStr)] = auth.TokenTypeConfig{
			JWKSURL:        config.JWKSURL,
			SigningMethods: config.SigningMethods,
			Issuer:         config.Issuer,
			Audience:       config.Audience,
			ClientIds:      config.ClientIds,
		}
	}

	// Add token types configuration
	builder = builder.WithTokenTypesConfig(
		tokenTypesConfig,
		a.Cfg.Auth.JWT.JWKSCache,
	)

	// Convert string token types to TokenType enum
	var tokenTypes []auth.TokenType
	for _, tokenTypeStr := range a.Cfg.Auth.JWT.TokenTypes {
		tokenTypes = append(tokenTypes, auth.TokenType(tokenTypeStr))
	}

	// Add token types if configured
	if len(tokenTypes) > 0 {
		builder = builder.WithTokenTypes(tokenTypes)
	}

	if err := builder.Build(); err != nil {
		return err
	}

	return nil
}
