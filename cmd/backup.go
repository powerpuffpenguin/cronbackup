package cmd

import (
	"log"
	"path/filepath"

	"github.com/powerpuffpenguin/cronbackup/core"
	"github.com/powerpuffpenguin/cronbackup/utils"
	"github.com/spf13/cobra"
)

func init() {
	var (
		basePath = utils.BasePath()
		backend, output,
		user, password, host string
		port                   uint16
		crontabs               []string
		immediate, description bool
	)
	cmd := &cobra.Command{
		Use:   `backup`,
		Short: `crontab incremental backup`,
		Run: func(cmd *cobra.Command, args []string) {
			c, e := core.New(
				core.WithOutput(output),
				core.WithServer(host, port), core.WithAuth(user, password),
				core.WithImmediate(immediate), core.WithDescription(description),
			)
			if e != nil {
				log.Fatalln(e)
				return
			}
			for _, crontab := range crontabs {
				_, e = c.Add(crontab)
				if e != nil {
					log.Fatalln(e)
				}
			}
			if len(c.Entries()) == 0 {
				if immediate {
					c.UnsafeJob()
				}
			} else {
				e = c.Serve()
				if e != nil {
					log.Fatalln(e)
				}
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&backend, `backend`, `b`,
		utils.Abs(basePath, filepath.Join(`backend`, `mariadb.js`)),
		`backend script`,
	)
	flags.StringVarP(&output, `output`, `o`,
		utils.Abs(basePath, `output`),
		`output path`,
	)
	flags.StringVarP(&user, `user`, `u`,
		`root`,
		`username for connecting to the server`,
	)
	flags.StringVarP(&password, `password`, `p`,
		``,
		`password for connecting to the server`,
	)
	flags.StringVarP(&host, `host`, `H`,
		`localhost`,
		`hostname for connecting to the server`,
	)
	flags.Uint16VarP(&port, `port`, `P`,
		3306,
		`port for connecting to the server`,
	)
	flags.StringSliceVarP(&crontabs, `contab`, `c`,
		nil,
		`minute hour dom mon dow`,
	)
	flags.BoolVarP(&immediate, `immediate`, `i`,
		false,
		`perform a backup immediately`,
	)
	flags.BoolVarP(&description, `description`, `d`,
		false,
		`generate description.json`,
	)
	rootCmd.AddCommand(cmd)
}
