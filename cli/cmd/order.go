package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
)

/////////////////////////////////////////////////
/////////           Stop Command            /////
////////////////////////////////////////////////

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Simulation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stop\n")
		sender.Post(nil, "/order/stop")

	},
}

func init() {
	orderCmd.AddCommand(stopCmd)
}

/////////////////////////////////////////////////
/////////           Start Command            /////
////////////////////////////////////////////////

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Simulation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("start\n")
		sender.Post(nil, "/order/start")

	},
}

func init() {
	orderCmd.AddCommand(startCmd)
}

/////////////////////////////////////////////////
/////////           Set Command            /////
////////////////////////////////////////////////

type AgentOptions struct {
	Num int `validate:"required,min=0,max=10", json:"num"`
}

type ClockOptions struct {
	Time int `validate:"required,min=0" json:"time"`
}

var (
	ao = &AgentOptions{}
	co = &ClockOptions{}
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set agent or clock or area",
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Set agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("set agent\n")
		aojson, _ := json.Marshal(ao)
		sender.Post(aojson, "/order/set/agent")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateParams(*ao)
	},
}

var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Set clock",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("set clock %v\n", co.Time)
		cojson, _ := json.Marshal(co)
		sender.Post(cojson, "/order/set/clock")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateParams(*co)
	},
}

func init() {
	agentCmd.Flags().IntVarP(&ao.Num, "num", "n", 0, "agent num (required)")
	clockCmd.Flags().IntVarP(&co.Time, "time", "t", 0, "clcok time (required)")
	setCmd.AddCommand(agentCmd)
	setCmd.AddCommand(clockCmd)
	orderCmd.AddCommand(setCmd)
}

/////////////////////////////////////////////////
//////////          Order Command          /////
////////////////////////////////////////////////
var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Start a provider",
	Long: `Start a provider with options 
For example:
    simulation order start   
	simulation order set-time   
	simulation order set-area   
`,
}

func init() {
	rootCmd.AddCommand(orderCmd)
}

/////////////////////////////////////////////////
//////////            Validation            /////
////////////////////////////////////////////////
var validate = validator.New()

func validateParams(p interface{}) error {

	errs := validate.Struct(p)

	return extractValidationErrors(errs)
}

func extractValidationErrors(err error) error {

	if err != nil {
		var errorText []string
		for _, err := range err.(validator.ValidationErrors) {
			errorText = append(errorText, validationErrorToText(err))
		}
		return fmt.Errorf("Parameter error: %s", strings.Join(errorText, "\n"))
	}

	return nil
}

func validationErrorToText(e validator.FieldError) string {

	f := e.Field()
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", f)
	case "max":
		return fmt.Sprintf("%s cannot be greater than %s", f, e.Param())
	case "min":
		return fmt.Sprintf("%s must be greater than %s", f, e.Param())
	}
	return fmt.Sprintf("%s is not valid %s", e.Field(), e.Value())
}
