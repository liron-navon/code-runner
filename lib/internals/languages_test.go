package internals

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func normalizeCliOutput(text string) string {
	symbolsToRemoveFromOutput := []string{"\x1b[93m", "\x1b[94m"}
	for _, s := range symbolsToRemoveFromOutput {
		text = strings.ReplaceAll(text, s, "")
	}
	return strings.TrimSpace(text)
}

func TestExecuteCode(t *testing.T) {
	tests := []struct {
		language         string
		code             string
		expectedLastLine string
	}{
		{
			"java",
			strings.TrimSpace(`
public class HelloWorld {
	public static void main(String[] args) {
		System.out.println("hello from java");
	}
}
`),
			"hello from java",
		},
		{
			"go",
			strings.TrimSpace(`
package main
import (
	"fmt"
)
func main() {
	fmt.Println("hello from go")
}
`),
			"hello from go",
		},
		{
			"node",
			strings.TrimSpace(`
console.log("hello from node.js")
`),
			"hello from node.js",
		},
		{
			"python3",
			strings.TrimSpace(`
print("hello from python3")
`),
			"hello from python3",
		},
	}

	for _, test := range tests {
		context := ProgrammingLanguages[test.language]
		output, err := context.TestExecuteCode(test.code)

		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, test.expectedLastLine, output[len(output)-1], test)
	}
}
