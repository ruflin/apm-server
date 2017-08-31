package beater

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elastic/apm-server/tests"
	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher/pipeline"
	"github.com/elastic/go-ucfg/yaml"
)

var pip *pipeline.Pipeline
var pub *publisher

func SetupServer(b *testing.B) *http.ServeMux {
	bt, err := instance.NewBeat("apm-server-test", "")
	if err != nil {
		b.Fatal(err)
	}

	config, err := yaml.NewConfig([]byte(`
output:
  console:
    pretty: false`))
	if err != nil {
		b.Fatal(err)
	}

	err = config.Unpack(&bt.Config)
	if err != nil {
		b.Fatalf("error unpacking config data: %v", err)
	}

	bt.RawConfig = (*common.Config)(config)
	bt.Beat.Config = &bt.Config.BeatConfig
	bt.Beat.BeatConfig, err = bt.BeatConfig()
	if err != nil {
		b.Fatal(err)
	}
	if pip == nil {
		var err error
		pip, err = pipeline.Load(bt.Info, bt.Config.Pipeline, bt.Config.Output)
		if err != nil {
			b.Fatalf("error initializing publisher: %v", err)
		}
	}
	// bt.Publisher = pipeline

	if pub == nil {
		var err error
		pub, err = newPublisher(pip, 20)

		if err != nil {
			b.Fatal(err)
		}
		defer pub.Stop()
	}

	return newMuxer(defaultConfig, pub.Send)
}

func BenchmarkEndToEndLargePayload(b *testing.B) {
	mux := SetupServer(b)
	data, _ := tests.LoadValidData("transaction")

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewReader(data))
	if err != nil {

	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}

	// beatt, err := beater.New(&bt.Beat, bt.Beat.BeatConfig)
	// fmt.Printf("config: %v", bt.Config.Output)
	// beatt.Run(bt)
	// bt.Config.Output = common.ConfigNamespace{
	// 	name: "console",
	// 	common.Config
	// }
	// assert.IsNil(b, err)

	// processor := NewProcessor()
	// for i := 0; i < b.N; i++ {
	// 	data, _ := tests.LoadValidData("error")
	// 	err := processor.Validate(data)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	processor.Transform(data)
	// }
}
