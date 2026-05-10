package check

import "nhatp.com/go/sugar"

func New() sugar.Plugin {
	return &pluginImpl{}
}

const PluginID = "nhatp.com/go/sugar/sugars/check"

type pluginImpl struct {
}

func (i *pluginImpl) ID() string {
	return PluginID
}

func (i *pluginImpl) LexicalParser() sugar.LexicalParser {
	return LexicalParser()
}

var _ sugar.Plugin = (*pluginImpl)(nil)
