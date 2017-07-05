package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gesquive/cli"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var displayVersion string

var logDebug bool
var showVersion bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "gack [flags] [packages]",
	Short:            "Package your executables",
	Long:             `Package your multi-os/arch executables`,
	PersistentPreRun: preRun,
	Run:              run,
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute(version string) {
	displayVersion = version
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/%s\n",
		RootCmd.HelpTemplate(), displayVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"config file (default .gack.yml)")
	RootCmd.PersistentFlags().BoolVarP(&logDebug, "debug", "D", false,
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

	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.BindPFlag("input", RootCmd.PersistentFlags().Lookup("input"))
	viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("files", RootCmd.PersistentFlags().Lookup("files"))
	viper.BindPFlag("archive", RootCmd.PersistentFlags().Lookup("archive"))
	viper.BindPFlag("os", RootCmd.PersistentFlags().Lookup("os"))
	viper.BindPFlag("arch", RootCmd.PersistentFlags().Lookup("arch"))

	viper.SetDefault("input", "{{.Dir}}_{{.OS}}_{{.Arch}}")
	viper.SetDefault("output", "{{.Dir}}_{{.OS}}_{{.Arch}}.{{.Archive}}")
	viper.SetDefault("archive", DefaultArchiveList)
	viper.SetDefault("os", OSList)
	viper.SetDefault("arch", ArchList)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		homeConfig := path.Join(home, ".config/gack")

		viper.SetConfigName(".gack")    // name of config file (without extension)
		viper.AddConfigPath(".")        // adding current directory as first search path
		viper.AddConfigPath(homeConfig) // adding home directory as next search path

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cli.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if logDebug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
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

	packages, err := GetUserDefinedPackages(archList, osList, archiveList)
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
		cli.Info("--> %60s", pkg.ArchivePath)
		err = Archive(pkg.ArchivePath, pkg.Archive, pkg.FileList)
		if err != nil {
			cli.Error("error: %s", err)
		}
	}

type TemplateData struct {
	Dir     string
	OS      string
	Arch    string
	Archive string
}
