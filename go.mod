module go.thethings.network/lorawan-stack

go 1.14

// Use our fork of grpc-gateway.
replace github.com/grpc-ecosystem/grpc-gateway => github.com/TheThingsIndustries/grpc-gateway v1.14.3-gogo

// Use our fork of otto.
replace github.com/robertkrimen/otto => github.com/TheThingsIndustries/otto v0.0.0-20181129100957-6ddbbb60554a

// github.com/blang/semver doesn't have a v3 semantic import.
replace github.com/blang/semver => github.com/blang/semver v0.0.0-20190414182527-1a9109f8c4a1

// Dependency of Goreleaser that causes problems with module management.
// See https://github.com/Azure/go-autorest/issues/414.
replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible

// Do not upgrade Echo beyond v4.1.2.
// See https://github.com/TheThingsNetwork/lorawan-stack/issues/977.
replace github.com/labstack/echo/v4 => github.com/labstack/echo/v4 v4.1.2

// Do not upgrade go-sqlmock beyond v1.3.0.
// See https://github.com/heptiolabs/healthcheck/issues/23.
replace gopkg.in/DATA-DOG/go-sqlmock.v1 => gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0

// Dependency of Hugo that causes problems with module management.
replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

// Dependency of Hugo that causes problems with module management.
replace github.com/nicksnyder/go-i18n => github.com/nicksnyder/go-i18n v1.10.0

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1
	github.com/TheThingsIndustries/magepkg v0.0.0-20190214092847-6c0299b7c3ed
	github.com/TheThingsIndustries/mystique v0.0.0-20200127144137-4aa959111fe7
	github.com/TheThingsNetwork/go-cayenne-lib v1.0.0
	github.com/aws/aws-sdk-go v1.29.26
	github.com/blang/semver v0.0.0-00010101000000-000000000000
	github.com/chrj/smtpd v0.1.2
	github.com/client9/misspell v0.3.4
	github.com/disintegration/imaging v1.6.2
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/envoyproxy/protoc-gen-validate v0.3.0-java
	github.com/fsnotify/fsnotify v1.4.9
	github.com/getsentry/sentry-go v0.5.1
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/gogo/protobuf v1.3.1
	github.com/gohugoio/hugo v0.67.1
	github.com/golang/gddo v0.0.0-20200310004957-95ce5a452273
	github.com/golang/protobuf v1.3.5
	github.com/goreleaser/goreleaser v0.129.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/gotnospirit/makeplural v0.0.0-20180622080156-a5f48d94d976 // indirect
	github.com/gotnospirit/messageformat v0.0.0-20190719172517-c1d0bdacdea2
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/jacobsa/crypto v0.0.0-20190317225127-9f44e2d11115
	github.com/jacobsa/oglematchers v0.0.0-20150720000706-141901ea67cd // indirect
	github.com/jacobsa/oglemock v0.0.0-20150831005832-e94d794d06ff // indirect
	github.com/jacobsa/ogletest v0.0.0-20170503003838-80d50a735a11 // indirect
	github.com/jacobsa/reqtrace v0.0.0-20150505043853-245c9e0234cb // indirect
	github.com/jarcoal/httpmock v1.0.5
	github.com/jaytaylor/html2text v0.0.0-20200220170450-61d9dc4d7195
	github.com/jinzhu/gorm v1.9.12
	github.com/kr/pretty v0.2.0
	github.com/labstack/echo/v4 v4.1.11
	github.com/labstack/gommon v0.3.0
	github.com/lib/pq v1.3.0
	github.com/magefile/mage v1.9.0
	github.com/mattn/go-isatty v0.0.12
	github.com/mattn/goveralls v0.0.5
	github.com/mdempsky/unconvert v0.0.0-20200228143138-95ecdbfc0b5f
	github.com/mgechev/revive v1.0.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nats-io/nats-server/v2 v2.1.4
	github.com/nats-io/nats.go v1.9.1
	github.com/oklog/ulid/v2 v2.0.2
	github.com/openshift/osin v1.0.1
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/robertkrimen/otto v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/rest v2.4.1+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.5.0+incompatible
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	github.com/smartystreets/assertions v1.0.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	go.opencensus.io v0.22.3
	go.packetbroker.org/api/v2 v2.0.0
	go.thethings.network/lorawan-stack-legacy v0.0.0-20190118141410-68812c833a78
	gocloud.dev v0.19.0
	gocloud.dev/pubsub/natspubsub v0.19.0
	golang.org/x/crypto v0.0.0-20200317142112-1b76d66859c6
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/tools v0.0.0-20200318132943-540150da7399
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/api v0.20.0
	google.golang.org/genproto v0.0.0-20200323114720-3f67cca34472
	google.golang.org/grpc v1.28.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.0.0-00010101000000-000000000000 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.2.8
)
