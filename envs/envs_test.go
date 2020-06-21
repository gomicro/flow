package envs_test

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/gomicro/flow/envs"
	. "github.com/onsi/gomega"
)

func TestEnvs(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Parsing File", func() {
		g.It("should parse envs from a file", func() {
			es, err := envs.ParseFile("./basic.env")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(es)).To(Equal(3))
			Expect(es[0].Key).To(Equal("FLOW_CONFIG1"))
			Expect(es[0].Value).To(Equal("foo"))
			Expect(es[1].Key).To(Equal("FLOW_CONFIG2"))
			Expect(es[1].Value).To(Equal("bar"))
			Expect(es[2].Key).To(Equal("FLOW_CONFIG3"))
			Expect(es[2].Value).To(Equal("baz"))
		})

		g.It("should parse only good envs from a file", func() {
			es, err := envs.ParseFile("./bads.env")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(es)).To(Equal(2))
			Expect(es[0].Key).To(Equal("FLOW_CONFIG1"))
			Expect(es[0].Value).To(Equal("foo"))
			Expect(es[1].Key).To(Equal("FLOW_CONFIG2"))
			Expect(es[1].Value).To(Equal("avaluewith=init"))
		})

		g.It("should return an error on file errors", func() {
			_, err := envs.ParseFile("./nonexistent.env")
			Expect(err).To(HaveOccurred())
		})
	})

	g.Describe("Parsing Slice", func() {
		g.It("should parse multiple envs", func() {
			slice := []string{
				"FOO=bar",
				"BAZ=biz",
			}

			es := envs.ParseSlice(slice)
			Expect(len(es)).To(Equal(2))
			Expect(es[0].Key).To(Equal("FOO"))
			Expect(es[0].Value).To(Equal("bar"))
			Expect(es[1].Key).To(Equal("BAZ"))
			Expect(es[1].Value).To(Equal("biz"))
		})

		g.It("should ignore incomplete key pairs", func() {
			slice := []string{
				"FOO=bar",
				"MISSING=",
			}

			es := envs.ParseSlice(slice)
			Expect(len(es)).To(Equal(1))
			Expect(es[0].Key).To(Equal("FOO"))
			Expect(es[0].Value).To(Equal("bar"))
		})
	})
}
