package infrastructure_test

import (
	"io/ioutil"
	"os"

	. "github.com/Tinker-Ware/digital-ocean-service/infrastructure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("infrastructure", func() {

	Describe("Read a YAML", func() {
		Context("Create a YAML on the fly", func() {

			yaml := `---
clientID: asdfg
clientSecret: aoihcou
redirectURI: http://localhost/oauth
port: 1000
scopes:
  - read
  - write`

			err := ioutil.WriteFile("conf.yaml", []byte(yaml), 0644)

			It("Should Return a Configuration struct with the correct values", func() {
				Ω(err).Should(BeNil())

				conf, err := GetConfiguration("conf.yaml")
				Ω(err).Should(BeNil())

				Ω(conf.Port).Should(Equal("1000"))
				Ω(conf.ClientID).Should(Equal("asdfg"))
				Ω(conf.ClientSecret).Should(Equal("aoihcou"))
				Ω(conf.RedirectURI).Should(Equal("http://localhost/oauth"))

				Ω(conf.Scopes).Should(ContainElement("read"))
				Ω(conf.Scopes).Should(ContainElement("write"))

				os.Remove("conf.yaml")

			})
		})
	})

})
