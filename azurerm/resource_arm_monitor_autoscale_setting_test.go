package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMMonitorAutoScaleSetting_basic(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMonitorAutoScaleSetting_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "tags.$type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorAutoScaleSetting_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_monitor_autoscale_setting"),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleProfiles(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMonitorAutoScaleSetting_multipleProfiles(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "primary"),
					resource.TestCheckResourceAttr(resourceName, "profile.1.name", "secondary"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_update(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(ri, location, 1, 3, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.minimum", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.maximum", "3"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.default", "2"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(ri, location, 2, 4, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.minimum", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.maximum", "4"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.default", "3"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_capacity(ri, location, 2, 45, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.minimum", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.maximum", "45"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.capacity.0.default", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_multipleRules(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.0.scale_action.0.direction", "Increase"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_multipleRules(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "metricRules"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.0.scale_action.0.direction", "Increase"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.rule.1.scale_action.0.direction", "Decrease"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_customEmails(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_email(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.0.custom_emails.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.0.custom_emails.0", fmt.Sprintf("acctest1-%d@example.com", ri)),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_emailUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.0.custom_emails.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.0.custom_emails.0", fmt.Sprintf("acctest1-%d@example.com", ri)),
					resource.TestCheckResourceAttr(resourceName, "notification.0.email.0.custom_emails.1", fmt.Sprintf("acctest2-%d@example.com", ri)),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_recurrence(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMonitorAutoScaleSetting_recurrence(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "recurrence"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_recurrenceUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_recurrence(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.1", "Wednesday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.2", "Friday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.hours.0", "18"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.minutes.0", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorAutoScaleSetting_recurrenceUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.0", "Monday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.1", "Tuesday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.days.2", "Wednesday"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.hours.0", "20"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.recurrence.0.minutes.0", "15"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorAutoScaleSetting_fixedDate(t *testing.T) {
	resourceName := "azurerm_monitor_autoscale_setting.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	config := testAccAzureRMMonitorAutoScaleSetting_fixedDate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorAutoScaleSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorAutoScaleSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "profile.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.name", "fixedDate"),
					resource.TestCheckResourceAttr(resourceName, "profile.0.fixed_date.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMMonitorAutoScaleSettingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		autoscaleSettingName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Monitor AutoScale Setting: %s", autoscaleSettingName)
		}

		conn := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, autoscaleSettingName)
		if err != nil {
			return fmt.Errorf("Bad: Get on Monitor AutoScale Setting: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: AutoScale Setting %q (Resource Group: %q) does not exist", autoscaleSettingName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMMonitorAutoScaleSettingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).autoscaleSettingsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_autoscale_setting" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("AutoScale Setting still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMonitorAutoScaleSetting_basic(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "import" {
  name                = "${azurerm_monitor_autoscale_setting.test.name}"
  resource_group_name = "${azurerm_monitor_autoscale_setting.test.resource_group_name}"
  location            = "${azurerm_monitor_autoscale_setting.test.location}"
  target_resource_id  = "${azurerm_monitor_autoscale_setting.test.target_resource_id}"

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleProfiles(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "primary"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  profile {
    name = "secondary"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Wednesday",
        "Friday",
      ]

      hours   = [18]
      minutes = [0]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_capacity(rInt int, location string, min int, max int, defaultVal int) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  enabled             = false

  profile {
    name = "metricRules"

    capacity {
      default = %d
      minimum = %d
      maximum = %d
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, rInt, defaultVal, min, max)
}

func testAccAzureRMMonitorAutoScaleSetting_multipleRules(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"
  enabled             = true

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "LessThan"
        threshold          = 25
      }

      scale_action {
        direction = "Decrease"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_email(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
      custom_emails                         = ["acctest1-%d@example.com"]
    }
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_emailUpdated(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "metricRules"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    rule {
      metric_trigger {
        metric_name        = "Percentage CPU"
        metric_resource_id = "${azurerm_virtual_machine_scale_set.test.id}"
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
      custom_emails                         = ["acctest1-%d@example.com", "acctest2-%d@example.com"]
    }
  }
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_recurrence(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "recurrence"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Wednesday",
        "Friday",
      ]

      hours   = [18]
      minutes = [0]
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_recurrenceUpdated(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "recurrence"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    recurrence {
      timezone = "Pacific Standard Time"

      days = [
        "Monday",
        "Tuesday",
        "Wednesday",
      ]

      hours   = [20]
      minutes = [15]
    }
  }

  notification {
    email {
      send_to_subscription_administrator    = false
      send_to_subscription_co_administrator = false
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_fixedDate(rInt int, location string) string {
	template := testAccAzureRMMonitorAutoScaleSetting_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  target_resource_id  = "${azurerm_virtual_machine_scale_set.test.id}"

  profile {
    name = "fixedDate"

    capacity {
      default = 1
      minimum = 1
      maximum = 10
    }

    fixed_date {
      timezone = "Pacific Standard Time"
      start    = "2020-06-18T00:00:00Z"
      end      = "2020-06-18T23:59:59Z"
    }
  }
}
`, template, rInt)
}

func testAccAzureRMMonitorAutoScaleSetting_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 2
  }

  os_profile {
    computer_name_prefix = "testvm-%d"
    admin_username       = "myadmin"
    admin_password       = "Passwword1234"
  }

  network_profile {
    name    = "TestNetworkProfile-%d"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      subnet_id = "${azurerm_subnet.test.id}"
      primary   = true
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "StandardSSD_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
