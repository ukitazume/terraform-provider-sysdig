package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

func resourceSysdigSecureNotificationChannelVictorOps() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelVictorOpsCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelVictorOpsUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelVictorOpsRead,
		DeleteContext: resourceSysdigSecureNotificationChannelVictorOpsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelVictorOpsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelVictorOpsFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))
	_ = d.Set("version", notificationChannel.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigSecureNotificationChannelVictorOpsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	err = secureNotificationChannelVictorOpsToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelVictorOpsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelVictorOpsFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelVictorOpsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteNotificationChannel(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func secureNotificationChannelVictorOpsFromResourceData(d *schema.ResourceData) (nc secure.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_VICTOROPS
	nc.Options.APIKey = d.Get("api_key").(string)
	nc.Options.RoutingKey = d.Get("routing_key").(string)
	return
}

func secureNotificationChannelVictorOpsToResourceData(nc *secure.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("api_key", nc.Options.APIKey)
	_ = d.Set("routing_key", nc.Options.RoutingKey)

	return
}
