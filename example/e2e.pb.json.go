package e2e

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
)

// BasicJSONMarshaler describes the default jsonpb.Marshaler used by all
// instances of Basic. This struct is safe to replace or modify but
// should not be done so concurrently.
var BasicJSONMarshaler = protojson.MarshalOptions{}

// MarshalJSON satisfies the encoding/json Marshaler interface. This method
// uses the more correct jsonpb package to correctly marshal the message.
func (m *Basic) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}
	return BasicJSONMarshaler.Marshal(m)
}

var _ json.Marshaler = (*Basic)(nil)

// BasicJSONUnmarshaler describes the default jsonpb.Unmarshaler used by all
// instances of Basic. This struct is safe to replace or modify but
// should not be done so concurrently.
var BasicJSONUnmarshaler = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface. This method
// uses the more correct jsonpb package to correctly unmarshal the message.
func (m *Basic) UnmarshalJSON(b []byte) error {
	return BasicJSONUnmarshaler.Unmarshal(b, m)
}

var _ json.Unmarshaler = (*Basic)(nil)
