package modelgen

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen/out"
	"github.com/stretchr/testify/require"
)

func TestModelGeneration(t *testing.T) {
	cfg, err := config.LoadConfig("testdata/gqlgen.yml")
	require.NoError(t, err)
	require.NoError(t, cfg.Init())
	p := Plugin{
		MutateHook: mutateHook,
	}
	require.NoError(t, p.MutateConfig(cfg))

	require.True(t, cfg.Models.UserDefined("MissingTypeNotNull"))
	require.True(t, cfg.Models.UserDefined("MissingTypeNullable"))
	require.True(t, cfg.Models.UserDefined("MissingEnum"))
	require.True(t, cfg.Models.UserDefined("MissingUnion"))
	require.True(t, cfg.Models.UserDefined("MissingInterface"))
	require.True(t, cfg.Models.UserDefined("TypeWithDescription"))
	require.True(t, cfg.Models.UserDefined("EnumWithDescription"))
	require.True(t, cfg.Models.UserDefined("InterfaceWithDescription"))
	require.True(t, cfg.Models.UserDefined("UnionWithDescription"))
	require.True(t, cfg.Models.UserDefined("FieldOverrides"))

	// TODO: not sure if this is a sufficient test, should examine ./out/generated.go
	//
	require.Equal(t, "NotName", cfg.Models["FieldOverrides"].Fields["name"].FieldName)

	t.Run("no pointer pointers", func(t *testing.T) {
		generated, err := ioutil.ReadFile("./out/generated.go")
		require.NoError(t, err)
		require.NotContains(t, string(generated), "**")
	})

	t.Run("description is generated", func(t *testing.T) {
		node, err := parser.ParseFile(token.NewFileSet(), "./out/generated.go", nil, parser.ParseComments)
		require.NoError(t, err)
		for _, commentGroup := range node.Comments {
			text := commentGroup.Text()
			words := strings.Split(text, " ")
			require.True(t, len(words) > 1, "expected description %q to have more than one word", text)
		}
	})

	t.Run("tags are applied", func(t *testing.T) {
		file, err := ioutil.ReadFile("./out/generated.go")
		require.NoError(t, err)

		fileText := string(file)

		expectedTags := []string{
			`json:"missing2" database:"MissingTypeNotNullMissing2"`,
			`json:"name" database:"MissingInputName"`,
			`json:"missing2" database:"MissingTypeNullableMissing2"`,
			`json:"name" database:"TypeWithDescriptionName"`,
		}

		for _, tag := range expectedTags {
			require.True(t, strings.Contains(fileText, tag))
		}
	})

	t.Run("concrete types implement interface", func(t *testing.T) {
		var _ out.FooBarer = out.FooBarr{}
	})
}

func mutateHook(b *ModelBuild) *ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			field.Tag += ` database:"` + model.Name + field.Name + `"`
		}
	}

	return b
}
