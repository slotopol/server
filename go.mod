module github.com/slotopol/server

go 1.26.0

require (
	github.com/bmatcuk/doublestar/v4 v4.10.3
	github.com/gin-gonic/gin v1.10.1
	github.com/go-sql-driver/mysql v1.10.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/klauspost/cpuid/v2 v2.4.1
	github.com/lib/pq v1.12.3
	github.com/mattn/go-isatty v0.0.24
	github.com/mattn/go-sqlite3 v1.14.52
	github.com/schwarzlichtbezirk/go-disk-usage v1.3.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	go.uber.org/atomic v1.12.0
	go.yaml.in/yaml/v3 v3.0.5
	golang.org/x/sync v0.24.0
	xorm.io/xorm v1.4.3
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.4 // indirect
	github.com/bytedance/sonic/loader v0.5.2 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/gin-contrib/sse v1.1.2 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/go-playground/validator/v10 v10.30.5 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.11.2 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.14 // indirect
	golang.org/x/arch v0.32.0 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	xorm.io/builder v0.3.13 // indirect
)

// Advice: add to file `.vscode/settings.json` followed lines:
// {
//     "gopls": {
//         "buildFlags": ["-tags", "full keno agt aristocrat betsoft ct igt megajack netent novomatic playngo playtech"],
//     },
//     "Lua.diagnostics.disable": ["lowercase-global"]
// }

// Advice: upgrade this file by command line
//   go get -u ./...
//   go mod tidy
