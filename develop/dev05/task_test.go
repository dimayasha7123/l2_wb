package main

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const directory = "test_data"

func Test_e2e_grep(t *testing.T) {
	grepPath := path.Join(t.TempDir(), "task")
	compileErr := exec.Command("go1.20.1", "build", "-o", grepPath, "task.go").Run()
	require.NoError(t, compileErr, "can't compile grep: %v\n", compileErr)

	tests := []struct {
		name           string
		args           []string
		stdinFile      string
		stdin          string
		wantStdoutFile string
		wantStdout     string
		wantStderr     string
		errFunc        assert.ErrorAssertionFunc
	}{
		{
			name:       "empty args",
			wantStderr: "ШАБЛОН не может быть пустым\n",
			errFunc:    assert.Error,
		},
		{
			name:       "only pattern",
			args:       []string{"o."},
			wantStderr: "Не указаны файлы или \"-\" для чтения из STDIN\n",
			errFunc:    assert.Error,
		},
		{
			name:       "validation check",
			args:       []string{"-c", "-1", "o.", "-"},
			wantStderr: "Не верно указаны флаги: COUNT не может быть меньше нуля\n",
			errFunc:    assert.Error,
		},
		{
			name:       "bad regex",
			args:       []string{"\\", "-"},
			stdin:      "London",
			wantStderr: "Не могу грепнуть: не могу скомпилировать регулярное выражение \\: error parsing regexp: trailing backslash at end of expression: ``\n",
			errFunc:    assert.Error,
		},
		{
			name:       "after less",
			args:       []string{"-A", "2", "o.", "-"},
			stdinFile:  "after_1.txt",
			wantStdout: "London\n1234\n",
			errFunc:    assert.NoError,
		},
		{
			name:       "after equal",
			args:       []string{"-A", "1", "o.", "-"},
			stdinFile:  "after_1.txt",
			wantStdout: "London\n1234\n",
			errFunc:    assert.NoError,
		},
		{
			name:       "after more",
			args:       []string{"-A", "1", "o.", "-"},
			stdinFile:  "after_2.txt",
			wantStdout: "London\n1234\n",
			errFunc:    assert.NoError,
		},
		{
			name:       "before less",
			args:       []string{"-B", "2", "o.", "-"},
			stdinFile:  "before_1.txt",
			wantStdout: "abcd\nLondon\n",
			errFunc:    assert.NoError,
		},
		{
			name:       "before equal",
			args:       []string{"-B", "1", "o.", "-"},
			stdinFile:  "before_1.txt",
			wantStdout: "abcd\nLondon\n",
			errFunc:    assert.NoError,
		},
		{
			name:       "before more",
			args:       []string{"-B", "1", "o.", "-"},
			stdinFile:  "before_2.txt",
			wantStdout: "efgh\nLondon\n",
			errFunc:    assert.NoError,
		},
		{
			name:           "after and before",
			args:           []string{"-A", "1", "-B", "2", "o.", "-"},
			stdinFile:      "after_and_before.txt",
			wantStdoutFile: "after_and_before_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "C prior than A and B",
			args:           []string{"-A", "2", "-B", "3", "-C", "1", "o.", "-"},
			stdinFile:      "c.txt",
			wantStdoutFile: "c_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "invert and ABC check",
			args:           []string{"-A", "1", "-B", "2", "-v", "o.", "-"},
			stdinFile:      "invert_and_ABC.txt",
			wantStdoutFile: "invert_and_ABC_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "count",
			args:           []string{"-c", "2", "o.", "-"},
			stdinFile:      "count.txt",
			wantStdoutFile: "count_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "count_invert",
			args:           []string{"-c", "3", "-v", "o.", "-"},
			stdinFile:      "count_invert.txt",
			wantStdoutFile: "count_invert_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "count_and_ABC",
			args:           []string{"-c", "2", "-A", "1", "-B", "2", "o.", "-"},
			stdinFile:      "count_and_ABC.txt",
			wantStdoutFile: "count_and_ABC_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "ignore_case_false",
			args:           []string{"o.", "-"},
			stdinFile:      "ignore_case.txt",
			wantStdoutFile: "ignore_case_o1.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "ignore_case_true",
			args:           []string{"-i", "o.", "-"},
			stdinFile:      "ignore_case.txt",
			wantStdoutFile: "ignore_case_o2.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "ignore_case_fixed",
			args:           []string{"-i", "-F", "o.", "-"},
			stdinFile:      "ignore_case_fixed.txt",
			wantStdoutFile: "ignore_case_fixed_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "line_num",
			args:           []string{"-n", "o.", "-"},
			stdinFile:      "line_num.txt",
			wantStdoutFile: "line_num_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "highlighting",
			args:           []string{"-h", "o.", "-"},
			stdinFile:      "highlighting.txt",
			wantStdoutFile: "highlighting_o.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "few files no names",
			args:           []string{"o.", path.Join(directory, "file1.txt"), path.Join(directory, "file2.txt")},
			wantStdoutFile: "file_o1.txt",
			errFunc:        assert.NoError,
		},
		{
			name:           "few files with names",
			args:           []string{"-f", "o.", path.Join(directory, "file1.txt"), path.Join(directory, "file2.txt")},
			wantStdoutFile: "file_o2.txt",
			errFunc:        assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.errFunc, "error assertion function can't be nil")

			cmd := exec.Command(grepPath, tt.args...)
			stdin := tt.stdin
			if len(tt.stdinFile) != 0 {
				data, err := os.ReadFile(path.Join(directory, tt.stdinFile))
				require.NoError(t, err, "can't read stdin file %s:%v\n", tt.stdinFile, err)
				stdin = string(data)
			}
			cmd.Stdin = strings.NewReader(stdin)
			output := strings.Builder{}
			cmd.Stdout = &output
			errOutput := strings.Builder{}
			cmd.Stderr = &errOutput
			wantStdout := tt.wantStdout
			if len(tt.wantStdoutFile) != 0 {
				data, err := os.ReadFile(path.Join(directory, tt.wantStdoutFile))
				require.NoError(t, err, "can't read want stdout file %s:%v\n", tt.wantStdoutFile, err)
				wantStdout = string(data)
			}

			err := cmd.Run()

			tt.errFunc(t, err)
			assert.Equal(t, wantStdout, output.String())
			assert.Equal(t, tt.wantStderr, errOutput.String())
		})
	}
}
