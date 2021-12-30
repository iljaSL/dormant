package cmd

import (
	"fmt"

	"github.com/iljaSL/dormant/lib"
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
	dormant inspect <GO.MOD>
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
	fmt.Println("TEST FLAG", inactivityDuration)

	deps, err := lib.ReadFile(args[0])
	if err != nil {
		return err
	}

	activityInfo, _ := lib.GetAPILastActivityInfo(deps)

	depsActivityList := lib.CalculateDepsActivity(activityInfo)
	fmt.Println(depsActivityList)

	// Todo: Next Step using gitHubs and gitLabs REST API to get the status
	// ! Before doing any CURL action, check if inactivityDuration is not 0
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
