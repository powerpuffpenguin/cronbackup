package cmd

import (
	"log"

	"github.com/powerpuffpenguin/cronbackup/core"
	"github.com/powerpuffpenguin/cronbackup/utils"
	"github.com/spf13/cobra"
)

func init() {
	var (
		basePath = utils.BasePath()
		output   string
	)
	cmd := &cobra.Command{
		Use:   `description`,
		Short: `generate description.json`,
		Run: func(cmd *cobra.Command, args []string) {
			c, e := core.New(
				core.WithOutput(output),
			)
			if e != nil {
				log.Fatalln(e)
				return
			}
			e = c.GenerateDescription()
			if e != nil {
				log.Fatalln(e)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&output, `output`, `o`,
		utils.Abs(basePath, `output`),
		`output path`,
	)
	rootCmd.AddCommand(cmd)
}
