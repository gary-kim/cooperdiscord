package cooper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCodes(t *testing.T) {
	assert.Equal(t, []string{"ESC 000.1", "ESC 000.2", "ESC 000.3", "ESC 000.4"}, splitCodes("ESC 000.1-000.4"), "splitting ESC 000.1-000.4")
	assert.Equal(t, []string{"EID 101"}, splitCodes("EID 101"), "splitting EID 101")
	assert.Equal(t, []string{"FA 342A", "FA 342B"}, splitCodes("FA 342A, FA 342B"), "Splitting FA 342A, FA 342B")
	assert.Equal(t, []string{"EID 320", "EID 321", "EID 322", "EID 323"}, splitCodes("EID 320 - 323"), "Splitting EID 320 - 323")
	assert.Equal(t, []string{"Arch 154 A", "Arch 154 B"}, splitCodes("Arch 154 A-B"), "Splitting Arch 154 A-B")
}
