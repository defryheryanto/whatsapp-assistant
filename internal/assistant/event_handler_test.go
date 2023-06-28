package assistant

import (
	"testing"

	"github.com/onsi/gomega"
)

func Test_extractCommands(t *testing.T) {
	m := gomega.NewWithT(t)
	t.Run("should extract commands delimited by enter", func(t *testing.T) {
		commands := extractCommands("%assign\ntest @628123456789 @628213458697")
		m.Expect(commands).To(gomega.BeEquivalentTo([]string{"assign"}))
	})
	t.Run("should extract commands delimited by space", func(t *testing.T) {
		commands := extractCommands("%assign test @628123456789 @628213458697")
		m.Expect(commands).To(gomega.BeEquivalentTo([]string{"assign"}))
	})
}
