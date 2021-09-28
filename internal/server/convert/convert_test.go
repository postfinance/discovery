package convert

import (
	"testing"
	"time"

	"github.com/postfinance/discovery"
	"github.com/postfinance/discovery/internal/exporter"
	"github.com/stretchr/testify/assert"
)

func TestConvertServer(t *testing.T) {
	expected := discovery.NewServer("name", discovery.Labels{"key": "val"})
	expected.State = discovery.Joining
	pb := ServerToPB(expected)
	s := ServerFromPB(pb)

	assert.True(t, expected.Modified.Equal(s.Modified))
	expected.Modified, s.Modified = time.Time{}, time.Time{}
	assert.Equal(t, expected, s)
}

func TestConvertNamespace(t *testing.T) {
	expected := discovery.DefaultNamespace()
	pb := NamespaceToPB(expected)
	n := NamespaceFromPB(pb)

	assert.True(t, expected.Modified.Equal(n.Modified))
	expected.Modified, n.Modified = time.Time{}, time.Time{}
	assert.Equal(t, expected, n)
}

func TestConvertTargetGroup(t *testing.T) {
	expected := exporter.TargetGroup{
		Targets: []string{"target1", "target2"},
		Labels: map[string]string{
			"label1": "value1",
		},
	}

	pb := TargetGroupToPB(&expected)
	tg := TargetGroupFromPB(pb)

	assert.EqualValues(t, expected, *tg)
}
