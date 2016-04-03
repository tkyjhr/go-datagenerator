package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
)

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
	maxFileSize = 1 * gb
)

type dataOption struct {
	Key         string
	Description string
	reader      io.Reader
}

func main() {

	log.SetFlags(0)

	dataOptionList := []dataOption{
		{"0", "Fill with 0x00.", ZeroReader},
		{"f", "Fill with 0xFF.", FfReader},
		{"r", "Fiil with random bytes.", RandReader},
		{"ra", "Fill with random alphabet character.", RandomAlphabetReader},
		{"ran", "Fill with random alphabet and numeric characters.", RandomAlphabetNumericCharacterReader},
		{"ctr", "Fill with bytes increasing from 0x00 to 0xFF. After 0xFF it backs to 0x00.", &ByteCounterReader{}},
		{"ctr2", "Fill with uint16(BigEndian) increasing from 0x0000 to 0xFFFF.\n\t       After 0xFFFF, the value starts from 0x0000 again.\n\t       The size must be a multiple of 2.", &Uint16CounterReader{}},
		{"ctr4", "Fill with uint32(BigEndian) increasing from 0x00000000 to 0xFFFFFFFF.\n\t       After 0xFFFFFFFF, the value starts from 0x0000 again.\n\t       The size must be a multiple of 4.", &Uint32CounterReader{}},
		{"ctr8", "Fill with uint64(BigEndian) increasing from 0x0000000000000000 to 0xFFFFFFFFFFFFFFFF.\n\t       After 0xFFFFFFFFFFFFFFFF, the value starts from 0x0000000000000000 again.\n\t       The size must be a multiple of 8.", &Uint64CounterReader{}},
	}

	// Generate the help text for the option flag using text/template
	dataOptionHelpTemplate := template.New("DataOptionHelp")
	template.Must(dataOptionHelpTemplate.Parse("{{range $i, $value := .}}\t{{printf \"%-4s\" $value.Key}} : {{$value.Description}}\n{{end}}"))
	dataOptionHelpBuffer := &bytes.Buffer{}
	dataOptionHelpTemplate.Execute(dataOptionHelpBuffer, dataOptionList)

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.UsageText}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = "datagenerator"
	app.Usage = "Output data"
	app.HideVersion = true
	app.UsageText = `datagenerator (-d [Data]) (-o [OutputPath]) [Size]

   [Data]       : Data to generate. See below.
   [OutputPath] : Output file path. Directries are created if necessary.
                  If ommitted, data is output to stdout.
   [Size]       : Size in bytes. KB and MB can be used as suffix for (* 1024) and (* 1024 * 1024) respectively.
                  Hexadecimal value can be used with 0x prefix.
                  Maximum size is 1GB.
`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "data, d",
			Usage: "Specify content of output data. 0 is the default.\n" + dataOptionHelpBuffer.String(),
		},
		cli.StringFlag{
			Name: "o",
			Usage: "Output file path.",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() < 1 {
			log.Fatalf("[Error] Insufficient aguments. You must pass [Size] argument. Use -h option to see help.")
		}

		// Get and check file size.
		sizeArg := c.Args().Get(0)
		factor := 1
		if strings.HasSuffix(sizeArg, "KB") || strings.HasSuffix(sizeArg, "kb") {
			factor = kb
			sizeArg = sizeArg[:len(sizeArg) - 2]
		} else if strings.HasSuffix(sizeArg, "MB") || strings.HasSuffix(sizeArg, "mb") {
			factor = mb
			sizeArg = sizeArg[:len(sizeArg) - 2]
		}
		size, err := strconv.ParseInt(sizeArg, 0, 64)
		if err != nil || size < 0 {
			log.Fatalf("[Error] '%s' is invalid as a size argument. Use -h option to see help and check the acceptable format of the [Size] option.\n", sizeArg)
		}
		size *= int64(factor)

		if size > maxFileSize {
			log.Fatalf("[Error] %d bytes(%.2fGB) is too large. Currently the maximum file size is %d bytes (%.2fGB).", size, float64(size) / gb, maxFileSize, float64(maxFileSize) / gb)
		}

		// Get and check --data option
		var reader io.Reader
		dataOption := c.GlobalString("data")
		if dataOption == "" {
			// Default
			reader = ZeroReader
		} else {
			for _, d := range dataOptionList {
				if d.Key == dataOption {
					reader = d.reader
				}
			}
		}

		if reader == nil {
			log.Fatalf("[Error] Invalid --data option.\n")
		}

		// Get and check output path.
		var outputFile *os.File
		outputFilePath := c.GlobalString("o")
		if outputFilePath == "" {
			outputFile = os.Stdout
		} else {
			dir, file := filepath.Split(outputFilePath)
			if file == "" {
				log.Fatalf("[Error] Invalid file name.\n")
			}
			if dir != "" {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					log.Fatalf("[Error] Failed to create directries \"%s\".\n%s\n", dir, err.Error())
				}
			}

			outputFile, err = os.Create(outputFilePath)
			if err != nil {
				log.Fatalf("[Error] Failed to create a file \"%s\".\n%s\n", outputFilePath, err.Error())
			}
			defer func() {
				if err := outputFile.Close(); err != nil {
					log.Fatalf("[Error] Failed to close file.\n%s\n", err.Error())
				}
			}()
		}

		// Initialize rand for Readers that use rand.
		rand.Seed(time.Now().UnixNano())

		writer := bufio.NewWriter(outputFile)
		_, err = io.CopyN(writer, reader, size)
		if err != nil {
			log.Fatalf("[Error] Failed to write data.\n%s\n", err.Error())
		}

		err = writer.Flush()
		if err != nil {
			log.Fatalf("[Error] Failed to flush data.\n%s\n", err.Error())
		}
	}

	app.Run(os.Args)
}
