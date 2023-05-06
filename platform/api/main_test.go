package api

import (
	"testing"

	db "github.com/UitsHabib/ecommerce-microservice/platform/db/sqlc"
	"github.com/UitsHabib/ecommerce-microservice/platform/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
