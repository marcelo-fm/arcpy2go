/*
Copyright © 2024 marcelo-fm

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
	"github.com/marcelo-fm/arcpy2go/arcpy-scraper/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arcpy2go",
	Short: "",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		url := args[0]
		file := os.Stdout
		if len(args) == 2 {
			path := args[1]
			dirpath := filepath.Dir(path)
			err = os.MkdirAll(dirpath, 0755)
			cobra.CheckErr(err)
			file, err = os.Create(args[1])
			cobra.CheckErr(err)
		}
		c := colly.NewCollector(
			// colly.AllowedDomains("https://pro.arcgis.com"),
			colly.CacheDir(filepath.Join(viper.GetString("appConfigDir"), "cache")),
		)
		data, err := web.Parse(c, url)
		cobra.CheckErr(err)
		err = data.Render(file)
		cobra.CheckErr(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find homeConfigDir directory.
		homeConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".arcpy2go" (without extension).
		appConfigDir := filepath.Join(homeConfigDir, "arcpy2go")
		err = os.MkdirAll(filepath.Join(appConfigDir, "cache"), 0755)
		cobra.CheckErr(err)
		viper.Set("appConfigDir", appConfigDir)
		viper.AddConfigPath(appConfigDir)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
