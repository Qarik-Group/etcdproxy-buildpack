package supply

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
}

type Manifest interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/manifest.go
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
}

type Installer interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/installer.go
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/command.go
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Supplier struct {
	Manifest  Manifest
	Installer Installer
	Stager    Stager
	Command   Command
	Log       *libbuildpack.Logger
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying etcd")

	etcd, err := s.Manifest.DefaultVersion("etcd")
	if err != nil {
		return err
	}
	s.Log.Info("Using etcd version %s", etcd.Version)

	if err := s.Installer.InstallDependency(etcd, s.Stager.DepDir()); err != nil {
		return err
	}

	// Find unpacked "etcd-v3.3.13-linux-amd64/etcd" file
	etcdFile, err := filepath.Glob(filepath.Join(s.Stager.DepDir(), "etcd*linux-amd64", "etcd"))
	if err != nil {
		return err
	}

	// Rename and move into bin/etcd, which will be in $PATH
	err = os.Rename(etcdFile[0], filepath.Join(s.Stager.DepDir(), "bin", "etcd"))
	if err != nil {
		return err
	}

	// Find unpacked "etcd-v3.3.13-linux-amd64/etcdctl" file
	etcdctlFile, err := filepath.Glob(filepath.Join(s.Stager.DepDir(), "etcd*linux-amd64", "etcdctl"))
	if err != nil {
		return err
	}

	// Rename and move into bin/etcdctl, which will be in $PATH
	err = os.Rename(etcdctlFile[0], filepath.Join(s.Stager.DepDir(), "bin", "etcdctl"))
	if err != nil {
		return err
	}

	return nil
}
