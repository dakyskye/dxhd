package parser

import (
	"bufio"
	"errors"
	"os"

	"github.com/dakyskye/dxhd/logger"
)

// Parser holds everything our parser needs to work.
type Parser struct {
	fileName    string
	reader      *bufio.Reader
	res         parser
	finalResult ParseResult
	finished    bool
}

type parser struct {
	lineNumber int
	line       string
	isPrefix   bool
	parseRes   ParseResult
}

// ParseResult is where our parser packs the parsed data.
type ParseResult struct {
	Shell string
}

var (
	// ErrParsingNotFinished is returned collection of parser results is requested before it finishes the job.
	ErrParsingNotFinished = errors.New("the parser has not finished the job yet")
	// ErrParsingHasFinished is returned when re-parse is requested, instead Collect should be requested.
	ErrParsingHasFinished = errors.New("the parser has already finished the job")
	// ErrParserHasNoReader is likely to be returned when New is not used to initialise a parser.
	ErrParserHasNoReader = errors.New("the parser does not have any data to parse")
)

// New returns a new parser for a given file.
func New(fileName string) (Parser, error) {
	logger.L().WithField("file", fileName).Debug("made a new parser")

	file, err := os.Open(fileName)
	if err != nil {
		return Parser{}, err
	}

	p := Parser{
		fileName: fileName,
		reader:   bufio.NewReader(file),
	}
	err = file.Close()

	return p, err
}

// Collect returns the parsed result.
func (p *Parser) Collect() (ParseResult, error) {
	if !p.finished {
		return ParseResult{}, ErrParsingNotFinished
	}
	return p.finalResult, nil
}
