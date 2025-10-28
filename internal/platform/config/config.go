package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Env  string `mapstructure:"env"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`
	Logging struct {
		Enabled        bool   `mapstructure:"enabled"`
		Level          string `mapstructure:"level"`
		Format         string `mapstructure:"format"`
		Output         string `mapstructure:"output"`           // "console", "file", or "both"
		FilePath       string `mapstructure:"file_path"`        // Path to log file
		FileMaxSize    int    `mapstructure:"file_max_size"`    // Max size in MB before rotation
		FileMaxBackups int    `mapstructure:"file_max_backups"` // Max number of old log files
		FileMaxAge     int    `mapstructure:"file_max_age"`     // Max days to retain old log files
	} `mapstructure:"logging"`
	Auth struct {
		DevMode bool `mapstructure:"dev_mode"`
		JWT     struct {
			Secret   string   `mapstructure:"secret"`
			Issuer   []string `mapstructure:"issuer"`
			Audience []string `mapstructure:"audience"`
			Cookie   struct {
				Name     string `mapstructure:"name"`
				Secure   bool   `mapstructure:"secure"`
				HTTPOnly bool   `mapstructure:"httponly"`
				MaxAge   int    `mapstructure:"max_age"`
			} `mapstructure:"cookie"`
		} `mapstructure:"jwt"`
	} `mapstructure:"auth"`
	Database struct {
		MongoDB struct {
			URI         string `mapstructure:"uri"`
			Database    string `mapstructure:"database"`
			Collections struct {
				Patients      string `mapstructure:"patients"`
				Addresses     string `mapstructure:"addresses"`
				Prescriptions string `mapstructure:"prescriptions"`
			} `mapstructure:"collections"`
			Connection struct {
				MaxPoolSize    uint64 `mapstructure:"max_pool_size"`
				MinPoolSize    uint64 `mapstructure:"min_pool_size"`
				MaxIdleTime    string `mapstructure:"max_idle_time"`
				ConnectTimeout string `mapstructure:"connect_timeout"`
				SocketTimeout  string `mapstructure:"socket_timeout"`
			} `mapstructure:"connection"`
			Options struct {
				RetryWrites bool `mapstructure:"retry_writes"`
				RetryReads  bool `mapstructure:"retry_reads"`
			} `mapstructure:"options"`
		} `mapstructure:"mongodb"`
	} `mapstructure:"database"`
	External struct {
		Stargate struct {
			UseMock      bool              `mapstructure:"use_mock"`
			Timeout      string            `mapstructure:"timeout"`
			ClientID     string            `mapstructure:"client_id"`
			ClientSecret string            `mapstructure:"client_secret"`
			Scope        string            `mapstructure:"scope"`
			Endpoints    StargateEndpoints `mapstructure:"endpoints"`
		} `mapstructure:"stargate"`
		Pharmacy struct {
			UseMock   bool              `mapstructure:"use_mock"`
			Timeout   string            `mapstructure:"timeout"`
			Endpoints PharmacyEndpoints `mapstructure:"endpoints"`
		} `mapstructure:"pharmacy"`
		Billing struct {
			UseMock   bool             `mapstructure:"use_mock"`
			Timeout   string           `mapstructure:"timeout"`
			Endpoints BillingEndpoints `mapstructure:"endpoints"`
		} `mapstructure:"billing"`
	} `mapstructure:"external"`
	Cache CacheConfig `mapstructure:"cache"`
}

// StargateEndpoints holds the full URLs for Stargate authentication endpoints
type StargateEndpoints struct {
	Token        string `mapstructure:"token"`
	RefreshToken string `mapstructure:"refresh_token"`
}

// PharmacyEndpoints holds the full URLs for pharmacy API endpoints
type PharmacyEndpoints struct {
	GetPrescription string `mapstructure:"get_prescription"`
}

// BillingEndpoints holds the full URLs for billing API endpoints
type BillingEndpoints struct {
	GetInvoice           string `mapstructure:"get_invoice"`
	GetInvoicesByPatient string `mapstructure:"get_invoices_by_patient"`
	CreateInvoice        string `mapstructure:"create_invoice"`
	AcknowledgeInvoice   string `mapstructure:"acknowledge_invoice"`
	GetInvoicePayment    string `mapstructure:"get_invoice_payment"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	MongoDB CacheMongoDBConfig `mapstructure:"mongodb"`
	Memory  MemoryCacheConfig  `mapstructure:"memory"`
}

// CacheMongoDBConfig holds MongoDB cache configuration
type CacheMongoDBConfig struct {
	URI        string `mapstructure:"uri"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
	Connection struct {
		MaxPoolSize    uint64 `mapstructure:"max_pool_size"`
		MinPoolSize    uint64 `mapstructure:"min_pool_size"`
		MaxIdleTime    string `mapstructure:"max_idle_time"`
		ConnectTimeout string `mapstructure:"connect_timeout"`
		SocketTimeout  string `mapstructure:"socket_timeout"`
	} `mapstructure:"connection"`
}

// MemoryCacheConfig holds in-memory cache configuration
type MemoryCacheConfig struct {
	MaxCost     int64  `mapstructure:"max_cost"`
	BufferItems int64  `mapstructure:"buffer_items"`
	Metrics     bool   `mapstructure:"metrics"`
	DefaultTTL  string `mapstructure:"default_ttl"`
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/configs")
	_ = v.ReadInConfig()

	// optional env-specific file: RX_APP_ENV=prod -> app.prod.yaml
	if env := os.Getenv("RX_APP_ENV"); env != "" {
		v2 := viper.New()
		v2.SetConfigName("app." + env)
		v2.AddConfigPath("./internal/configs")
		v2.SetConfigType("yaml")
		if err := v2.ReadInConfig(); err == nil {
			v.MergeConfigMap(v2.AllSettings())
		}
	}

	v.SetEnvPrefix("RX")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{}
	_ = v.Unmarshal(cfg)
	if cfg.App.Port == 0 {
		cfg.App.Port = 8080
	}
	if cfg.App.Name == "" {
		cfg.App.Name = "PharmacyModernization"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "dev"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "debug"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "console"
	}
	if cfg.Logging.Output == "" {
		cfg.Logging.Output = "console"
	}
	if cfg.Logging.FilePath == "" {
		cfg.Logging.FilePath = "logs/app.log"
	}
	if cfg.Logging.FileMaxSize == 0 {
		cfg.Logging.FileMaxSize = 100 // 100MB default
	}
	if cfg.Logging.FileMaxBackups == 0 {
		cfg.Logging.FileMaxBackups = 3
	}
	if cfg.Logging.FileMaxAge == 0 {
		cfg.Logging.FileMaxAge = 28 // 28 days default
	}
	// Logging enabled by default if not specified
	if !v.IsSet("logging.enabled") {
		cfg.Logging.Enabled = true
	}
	// Auth defaults
	// JWT Secret is REQUIRED via RX_AUTH_JWT_SECRET environment variable
	// No default provided for security reasons
	if len(cfg.Auth.JWT.Issuer) == 0 {
		cfg.Auth.JWT.Issuer = []string{"PharmacyModernization"}
	}
	if len(cfg.Auth.JWT.Audience) == 0 {
		cfg.Auth.JWT.Audience = []string{"PharmacyModernization"}
	}
	if cfg.Auth.JWT.Cookie.Name == "" {
		cfg.Auth.JWT.Cookie.Name = "auth_token"
	}
	if cfg.Auth.JWT.Cookie.MaxAge == 0 {
		cfg.Auth.JWT.Cookie.MaxAge = 3600 // 1 hour
	}
	cfg.Auth.JWT.Cookie.HTTPOnly = true // Always true for security
	return cfg
}
