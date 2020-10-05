package e2e

import (
	"encoding/json"
	"testing"
)

func TestJSONUnmarshal(t *testing.T) {
	buf := `{
    "b": "what",
    "a": "Hello",
    "meta": {
      "hello": "world"
    }
  }
  `
	s := new(Basic)
	err := json.Unmarshal([]byte(buf), s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", s)

	buf2, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", buf2)
}

// func TestJSONMarshal(t *testing.T) {
// 	s := &e2e.Basic{
// 		A: "Hello",
// 		Meta: &structpb.Value{
// 			Kind: &structpb.Value_StructValue{
// 				StructValue: &structpb.Struct{
// 					Fields: map[string]*structpb.Value{
// 						"Hello": {
// 							Kind: &structpb.Value_StringValue{
// 								StringValue: "World",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	buf, err := json.MarshalIndent(s, "", "  ")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("%s", buf)
// }
