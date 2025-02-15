package memberships

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amiulam/music-catalog/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Logout(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		withToken          bool
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			withToken:          true,
		},
		{
			name:               "unauthorized",
			expectedStatusCode: http.StatusUnauthorized,
			withToken:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := gin.New()

			h := &Handler{
				Engine: api,
			}

			h.RegisterRoute()
			w := httptest.NewRecorder()
			endpoint := `/memberships/logout`

			req, err := http.NewRequest(http.MethodPost, endpoint, nil)
			assert.NoError(t, err)

			if tt.withToken {
				token, err := jwt.CreateToken(1, "username", "secret")
				assert.NoError(t, err)
				req.Header.Set("Authorization", token)
			}

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
