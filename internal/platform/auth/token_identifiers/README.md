# Token Identifiers

This folder contains the polymorphic token identifier implementations for different JWT token types.

## Structure

- `interfaces.go` - Defines the `TokenTypeIdentifier` interface that all token identifiers must implement
- `auth_pass_token.go` - Implementation for auth_pass/custom token format
- `azure_b2c_token.go` - Implementation for Azure B2C tokens

## Adding New Token Types

To add a new token type (e.g., Google OAuth):

1. **Create a new file** (e.g., `google_token.go`)
2. **Implement the interface**:
   ```go
   type GoogleTokenIdentifier struct {
       jwksURL   string
       jwksKeyfunc keyfunc.Keyfunc
       config    types.JWTConfig
   }
   
   func (gti *GoogleTokenIdentifier) DetectTokenType(ctx context.Context, tokenString string) (types.TokenType, error) {
       // Implementation
   }
   
   func (gti *GoogleTokenIdentifier) IsValidToken(ctx context.Context, tokenString string) (types.TokenType, error) {
       // Implementation
   }
   
   func (gti *GoogleTokenIdentifier) GetTokenType() types.TokenType {
       return types.TokenTypeGoogle
   }
   
   func (gti *GoogleTokenIdentifier) GetJWKSURL() string {
       return gti.jwksURL
   }
   
   func (gti *GoogleTokenIdentifier) ExtractUser(token *jwt.Token) (*types.User, error) {
       // Implementation
   }
   ```

3. **Add the token type constant** in `types/types.go`:
   ```go
   const (
       TokenTypeAuthPass  TokenType = "auth_pass"
       TokenTypeAzureB2C TokenType = "azure_b2c"
       TokenTypeGoogle   TokenType = "google"  // NEW
   )
   ```

4. **Register the identifier** in `jwt.go`:
   ```go
   case types.TokenTypeGoogle:
       identifier, err := token_identifiers.NewGoogleTokenIdentifier(jwksURL, config)
       if err != nil {
           return fmt.Errorf("failed to create Google token identifier: %w", err)
       }
       tokenManager.RegisterIdentifier(identifier)
   ```

5. **Update configuration**:
   ```yaml
   jwks_urls:
     auth_pass: "https://your-auth-provider.com/.well-known/jwks.json"
     azure_b2c: "https://yourtenant.b2clogin.com/yourtenant.onmicrosoft.com/B2C_1_signupsignin/discovery/v2.0/keys"
     google: "https://www.googleapis.com/oauth2/v3/certs"
   token_types: ["auth_pass", "azure_b2c", "google"]
   ```

## Benefits

- **Separation of Concerns**: Each token type has its own file and logic
- **Easy to Extend**: Adding new token types doesn't require modifying existing code
- **Polymorphic Design**: All token types implement the same interface
- **Independent JWKS URLs**: Each token type can have its own JWKS endpoint
- **Clean Organization**: Related functionality is grouped together
