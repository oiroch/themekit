package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Shopify/themekit/src/env"
	"github.com/Shopify/themekit/src/shopify"
)

func TestGet(t *testing.T) {
	ctx, client, conf, _, _ := createTestCtx()
	conf.On("Set", "development", env.Env{}).Return(nil, nil)
	conf.On("Save").Return(nil)
	client.On("GetAllAssets").Return([]string{}, nil)
	assert.Error(t, getTheme(ctx), "No files to download")

	ctx, _, conf, _, _ = createTestCtx()
	conf.On("Set", "development", env.Env{}).Return(nil, fmt.Errorf("invalid conf"))
	conf.On("Save").Return(nil)
	err := getTheme(ctx)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "invalid conf")
	}

	ctx, _, conf, _, _ = createTestCtx()
	conf.On("Set", "development", env.Env{}).Return(nil, nil)
	conf.On("Save").Return(fmt.Errorf("no file"))
	err = getTheme(ctx)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "no file")
	}

	ctx, client, conf, _, _ = createTestCtx()
	ctx.Flags.Environments = []string{"test"}
	ctx.Env.Domain = "my.domain.com"
	conf.On("Set", "test", env.Env{Domain: "my.domain.com"}).Return(nil, nil)
	conf.On("Save").Return(nil)
	client.On("GetAllAssets").Return([]string{}, nil)
	assert.Error(t, getTheme(ctx), "No files to download")

	ctx, client, conf, _, _ = createTestCtx()
	ctx.Flags.List = true
	client.On("Themes").Return([]shopify.Theme{}, nil)
	assert.EqualError(t, getTheme(ctx), errNoThemes.Error())

	ctx, client, conf, _, _ = createTestCtx()
	ctx.Flags.List = true
	client.On("Themes").Return([]shopify.Theme{}, fmt.Errorf("server error"))
	assert.EqualError(t, getTheme(ctx), "server error")

	ctx, client, conf, stdOut, _ := createTestCtx()
	ctx.Flags.List = true
	client.On("Themes").Return([]shopify.Theme{{ID: 1234, Role: "main", Name: "test"}}, nil)
	assert.Nil(t, getTheme(ctx))
	assert.Contains(t, stdOut.String(), "[1234][live] test")
}
