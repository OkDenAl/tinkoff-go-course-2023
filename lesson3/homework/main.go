package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"lecture03_homework/entity"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrOffsetTooBig          = errors.New("offset is bigger then buf")
	ErrOffsetIsLessThenZero  = errors.New("offset is less then 0")
	ErrUpperAndLowerTogether = errors.New("cant use upper_case and lower_case conversation together")
	ErrUnknownConv           = errors.New("unknown conversation")
	ErrFileExists            = errors.New("file already exists")
)

type offsetReader struct {
	offset         int
	IsOffsetUseful bool
}

func (r *offsetReader) Read(buf []byte) (n int, err error) {
	if len(buf) < r.offset {
		r.offset -= len(buf)
		return 0, ErrOffsetTooBig
	}
	r.IsOffsetUseful = true
	return r.offset, nil
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
	Conv      []string
}

func ParseFlags() (*Options, error) {
	var opts Options
	var convInStr string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "number of bytes to be skipped")
	flag.IntVar(&opts.Limit, "limit", -1, "max number of bytes to be read. by default - all file")
	flag.IntVar(&opts.BlockSize, "block-size", 10, "the size of one block in bytes when reading and writing")
	flag.StringVar(&convInStr, "conv", "", "one or more of the possible "+
		"transformations over the text, separated by commas.")
	flag.Parse()
	opts.Conv = strings.Split(convInStr, ",")

	if err := validateOffset(opts.Offset); err != nil {
		return nil, fmt.Errorf("invalid offset: %w", err)
	}
	if err := validateConv(opts.Conv); err != nil {
		return nil, fmt.Errorf("invalid conv: %w", err)
	}

	return &opts, nil
}

func setupReader(opts *Options) (*bufio.Reader, *os.File, error) {
	var (
		reader *bufio.Reader
		f      *os.File
		err    error
	)
	if opts.From != "" {
		f, err = os.Open(opts.From)
		if err != nil {
			return nil, nil, fmt.Errorf("cant open input file: %w", err)
		}
		reader = bufio.NewReader(f)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}
	return reader, f, nil
}

func setupWriter(opts *Options) (*bufio.Writer, *os.File, error) {
	var (
		writer *bufio.Writer
		f      *os.File
		err    error
	)
	if opts.To != "" {
		if _, err := os.Stat(opts.To); !os.IsNotExist(err) {
			return nil, nil, ErrFileExists
		}
		f, err = os.Create(opts.To)
		if err != nil {
			return nil, nil, fmt.Errorf("cant create the output file: err: %w", err)
		}
		writer = bufio.NewWriter(f)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}
	return writer, f, nil
}

func validateConv(inputConv []string) error {
	if len(inputConv) == 1 && inputConv[0] == "" {
		return nil
	}
	usedConvs := make(map[entity.Conversation]struct{})
	for _, name := range inputConv {
		switch entity.Conversation(name) {
		case entity.LowerCase:
			if _, ok := usedConvs[entity.UpperCase]; ok {
				return ErrUpperAndLowerTogether
			}
			usedConvs[entity.LowerCase] = struct{}{}
		case entity.UpperCase:
			if _, ok := usedConvs[entity.LowerCase]; ok {
				return ErrUpperAndLowerTogether
			}
			usedConvs[entity.UpperCase] = struct{}{}
		case entity.TrimSpaces:
		default:
			return ErrUnknownConv
		}
	}
	return nil
}

func validateOffset(offset int) error {
	if offset < 0 {
		return ErrOffsetIsLessThenZero
	}
	return nil
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
	logger := log.New(os.Stderr, "ERROR:\t", log.LstdFlags)
	opts, err := ParseFlags()
	if err != nil {
		logger.Fatal("can not parse flags: ", err)
	}

	reader, _, err := setupReader(opts)
	//if input != nil {
	//	defer input.Close()
	//}
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
	var outputStorageBuf []byte

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
		for _, name := range opts.Conv {
			switch entity.Conversation(name) {
			case entity.LowerCase:
				convLine = strings.ToLower(convLine)
			case entity.UpperCase:
				convLine = strings.ToUpper(convLine)
			case entity.TrimSpaces:
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
