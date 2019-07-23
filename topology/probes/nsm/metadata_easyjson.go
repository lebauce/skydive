// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package nsm

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(in *jlexer.Lexer, out *RemoteNSMMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "SourceCrossConnectID":
			out.SourceCrossConnectID = string(in.String())
		case "DestinationCrossConnectID":
			out.DestinationCrossConnectID = string(in.String())
		case "Via":
			(out.Via).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(out *jwriter.Writer, in RemoteNSMMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"SourceCrossConnectID\":"
		out.RawString(prefix[1:])
		out.String(string(in.SourceCrossConnectID))
	}
	{
		const prefix string = ",\"DestinationCrossConnectID\":"
		out.RawString(prefix)
		out.String(string(in.DestinationCrossConnectID))
	}
	{
		const prefix string = ",\"Via\":"
		out.RawString(prefix)
		(in.Via).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoteNSMMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoteNSMMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoteNSMMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoteNSMMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(in *jlexer.Lexer, out *RemoteConnectionMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "SourceNSM":
			out.SourceNSM = string(in.String())
		case "DestinationNSM":
			out.DestinationNSM = string(in.String())
		case "NetworkServiceEndpoint":
			out.NetworkServiceEndpoint = string(in.String())
		case "MechanismType":
			out.MechanismType = string(in.String())
		case "MechanismParameters":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.MechanismParameters = make(map[string]string)
				} else {
					out.MechanismParameters = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.MechanismParameters)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "Labels":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Labels = make(map[string]string)
				} else {
					out.Labels = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 string
					v2 = string(in.String())
					(out.Labels)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(out *jwriter.Writer, in RemoteConnectionMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"SourceNSM\":"
		out.RawString(prefix[1:])
		out.String(string(in.SourceNSM))
	}
	{
		const prefix string = ",\"DestinationNSM\":"
		out.RawString(prefix)
		out.String(string(in.DestinationNSM))
	}
	{
		const prefix string = ",\"NetworkServiceEndpoint\":"
		out.RawString(prefix)
		out.String(string(in.NetworkServiceEndpoint))
	}
	{
		const prefix string = ",\"MechanismType\":"
		out.RawString(prefix)
		out.String(string(in.MechanismType))
	}
	{
		const prefix string = ",\"MechanismParameters\":"
		out.RawString(prefix)
		if in.MechanismParameters == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v3First := true
			for v3Name, v3Value := range in.MechanismParameters {
				if v3First {
					v3First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v3Name))
				out.RawByte(':')
				out.String(string(v3Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"Labels\":"
		out.RawString(prefix)
		if in.Labels == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v4First := true
			for v4Name, v4Value := range in.Labels {
				if v4First {
					v4First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v4Name))
				out.RawByte(':')
				out.String(string(v4Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoteConnectionMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoteConnectionMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoteConnectionMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoteConnectionMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm1(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(in *jlexer.Lexer, out *LocalNSMMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "CrossConnectID":
			out.CrossConnectID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(out *jwriter.Writer, in LocalNSMMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"CrossConnectID\":"
		out.RawString(prefix[1:])
		out.String(string(in.CrossConnectID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LocalNSMMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LocalNSMMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LocalNSMMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LocalNSMMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm2(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(in *jlexer.Lexer, out *LocalConnectionMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "IP":
			out.IP = string(in.String())
		case "MechanismType":
			out.MechanismType = string(in.String())
		case "MechanismParameters":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.MechanismParameters = make(map[string]string)
				} else {
					out.MechanismParameters = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v5 string
					v5 = string(in.String())
					(out.MechanismParameters)[key] = v5
					in.WantComma()
				}
				in.Delim('}')
			}
		case "Labels":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Labels = make(map[string]string)
				} else {
					out.Labels = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v6 string
					v6 = string(in.String())
					(out.Labels)[key] = v6
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(out *jwriter.Writer, in LocalConnectionMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"IP\":"
		out.RawString(prefix[1:])
		out.String(string(in.IP))
	}
	{
		const prefix string = ",\"MechanismType\":"
		out.RawString(prefix)
		out.String(string(in.MechanismType))
	}
	{
		const prefix string = ",\"MechanismParameters\":"
		out.RawString(prefix)
		if in.MechanismParameters == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v7First := true
			for v7Name, v7Value := range in.MechanismParameters {
				if v7First {
					v7First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v7Name))
				out.RawByte(':')
				out.String(string(v7Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"Labels\":"
		out.RawString(prefix)
		if in.Labels == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v8First := true
			for v8Name, v8Value := range in.Labels {
				if v8First {
					v8First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v8Name))
				out.RawByte(':')
				out.String(string(v8Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LocalConnectionMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LocalConnectionMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LocalConnectionMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LocalConnectionMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm3(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(in *jlexer.Lexer, out *EdgeMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "SourceCrossConnectID":
			out.SourceCrossConnectID = string(in.String())
		case "DestinationCrossConnectID":
			out.DestinationCrossConnectID = string(in.String())
		case "Via":
			(out.Via).UnmarshalEasyJSON(in)
		case "CrossConnectID":
			out.CrossConnectID = string(in.String())
		case "NetworkService":
			out.NetworkService = string(in.String())
		case "Payload":
			out.Payload = string(in.String())
		case "Source":
			(out.Source).UnmarshalEasyJSON(in)
		case "Destination":
			(out.Destination).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(out *jwriter.Writer, in EdgeMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"SourceCrossConnectID\":"
		out.RawString(prefix[1:])
		out.String(string(in.SourceCrossConnectID))
	}
	{
		const prefix string = ",\"DestinationCrossConnectID\":"
		out.RawString(prefix)
		out.String(string(in.DestinationCrossConnectID))
	}
	{
		const prefix string = ",\"Via\":"
		out.RawString(prefix)
		(in.Via).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"CrossConnectID\":"
		out.RawString(prefix)
		out.String(string(in.CrossConnectID))
	}
	{
		const prefix string = ",\"NetworkService\":"
		out.RawString(prefix)
		out.String(string(in.NetworkService))
	}
	{
		const prefix string = ",\"Payload\":"
		out.RawString(prefix)
		out.String(string(in.Payload))
	}
	{
		const prefix string = ",\"Source\":"
		out.RawString(prefix)
		(in.Source).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"Destination\":"
		out.RawString(prefix)
		(in.Destination).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EdgeMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EdgeMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EdgeMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EdgeMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm4(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(in *jlexer.Lexer, out *BaseNSMMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "NetworkService":
			out.NetworkService = string(in.String())
		case "Payload":
			out.Payload = string(in.String())
		case "Source":
			(out.Source).UnmarshalEasyJSON(in)
		case "Destination":
			(out.Destination).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(out *jwriter.Writer, in BaseNSMMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"NetworkService\":"
		out.RawString(prefix[1:])
		out.String(string(in.NetworkService))
	}
	{
		const prefix string = ",\"Payload\":"
		out.RawString(prefix)
		out.String(string(in.Payload))
	}
	{
		const prefix string = ",\"Source\":"
		out.RawString(prefix)
		(in.Source).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"Destination\":"
		out.RawString(prefix)
		(in.Destination).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BaseNSMMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BaseNSMMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BaseNSMMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BaseNSMMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm5(l, v)
}
func easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(in *jlexer.Lexer, out *BaseConnectionMetadata) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "MechanismType":
			out.MechanismType = string(in.String())
		case "MechanismParameters":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.MechanismParameters = make(map[string]string)
				} else {
					out.MechanismParameters = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v9 string
					v9 = string(in.String())
					(out.MechanismParameters)[key] = v9
					in.WantComma()
				}
				in.Delim('}')
			}
		case "Labels":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Labels = make(map[string]string)
				} else {
					out.Labels = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v10 string
					v10 = string(in.String())
					(out.Labels)[key] = v10
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(out *jwriter.Writer, in BaseConnectionMetadata) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"MechanismType\":"
		out.RawString(prefix[1:])
		out.String(string(in.MechanismType))
	}
	{
		const prefix string = ",\"MechanismParameters\":"
		out.RawString(prefix)
		if in.MechanismParameters == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v11First := true
			for v11Name, v11Value := range in.MechanismParameters {
				if v11First {
					v11First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v11Name))
				out.RawByte(':')
				out.String(string(v11Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"Labels\":"
		out.RawString(prefix)
		if in.Labels == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v12First := true
			for v12Name, v12Value := range in.Labels {
				if v12First {
					v12First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v12Name))
				out.RawByte(':')
				out.String(string(v12Value))
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BaseConnectionMetadata) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BaseConnectionMetadata) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa0ee0e3EncodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BaseConnectionMetadata) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BaseConnectionMetadata) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa0ee0e3DecodeGithubComSkydiveProjectSkydiveTopologyProbesNsm6(l, v)
}