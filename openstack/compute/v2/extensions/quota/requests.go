package quotasets

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

// CreateOptsExt adds a Quota option to the base CreateOpts.
type CreateOptsExt struct {
	servers.CreateOptsBuilder
	KeyName string `json:"key_name,omitempty"`
}

// ToServerCreateMap adds the key_name and, optionally, key_data options to
// the base server creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	if opts.KeyName == "" {
		return base, nil
	}

	serverMap := base["server"].(map[string]interface{})
	serverMap["key_name"] = opts.KeyName

	return base, nil
}

// List returns a Pager that allows you to iterate over a collection of Quotas.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return QuotaSetPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call. Notable, the
// CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	ToQuotaCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies quota creation or import parameters.
type CreateOpts struct {
	// Name [required] is a friendly name to refer to this Quota in other services.
	Name string

	// PublicKey [optional] is a pregenerated OpenSSH-formatted public key. If provided, this key
	// will be imported and no new key will be created.
	PublicKey string
}

// ToQuotaCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToQuotaCreateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		return nil, errors.New("Missing field required for quota creation: Name")
	}

	quota := make(map[string]interface{})
	quota["name"] = opts.Name
	if opts.PublicKey != "" {
		quota["public_key"] = opts.PublicKey
	}

	return map[string]interface{}{"quota": quota}, nil
}

// Create requests the creation of a new quota on the server, or to import a pre-existing
// quota.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToQuotaCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(createURL(client), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Get returns public data about a previously uploaded Quota.
func Get(client *gophercloud.ServiceClient, name string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getURL(client, name), &res.Body, nil)
	return res
}

// Delete requests the deletion of a previous stored Quota from the server.
func Delete(client *gophercloud.ServiceClient, name string) DeleteResult {
	var res DeleteResult
	_, res.Err = client.Delete(deleteURL(client, name), nil)
	return res
}
