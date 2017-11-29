package travis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otiai10/marmoset"
	. "github.com/otiai10/mint"
)

func TestNewClient(t *testing.T) {
	client := NewClient("my token")
	Expect(t, client).TypeOf("*travis.Client")
}

func TestClient_Trigger(t *testing.T) {

	client := NewClient("my token")
	client.Host = dummyTravisAPIServer().URL
	client.Output = httptest.NewRecorder()
	err := client.Trigger("case/01/ok")
	Expect(t, err).ToBe(nil)

	When(t, "failed to decode response", func(t *testing.T) {
		err := client.Trigger("case/02/EOF")
		Expect(t, err).Not().ToBe(nil)
	})

	When(t, "HTTP Client is not set", func(t *testing.T) {
		c0 := NewClient("my token")
		c0.HTTPClient = nil
		err := c0.Trigger("case/01/ok")
		Expect(t, err).Not().ToBe(nil)
	})

	When(t, "HTTP failed", func(t *testing.T) {
		c1 := NewClient("my token")
		c1.Host = ""
		err := c1.Trigger("case/01/ok")
		Expect(t, err).Not().ToBe(nil)
	})

	When(t, "response.@type is `error`", func(t *testing.T) {
		c2 := NewClient("my token")
		c2.Host = dummyTravisAPIServer().URL
		err := c2.Trigger("case/03/error")
		Expect(t, err).Not().ToBe(nil)
	})
}

func dummyTravisAPIServer() *httptest.Server {
	r := marmoset.NewRouter()
	r.POST("/repo/(?P<slug>.+)/requests", func(w http.ResponseWriter, r *http.Request) {
		switch r.FormValue("slug") {
		case "case/01/ok":
			json.NewEncoder(w).Encode(map[string]interface{}{"@type": "pending"})
		case "case/02/EOF":
			w.Write(nil)
		case "case/03/error":
			json.NewEncoder(w).Encode(map[string]interface{}{"@type": "error"})
		default:
			json.NewEncoder(w).Encode(map[string]interface{}{"@type": "pending"})
		}
	})
	return httptest.NewServer(r)
}
