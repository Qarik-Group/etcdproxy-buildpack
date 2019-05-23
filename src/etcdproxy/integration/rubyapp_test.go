package integration_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rubyapp Integration Test", func() {
	var app *cutlass.App
	BeforeEach(func() {
		Expect(ApiHasSidecar()).To(BeTrue())
	})
	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	It("app deploys", func() {
		app = cutlass.New(filepath.Join(bpDir, "fixtures", "rubyapp"))
		app.Manifest = filepath.Join(bpDir, "fixtures", "rubyapp", "manifest.cfdev.yml")
		V3PushAppAndConfirm(app)
		Expect(app.GetBody("/")).To(ContainSubstring("Hi, I'm an app with an etcd grpc-proxy sidecar!"))
	})
})
