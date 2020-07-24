package limelight

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLimelightIPRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLimelightIPRangesRead,
		Schema: map[string]*schema.Schema{
			"ip_ranges": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLimelightIPRangesRead(d *schema.ResourceData, m interface{}) error {
	l := getConfigurationClient(m)
	log.Printf("[INFO] Fetching IP ranges")
	ipList, _, err := l.GetIPAllowList()
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(hashcode.String(strings.Join(ipList.IPRanges, ","))))
	d.Set("ip_ranges", ipList.IPRanges)
	d.Set("version", ipList.Version)
	return nil
}
