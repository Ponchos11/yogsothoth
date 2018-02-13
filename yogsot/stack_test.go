package yogsot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *godo.Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = godo.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func newTestClient() *YogClient {
	yogClient := YogClient{client}
	return &yogClient
}

func TestCreateStack(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "name",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":1}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	template := []byte(`
  Parameters:
    StackName:
      Description: The name of the stack to deploy
      Type: String
      Default: FurnaceStack
    Port:
      Description: Test port
      Type: Number
      Default: 80

  Resources:
    Droplet:
      Name: name
      Region: region
      Size: size
      Image: 1
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"`)
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}

func TestCreateStackMoreThanFiveDroplets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "name",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":1}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	// Normally, the none uniqueness of these names should raise an error.
	// But for unit testing purposes, I'm ignore those for now.
	template := []byte(`
  Parameters:
    StackName:
      Description: The name of the stack to deploy
      Type: String
      Default: FurnaceStack
    Port:
      Description: Test port
      Type: Number
      Default: 80

  Resources:
    Droplet1:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet2:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet3:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet4:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet5:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet6:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    Droplet7:
      Name: name
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"`)
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
}

func TestCreateStackMultipleResources(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":               "MyDroplet",
			"region":             "region",
			"size":               "size",
			"image":              "ubuntu-14-04-x64",
			"ssh_keys":           nil,
			"backups":            false,
			"ipv6":               false,
			"private_networking": false,
			"monitoring":         false,
			"tags":               []interface{}{"TestStack"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"droplet":{"id":987}, "links":{"actions": [{"id": 1, "href": "http://example.com", "rel": "create"}]}}`)
	})

	mux.HandleFunc("/v2/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"region":     "nyc3",
			"droplet_id": float64(987),
		}
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body = %+v, expected = %+v", v, expected)
		}

		fmt.Fprint(w, `{"floating_ip":{"region":{"slug":"nyc3"},"droplet":{"id":987},"ip":"192.168.0.1"}}`)
	})

	template := []byte(`
  Parameters:
    StackName:
      Description: The name of the stack to deploy
      Type: String
      Default: FurnaceStack
    Port:
      Description: Test port
      Type: Number
      Default: 80

  Resources:
    Droplet1:
      Name: MyDroplet
      Region: region
      Size: size
      Backups: false
      IPv6: false
      PrivateNetworking: false
      Monitoring: false
      Type: Droplet
      Image:
        Slug: "ubuntu-14-04-x64"
    FloatingIP:
      Type: FloatingIP
      Region: nyc3
      DropletID: MyDroplet`)
	request := CreateStackRequest{TemplateBody: template, StackName: "TestStack"}
	yogClient := newTestClient()
	response, err := yogClient.CreateStack(request)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}
	if len(response.Resources) < 1 {
		t.Fatal("should have contained one created resource")
	}
	for _, v := range response.Resources {
		if f, ok := v.(*FloatingIP); ok {
			if f.Request.DropletID != 987 {
				t.Fatalf("floatingip request droplet id should have equaled 987. Was instead: %d\n", f.Request.DropletID)
			}
		}
	}
}
