package push

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/store"

	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/push"
	"github.com/solo-io/wasme/pkg/resolver"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	ref        string
	storageDir string

	*opts.GeneralOptions
}

func PushCmd(ctx *context.Context, generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts pushOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "push name[:tag|@digest]",
		Short: "Push a wasm filter to remote registry",
		Long: `Push wasm filter to remote registry. E.g.:

wasme push webassemblyhub.io/my/filter:v1
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ref = args[0]
			return runPush(*ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")

	return cmd
}

func runPush(ctx context.Context, opts pushOptions) error {
	logrus.Infof("Pushing image %v", opts.ref)

	image, err := store.NewStore(opts.storageDir).Get(opts.ref)
	if err != nil {
		return errors.Wrap(err, "image not found. run `wasme list` to see locally cahced images")
	}

	resolver, authorizer := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.Configs...)
	pusher := push.NewPusher(resolver, authorizer)
	if err := pusher.Push(ctx, image); err != nil {
		return err
	}

	return nil
}
