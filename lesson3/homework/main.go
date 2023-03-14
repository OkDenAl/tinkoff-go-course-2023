package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
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
	flag.IntVar(&opts.BlockSize, "block-size", 10, "the size of one block in bytes when reading and writing")
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
	outputLine := ""
	spaces := make([]byte, 0)
	buf2 := make([]byte, 0)

	for bytesRead <= opts.Limit || opts.Limit == -1 {
		buf := make([]byte, opts.BlockSize)
		n, errorReader := reader.Read(buf)
		if n == 0 {
			if errorReader == io.EOF {
				if !isOffsetUseful {
					logger.Fatal("offset is bigger then file is")
				}
				break
			} else {
				logger.Fatal("cant read the input file", errorReader)
			}
		}
		buf = buf[:n]
		if len(buf2) != 0 {
			buf2 = append(buf2, buf...)
			buf = buf2
			buf2 = make([]byte, 0)
		}
		i := len(buf) - 1
		k := 0
		var rune1 rune
		for i >= 0 && i >= len(buf)-5 {
			rune1, k = utf8.DecodeRune(buf[i:])
			if rune1 != utf8.RuneError {
				break
			}
			i--
		}
		buf2 = append(buf2, buf[i+k:]...)
		buf = buf[:i+k]
		if !isOffsetUseful {
			if len(buf) < opts.Offset {
				opts.Offset -= len(buf)
				continue
			} else {
				isOffsetUseful = true
				buf = buf[opts.Offset:]
			}
		}
		bytesRead += len(buf)
		if bytesRead > opts.Limit && opts.Limit != -1 {
			buf = buf[:len(buf)-(bytesRead-opts.Limit)]
		}
		splited := strings.Split(opts.Conv, ",")
		convLine := string(buf)
		for _, name := range splited {
			switch Conversation(name) {
			case LowerCase:
				convLine = strings.ToLower(convLine)
			case UpperCase:
				convLine = strings.ToUpper(convLine)
			case TrimSpaces:
				if !isTrimmed {
					convLine = strings.TrimLeft(convLine, " \n\r\v\t ")
				}
				if len(convLine) != 0 {
					isTrimmed = true
					buf = buf[len(buf)-len(convLine):]
				}
				if isTrimmed {
					convLine = ""
					for _, b := range string(buf) {
						if unicode.IsSpace(b) {
							spaces = append(spaces, []byte(string(b))...)
						} else {
							convLine += string(spaces) + string(b)
							spaces = make([]byte, 0)
						}
					}
				}
			}
			outputLine = convLine
		}
		fmt.Fprint(writer, outputLine)
	}
}
