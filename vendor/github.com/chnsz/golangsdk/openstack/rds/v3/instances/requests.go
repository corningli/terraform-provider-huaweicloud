package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	Name                string             `json:"name"  required:"true"`
	Datastore           *Datastore         `json:"datastore" required:"true"`
	Ha                  *Ha                `json:"ha,omitempty"`
	ConfigurationId     string             `json:"configuration_id,omitempty"`
	Port                string             `json:"port,omitempty"`
	Password            string             `json:"password,omitempty"`
	BackupStrategy      *BackupStrategy    `json:"backup_strategy,omitempty"`
	EnterpriseProjectId string             `json:"enterprise_project_id,omitempty"`
	DiskEncryptionId    string             `json:"disk_encryption_id,omitempty"`
	FlavorRef           string             `json:"flavor_ref" required:"true"`
	Volume              *Volume            `json:"volume" required:"true"`
	Region              string             `json:"region" required:"true"`
	AvailabilityZone    string             `json:"availability_zone" required:"true"`
	VpcId               string             `json:"vpc_id" required:"true"`
	SubnetId            string             `json:"subnet_id" required:"true"`
	SecurityGroupId     string             `json:"security_group_id" required:"true"`
	RestorePoint        *RestorePoint      `json:"restore_point,omitempty"`
	ChargeInfo          *ChargeInfo        `json:"charge_info,omitempty"`
	TimeZone            string             `json:"time_zone,omitempty"`
	DssPoolId           string             `json:"dsspool_id,omitempty"`
	FixedIp             string             `json:"data_vip,omitempty"`
	Collation           string             `json:"collation,omitempty"`
	UnchangeableParam   *UnchangeableParam `json:"unchangeable_param,omitempty"`
	Tags                []tags.ResourceTag `json:"tags,omitempty"`
}

type CreateReplicaOpts struct {
	Name                string      `json:"name"  required:"true"`
	ReplicaOfId         string      `json:"replica_of_id" required:"true"`
	EnterpriseProjectId string      `json:"enterprise_project_id,omitempty"`
	DiskEncryptionId    string      `json:"disk_encryption_id,omitempty"`
	FlavorRef           string      `json:"flavor_ref" required:"true"`
	Volume              *Volume     `json:"volume" required:"true"`
	Region              string      `json:"region,omitempty"`
	AvailabilityZone    string      `json:"availability_zone" required:"true"`
	ChargeInfo          *ChargeInfo `json:"charge_info,omitempty"`
}

type Datastore struct {
	Type            string `json:"type" required:"true"`
	Version         string `json:"version" required:"true"`
	CompleteVersion string `json:"complete_version,omitempty"`
}

type Ha struct {
	Mode            string `json:"mode" required:"true"`
	ReplicationMode string `json:"replication_mode,omitempty"`
}

type UnchangeableParam struct {
	LowerCaseTableNames string `json:"lower_case_table_names"`
}

type BackupStrategy struct {
	StartTime string `json:"start_time" required:"true"`
	KeepDays  int    `json:"keep_days,omitempty"`
}

type Volume struct {
	Type string `json:"type" required:"true"`
	Size int    `json:"size" required:"true"`
}

type RestorePoint struct {
	InstanceId   string            `json:"instance_id" required:"true"`
	Type         string            `json:"type" required:"true"`
	BackupId     string            `json:"backup_id,omitempty"`
	RestoreTime  string            `json:"restore_time,omitempty"`
	DatabaseName map[string]string `json:"database_name,omitempty"`
}

type ChargeInfo struct {
	ChargeMode  string `json:"charge_mode" required:"true"`
	PeriodType  string `json:"period_type,omitempty"`
	PeriodNum   int    `json:"period_num,omitempty"`
	IsAutoRenew bool   `json:"is_auto_renew,omitempty"`
	IsAutoPay   bool   `json:"is_auto_pay,omitempty"`
}

type CreateRdsBuilder interface {
	ToInstancesCreateMap() (map[string]interface{}, error)
}

type RestRootPasswordOpts struct {
	DbUserPwd string `json:"db_user_pwd" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func (opts CreateOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateRdsBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

func (opts CreateReplicaOpts) ToInstancesCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CreateReplica(client *golangsdk.ServiceClient, opts CreateRdsBuilder) (r CreateResult) {
	b, err := opts.ToInstancesCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})

	return
}

type DeleteOpts struct {
	InstanceId string `json:"instance_id" required:"true"`
}

type DeleteInstanceBuilder interface {
	ToInstancesDeleteMap() (map[string]interface{}, error)
}

func (opts DeleteOpts) ToInstancesDeleteMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Delete(client *golangsdk.ServiceClient, instanceId string) (r DeleteResult) {

	url := deleteURL(client, instanceId)

	_, r.Err = client.Delete(url, &golangsdk.RequestOpts{
		JSONResponse: &r.Body,
		MoreHeaders:  map[string]string{"Content-Type": "application/json"},
	})
	return
}

type ListOpts struct {
	Id            string `q:"id"`
	Name          string `q:"name"`
	Type          string `q:"type"`
	DataStoreType string `q:"datastore_type"`
	VpcId         string `q:"vpc_id"`
	SubnetId      string `q:"subnet_id"`
	Offset        int    `q:"offset"`
	Limit         int    `q:"limit"`
}

type ListRdsBuilder interface {
	ToRdsListDetailQuery() (string, error)
}

func (opts ListOpts) ToRdsListDetailQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListRdsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToRdsListDetailQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RdsPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}

type ActionInstanceBuilder interface {
	ToActionInstanceMap() (map[string]interface{}, error)
}

func toActionInstanceMap(opts ActionInstanceBuilder) (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

type RestartInstanceOpts struct {
	Restart string `json:"restart" required:"true"`
}

func (opts RestartInstanceOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

func Restart(client *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r RestartResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

type RenameInstanceOpts struct {
	Name string `json:"name" required:"true"`
}

func (opts RenameInstanceOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

func Rename(client *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r RenameResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, instanceId, "name"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

type SingleToHaRdsOpts struct {
	SingleToHa *SingleToHaRds `json:"single_to_ha" required:"true"`
}

type SingleToHaRds struct {
	AzCodeNewNode string `json:"az_code_new_node" required:"true"`
	Password      string `json:"password,omitempty"`
}

func (opts SingleToHaRdsOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

func SingleToHa(client *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r SingleToHaResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type SpecCode struct {
	Speccode  string `json:"spec_code" required:"true"`
	IsAutoPay bool   `json:"is_auto_pay,omitempty"`
}

type ResizeFlavorOpts struct {
	ResizeFlavor *SpecCode `json:"resize_flavor" required:"true"`
}

func (opts ResizeFlavorOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

func Resize(client *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ResizeFlavorResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type EnlargeVolumeOpts struct {
	EnlargeVolume *EnlargeVolumeSize `json:"enlarge_volume" required:"true"`
}

type EnlargeVolumeSize struct {
	Size int `json:"size" required:"true"`
}

func (opts EnlargeVolumeOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

func EnlargeVolume(client *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r EnlargeVolumeResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(updateURL(client, instanceId, "action"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

type DbErrorlogOpts struct {
	StartDate string `q:"start_date"`
	EndDate   string `q:"end_date"`
	Offset    string `q:"offset"`
	Limit     string `q:"limit"`
	Level     string `q:"level"`
}

type DbErrorlogBuilder interface {
	DbErrorlogQuery() (string, error)
}

func (opts DbErrorlogOpts) DbErrorlogQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func ListErrorLog(client *golangsdk.ServiceClient, opts DbErrorlogBuilder, instanceID string) pagination.Pager {
	url := updateURL(client, instanceID, "errorlog")
	if opts != nil {
		query, err := opts.DbErrorlogQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ErrorLogPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}

type DbSlowLogOpts struct {
	StartDate string `q:"start_date"`
	EndDate   string `q:"end_date"`
	Offset    string `q:"offset"`
	Limit     string `q:"limit"`
	Level     string `q:"level"`
}

type DbSlowLogBuilder interface {
	ToDbSlowLogListQuery() (string, error)
}

func (opts DbSlowLogOpts) ToDbSlowLogListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func ListSlowLog(client *golangsdk.ServiceClient, opts DbSlowLogBuilder, instanceID string) pagination.Pager {
	url := updateURL(client, instanceID, "slowlog")
	if opts != nil {
		query, err := opts.ToDbSlowLogListQuery()

		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ErrorLogPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}

type RDSJobOpts struct {
	JobID string `q:"id"`
}

type RDSJobBuilder interface {
	ToRDSJobQuery() (string, error)
}

func (opts RDSJobOpts) ToRDSJobQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func GetRDSJob(client *golangsdk.ServiceClient, opts RDSJobBuilder) (r RDSJobResult) {
	url := jobURL(client)
	if opts != nil {
		query, err := opts.ToRDSJobQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

func ListEngine(client *golangsdk.ServiceClient, dbName string) (*Engine, error) {
	var rst golangsdk.Result
	_, err := client.Get(engineURL(client, dbName), &rst.Body, nil)
	if err == nil {
		var s Engine
		err := rst.ExtractInto(&s)
		return &s, err
	}
	return nil, err
}

func RestRootPassword(c *golangsdk.ServiceClient, instanceID string, opts RestRootPasswordOpts) (*ErrorResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r ErrorResponse
	_, err = c.Post(resetRootPasswordURL(c, instanceID), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}

type ModifyConfigurationOpts struct {
	Values map[string]string `json:"values" required:"true"`
}

func ModifyConfiguration(c *golangsdk.ServiceClient, instanceID string, opts ModifyConfigurationOpts) (r ModifyConfigurationResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Put(configurationsURL(c, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

func GetConfigurations(c *golangsdk.ServiceClient, instanceID string) (r GetConfigurationResult) {
	_, r.Err = c.Get(configurationsURL(c, instanceID), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}

func RebootInstance(c *golangsdk.ServiceClient, instanceID string) (r RebootResult) {
	b, err := golangsdk.BuildRequestBody(struct{}{}, "restart")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(actionURL(c, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

var (
	enableAutoExpand  bool = true
	disableAutoExpand bool = false
)

// EnableAutoExpandOpts is the structure used to enable the volume automatic expansion of RDS instance.
type EnableAutoExpandOpts struct {
	// The instnace ID.
	InstanceId string `json:"-" required:"true"`
	// The upper limit of automatic expansion of storage, in GB.
	// This parameter is mandatory when switch_option is set to true.
	// The value ranges from 40 GB to 4,000 GB and must be no less than the current storage of the instance.
	LimitSize int `json:"limit_size" required:"true"`
	// The threshold to trigger automatic expansion.
	// If the available storage drops to this threshold or 10 GB, the automatic expansion is triggered.
	// This parameter is mandatory when switch_option is set to true.
	// The valid values are as follows:
	// + 10
	// + 15
	// + 20
	TriggerThreshold int `json:"trigger_threshold" required:"true"`
}

// autoExpandOpts is the structure used to configure the volume automatic expansion of RDS instance.
type autoExpandOpts struct {
	// Whether the auto-expansion is enabled.
	SwitchOption *bool `json:"switch_option" required:"true"`
	// The upper limit of automatic expansion of storage, in GB.
	// This parameter is mandatory when switch_option is set to true.
	// The value ranges from 40 GB to 4,000 GB and must be no less than the current storage of the instance.
	LimitSize int `json:"limit_size,omitempty"`
	// The threshold to trigger automatic expansion.
	// If the available storage drops to this threshold or 10 GB, the automatic expansion is triggered.
	// This parameter is mandatory when switch_option is set to true.
	// The valid values are as follows:
	// + 10
	// + 15
	// + 20
	TriggerThreshold int `json:"trigger_threshold,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// EnableAutoExpand is a method used to configure the volume automatic expansion of RDS instance.
func EnableAutoExpand(c *golangsdk.ServiceClient, opts EnableAutoExpandOpts) error {
	enableOpts := autoExpandOpts{
		SwitchOption:     &enableAutoExpand,
		LimitSize:        opts.LimitSize,
		TriggerThreshold: opts.TriggerThreshold,
	}
	b, err := golangsdk.BuildRequestBody(enableOpts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(autoExpandURL(c, opts.InstanceId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// DisableAutoExpand is a method used to remove the volume automatic expansion configuration of RDS instance.
func DisableAutoExpand(c *golangsdk.ServiceClient, instanceId string) error {
	autoExpandOpts := autoExpandOpts{
		SwitchOption: &disableAutoExpand,
	}
	b, err := golangsdk.BuildRequestBody(autoExpandOpts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(autoExpandURL(c, instanceId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// GetAutoExpand is a method used to obtain the automatic expansion configuarion of instance storage.
func GetAutoExpand(c *golangsdk.ServiceClient, instanceId string) (*AutoExpansion, error) {
	var r AutoExpansion
	_, err := c.Get(autoExpandURL(c, instanceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

type ModifyAliasOpts struct {
	Alias string `json:"alias,omitempty"`
}

func (opts ModifyAliasOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

// ModifyAlias is a method used to modify the alias.
func ModifyAlias(c *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ModifyAliasResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, instanceId, "alias"), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

type ModifyMaintainWindowOpts struct {
	StartTime string `json:"start_time" required:"true"`
	EndTime   string `json:"end_time" required:"true"`
}

func (opts ModifyMaintainWindowOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

// ModifyMaintainWindow is a method used to modify maintain window.
func ModifyMaintainWindow(c *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ModifyMaintainWindowResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, instanceId, "ops-window"), b, nil, &golangsdk.RequestOpts{})
	return
}

type ModifyReplicationModeOpts struct {
	Mode string `json:"mode" required:"true"`
}

func (opts ModifyReplicationModeOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

// ModifyReplicationMode is a method used to modify replication mode.
func ModifyReplicationMode(c *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ModifyReplicationModeResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, instanceId, "failover/mode"), b, &r.Body, &golangsdk.RequestOpts{})
	return
}

type ModifySwitchStrategyOpts struct {
	RepairStrategy string `json:"repairStrategy" required:"true"`
}

func (opts ModifySwitchStrategyOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

// ModifySwitchStrategy is a method used to modify replication mode.
func ModifySwitchStrategy(c *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ModifySwitchStrategyResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, instanceId, "failover/strategy"), b, &r.Body, &golangsdk.RequestOpts{})
	return
}

type ModifyCollationOpts struct {
	Collation string `json:"collation" required:"true"`
}

func (opts ModifyCollationOpts) ToActionInstanceMap() (map[string]interface{}, error) {
	return toActionInstanceMap(opts)
}

// ModifyCollation is a method used to modify collation.
func ModifyCollation(c *golangsdk.ServiceClient, opts ActionInstanceBuilder, instanceId string) (r ModifyCollationResult) {
	b, err := opts.ToActionInstanceMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, instanceId, "collations"), b, &r.Body, &golangsdk.RequestOpts{})
	return
}
