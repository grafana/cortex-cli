module github.com/grafana/cortextool

go 1.12

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

require (
	cloud.google.com/go v0.35.0
	github.com/alecthomas/chroma v0.7.0
	github.com/alecthomas/repr v0.0.0-20181024024818-d37bc2a10ba1 // indirect
	github.com/cortexproject/cortex v0.1.1-0.20190808112445-606262b7a637
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/hashicorp/memberlist v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/alertmanager v0.13.0
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/common v0.7.0
	github.com/prometheus/prometheus v0.0.0-20190731144842-63ed2e28f1ac
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	google.golang.org/api v0.4.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
