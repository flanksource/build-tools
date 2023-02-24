module github.com/flanksource/build-tools

go 1.13

require (
	github.com/aws/aws-sdk-go v1.29.25
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/flanksource/commons v1.2.0
	github.com/google/go-github/v32 v32.0.0
	github.com/joshdk/go-junit v0.0.0-20200312181801-e5d93c0f31a8
	github.com/palantir/stacktrace v0.0.0-20161112013806-78658fd2d177
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
)

replace gopkg.in/hairyhenderson/yaml.v2 => github.com/maxaudron/yaml v0.0.0-20190411130442-27c13492fe3c
