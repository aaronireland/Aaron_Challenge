// Cobra CLI command to validate batches of ABCD Bank Credit Card numbers
//
// Implementation of the HackerRank code challenge.
//
//	Validates lines of data provided either by unix pipe or file path.
//
// Usage:
//
//	validate [flags]
//
// Flags:
//
//	-f, --file string   path to a file containing the data to validate
//	-h, --help          help for validate
//	-v, --verbose       output full validation result with input text and the error messages
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const MAXLINES = 99

var (
	cmdFs         *afero.Afero
	fs            afero.Fs
	filePath      string
	verbose       bool
	validationCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate a batch of credit card numbers",
		Long: `Implementation of the HackerRank code challenge.1
		Validates lines of data provided either by unix pipe or file path.`,
		RunE: RunValidation,
	}
)

func init() {
	fs = afero.NewOsFs()
	cmdFs = &afero.Afero{Fs: fs}
	InitValidationCmdFlags(validationCmd)
}

func InitValidationCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&filePath, "file", "f", "", "path to a file containing the data to validate",
	)
	cmd.Flags().BoolVarP(
		&verbose, "verbose", "v", false, "output full validation result with input text and the error messages",
	)
}

func RunValidation(cmd *cobra.Command, args []string) error {

	validator := NewAbcdBankValidator()

	r, err := dataReader(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	var count int
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		input := scanner.Text()
		switch count {
		case 0:
			expected, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("malformed input: first line must be the size of the batch, got: %s", input)
			} else if expected > MAXLINES || expected < 1 {
				return fmt.Errorf("invalid input: %v, must be between 1 and %v", expected, MAXLINES)
			}
		case MAXLINES + 1:
			return nil
		default:
			isValid, errs := validator.Validate(input)
			cmd.Println(validationToString(input, isValid, errs...))
		}
		count++
	}

	return nil
}

func Execute() {
	if err := validationCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func validationToString(input string, isValid bool, errs ...error) string {
	validStr := "INVALID"
	if isValid {
		validStr = "VALID"
	}

	if !verbose {
		return validStr
	}

	msg := fmt.Sprintf("%s: %s", input, validStr)
	if len(errs) > 0 {
		errMsg := fmt.Sprintf("%d validation errors...", len(errs))
		for _, err := range errs {
			errStr := fmt.Sprintf("\n\t- %s", err)
			errMsg += errStr
		}
		msg = fmt.Sprintf("%s with %s", msg, errMsg)
	}

	return msg
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func dataReader(fp string) (io.ReadCloser, error) {
	if fp == "" && isInputFromPipe() {
		return os.Stdin, nil
	}

	if !fileExists(fp) {
		return nil, fmt.Errorf("the file provided does not exist")
	}

	file, err := fs.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read the file %s: %w", fp, err)
	}
	return file, nil
}

func fileExists(filepath string) bool {
	info, err := fs.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
