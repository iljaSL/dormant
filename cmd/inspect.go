package cmd

import (
	"fmt"
	"time"

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
	// ! inactivityDuration Need to replace hardcoded months
	fmt.Println("TEST FLAG", inactivityDuration)

	deps, err := lib.ReadFile(args[0])
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
	d := pterm.TableData{{"Dependencies", "Status", "Size"}}

	// ! Check if activityDuration is not 0!!!
	// ! REMOVE HARDCODED MONTHS
	for _, v := range depsLastActivityList {
		p.UpdateTitle("Analyzing " + v.URL) // Update the title of the progressbar.
		// pterm.Success.Println("Analyzing " + v.URL)
		if v.Month <= 6 && v.Month >= 4 {
			d = append(d, []string{pterm.LightYellow(v.URL), pterm.LightYellow("Sporadic")})
		} else if v.Month <= 3 {
			d = append(d, []string{pterm.LightGreen(v.URL), pterm.LightGreen("Active")})
		} else {
			d = append(d, []string{pterm.LightRed(v.URL), pterm.LightRed("Inactive")})
		}
		p.Increment() // Increment the progressbar by one. Use Add(x int) to increment by a custom amount.
		time.Sleep(time.Millisecond * 250)
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
