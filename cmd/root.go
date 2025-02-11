package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string 
var rootCMD = &cobra.Command{
	Use: "core-api",
	Short: "this api news go",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(startCMD, nil)
	},
}

func Execute() {
	cobra.CheckErr(rootCMD.Execute())
}

func init(){
	cobra.OnInitialize(initConfig)

	rootCMD.PersistentFlags().StringVar(&cfgFile, "config","","config file (default is .env)")
	rootCMD.Flags().BoolP("toogle", "t", false, "help message for toogle")
}

func initConfig(){
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(`.env`)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(os.Stdout, "using config file:", viper.ConfigFileUsed())
	}
}