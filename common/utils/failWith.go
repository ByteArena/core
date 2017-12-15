package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ttacon/chalk"
	bettererrors "github.com/xtuc/better-errors"
	bettererrorstree "github.com/xtuc/better-errors/printer/tree"
)

// TODO(sven): we should disable the colors when the terminal has no frontend
// and/or expliclty pass an --no-colors argument.
var (
	WarnColor = chalk.Red.Color
)

func FailWith(err error) {
	if bettererrors.IsBetterError(err) {

		command := strings.Join(os.Args, " ")

		berror := bettererrors.
			New(command).
			SetContext("version", GetVersion()).
			With(err)

		msg := bettererrorstree.PrintChain(berror)

		urlOptions := url.Values{}
		urlOptions.Set("body", wrapInMarkdownCode(msg))

		fmt.Println("")
		fmt.Println(WarnColor("❌  An error occurred."))
		fmt.Println("")

		fmt.Print(WarnColor(msg))

		fmt.Println("")

		fmt.Println("Please report this error here: https://github.com/ByteArena/cli/issues/new?" + urlOptions.Encode())

		os.Exit(1)
	} else {
		panic(err)
	}
}

func wrapInMarkdownCode(str string) string {
	return fmt.Sprintf("```sh\n%s\n```", str)
}

func WarnWith(err error) {
	if bettererrors.IsBetterError(err) {
		msg := bettererrorstree.PrintChain(err.(*bettererrors.Chain))

		fmt.Println("")
		fmt.Println(WarnColor("⚠️  Warning"))
		fmt.Println("")

		fmt.Print(WarnColor(msg))

		fmt.Println("")
	} else {
		fmt.Println(err.Error())
	}
}
