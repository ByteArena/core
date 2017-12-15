package dockerfile

import (
	"io"
	"strconv"
	"strings"

	dockerfileparser "github.com/docker/docker/builder/dockerfile/parser"
)

type Instruction string
type SourceLoc int
type Whitelist map[Instruction]interface{}

func (s Instruction) String() string {
	return string(s)
}

func (s SourceLoc) String() string {
	return strconv.Itoa(int(s))
}

var dockerfileInstructionWhitelist = Whitelist{
	"COPY":       nil,
	"FROM":       nil,
	"WORKDIR":    nil,
	"RUN":        nil,
	"CMD":        nil,
	"ENTRYPOINT": nil,
	"ENV":        nil,
}

func DockerfileParserGetFroms(source io.Reader) ([]string, error) {
	result, err := dockerfileparser.Parse(source)
	if err != nil {
		return nil, err
	}

	fromValues := make([]string, 0)
	visitRow(result.AST, func(node *dockerfileparser.Node) {
		if node.Value == "from" {
			fromValues = append(fromValues, node.Next.Value)
		}
	})

	return fromValues, nil
}

func DockerfileFindForbiddenInstructions(source io.Reader) (map[Instruction]interface{}, error) {
	res := make(map[Instruction]interface{})

	result, err := dockerfileparser.Parse(source)
	if err != nil {
		return nil, err
	}

	visitRow(result.AST, func(node *dockerfileparser.Node) {
		instruction := Instruction(strings.ToUpper(node.Value))
		_, isWhitelisted := dockerfileInstructionWhitelist[instruction]

		if !isWhitelisted {
			res[instruction] = nil
		}
	})

	return res, nil
}

func visitRow(node *dockerfileparser.Node, cbk func(n *dockerfileparser.Node)) {
	for _, n := range node.Children {
		cbk(n)
	}
}
