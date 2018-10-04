// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	compute "google.golang.org/api/compute/v1"
)

// Whether the IP CIDR change shrinks the block.
func isShrinkageIpCidr(old, new, _ interface{}) bool {
	_, oldCidr, oldErr := net.ParseCIDR(old.(string))
	_, newCidr, newErr := net.ParseCIDR(new.(string))

	if oldErr != nil || newErr != nil {
		// This should never happen. The ValidateFunc on the field ensures it.
		return false
	}

	oldStart, oldEnd := cidr.AddressRange(oldCidr)

	if newCidr.Contains(oldStart) && newCidr.Contains(oldEnd) {
		// This is a CIDR range expansion, no need to ForceNew, we have an update method for it.
		return false
	}

	return true
}

func splitSubnetID(id string) (region string, name string) {
	parts := strings.Split(id, "/")
	region = parts[0]
	name = parts[1]
	return
}

func resourceComputeSubnetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSubnetworkCreate,
		Read:   resourceComputeSubnetworkRead,
		Update: resourceComputeSubnetworkUpdate,
		Delete: resourceComputeSubnetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeSubnetworkImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(360 * time.Second),
			Update: schema.DefaultTimeout(360 * time.Second),
			Delete: schema.DefaultTimeout(360 * time.Second),
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("ip_cidr_range", isShrinkageIpCidr),
			resourceComputeSubnetworkSecondaryIpRangeSetStyleDiff,
		),

		Schema: map[string]*schema.Schema{
			"ip_cidr_range": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIpCidrRange,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_flow_logs": {
				Type:     schema.TypeBool,
				Optional: true,
				Deprecated: `This field is in beta and will be removed from this provider.
Use the terraform-provider-google-beta provider to continue using it.
See https://terraform.io/docs/provider/google/provider_versions.html for more details on beta fields.`,
			},
			"private_ip_google_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"secondary_ip_range": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Deprecated: `This field is in beta and will be removed from this provider.
Use the terraform-provider-google-beta provider to continue using it.
See https://terraform.io/docs/provider/google/provider_versions.html for more details on beta fields.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_cidr_range": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIpCidrRange,
						},
						"range_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateGCPName,
						},
					},
				},
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
				Deprecated: `This field is in beta and will be removed from this provider.
Use the terraform-provider-google-beta provider to continue using it.
See https://terraform.io/docs/provider/google/provider_versions.html for more details on beta fields.`,
			},
			"gateway_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceComputeSubnetworkSecondaryIpRangeSetStyleDiff(diff *schema.ResourceDiff, meta interface{}) error {
	keys := diff.GetChangedKeysPrefix("secondary_ip_range")
	if len(keys) == 0 {
		return nil
	}
	oldCount, newCount := diff.GetChange("secondary_ip_range.#")
	var count int
	// There could be duplicates - worth continuing even if the counts are unequal.
	if oldCount.(int) < newCount.(int) {
		count = newCount.(int)
	} else {
		count = oldCount.(int)
	}

	if count < 1 {
		return nil
	}
	old := make([]interface{}, count)
	new := make([]interface{}, count)
	for i := 0; i < count; i++ {
		o, n := diff.GetChange(fmt.Sprintf("secondary_ip_range.%d", i))

		if o != nil {
			old = append(old, o)
		}
		if n != nil {
			new = append(new, n)
		}
	}

	oldSet := schema.NewSet(schema.HashResource(resourceComputeSubnetwork().Schema["secondary_ip_range"].Elem.(*schema.Resource)), old)
	newSet := schema.NewSet(schema.HashResource(resourceComputeSubnetwork().Schema["secondary_ip_range"].Elem.(*schema.Resource)), new)

	if oldSet.Equal(newSet) {
		if err := diff.Clear("secondary_ip_range"); err != nil {
			return err
		}
	}

	return nil
}

func resourceComputeSubnetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeSubnetworkDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	ipCidrRangeProp, err := expandComputeSubnetworkIpCidrRange(d.Get("ip_cidr_range"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ip_cidr_range"); !isEmptyValue(reflect.ValueOf(ipCidrRangeProp)) && (ok || !reflect.DeepEqual(v, ipCidrRangeProp)) {
		obj["ipCidrRange"] = ipCidrRangeProp
	}
	nameProp, err := expandComputeSubnetworkName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	networkProp, err := expandComputeSubnetworkNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("network"); !isEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
		obj["network"] = networkProp
	}
	enableFlowLogsProp, err := expandComputeSubnetworkEnableFlowLogs(d.Get("enable_flow_logs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_flow_logs"); ok || !reflect.DeepEqual(v, enableFlowLogsProp) {
		obj["enableFlowLogs"] = enableFlowLogsProp
	}
	secondaryIpRangesProp, err := expandComputeSubnetworkSecondaryIpRange(d.Get("secondary_ip_range"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("secondary_ip_range"); !isEmptyValue(reflect.ValueOf(secondaryIpRangesProp)) && (ok || !reflect.DeepEqual(v, secondaryIpRangesProp)) {
		obj["secondaryIpRanges"] = secondaryIpRangesProp
	}
	privateIpGoogleAccessProp, err := expandComputeSubnetworkPrivateIpGoogleAccess(d.Get("private_ip_google_access"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("private_ip_google_access"); !isEmptyValue(reflect.ValueOf(privateIpGoogleAccessProp)) && (ok || !reflect.DeepEqual(v, privateIpGoogleAccessProp)) {
		obj["privateIpGoogleAccess"] = privateIpGoogleAccessProp
	}
	regionProp, err := expandComputeSubnetworkRegion(d.Get("region"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("region"); !isEmptyValue(reflect.ValueOf(regionProp)) && (ok || !reflect.DeepEqual(v, regionProp)) {
		obj["region"] = regionProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Subnetwork: %#v", obj)
	res, err := sendRequest(config, "POST", url, obj)
	if err != nil {
		return fmt.Errorf("Error creating Subnetwork: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{region}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	waitErr := computeOperationWaitTime(
		config.clientCompute, op, project, "Creating Subnetwork",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Subnetwork: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Subnetwork %q: %#v", d.Id(), res)

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeSubnetwork %q", d.Id()))
	}

	if err := d.Set("creation_timestamp", flattenComputeSubnetworkCreationTimestamp(res["creationTimestamp"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("description", flattenComputeSubnetworkDescription(res["description"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("gateway_address", flattenComputeSubnetworkGatewayAddress(res["gatewayAddress"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("ip_cidr_range", flattenComputeSubnetworkIpCidrRange(res["ipCidrRange"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("name", flattenComputeSubnetworkName(res["name"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("network", flattenComputeSubnetworkNetwork(res["network"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("enable_flow_logs", flattenComputeSubnetworkEnableFlowLogs(res["enableFlowLogs"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("fingerprint", flattenComputeSubnetworkFingerprint(res["fingerprint"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("secondary_ip_range", flattenComputeSubnetworkSecondaryIpRange(res["secondaryIpRanges"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("private_ip_google_access", flattenComputeSubnetworkPrivateIpGoogleAccess(res["privateIpGoogleAccess"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("region", flattenComputeSubnetworkRegion(res["region"])); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Subnetwork: %s", err)
	}

	return nil
}

func resourceComputeSubnetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("ip_cidr_range") {
		obj := make(map[string]interface{})
		ipCidrRangeProp, err := expandComputeSubnetworkIpCidrRange(d.Get("ip_cidr_range"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("ip_cidr_range"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, ipCidrRangeProp)) {
			obj["ipCidrRange"] = ipCidrRangeProp
		}

		url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks/{{name}}/expandIpCidrRange")
		if err != nil {
			return err
		}
		res, err := sendRequest(config, "POST", url, obj)
		if err != nil {
			return fmt.Errorf("Error updating Subnetwork %q: %s", d.Id(), err)
		}

		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		op := &compute.Operation{}
		err = Convert(res, op)
		if err != nil {
			return err
		}

		err = computeOperationWaitTime(
			config.clientCompute, op, project, "Updating Subnetwork",
			int(d.Timeout(schema.TimeoutUpdate).Minutes()))

		if err != nil {
			return err
		}

		d.SetPartial("ip_cidr_range")
	}
	if d.HasChange("enable_flow_logs") || d.HasChange("fingerprint") || d.HasChange("secondary_ip_range") {
		obj := make(map[string]interface{})
		enableFlowLogsProp, err := expandComputeSubnetworkEnableFlowLogs(d.Get("enable_flow_logs"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("enable_flow_logs"); ok || !reflect.DeepEqual(v, enableFlowLogsProp) {
			obj["enableFlowLogs"] = enableFlowLogsProp
		}
		fingerprintProp := d.Get("fingerprint")
		obj["fingerprint"] = fingerprintProp
		secondaryIpRangesProp, err := expandComputeSubnetworkSecondaryIpRange(d.Get("secondary_ip_range"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("secondary_ip_range"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, secondaryIpRangesProp)) {
			obj["secondaryIpRanges"] = secondaryIpRangesProp
		}

		url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks/{{name}}")
		if err != nil {
			return err
		}
		res, err := sendRequest(config, "PATCH", url, obj)
		if err != nil {
			return fmt.Errorf("Error updating Subnetwork %q: %s", d.Id(), err)
		}

		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		op := &compute.Operation{}
		err = Convert(res, op)
		if err != nil {
			return err
		}

		err = computeOperationWaitTime(
			config.clientCompute, op, project, "Updating Subnetwork",
			int(d.Timeout(schema.TimeoutUpdate).Minutes()))

		if err != nil {
			return err
		}

		d.SetPartial("enable_flow_logs")
		d.SetPartial("fingerprint")
		d.SetPartial("secondary_ip_range")
	}
	if d.HasChange("private_ip_google_access") {
		obj := make(map[string]interface{})
		privateIpGoogleAccessProp, err := expandComputeSubnetworkPrivateIpGoogleAccess(d.Get("private_ip_google_access"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("private_ip_google_access"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, privateIpGoogleAccessProp)) {
			obj["privateIpGoogleAccess"] = privateIpGoogleAccessProp
		}

		url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks/{{name}}/setPrivateIpGoogleAccess")
		if err != nil {
			return err
		}
		res, err := sendRequest(config, "POST", url, obj)
		if err != nil {
			return fmt.Errorf("Error updating Subnetwork %q: %s", d.Id(), err)
		}

		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		op := &compute.Operation{}
		err = Convert(res, op)
		if err != nil {
			return err
		}

		err = computeOperationWaitTime(
			config.clientCompute, op, project, "Updating Subnetwork",
			int(d.Timeout(schema.TimeoutUpdate).Minutes()))

		if err != nil {
			return err
		}

		d.SetPartial("private_ip_google_access")
	}

	d.Partial(false)

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/beta/projects/{{project}}/regions/{{region}}/subnetworks/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Subnetwork %q", d.Id())
	res, err := sendRequest(config, "DELETE", url, obj)
	if err != nil {
		return handleNotFoundError(err, d, "Subnetwork")
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Deleting Subnetwork",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Subnetwork %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeSubnetworkImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+)", "(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{region}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeSubnetworkCreationTimestamp(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkDescription(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkGatewayAddress(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkIpCidrRange(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkName(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkNetwork(v interface{}) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeSubnetworkEnableFlowLogs(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkFingerprint(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkSecondaryIpRange(v interface{}) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		transformed = append(transformed, map[string]interface{}{
			"range_name":    flattenComputeSubnetworkSecondaryIpRangeRangeName(original["rangeName"]),
			"ip_cidr_range": flattenComputeSubnetworkSecondaryIpRangeIpCidrRange(original["ipCidrRange"]),
		})
	}
	return transformed
}
func flattenComputeSubnetworkSecondaryIpRangeRangeName(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkSecondaryIpRangeIpCidrRange(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkPrivateIpGoogleAccess(v interface{}) interface{} {
	return v
}

func flattenComputeSubnetworkRegion(v interface{}) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func expandComputeSubnetworkDescription(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkIpCidrRange(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkNetwork(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("networks", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeSubnetworkEnableFlowLogs(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkSecondaryIpRange(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedRangeName, err := expandComputeSubnetworkSecondaryIpRangeRangeName(original["range_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRangeName); val.IsValid() && !isEmptyValue(val) {
			transformed["rangeName"] = transformedRangeName
		}

		transformedIpCidrRange, err := expandComputeSubnetworkSecondaryIpRangeIpCidrRange(original["ip_cidr_range"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIpCidrRange); val.IsValid() && !isEmptyValue(val) {
			transformed["ipCidrRange"] = transformedIpCidrRange
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeSubnetworkSecondaryIpRangeRangeName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkSecondaryIpRangeIpCidrRange(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkPrivateIpGoogleAccess(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeSubnetworkRegion(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("regions", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for region: %s", err)
	}
	return f.RelativeLink(), nil
}
