// Code generated by ogen, DO NOT EDIT.

package apiv2

import (
	"io"

	"github.com/go-faster/jx"
)

// Ref: #/components/schemas/Error
type Error map[string]jx.Raw

func (s *Error) init() Error {
	m := *s
	if m == nil {
		m = map[string]jx.Raw{}
		*s = m
	}
	return m
}

type GetConfigurationBadGateway struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s GetConfigurationBadGateway) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

// GetConfigurationBadGatewayHeaders wraps GetConfigurationBadGateway with response headers.
type GetConfigurationBadGatewayHeaders struct {
	ContentType string
	Response    GetConfigurationBadGateway
}

// GetContentType returns the value of ContentType.
func (s *GetConfigurationBadGatewayHeaders) GetContentType() string {
	return s.ContentType
}

// GetResponse returns the value of Response.
func (s *GetConfigurationBadGatewayHeaders) GetResponse() GetConfigurationBadGateway {
	return s.Response
}

// SetContentType sets the value of ContentType.
func (s *GetConfigurationBadGatewayHeaders) SetContentType(val string) {
	s.ContentType = val
}

// SetResponse sets the value of Response.
func (s *GetConfigurationBadGatewayHeaders) SetResponse(val GetConfigurationBadGateway) {
	s.Response = val
}

func (*GetConfigurationBadGatewayHeaders) getConfigurationRes() {}

type GetConfigurationBadRequest Error

func (*GetConfigurationBadRequest) getConfigurationRes() {}

type GetConfigurationGatewayTimeout struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s GetConfigurationGatewayTimeout) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

// GetConfigurationGatewayTimeoutHeaders wraps GetConfigurationGatewayTimeout with response headers.
type GetConfigurationGatewayTimeoutHeaders struct {
	ContentType string
	Response    GetConfigurationGatewayTimeout
}

// GetContentType returns the value of ContentType.
func (s *GetConfigurationGatewayTimeoutHeaders) GetContentType() string {
	return s.ContentType
}

// GetResponse returns the value of Response.
func (s *GetConfigurationGatewayTimeoutHeaders) GetResponse() GetConfigurationGatewayTimeout {
	return s.Response
}

// SetContentType sets the value of ContentType.
func (s *GetConfigurationGatewayTimeoutHeaders) SetContentType(val string) {
	s.ContentType = val
}

// SetResponse sets the value of Response.
func (s *GetConfigurationGatewayTimeoutHeaders) SetResponse(val GetConfigurationGatewayTimeout) {
	s.Response = val
}

func (*GetConfigurationGatewayTimeoutHeaders) getConfigurationRes() {}

type GetConfigurationInternalServerError struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s GetConfigurationInternalServerError) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

// GetConfigurationInternalServerErrorHeaders wraps GetConfigurationInternalServerError with response headers.
type GetConfigurationInternalServerErrorHeaders struct {
	ContentType string
	Response    GetConfigurationInternalServerError
}

// GetContentType returns the value of ContentType.
func (s *GetConfigurationInternalServerErrorHeaders) GetContentType() string {
	return s.ContentType
}

// GetResponse returns the value of Response.
func (s *GetConfigurationInternalServerErrorHeaders) GetResponse() GetConfigurationInternalServerError {
	return s.Response
}

// SetContentType sets the value of ContentType.
func (s *GetConfigurationInternalServerErrorHeaders) SetContentType(val string) {
	s.ContentType = val
}

// SetResponse sets the value of Response.
func (s *GetConfigurationInternalServerErrorHeaders) SetResponse(val GetConfigurationInternalServerError) {
	s.Response = val
}

func (*GetConfigurationInternalServerErrorHeaders) getConfigurationRes() {}

type GetConfigurationNotFound Error

func (*GetConfigurationNotFound) getConfigurationRes() {}

type GetConfigurationOK struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s GetConfigurationOK) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

// GetConfigurationOKHeaders wraps GetConfigurationOK with response headers.
type GetConfigurationOKHeaders struct {
	ConfigurationVersion OptString
	ContentType          string
	Response             GetConfigurationOK
}

// GetConfigurationVersion returns the value of ConfigurationVersion.
func (s *GetConfigurationOKHeaders) GetConfigurationVersion() OptString {
	return s.ConfigurationVersion
}

// GetContentType returns the value of ContentType.
func (s *GetConfigurationOKHeaders) GetContentType() string {
	return s.ContentType
}

// GetResponse returns the value of Response.
func (s *GetConfigurationOKHeaders) GetResponse() GetConfigurationOK {
	return s.Response
}

// SetConfigurationVersion sets the value of ConfigurationVersion.
func (s *GetConfigurationOKHeaders) SetConfigurationVersion(val OptString) {
	s.ConfigurationVersion = val
}

// SetContentType sets the value of ContentType.
func (s *GetConfigurationOKHeaders) SetContentType(val string) {
	s.ContentType = val
}

// SetResponse sets the value of Response.
func (s *GetConfigurationOKHeaders) SetResponse(val GetConfigurationOK) {
	s.Response = val
}

func (*GetConfigurationOKHeaders) getConfigurationRes() {}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}
