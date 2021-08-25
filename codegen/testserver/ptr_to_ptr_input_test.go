package testserver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

type UpdatePtrToPtrResults struct {
	UpdatedPtrToPtr PtrToPtrOuter `json:"updatePtrToPtr"`
}

func TestPtrToPtr(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.MutationResolver.UpdatePtrToPtr = func(ctx context.Context, in UpdatePtrToPtrOuter) (ret *PtrToPtrOuter, err error) {
		ret = &PtrToPtrOuter{
			Name: "oldName",
			Inner: &PtrToPtrInner{
				Key:   "oldKey",
				Value: "oldValue",
			},
		}

		if in.Name != nil {
			ret.Name = *in.Name
		}

		if in.Inner != nil {
			inner := *in.Inner
			if inner == nil {
				ret.Inner = nil
			} else {
				if in.Inner == nil {
					ret.Inner = &PtrToPtrInner{}
				}
				if inner.Key != nil {
					ret.Inner.Key = *inner.Key
				}
				if inner.Value != nil {
					ret.Inner.Value = *inner.Value
				}
			}
		}
		return
	}

	t.Run("pointer to pointer input missing", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation { updatePtrToPtr(input: { name: "newName" }) { name, inner { key, value }}}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "newName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "oldKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "oldValue")
	})

	t.Run("pointer to pointer input non-null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation {
			updatePtrToPtr(input: {
				inner: {
					key: "newKey"
					value: "newValue"
				}
			})
			{ name, inner { key, value }}
		}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.NotNil(t, resp.UpdatedPtrToPtr.Inner)
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Key, "newKey")
		require.Equal(t, resp.UpdatedPtrToPtr.Inner.Value, "newValue")
	})

	t.Run("pointer to pointer input null", func(t *testing.T) {
		var resp UpdatePtrToPtrResults

		err := c.Post(`mutation { updatePtrToPtr(input: { inner: null }) { name, inner { key, value }}}`, &resp)
		require.NoError(t, err)

		require.Equal(t, resp.UpdatedPtrToPtr.Name, "oldName")
		require.Nil(t, resp.UpdatedPtrToPtr.Inner)
	})
}
