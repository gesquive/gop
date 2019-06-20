package main

import (
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

func archive(archivePath string, archiveType string, files []string) error {
	switch strings.ToLower(archiveType) {
	case "zip":
		zip := archiver.NewZip()
		if err := zip.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving zip")
		}
	case "tar":
		tar := archiver.NewTar()
		if err := tar.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar")
		}
	case "tbz2", "tar.bz2":
		tarbz2 := archiver.NewTarBz2()
		if err := tarbz2.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar.bz2")
		}
	case "tgz", "tar.gz":
		targz := archiver.NewTarGz()
		if err := targz.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar.gz")
		}
	case "tlz4", "tar.lz4":
		tarlz4 := archiver.NewTarLz4()
		if err := tarlz4.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar.lz4")
		}
	case "tsz", "tar.sz":
		tarsz := archiver.NewTarSz()
		if err := tarsz.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar.sz")
		}
	case "txz", "tar.xz":
		tarxz := archiver.NewTarXz()
		if err := tarxz.Archive(files, archivePath); err != nil {
			return errors.Wrap(err, "archving tar.xz")
		}
	default:
		return errors.Errorf("unknown archving format '%s'", archiveType)
	}
	return nil
}
