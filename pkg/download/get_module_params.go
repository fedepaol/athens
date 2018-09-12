package download

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/paths"
)

func getModuleParams(c buffalo.Context, op errors.Op) (mod string, vers string, err error) {
	params, err := paths.GetAllParams(c)
	if err != nil {
		return "", "", errors.E(op, err, errors.KindBadRequest)
	}

	return params.Module, params.Version, nil
}
