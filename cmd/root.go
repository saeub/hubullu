package cmd

import (
	"fmt"
	"os"

	"github.com/saeub/hubullu/internal/ui"
	"github.com/saeub/hubullu/internal/vocabulary"
	"github.com/saeub/hubullu/pkg/translation"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var vocabFile string

var rootCmd = &cobra.Command{
	Use:   "hubullu",
	Short: "A CLI for looking up words",
	Args:  cobra.ExactArgs(2), // TODO: custom validator
	Run: func(cmd *cobra.Command, args []string) {
		var vocab *vocabulary.Vocabulary
		if vocabFile == "" {
			vocab = vocabulary.NewVocabulary()
		} else {
			var err error
			vocab, err = vocabulary.LoadVocabulary(vocabFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		srcLang, err := language.Parse(args[0])
		if err != nil {
			fmt.Printf("invalid language tag: %s\n", srcLang)
			os.Exit(1)
		}
		trgLang, err := language.Parse(args[1])
		if err != nil {
			fmt.Printf("invalid language tag: %s\n", err)
			os.Exit(1)
		}
		sk, err := translation.NewSanakirjaOrg(srcLang, trgLang)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		trltrs := []translation.Translator{sk}
		ui := ui.NewPrompt(vocab, trltrs)
		ui.Run()
	},
}

// Execute initializes and runs the root command.
// If an error happens, this writes it to standard output and exits with status 1.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.hubullu.yaml)")
	rootCmd.Flags().StringVarP(&vocabFile, "vocab", "v", "",
		"file to load/save vocabulary")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".hubullu")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initVocab() {
}
