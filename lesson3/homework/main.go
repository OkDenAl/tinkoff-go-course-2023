package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Conversation string

const (
	UpperCase  Conversation = "upper_case"
	LowerCase  Conversation = "lower_case"
	TrimSpaces Conversation = "trim_spaces"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	BlockSize int
	Conv      string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "number of bytes to be skipped")
	flag.IntVar(&opts.Limit, "limit", -1, "max number of bytes to be read. by default - all file")
	flag.IntVar(&opts.BlockSize, "block-size", 0, "the size of one block in bytes when reading and writing")
	flag.StringVar(&opts.Conv, "conv", "", "one or more of the possible "+
		"transformations over the text, separated by commas.")
	flag.Parse()
	return &opts, nil
}

func setupReader(opts *Options) (*bufio.Reader, *os.File, error) {
	var reader *bufio.Reader
	var f *os.File
	if opts.From != "" {
		f, err := os.Open(opts.From)
		if err != nil {
			return nil, nil, fmt.Errorf("cant open input file: %v", err)
		}
		reader = bufio.NewReader(f)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	return reader, f, nil
}

func setupWriter(opts *Options) (*bufio.Writer, *os.File, error) {
	var writer *bufio.Writer
	var f *os.File
	if opts.To != "" {
		if _, err := os.Stat(opts.To); err == nil {
			return nil, nil, fmt.Errorf("output file already exist")
		}
		f, err := os.Create(opts.To)
		if err != nil {
			return nil, nil, fmt.Errorf("cant create the output file: err: %v", err)
		}
		writer = bufio.NewWriter(f)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}
	return writer, f, nil
}

func validateConv(inputConv string) error {
	if inputConv == "" {
		return nil
	}
	splited := strings.Split(inputConv, ",")
	usedConvs := make(map[Conversation]struct{})
	for _, name := range splited {
		switch Conversation(name) {
		case LowerCase:
			if _, ok := usedConvs[UpperCase]; ok {
				return fmt.Errorf("cant use upper_case and lower_case conversation together")
			}
			usedConvs[LowerCase] = struct{}{}
		case UpperCase:
			if _, ok := usedConvs[LowerCase]; ok {
				return fmt.Errorf("cant use upper_case and lower_case conversation together")
			}
			usedConvs[UpperCase] = struct{}{}
		case TrimSpaces:
		default:
			return fmt.Errorf("unknown conversation")
		}
	}
	return nil
}

func validateOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset is less then 0")
	} else {
		return nil
	}
}

func shutDown() {

}

func main() {
	logger := log.New(os.Stderr, "ERROR:\t", 3)
	opts, err := ParseFlags()
	if err != nil {
		logger.Fatal("can not parse flags: ", err)
	}
	if err = validateOffset(opts.Offset); err != nil {
		logger.Fatal("invalid offset: ", err)
	}
	if err = validateConv(opts.Conv); err != nil {
		logger.Fatal("invalid conv: ", err)
	}

	reader, input, err := setupReader(opts)
	if err != nil {
		logger.Fatal("cant setup reader: ", err)
	}
	if input != nil {
		defer input.Close()
	}

	writer, output, err := setupWriter(opts)
	if err != nil {
		logger.Fatal("cant setup writer: ", err)
	}
	if output != nil {
		defer output.Close()
	}
	defer writer.Flush()

	isOffsetUseful := false
	bytesRead := 0
	isTrimmed := false

	for bytesRead <= opts.Limit || opts.Limit == -1 {
		line, errorReader := reader.ReadBytes('\n')
		if errorReader != nil && len(line) == 0 {
			if errorReader == io.EOF {
				if !isOffsetUseful {
					logger.Fatal("offset is bigger then file is")
				}
				break
			} else {
				logger.Fatal("cant read the input file", errorReader)
			}
		}
		if !isOffsetUseful {
			if len(line) < opts.Offset {
				opts.Offset -= len(line)
				continue
			} else {
				isOffsetUseful = true
				line = line[opts.Offset:]
			}
		}
		bytesRead += len(line)
		if bytesRead > opts.Limit && opts.Limit != -1 {
			line = line[:len(line)-(bytesRead-opts.Limit)]
		}
		outputLine := string(line)
		splited := strings.Split(opts.Conv, ",")
		for _, name := range splited {
			switch Conversation(name) {
			case LowerCase:
				outputLine = strings.ToLower(outputLine)
			case UpperCase:
				outputLine = strings.ToUpper(outputLine)
			case TrimSpaces:
				if !isTrimmed {
					outputLine = strings.TrimLeft(outputLine, " \n\r\v\t")
				}
				if errorReader != nil {
					outputLine = strings.TrimRight(outputLine, " \n\r\v\t")
				}
				if len(outputLine) != 0 {
					isTrimmed = true
				}
			}
		}
		fmt.Fprint(writer, outputLine)

		if errorReader != nil {
			if errorReader == io.EOF {
				if !isOffsetUseful {
					logger.Fatal("offset is bigger then file is")
				}
				break
			} else {
				logger.Fatal("cant read the input file", errorReader)
			}
		}
	}
}
