package cmd

import (
	"github.com/iljaSL/dormant/lib"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// type inspectOptions struct {
// 	inspectOptions string
// }

func init() {
	rootCmd.AddCommand(inspectModFileCmd())
}

//
//
func inspectModFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a go.mod file for inactive packages.",
		Long:  `Inspect a go.mod file for inactive packages.`,
		Example: `
	dormant inspect go.mod
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectModFile(args)
		},
	}

	// cmd.AddCommand(aNewCommand())

	return cmd
}

func inspectModFile(args []string) error {
	deps, err := lib.ReadGoModFile(args[0])
	if err != nil {
		return err
	}

	activityInfo, err := lib.GetAPILastActivityInfo(deps)
	if err != nil {
		return err
	}

	depsLastActivityList, err := lib.CalculateDepsActivity(activityInfo)
	if err != nil {
		return err
	}

	p, _ := pterm.DefaultProgressbar.WithTotal(len(depsLastActivityList)).Start()
	d := pterm.TableData{{"Dependencies", "Status", "Last Activity (months)"}}

	for _, v := range depsLastActivityList {
		p.UpdateTitle("Analyzing " + v.URL) // Update the title of the progressbar.
		// pterm.Success.Println("Analyzing " + v.URL)
		if v.Month <= inactivityDuration && v.Month >= sporadicDuration {
			d = append(d, []string{pterm.LightYellow(v.URL), pterm.LightYellow("Sporadic"), pterm.LightYellow(v.Month)})
		} else if v.Month < sporadicDuration && !(sporadicDuration > inactivityDuration) {
			d = append(d, []string{pterm.LightGreen(v.URL), pterm.LightGreen("Active"), pterm.LightGreen(v.Month)})
		} else {
			d = append(d, []string{pterm.LightRed(v.URL), pterm.LightRed("Inactive"), pterm.LightRed(v.Month)})
		}
		p.Increment() // Increment the progress bar by one. Use Add(x int) to increment by a custom amount.
	}
	pterm.DefaultTable.WithHasHeader().WithData(d).Render()

	return err
}

// func stdout(data interface{}) error {
// 	encoded, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = os.Stdout.Write(encoded)
// 	return err
// }

//
//
// func NewCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "test",
// 		Short: "test",
// 		Long:  `test`,
// 		Example: `
// 	test
// 		`,
// 		Args:         cobra.ExactArgs(1),
// 		SilenceUsage: true,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return test(args)
// 		},
// 	}

// 	return cmd
// }
