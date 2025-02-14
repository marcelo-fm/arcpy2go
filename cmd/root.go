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
	Short: "parse an arcpy tool documentation from the url provided",
	Long: `arcpy2go scrapes arcpy tool documentation from the the <url> provided
in the first argument, and create a go file with the arcpy tool. passing a <path>
as the second argument saves the go file.
Usage:
  arcpy2go <url> <path>(optional)`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		url := args[0]
		file := os.Stdout
		c := colly.NewCollector(
			colly.CacheDir(filepath.Join(viper.GetString("appConfigDir"), "cache")),
		)
		data, err := web.Parse(c, url)
		if len(args) == 2 {
			path := args[1]
			filepathname := args[1]
			if viper.GetBool("package") {
				err = os.MkdirAll(args[1], 0755)
				cobra.CheckErr(err)
				filepathname = filepath.Join(args[1], fmt.Sprintf("%s.go", data.FunctionName))
			} else {
				dirpath := filepath.Dir(path)
				err = os.MkdirAll(dirpath, 0755)
				cobra.CheckErr(err)
			}
			file, err = os.Create(filepathname)
			cobra.CheckErr(err)
		}
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
	rootCmd.Flags().BoolP("package", "p", false, "saves the file in a package in <path> with the --package-name as folder. default (.)")
	viper.BindPFlag("package", rootCmd.Flags().Lookup("package"))
	rootCmd.Flags().String("package-name", "arcpy", "name of the package of the go file generated")
	viper.BindPFlag("packageName", rootCmd.Flags().Lookup("package-name"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	homeConfigDir, err := os.UserConfigDir()
	cobra.CheckErr(err)
	appConfigDir := filepath.Join(homeConfigDir, "arcpy2go")
	err = os.MkdirAll(filepath.Join(appConfigDir, "cache"), 0755)
	cobra.CheckErr(err)
	viper.Set("appConfigDir", appConfigDir)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(appConfigDir)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}
	viper.SetEnvPrefix("arcpy2go")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
