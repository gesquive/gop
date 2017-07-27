# gop - go packager
[![Travis CI](https://img.shields.io/travis/gesquive/gop/master.svg?style=flat-square)](https://travis-ci.org/gesquive/gop)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/gop/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/gop)

Package your multi-os/arch executables. `gop` is intended to compliment [gox](https://github.com/mitchellh/gox) and provide the packaging aspect of a deployment.

You should be able to use the same `arch` & `os` arguments used to run `gox`. The `ouput`/`input` arguments are also similar.

`gop` supports packaging into the following archive formats: zip, tar, tgz, tar.gz, tbz2, tar.bz2, txz, tar.xz, tlz4, tar.lz4, tsz, tar.sz


## Installing

### Compile
This project requires go 1.8+ to compile. Just run `go get -u github.com/gesquive/gop` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/gop/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/gop`.

## Configuration

You can add a `.gop.yml` file to your project which `gop` will use in place of command line arguments. A [sample config](config.sample.yml) is provided in this repo. Just rename the the sample to `.gop.yml`.


### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./.gop.yml`
 - `~/.config/gop/.gop.yml`

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix "GOP_" in front of the uppercased variable name. For example, the config variable `archive` would be the environment variable `GOP_ARCHIVE`.

## Usage

```console
Package your multi-os/arch executables

Usage:
  gop [flags] [packages]

Flags:
  -a, --arch stringSlice       List of architectures to package (default [386,amd64,amd64p32,arm,arm64,ppc64,ppc64le])
  -r, --archive stringSlice    List of package types to create (default [zip,tar.gz,tar.xz])
  -c, --config string          config file (default .gop.yml)
  -d, --delete                 Delete the packaged executables
  -f, --files stringSlice      Add additional file to package
  -h, --help                   help for gop
  -i, --input string           The input path template. (default "{{.Dir}}_{{.OS}}_{{.Arch}}")
  -s, --os stringSlice         List of operating systems to package (default [darwin,dragonfly,freebsd,linux,netbsd,openbsd,plan9,solaris,windows])
  -o, --output string          The output path template. (default "{{.Dir}}_{{.OS}}_{{.Arch}}.{{.Archive}}")
  -p, --packages stringSlice   List of os/arch/archive groups to package
  -V, --version                Show the version and exit
```
Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

## Documentation

This documentation can be found at github.com/gesquive/gop

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
