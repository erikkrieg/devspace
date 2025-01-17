package commands

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/loft-sh/devspace/pkg/devspace/pipeline/engine/types"
	"mvdan.cc/sh/v3/interp"
)

var errXArgsUsage = errors.New(`usage: xargs [utility [argument ...]]`)

type XArgsOptions struct {
	Delimiter string
}

func XArgs(ctx context.Context, args []string, handler types.ExecHandler) error {
	options := &XArgsOptions{
		Delimiter: " ",
	}

	args, err := parseXArgsOptions(args, options)
	if err != nil {
		return err
	} else if len(args) == 0 {
		return errXArgsUsage
	}

	hc := interp.HandlerCtx(ctx)
	out, err := ioutil.ReadAll(hc.Stdin)
	if err != nil {
		return err
	}

	addArgs := strings.Split(string(out), options.Delimiter)
	for _, addArg := range addArgs {
		addArg = strings.TrimSpace(addArg)
		if addArg == "" {
			continue
		}

		args = append(args, addArg)
	}
	return handler.ExecHandler(ctx, args)
}

func parseXArgsOptions(args []string, options *XArgsOptions) ([]string, error) {
	// check args for flags
	startAt := 0
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 0 && arg[0] == '-' {
			startAt++
			arg = arg[1:]

			switch arg {
			case "d", "-delimiter":
				if i+1 == len(args) {
					return nil, errXArgsUsage
				}

				i++
				startAt++
				options.Delimiter = args[i]
			default:
				return nil, errXArgsUsage
			}

			continue
		}

		break
	}

	return args[startAt:], nil
}
