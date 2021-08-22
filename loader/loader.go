package loader

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
)

type Loader struct {
	sysLanguageProvider common.SysLanguageProvider
	providers           []LoaderProvider
}

func New(sysLanguageProvider common.SysLanguageProvider) *Loader {
	result := &Loader{
		sysLanguageProvider: sysLanguageProvider,
		providers:           make([]LoaderProvider, 0),
	}

	return result
}

func (l *Loader) AddProvider(provider LoaderProvider) {
	l.providers = append(l.providers, provider)
}

func (l *Loader) Load(path string) (io.ReadSeekCloser, error) {
	if len(path) == 0 {
		return nil, errors.New("blank path provided")
	}

	path = strings.ReplaceAll(path, "\\", "/")

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	path = strings.ReplaceAll(path, "{LANG}", l.sysLanguageProvider.GetLanguageCode())

	for providerIdx := range l.providers {
		if !l.providers[providerIdx].Exists(path) {
			continue
		}

		return l.providers[providerIdx].Load(path)
	}

	return nil, fmt.Errorf("file not found: \"%s\"", path)
}
