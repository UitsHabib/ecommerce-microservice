package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/UitsHabib/ecommerce-microservice/platform/db/mock"
	db "github.com/UitsHabib/ecommerce-microservice/platform/db/sqlc"
	"github.com/UitsHabib/ecommerce-microservice/platform/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetFeatureAPI(t *testing.T) {
	feature := randomFeature()

	testCases := []struct {
		name          string
		featureID     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			featureID: feature.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFeature(gomock.Any(), feature.ID).
					Times(1).
					Return(feature, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFeature(t, recorder.Body, feature)
			},
		},
		{
			name:      "NotFound",
			featureID: feature.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFeature(gomock.Any(), feature.ID).
					Times(1).
					Return(db.Feature{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			featureID: feature.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFeature(gomock.Any(), feature.ID).
					Times(1).
					Return(db.Feature{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			featureID: "1",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetFeature(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/features/%s", tc.featureID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetFeaturesAPI(t *testing.T) {
	n := 5
	features := make([]db.Feature, n)
	for i := 0; i < n; i++ {
		features[i] = randomFeature()
	}

	type Query struct {
		page  int
		limit int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				page:  1,
				limit: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListFeaturesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListFeatures(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(features, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFeatures(t, recorder.Body, features)
			},
		},
		{
			name: "InternalError",
			query: Query{
				page:  1,
				limit: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListFeaturesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListFeatures(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Feature{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				page:  -1,
				limit: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFeature(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLimit",
			query: Query{
				page:  1,
				limit: -n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFeature(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/features"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// add query parameters to request	 url
			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("limit", fmt.Sprintf("%d", tc.query.limit))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomFeature() db.Feature {
	return db.Feature{
		ID:          util.RandomUUID(),
		Title:       util.RandomString(6),
		Slug:        util.RandomString(6),
		Description: util.RandomNullableString(),
	}
}

func requireBodyMatchFeature(t *testing.T, body *bytes.Buffer, feature db.Feature) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotFeature getFeatureResponse
	err = json.Unmarshal(data, &gotFeature)
	require.NoError(t, err)

	f := formatFeature(&feature)
	require.Equal(t, gotFeature, f)
}

func requireBodyMatchFeatures(t *testing.T, body *bytes.Buffer, features []db.Feature) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotFeatures []getFeatureResponse
	err = json.Unmarshal(data, &gotFeatures)
	require.NoError(t, err)

	tmp_features := FormatFeatures(&features)
	require.Equal(t, gotFeatures, tmp_features)
}
