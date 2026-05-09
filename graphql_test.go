package marketsurge

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGraphQLRequestMarshal(t *testing.T) {
	t.Parallel()

	t.Run("with operation name", func(t *testing.T) {
		t.Parallel()

		got, err := json.Marshal(GraphQLRequest[map[string]string]{
			OperationName: "GetStory",
			Variables:     map[string]string{"id": "123"},
			Query:         "query GetStory($id: ID!) { story(id: $id) { id } }",
		})
		if err != nil {
			t.Fatalf("json.Marshal(GraphQLRequest) error = %v", err)
		}

		want := []byte(
			`{"operationName":"GetStory","variables":{"id":"123"},` +
				`"query":"query GetStory($id: ID!) { story(id: $id) { id } }"}`,
		)
		if diff := cmp.Diff(string(want), string(got)); diff != "" {
			t.Fatalf("json.Marshal(GraphQLRequest) mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("without operation name", func(t *testing.T) {
		t.Parallel()

		got, err := json.Marshal(GraphQLRequest[map[string]string]{
			Variables: map[string]string{"id": "123"},
			Query:     "query { story(id: \"123\") { id } }",
		})
		if err != nil {
			t.Fatalf("json.Marshal(GraphQLRequest) error = %v", err)
		}

		var fields map[string]json.RawMessage
		if unmarshalErr := json.Unmarshal(got, &fields); unmarshalErr != nil {
			t.Fatalf("json.Unmarshal(GraphQLRequest JSON) error = %v", unmarshalErr)
		}
		if _, ok := fields["operationName"]; ok {
			t.Fatalf("json.Marshal(GraphQLRequest) unexpectedly included operationName key: %s", got)
		}

		want := []byte(`{"variables":{"id":"123"},"query":"query { story(id: \"123\") { id } }"}`)
		if diff := cmp.Diff(string(want), string(got)); diff != "" {
			t.Fatalf("json.Marshal(GraphQLRequest) mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestGraphQLResponseUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("data only", func(t *testing.T) {
		t.Parallel()

		var got GraphQLResponse[struct {
			Story struct {
				ID string `json:"id"`
			} `json:"story"`
		}]
		if err := json.Unmarshal([]byte(`{"data":{"story":{"id":"123"}}}`), &got); err != nil {
			t.Fatalf("json.Unmarshal(GraphQLResponse) error = %v", err)
		}

		want := GraphQLResponse[struct {
			Story struct {
				ID string `json:"id"`
			} `json:"story"`
		}]{
			Data: struct {
				Story struct {
					ID string `json:"id"`
				} `json:"story"`
			}{
				Story: struct {
					ID string `json:"id"`
				}{ID: "123"},
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("GraphQLResponse unmarshal mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("errors only", func(t *testing.T) {
		t.Parallel()

		var got GraphQLResponse[struct{}]
		if err := json.Unmarshal([]byte(
			`{"errors":[{"message":"bad request","path":["story"],`+
				`"extensions":{"code":"BAD_REQUEST"}}]}`,
		), &got); err != nil {
			t.Fatalf("json.Unmarshal(GraphQLResponse) error = %v", err)
		}

		want := GraphQLResponse[struct{}]{
			Errors: []GraphQLFieldError{{
				Message: "bad request",
				Path:    []string{"story"},
				Extensions: &GraphQLFieldErrorExtensions{
					Code: "BAD_REQUEST",
				},
			}},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("GraphQLResponse unmarshal mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("mixed data and errors", func(t *testing.T) {
		t.Parallel()

		var got GraphQLResponse[struct {
			Story struct {
				ID string `json:"id"`
			} `json:"story"`
		}]
		if err := json.Unmarshal(
			[]byte(`{"data":{"story":{"id":"123"}},"errors":[{"message":"partial failure"}]}`),
			&got,
		); err != nil {
			t.Fatalf("json.Unmarshal(GraphQLResponse) error = %v", err)
		}

		if got.Data.Story.ID != "123" {
			t.Fatalf("GraphQLResponse.Data.Story.ID = %q, want %q", got.Data.Story.ID, "123")
		}
		if diff := cmp.Diff([]GraphQLFieldError{{Message: "partial failure"}}, got.Errors); diff != "" {
			t.Fatalf("GraphQLResponse.Errors mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("null data", func(t *testing.T) {
		t.Parallel()

		var got GraphQLResponse[struct {
			Story struct {
				ID string `json:"id"`
			} `json:"story"`
		}]
		if err := json.Unmarshal([]byte(`{"data":null}`), &got); err != nil {
			t.Fatalf("json.Unmarshal(GraphQLResponse) error = %v", err)
		}

		if !reflect.DeepEqual(got.Data, struct {
			Story struct {
				ID string `json:"id"`
			} `json:"story"`
		}{}) {
			t.Fatalf("GraphQLResponse.Data = %#v, want zero value", got.Data)
		}
	})
}
