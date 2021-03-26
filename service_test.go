package discovery

import (
	"encoding/json"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	var tt = []struct {
		name           string
		url            string
		valid          bool
		expectedScheme string
		expectedHost   string
		expectedParams url.Values
		expectedPort   string
		expectedPath   string
	}{
		{
			"valid_url",
			"http://example.com/metrics?query=value",
			true,
			"http",
			"example.com",
			url.Values{
				"query": []string{"value"},
			},
			"",
			"/metrics",
		},
		{
			"space in name",
			"http://example.com/metrics?query=value",
			false,
			"http",
			"example.com",
			url.Values{
				"query": []string{"value"},
			},
			"",
			"/metrics",
		},
		{
			"valid_url_with_port",
			"http://example.com:9999/metrics",
			true,
			"http",
			"example.com",
			url.Values{},
			"9999",
			"/metrics",
		},
		{
			"not_valid: missing scheme",
			"example.com:9999/metrics",
			false,
			"",
			"example.com",
			url.Values{},
			"9999",
			"/metrics",
		},
		{
			"not_valid: empty endpoint",
			"",
			false,
			"",
			"",
			url.Values{},
			"",
			"",
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewService(tc.name, tc.url)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedHost, s.Endpoint.Hostname())
			require.Equal(t, tc.expectedParams, s.Endpoint.Query())
			require.Equal(t, tc.expectedPort, s.Endpoint.Port())
			require.Equal(t, tc.expectedPath, s.Endpoint.Path)
			require.Equal(t, tc.expectedScheme, s.Endpoint.Scheme)
		})
	}
}

func TestKey(t *testing.T) {
	var tt = []struct {
		name     string
		endpoint string
		valid    bool
	}{
		{
			"valid_with_port",
			"https://www.test.com:7002/metrics",
			true,
		},
		{
			"valid_without_port",
			"http://www.test.com/metrics",
			true,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewService(tc.name, tc.endpoint)
			require.NoError(t, err)
			if !tc.valid {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	s, err := NewService("name", "https://test.pnet.ch:7001/metrics")
	assert.NoError(t, err)
	t.Run("with url", func(t *testing.T) {
		d, err := json.Marshal(s)
		assert.NoError(t, err)
		s2 := Service{}
		err = json.Unmarshal(d, &s2)
		assert.NoError(t, err)
		// https://github.com/stretchr/testify/issues/502
		// assert.Equal(t, *s, s2)
		assert.Equal(t, s.Endpoint, s2.Endpoint)
		assert.Equal(t, s.Name, s2.Name)
		assert.Equal(t, s.Labels, s2.Labels)
		assert.Equal(t, s.ID, s2.ID)
		assert.Equal(t, s.Description, s2.Description)
		assert.True(t, s.Modified.Equal(s2.Modified))
	})
	t.Run("without url", func(t *testing.T) {
		s.Endpoint = nil
		d, err := json.Marshal(s)
		assert.NoError(t, err)
		s2 := Service{}
		err = json.Unmarshal(d, &s2)
		assert.NoError(t, err)
		assert.Equal(t, s.Endpoint, s2.Endpoint)
		assert.Equal(t, s.Name, s2.Name)
		assert.Equal(t, s.Labels, s2.Labels)
		assert.Equal(t, s.ID, s2.ID)
		assert.Equal(t, s.Description, s2.Description)
		assert.True(t, s.Modified.Equal(s2.Modified))
	})
}

func TestValid(t *testing.T) {
	var tests = []struct {
		testname string
		name     string
		id       string
	}{
		{
			"empty name",
			"",
			"name-123",
		},
		{
			"empty id",
			"name",
			"",
		},
		{
			"too short",
			"name",
			"name-",
		},
	}

	for _, test := range tests {
		s := Service{
			ID:   test.id,
			Name: test.name,
		}
		err := s.Validate()
		assert.Error(t, err)
	}
}

func TestTags(t *testing.T) {
	var tt = []struct {
		name     string
		tags     Labels
		expected string
	}{
		{
			"empty",
			Labels{},
			"",
		},
		{
			"one tag",
			Labels{
				"key": "value",
			},
			"key=value",
		},
		{
			"two tags",
			Labels{
				"key2": "value2",
				"key1": "value1",
			},
			"key1=value1,key2=value2",
		},
	}

	for _, tc := range tt {
		assert.Equal(t, tc.expected, tc.tags.String())
	}
}

func TestMarshalService(t *testing.T) {
	u, err := url.Parse("https://example.com")
	require.NoError(t, err)

	s := Service{
		Servers:     []string{"server1"},
		Selector:    "env=test",
		Name:        "job",
		ID:          "123",
		Labels:      Labels{"env": "test"},
		Endpoint:    u,
		Description: "description",
	}

	d, err := json.Marshal(s)
	require.NoError(t, err)

	ns := Service{}

	err = json.Unmarshal(d, &ns)
	require.NoError(t, err)
	assert.Equal(t, s, ns)
}

func TestFilterService(t *testing.T) {
	services := Services{
		Service{
			Name:     "service1",
			Endpoint: mustParseRequestURI("http://example1.pnet.ch/metrics"),
			Servers:  []string{"server1", "server2"},
		},
		Service{
			Name:     "service1",
			Endpoint: mustParseRequestURI("http://example2.pnet.ch/metrics"),
			Servers:  []string{"server2", "server3"},
		},
		Service{
			Name:     "service2",
			Endpoint: mustParseRequestURI("http://example3.pnet.ch/metrics"),
			Servers:  []string{"server1", "server2"},
		},
	}

	var tt = []struct {
		filters     []FilterFunc
		expectedLen int
	}{
		{
			[]FilterFunc{ServiceByName(regexp.MustCompile("1"))},
			2,
		},
		{
			[]FilterFunc{ServiceByName(regexp.MustCompile("service"))},
			3,
		},
		{
			[]FilterFunc{ServiceByName(regexp.MustCompile("service")), ServiceByServer(regexp.MustCompile("3"))},
			1,
		},
		{
			[]FilterFunc{ServiceByName(regexp.MustCompile("service")), ServiceByServer(regexp.MustCompile("3")), ServiceByEndpoint(regexp.MustCompile("2"))},
			1,
		},
		{
			[]FilterFunc{ServiceByName(regexp.MustCompile("service")), ServiceByServer(regexp.MustCompile("3")), ServiceByEndpoint(regexp.MustCompile("1"))},
			0,
		},
		{
			[]FilterFunc{},
			3,
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			assert.Len(t, services.Filter(tc.filters...), tc.expectedLen)
		})
	}

}

func mustParseRequestURI(rawURL string) *url.URL {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		panic(err)
	}

	return u
}
