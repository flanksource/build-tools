module github.com/flanksource/build-tools

go 1.13

require (
	github.com/flanksource/commons v1.2.0
	github.com/google/go-github/v31 v31.0.0
	github.com/joshdk/go-junit v0.0.0-20200312181801-e5d93c0f31a8
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gotest.tools v2.2.0+incompatible
)

replace gopkg.in/hairyhenderson/yaml.v2 => github.com/maxaudron/yaml v0.0.0-20190411130442-27c13492fe3c
