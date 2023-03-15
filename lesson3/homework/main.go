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

type offsetReader struct {
	offset         int
	IsOffsetUseful bool
}

func (r *offsetReader) Read(buf []byte) (n int, err error) {
	if len(buf) < r.offset {
		r.offset -= len(buf)
		return 0, fmt.Errorf("offset is bigger then buf")
	} else {
		r.IsOffsetUseful = true
		return r.offset, nil
	}
}

type limitReader struct {
	Limit      int
	BytesRead  int
	NeedToStop bool
}

func (r *limitReader) Read(buf []byte) (n int, err error) {
	r.BytesRead += len(buf)
	if r.BytesRead > r.Limit && r.Limit != -1 {
		r.NeedToStop = true
		return len(buf) - (r.BytesRead - r.Limit), nil
	}
	return len(buf), nil
}

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

func fixTruncBytes(inputBuf, storageBuf []byte, n int) ([]byte, []byte) {
	inputBuf = inputBuf[:n]
	if len(storageBuf) != 0 {
		storageBuf = append(storageBuf, inputBuf...)
		inputBuf = storageBuf
		storageBuf = make([]byte, 0)
	}
	i := len(inputBuf) - 1
	runeWidth := 0
	var rune1 rune
	for i >= 0 && i >= len(inputBuf)-5 {
		rune1, runeWidth = utf8.DecodeRune(inputBuf[i:])
		if rune1 != utf8.RuneError {
			break
		}
		i--
	}
	storageBuf = append(storageBuf, inputBuf[i+runeWidth:]...)
	inputBuf = inputBuf[:i+runeWidth]
	return inputBuf, storageBuf
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
	conversations := strings.Split(opts.Conv, ",")

	reader, input, err := setupReader(opts)
	if input != nil {
		defer input.Close()
	}
	if err != nil {
		logger.Fatal("cant setup reader: ", err)
	}

	writer, output, err := setupWriter(opts)
	if output != nil {
		defer output.Close()
	}
	if err != nil {
		logger.Fatal("cant setup writer: ", err)
	}

	offset := offsetReader{offset: opts.Offset, IsOffsetUseful: false}
	limit := limitReader{BytesRead: 0, Limit: opts.Limit, NeedToStop: false}

	isTrimmed := false
	spaces := make([]byte, 0)
	storageBuf := make([]byte, 0)

	for !limit.NeedToStop {
		inputBuf := make([]byte, opts.BlockSize)
		n, errorReader := reader.Read(inputBuf)
		if n == 0 {
			if errorReader == io.EOF {
				if !offset.IsOffsetUseful {
					logger.Fatal("offset is bigger then file is")
				}
				break
			} else {
				logger.Fatal("cant read the input file", errorReader)
			}
		}
		inputBuf, storageBuf = fixTruncBytes(inputBuf, storageBuf, n)

		if !offset.IsOffsetUseful {
			if n, err = offset.Read(inputBuf); err != nil {
				continue
			} else {
				inputBuf = inputBuf[n:]
			}
		}
		n, _ = limit.Read(inputBuf)
		inputBuf = inputBuf[:n]

		convLine := string(inputBuf)
		for _, name := range conversations {
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
					inputBuf = inputBuf[len(inputBuf)-len(convLine):]
				}
				if isTrimmed {
					convLine = ""
					for _, b := range string(inputBuf) {
						if unicode.IsSpace(b) {
							spaces = append(spaces, []byte(string(b))...)
						} else {
							convLine += string(spaces) + string(b)
							spaces = make([]byte, 0)
						}
					}
				}
			}
		}
		if len(convLine) < opts.BlockSize {
			fmt.Fprint(writer, convLine)
			writer.Flush()
		} else {
			convLineInBytes := []byte(convLine)
			outputStorageBuf := make([]byte, 0)
			for len(convLineInBytes) >= opts.BlockSize {
				inputBuf = convLineInBytes[:opts.BlockSize]
				inputBuf, outputStorageBuf = fixTruncBytes(inputBuf, outputStorageBuf, len(inputBuf))
				if len(inputBuf) != 0 {
					fmt.Fprint(writer, string(inputBuf))
					writer.Flush()
				}
				convLineInBytes = convLineInBytes[opts.BlockSize:]
			}
			if len(convLineInBytes) != 0 {
				inputBuf = convLineInBytes
				inputBuf, outputStorageBuf = fixTruncBytes(inputBuf, outputStorageBuf, len(inputBuf))
				fmt.Fprint(writer, string(inputBuf))
				writer.Flush()
			}
		}
	}
}
