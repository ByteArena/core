package assert

import (
	"github.com/bytearena/core/common/utils"
	bettererrors "github.com/xtuc/better-errors"
)

func Assert(cond bool, msg string) {

	if !cond {
		berror := bettererrors.
			New("Assertion error").
			With(bettererrors.New(msg))

		utils.FailWith(berror)
	}
}
func AssertBE(cond bool, err *bettererrors.Chain) {

	if !cond {
		berror := bettererrors.
			New("Assertion error").
			With(err)

		utils.FailWith(berror)
	}
}
