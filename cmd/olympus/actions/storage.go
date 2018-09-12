package actions

import (
	"fmt"

	"github.com/gomods/athens/pkg/config"
	"github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/storage"
	"github.com/gomods/athens/pkg/storage/fs"
	"github.com/gomods/athens/pkg/storage/mem"
	"github.com/gomods/athens/pkg/storage/mongo"
	"github.com/spf13/afero"
)

// GetStorage returns storage.Backend implementation
func GetStorage(storageType string, storageConfig *config.StorageConfig) (storage.Backend, error) {
	const op errors.Op = "actions.GetStorage"
	switch storageType {
	case "memory":
		return mem.NewStorage()
	case "disk":
		if storageConfig.Disk == nil {
			return nil, errors.E(op, "Invalid Disk Storage Configuration")
		}
		rootLocation := storageConfig.Disk.RootPath
		s, err := fs.NewStorage(rootLocation, afero.NewOsFs())
		if err != nil {
			return nil, fmt.Errorf("could not create new storage from os fs (%s)", err)
		}
		return s, nil
	case "mongo":
		if storageConfig.Mongo == nil {
			return nil, errors.E(op, "Invalid Mongo Storage Configuration")
		}
		return mongo.NewStorage(storageConfig.Mongo)
	default:
		return nil, fmt.Errorf("storage type %s is unknown", storageType)
	}
}
