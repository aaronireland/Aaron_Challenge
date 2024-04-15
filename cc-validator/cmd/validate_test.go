package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func writeLines(w io.Writer, lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return err
		}
	}
	return nil
}

func str(b []byte) string {
	return strings.TrimSpace(string(b))
}

func execute(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs(args)

	err := cmd.Execute()

	return str(b.Bytes()), err
}

func runWithPipe(t *testing.T, args []string, lines ...string) (string, error) {
	t.Helper()
	r, w, _ := os.Pipe()
	writeLines(w, lines...)
	_ = w.Close()

	os.Stdin = r
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)

	cmd := &cobra.Command{Use: "test", RunE: RunValidation}
	InitValidationCmdFlags(cmd)

	return execute(t, cmd, args...)
}

func TestExecute(t *testing.T) {
	// Only call the command execution if test called recursively
	if os.Getenv("TESTING_VALIDATION_COMMAND") == "1" {
		runE := validationCmd.RunE
		defer func() { validationCmd.RunE = runE }()

		validationCmd.RunE = func(cmd *cobra.Command, args []string) error {
			return nil
		}
		Execute()
		return
	} else if os.Getenv("TESTING_VALIDATION_COMMAND") == "boom" {
		runE := validationCmd.RunE
		defer func() { validationCmd.RunE = runE }()

		validationCmd.RunE = func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("boom!")
		}
		Execute()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecute")
	cmd.Env = append(os.Environ(), "TESTING_VALIDATION_COMMAND=1")

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			assert.Truef(t, exitErr.Success(), "validation execution failed: %s", exitErr)
		} else {
			assert.FailNow(t, fmt.Sprintf("validation command execution failed: %s", err))
		}
	}

	cmd = exec.Command(os.Args[0], "-test.run=TestExecute")
	cmd.Env = append(os.Environ(), "TESTING_VALIDATION_COMMAND=boom")
	err = cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		assert.Falsef(t, exitErr.Success(), "failed command should produce exit code 1: %s", exitErr)
	} else {
		assert.FailNow(t, "failed command should have generated an error")
	}
}

func TestValidationCommand(t *testing.T) {
	t.Run("Cobra command should run successfully with data piped directly from shell", func(t *testing.T) {
		args := []string{}
		output, err := runWithPipe(t, args, "3", "123456", "abc123", "4123456789123456")

		if assert.NoErrorf(t, err, "unexpected error: %s", err) {
			assert.Equal(t, "INVALID\nINVALID\nVALID", string(output), "validation command should generate expected validation results")
		}

	})

	t.Run("Cobra command should output validation errors and count if verbose flag set", func(t *testing.T) {
		args := []string{"-v"}
		output, err := runWithPipe(t, args, "3", "123456", "abc123", "4123456789123456")

		if assert.NoErrorf(t, err, "unexpected error: %s", err) {
			assert.Contains(t, string(output), "validation errors", "validation command with verbose flag set should show the errors")
		}
	})

	t.Run("Cobra command should run successfully with text file specified with path set on file flag", func(t *testing.T) {
		fs = afero.NewMemMapFs()
		cmdFs.Fs = fs
		dir := "data/tests/validate_from_file"
		file := "data/tests/validate_from_file/input.txt"
		data := "3\n123456\nabc123\n4123456789123456"
		fs.MkdirAll(dir, 0755)
		cmdFs.WriteFile(file, []byte(data), 0644)

		cmd := &cobra.Command{Use: "test", RunE: RunValidation}
		InitValidationCmdFlags(cmd)

		output, err := execute(t, cmd, fmt.Sprintf("-f=%s", file))

		if assert.NoErrorf(t, err, "unexpected error: %s", err) {
			assert.Equal(t, "INVALID\nINVALID\nVALID", string(output), "validation command should generate expected validation results")
		}

		t.Run("Cobra command should fail if data provided indicates a batch size of 0", func(t *testing.T) {
			args := []string{}
			_, err := runWithPipe(t, args, "0")
			assert.Error(t, err, "validation command should fail if the first line given is 0")

		})

		t.Run("Cobra command should fail if data provided indicates a batch size of greater than MAXSIZE", func(t *testing.T) {
			args := []string{}
			_, err := runWithPipe(t, args, "100", "1234-5678-1232-5554", "24acd", "ok ok ok")
			assert.Error(t, err, "validation command should fail if the first line given is a number greater than MAXSIZE")
		})

		t.Run("Cobra command should fail if there is invalid data provided in the first line", func(t *testing.T) {
			args := []string{}
			_, err := runWithPipe(t, args, "4131-5545-6709-0001")
			assert.Error(t, err, "validation command should fail if the first line given is not an integer value")
		})

	})
}
