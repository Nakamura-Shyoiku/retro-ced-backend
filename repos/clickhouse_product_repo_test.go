package repos

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepp(t *testing.T) {
	r, err := NewClickhouseProductsRepo("clickhouse://130.211.211.73:9000?user=default&password=default&database=retroced")
	require.NoError(t, err)
	p, err := r.Products("bags")
	require.NoError(t, err)

	assert.NotEmpty(t, p)
}
