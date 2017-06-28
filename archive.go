package main

import (
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

func Archive(archivePath string, archiveType string, files []string) error {
	switch strings.ToLower(archiveType) {
	case "zip":
		if err := archiver.Zip.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving zip")
		}
	case "tar":
		if err := archiver.Tar.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar")
		}
	case "tgz", "tar.gz":
		if err := archiver.TarGz.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar.gz")
		}
	case "tbz2", "tar.bz2":
		if err := archiver.TarBz2.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar.bz2")
		}
	case "txz", "tar.xz":
		if err := archiver.TarXZ.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar.xz")
		}
	case "tlz4", "tar.lz4":
		if err := archiver.TarLz4.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar.lz4")
		}
	case "tsz", "tar.sz":
		if err := archiver.TarSz.Make(archivePath, files); err != nil {
			return errors.Wrap(err, "archving tar.sz")
		}
	default:
		return errors.Errorf("unknown archving format '%s'", archiveType)
	}
	return nil
}
