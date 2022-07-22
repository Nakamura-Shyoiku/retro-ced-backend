package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ulventech/retro-ced-backend/models"
)

func TestProductRepoInsert(t *testing.T) {
	r, err := NewMysqlProductRepo("root:retroced@tcp(127.0.0.1:3306)/retroced")
	require.NoError(t, err)
	p := []models.Product{{
		SiteId:       1,
		Img:          "aaaa",
		Category:     "bags",
		Price:        37780,
		Title:        "gucci handbag",
		Measurements: "1x 4 x 7 cm",
		Description:  "desc",
		Brand:        "brand",
		UpdatedAt:    time.Now().UTC(),
	}}
	err = r.InsertProducts(p)
	require.NoError(t, err)
}
