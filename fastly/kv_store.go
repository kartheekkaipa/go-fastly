package fastly

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://developer.fastly.com/reference/api/services/resources/kv-store

// KVStore represents an KV Store response from the Fastly API.
type KVStore struct {
	CreatedAt *time.Time `mapstructure:"created_at"`
	ID        string     `mapstructure:"id"`
	Name      string     `mapstructure:"name"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// CreateKVStoreInput is used as an input to the CreateKVStore function.
type CreateKVStoreInput struct {
	// Name is the name of the store to create (required).
	Name string `json:"name"`
}

// CreateKVStore creates a new resource.
func (c *Client) CreateKVStore(i *CreateKVStoreInput) (*KVStore, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	const path = "/resources/stores/kv"
	resp, err := c.PostJSON(path, i, nil)
	if err != nil {
		return nil, err
	}

	var store *KVStore
	if err := decodeBodyMap(resp.Body, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// ListKVStoresInput is used as an input to the ListKVStores function.
type ListKVStoresInput struct {
	// Cursor is used for paginating through results.
	Cursor string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListKVStoresInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	if l.Limit == 0 && l.Cursor == "" {
		return nil
	}

	m := make(map[string]string)

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	return m
}

// ListKVStoresResponse retrieves all resources.
type ListKVStoresResponse struct {
	// Data is the list of returned kv stores
	Data []KVStore
	// Meta is the information for pagination
	Meta map[string]string
}

// ListKVStores retrieves all resources.
func (c *Client) ListKVStores(i *ListKVStoresInput) (*ListKVStoresResponse, error) {
	const path = "/resources/stores/kv"

	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var output *ListKVStoresResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListKVStoresPaginator is the opaque type for a ListKVStores call with pagination.
type ListKVStoresPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListKVStoresInput
	stores   []KVStore // stored response from previous api call
}

// NewListKVStoresPaginator creates a new paginator for the given ListKVStoresInput.
func (c *Client) NewListKVStoresPaginator(i *ListKVStoresInput) *ListKVStoresPaginator {
	return &ListKVStoresPaginator{
		client: c,
		input:  i,
	}
}

// Next advances the paginator and fetches the next set of kv stores.
func (l *ListKVStoresPaginator) Next() bool {
	if l.finished {
		l.stores = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListKVStores(l.input)
	if err != nil {
		l.err = err
		l.finished = true
	}

	l.stores = o.Data
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Stores returns the current partial list of kv stores.
func (l *ListKVStoresPaginator) Stores() []KVStore {
	return l.stores
}

// Err returns any error from the pagination.
func (l *ListKVStoresPaginator) Err() error {
	return l.err
}

// GetKVStoreInput is the input to the GetKVStore function.
type GetKVStoreInput struct {
	// ID is the ID of the store to fetch (required).
	ID string
}

// GetKVStore retrieves the specified resource.
func (c *Client) GetKVStore(i *GetKVStoreInput) (*KVStore, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := "/resources/stores/kv/" + i.ID
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var output *KVStore
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// DeleteKVStoreInput is the input to the DeleteKVStore function.
type DeleteKVStoreInput struct {
	// ID is the ID of the kv store to delete (required).
	ID string
}

// DeleteKVStore deletes the specified resource.
func (c *Client) DeleteKVStore(i *DeleteKVStoreInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := "/resources/stores/kv/" + i.ID
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}

// ListKVStoreKeysInput is the input to the ListKVStoreKeys function.
type ListKVStoreKeysInput struct {
	// Cursor is used for paginating through results.
	Cursor string
	// ID is the ID of the kv store to list keys for (required).
	ID string
	// Limit is the maximum number of items included the response.
	Limit int
}

func (l *ListKVStoreKeysInput) formatFilters() map[string]string {
	if l == nil {
		return nil
	}

	if l.Limit == 0 && l.Cursor == "" {
		return nil
	}

	m := make(map[string]string)

	if l.Limit != 0 {
		m["limit"] = strconv.Itoa(l.Limit)
	}

	if l.Cursor != "" {
		m["cursor"] = l.Cursor
	}

	return m
}

// ListKVStoreKeysResponse retrieves all resources.
type ListKVStoreKeysResponse struct {
	// Data is the list of keys
	Data []string
	// Meta is the information for pagination
	Meta map[string]string
}

// ListKVStoreKeys retrieves all resources.
func (c *Client) ListKVStoreKeys(i *ListKVStoreKeysInput) (*ListKVStoreKeysResponse, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := "/resources/stores/kv/" + i.ID + "/keys"
	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var output *ListKVStoreKeysResponse
	if err := decodeBodyMap(resp.Body, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// ListKVStoreKeysPaginator is the opaque type for a ListKVStoreKeys calls with pagination.
type ListKVStoreKeysPaginator struct {
	client   *Client
	cursor   string // == "" if no more pages
	err      error
	finished bool
	input    *ListKVStoreKeysInput
	keys     []string // stored response from previous api call
}

// NewListKVStoreKeysPaginator returns a new paginator for the provided LitKVStoreKeysInput.
func (c *Client) NewListKVStoreKeysPaginator(i *ListKVStoreKeysInput) *ListKVStoreKeysPaginator {
	return &ListKVStoreKeysPaginator{
		client: c,
		input:  i,
	}
}

// Next advanced the paginator.
func (l *ListKVStoreKeysPaginator) Next() bool {
	if l.finished {
		l.keys = nil
		return false
	}

	l.input.Cursor = l.cursor
	o, err := l.client.ListKVStoreKeys(l.input)
	if err != nil {
		l.err = err
		l.finished = true
	}

	l.keys = o.Data
	if next := o.Meta["next_cursor"]; next == "" {
		l.finished = true
	} else {
		l.cursor = next
	}

	return true
}

// Err returns any error from the paginator.
func (l *ListKVStoreKeysPaginator) Err() error {
	return l.err
}

// Keys returns the current set of keys retrieved by the paginator.
func (l *ListKVStoreKeysPaginator) Keys() []string {
	return l.keys
}

// GetKVStoreKeyInput is the input to the GetKVStoreKey function.
type GetKVStoreKeyInput struct {
	// ID is the ID of the kv store (required).
	ID string
	// Key is the key to fetch (required).
	Key string
}

// GetKVStoreKey retrieves the specified resource.
func (c *Client) GetKVStoreKey(i *GetKVStoreKeyInput) (string, error) {
	if i.ID == "" {
		return "", ErrMissingID
	}
	if i.Key == "" {
		return "", ErrMissingKey
	}

	path := "/resources/stores/kv/" + i.ID + "/keys/" + i.Key
	resp, err := c.Get(path, nil)
	if err != nil {
		return "", err
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// InsertKVStoreKeyInput is the input to the InsertKVStoreKey function.
type InsertKVStoreKeyInput struct {
	// ID is the ID of the kv store (required).
	ID string
	// Key is the key to add (required).
	Key string
	// Value is the value to insert (required).
	Value string
}

// InsertKVStoreKey inserts a key/value pair into an kv store.
func (c *Client) InsertKVStoreKey(i *InsertKVStoreKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}
	if i.Key == "" {
		return ErrMissingKey
	}

	path := "/resources/stores/kv/" + i.ID + "/keys/" + i.Key
	resp, err := c.Put(path, &RequestOptions{Body: io.NopCloser(strings.NewReader(i.Value))})
	if err != nil {
		return err
	}

	_, err = checkResp(resp, err)
	return err
}

// DeleteKVStoreKeyInput is the input to the DeleteKVStoreKey function.
type DeleteKVStoreKeyInput struct {
	// ID is the ID of the kv store (required).
	ID string
	// Key is the key to delete (required).
	Key string
}

// DeleteKVStoreKey deletes the specified resource.
func (c *Client) DeleteKVStoreKey(i *DeleteKVStoreKeyInput) error {
	if i.ID == "" {
		return ErrMissingID
	}
	if i.Key == "" {
		return ErrMissingKey
	}

	path := "/resources/stores/kv/" + i.ID + "/keys/" + i.Key
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}