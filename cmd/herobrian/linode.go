package main

import (
	"github.com/bdreece/herobrian/pkg/linode"
	"go.uber.org/config"
)

func createLinodeClient(provider config.Provider) (*linode.Client, error) {
    opts := new(linode.Options)
    if err := provider.Get("linode").Populate(opts); err != nil {
        return nil, err
    }

    return linode.NewClient(opts), nil
}
