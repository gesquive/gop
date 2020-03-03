package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/gesquive/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug bool
var showVersion bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gop [flags] [packages]",
	Short: "Package your executables",
	Long: `Package your multi-os/arch executables

If no specific operating systems, architectures or archives are specified, gop
  will search for all known builds and package any found.

Input/Output path template:

  The input & output path for the binaries/packages is specified with the
  "--input" and "--output" flags respectively. The value is a string that
  is a Go text template. The default values are "{{.Dir}}_{{.OS}}_{{.Arch}}"
  and "{{.Dir}}_{{.OS}}_{{.Arch}}.{{.Archive}}". The variables and
  their values should be self-explanatory.

Packages (OS/Arch/Archive):

  The operating systems, architectures, and archives to package may be
  specified with the "--arch", "--os" & "--archive" flags. These are space
  separated lists of values to build for, respectively. You may prefix an
  OS, Arch or Archive with "!" to negate and not package for that value.
  If the list is made up of only negations, then the negations will come from
  the default list.

  Additionally, the "--packages" flag may be used to specify complete
  os/arch/archive values that should be built or ignored. The syntax for
  this is what you would expect: "linux/amd64/zip" would be a valid package
  value. Multiple values can be space separated. An os/arch/archive definition
  can begin with "!" to not build for that platform.

  The "--packages" flag has the highest precedent when determing whether to
  build for a platform. If it is included in the "--packages" list, it will be
  built even if the specific os, arch or archive is negated in  the "--os",
  "--arch" and "--archive" flags respectively.

`,
	PersistentPreRun: preRun,
	Run:              run,
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/gop %s\n",
		RootCmd.HelpTemplate(), buildVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringP("config", "c", "",
		"config file (default .gop.yml)")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().StringP("input", "i", "{{.Dir}}_{{.OS}}_{{.Arch}}",
		"The input path template.")
	RootCmd.PersistentFlags().StringP("output", "o", "{{.Dir}}_{{.OS}}_{{.Arch}}.{{.Archive}}",
		"The output path template.")

	RootCmd.PersistentFlags().StringSliceP("files", "f", []string{},
		"Add additional file to package")
	RootCmd.PersistentFlags().StringSliceP("archive", "r", DefaultArchiveList,
		"List of package types to create")
	RootCmd.PersistentFlags().StringSliceP("os", "s", OSList,
		"List of operating systems to package")
	RootCmd.PersistentFlags().StringSliceP("arch", "a", ArchList,
		"List of architectures to package")
	RootCmd.PersistentFlags().StringSliceP("packages", "p", []string{},
		"List of os/arch/archive groups to package")
	RootCmd.PersistentFlags().BoolP("delete", "d", false,
		"Delete the packaged executables")

	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("gop")
	viper.AutomaticEnv()
	viper.BindEnv("config")
	viper.BindEnv("input")
	viper.BindEnv("output")
	viper.BindEnv("files")
	viper.BindEnv("archive")
	viper.BindEnv("os")
	viper.BindEnv("arch")
	viper.BindEnv("packages")
	viper.BindEnv("delete")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("input", RootCmd.PersistentFlags().Lookup("input"))
	viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("files", RootCmd.PersistentFlags().Lookup("files"))
	viper.BindPFlag("archive", RootCmd.PersistentFlags().Lookup("archive"))
	viper.BindPFlag("os", RootCmd.PersistentFlags().Lookup("os"))
	viper.BindPFlag("arch", RootCmd.PersistentFlags().Lookup("arch"))
	viper.BindPFlag("packages", RootCmd.PersistentFlags().Lookup("packages"))
	viper.BindPFlag("delete", RootCmd.PersistentFlags().Lookup("delete"))

	viper.SetDefault("input", "{{.Dir}}_{{.OS}}_{{.Arch}}")
	viper.SetDefault("output", "{{.Dir}}_{{.OS}}_{{.Arch}}.{{.Archive}}")
	viper.SetDefault("archive", DefaultArchiveList)
	viper.SetDefault("os", OSList)
	viper.SetDefault("arch", ArchList)
	viper.SetDefault("delete", false)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".gop")               // name of config file (without extension)
		viper.AddConfigPath(".")                  // adding current directory as first search path
		viper.AddConfigPath("$HOME/.config/.gop") // adding home directory as next search path
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		if !showVersion {
			fmt.Println("Error opening config: ", err)
		}
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/gop\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if debug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	cli.Debug("Running with debug turned on")
	cli.Debug("config: %s", viper.ConfigFileUsed())
}

func run(cmd *cobra.Command, args []string) {
	srcPackages := args
	if len(srcPackages) < 1 {
		srcPackages = []string{"."}
	}
	cli.Debug("cfg: packages=%v", srcPackages)

	inputTemplate := viper.GetString("input")
	cli.Debug("cfg: input=%s", inputTemplate)

	outputTemplate := viper.GetString("output")
	cli.Debug("cfg: output=%s", outputTemplate)

	fileList := viper.GetStringSlice("files")
	cli.Debug("cfg: files=%v", fileList)

	archList := viper.GetStringSlice("arch")
	cli.Debug("cfg: arch=%v", archList)

	osList := viper.GetStringSlice("os")
	cli.Debug("cfg: os=%v", osList)

	archiveList := viper.GetStringSlice("archive")
	cli.Debug("cfg: archive=%v", archiveList)

	// Get the packages that are in the given paths
	appDirs, err := GetAppDirs(srcPackages)
	if err != nil {
		cli.Fatal("error getting app dirs: %s", err)
	}

	userPackages := viper.GetStringSlice("packages")
	packages, err := AssemblePackageInfo(archList, osList, archiveList, userPackages)
	if err != nil {
		cli.Fatal("error getting package list: %s", err)
	}
	cli.Debug("packages found: %s", packages)

	packages, err = GetPackagePaths(packages, appDirs, inputTemplate, outputTemplate)
	if err != nil {
		cli.Fatal("error getting package paths: %s", err)
	}

	packages, err = GetPackageFiles(packages, fileList)
	if err != nil {
		cli.Fatal("error getting package files: %s", err)
	}

	cli.Info("Packaging archives:")

	for _, pkg := range packages {
		if _, err := os.Stat(pkg.ExePath); os.IsNotExist(err) {
			cli.Debug("xxx %60s", pkg.ArchivePath)
			continue
		}
		cli.Info("--> %60s", pkg.ArchivePath)
		err = archive(pkg.ArchivePath, pkg.Archive, pkg.FileList)
		if err != nil {
			cli.Error("error: %s", err)
		}
	}

	cli.Debug("cfg: delete=%t", viper.GetBool("delete"))
	if viper.GetBool("delete") {
		cli.Info("Cleaning up executables")
		for _, pkg := range packages {
			os.Remove(pkg.ExePath)
			dir := path.Dir(pkg.ExePath)
			if isEmpty, _ := IsEmpty(dir); isEmpty {
				os.Remove(dir)
			}
		}
	}
}
