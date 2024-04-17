// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/project-go-sdk/projectv1"
)

func DataSourceIbmProjectConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmProjectConfigRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique project ID.",
			},
			"project_config_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique configuration ID.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the configuration.",
			},
			"is_draft": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag that indicates whether the version of the configuration is draft, or active.",
			},
			"needs_attention_state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The needs attention state of a configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the event.",
						},
						"event": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the event.",
						},
						"severity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The severity of the event. This is a system generated field. For user triggered events the field is not present.",
						},
						"action_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An actionable URL that users can access in response to the event. This is a system generated field. For user triggered events the field is not present.",
						},
						"target": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration id and version for which the event occurred. This field is only available for user generated events. For system triggered events the field is not present.",
						},
						"triggered_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM id of the user that triggered the event. This field is only available for user generated events. For system triggered events the field is not present.",
						},
						"timestamp": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timestamp of the event.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"last_saved_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"outputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The outputs of a Schematics template property.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The variable name.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A short explanation of the output value.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeMap,
							Deprecated:  "This property will be deprecated, the new property will be of type String.",
							Computed:    true,
							Description: "This property can be any value - a string, number, boolean, array, or object.",
						},
						"value_json": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This property can be any value - a string, number, boolean, array, or object.",
						},
					},
				},
			},
			"project": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project that is referenced by this resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The definition of the project reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the project.",
									},
								},
							},
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An IBM Cloud resource name that uniquely identifies a resource.",
						},
					},
				},
			},
			"schematics": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A Schematics workspace that is associated to a project configuration, with scripts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An IBM Cloud resource name that uniquely identifies a resource.",
						},
						"validate_pre_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
						"validate_post_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
						"deploy_pre_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
						"deploy_post_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
						"undeploy_pre_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
						"undeploy_post_script": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A script to be run as part of a project configuration for a specific stage (pre or post) and action (validate, deploy, or undeploy).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the script.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to this script is within the current version source.",
									},
									"short_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The short description for this script.",
									},
								},
							},
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the configuration.",
			},
			"update_available": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag that indicates whether a configuration update is available.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"definition": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compliance_profile": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile that is required for compliance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID for the compliance profile.",
									},
									"instance_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for the instance of a compliance profile.",
									},
									"instance_location": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The location of the compliance instance.",
									},
									"attachment_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique ID for the attachment to a compliance profile.",
									},
									"profile_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the compliance profile.",
									},
								},
							},
						},
						"locator_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A project configuration description.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration name. It's unique within the account across projects and regions.",
						},
						"environment_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the project environment.",
						},
						"authorizations": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trusted profile ID.",
									},
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The authorization method. You can authorize by using a trusted profile or an API key in Secrets Manager.",
									},
									"api_key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
										Description: "The IBM Cloud API Key. It can be either raw or pulled from the catalog via a `CRN` or `JSON` blob.",
									},
								},
							},
						},
						"inputs": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The input variables that are used for configuration definition and environment.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"settings": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The Schematics environment variables to use to deploy the configuration. Settings are only available if they are specified when the configuration is initially created.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_crns": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The CRNs of the resources that are associated with this configuration.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"approved_version": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A summary of a project configuration version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of the definition in a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the project environment.",
									},
									"locator_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
									},
								},
							},
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the configuration.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the configuration.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
					},
				},
			},
			"deployed_version": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A summary of a project configuration version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of the definition in a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the project environment.",
									},
									"locator_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
									},
								},
							},
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the configuration.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the configuration.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmProjectConfigRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getConfigOptions := &projectv1.GetConfigOptions{}

	getConfigOptions.SetProjectID(d.Get("project_id").(string))
	getConfigOptions.SetID(d.Get("project_config_id").(string))

	projectConfig, _, err := projectClient.GetConfigWithContext(context, getConfigOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigWithContext failed: %s", err.Error()), "(Data) ibm_project_config", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getConfigOptions.ProjectID, *getConfigOptions.ID))

	if err = d.Set("version", flex.IntValue(projectConfig.Version)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("is_draft", projectConfig.IsDraft); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting is_draft: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	needsAttentionState := []map[string]interface{}{}
	if projectConfig.NeedsAttentionState != nil {
		for _, modelItem := range projectConfig.NeedsAttentionState {
			modelMap, err := dataSourceIbmProjectConfigProjectConfigNeedsAttentionStateToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
				return tfErr.GetDiag()
			}
			needsAttentionState = append(needsAttentionState, modelMap)
		}
	}
	if err = d.Set("needs_attention_state", needsAttentionState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting needs_attention_state: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(projectConfig.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("modified_at", flex.DateTimeToString(projectConfig.ModifiedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting modified_at: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("last_saved_at", flex.DateTimeToString(projectConfig.LastSavedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting last_saved_at: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	outputs := []map[string]interface{}{}
	if projectConfig.Outputs != nil {
		for _, modelItem := range projectConfig.Outputs {
			modelMap, err := dataSourceIbmProjectConfigOutputValueToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
				return tfErr.GetDiag()
			}
			outputs = append(outputs, modelMap)
		}
	}
	if err = d.Set("outputs", outputs); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting outputs: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	project := []map[string]interface{}{}
	if projectConfig.Project != nil {
		modelMap, err := dataSourceIbmProjectConfigProjectReferenceToMap(projectConfig.Project)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
			return tfErr.GetDiag()
		}
		project = append(project, modelMap)
	}
	if err = d.Set("project", project); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting project: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	schematics := []map[string]interface{}{}
	if projectConfig.Schematics != nil {
		modelMap, err := dataSourceIbmProjectConfigSchematicsMetadataToMap(projectConfig.Schematics)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
			return tfErr.GetDiag()
		}
		schematics = append(schematics, modelMap)
	}
	if err = d.Set("schematics", schematics); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting schematics: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", projectConfig.State); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("update_available", projectConfig.UpdateAvailable); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting update_available: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", projectConfig.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	definition := []map[string]interface{}{}
	if projectConfig.Definition != nil {
		modelMap, err := dataSourceIbmProjectConfigProjectConfigDefinitionResponseToMap(projectConfig.Definition)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
			return tfErr.GetDiag()
		}
		definition = append(definition, modelMap)
	}
	if err = d.Set("definition", definition); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting definition: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	approvedVersion := []map[string]interface{}{}
	if projectConfig.ApprovedVersion != nil {
		modelMap, err := dataSourceIbmProjectConfigProjectConfigVersionSummaryToMap(projectConfig.ApprovedVersion)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
			return tfErr.GetDiag()
		}
		approvedVersion = append(approvedVersion, modelMap)
	}
	if err = d.Set("approved_version", approvedVersion); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting approved_version: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	deployedVersion := []map[string]interface{}{}
	if projectConfig.DeployedVersion != nil {
		modelMap, err := dataSourceIbmProjectConfigProjectConfigVersionSummaryToMap(projectConfig.DeployedVersion)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_project_config", "read")
			return tfErr.GetDiag()
		}
		deployedVersion = append(deployedVersion, modelMap)
	}
	if err = d.Set("deployed_version", deployedVersion); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting deployed_version: %s", err), "(Data) ibm_project_config", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmProjectConfigProjectConfigNeedsAttentionStateToMap(model *projectv1.ProjectConfigNeedsAttentionState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["event_id"] = model.EventID
	modelMap["event"] = model.Event
	if model.Severity != nil {
		modelMap["severity"] = model.Severity
	}
	if model.ActionURL != nil {
		modelMap["action_url"] = model.ActionURL
	}
	if model.Target != nil {
		modelMap["target"] = model.Target
	}
	if model.TriggeredBy != nil {
		modelMap["triggered_by"] = model.TriggeredBy
	}
	modelMap["timestamp"] = model.Timestamp
	return modelMap, nil
}

func dataSourceIbmProjectConfigOutputValueToMap(model *projectv1.OutputValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Value != nil {
		modelMap["value_json"] = stringify(model.Value)
		jsonStr, err := json.Marshal(model.Value)
		if err != nil {
			b := []byte(jsonStr)
			var f interface{}
			json.Unmarshal(b, &f)
			valueMap := f.(map[string]interface{})
			modelMap["value"] = valueMap
		}
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["href"] = model.Href
	definitionMap, err := dataSourceIbmProjectConfigProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmProjectConfigSchematicsMetadataToMap(model *projectv1.SchematicsMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.WorkspaceCrn != nil {
		modelMap["workspace_crn"] = model.WorkspaceCrn
	}
	if model.ValidatePreScript != nil {
		validatePreScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.ValidatePreScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["validate_pre_script"] = []map[string]interface{}{validatePreScriptMap}
	}
	if model.ValidatePostScript != nil {
		validatePostScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.ValidatePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["validate_post_script"] = []map[string]interface{}{validatePostScriptMap}
	}
	if model.DeployPreScript != nil {
		deployPreScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.DeployPreScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["deploy_pre_script"] = []map[string]interface{}{deployPreScriptMap}
	}
	if model.DeployPostScript != nil {
		deployPostScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.DeployPostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["deploy_post_script"] = []map[string]interface{}{deployPostScriptMap}
	}
	if model.UndeployPreScript != nil {
		undeployPreScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.UndeployPreScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["undeploy_pre_script"] = []map[string]interface{}{undeployPreScriptMap}
	}
	if model.UndeployPostScript != nil {
		undeployPostScriptMap, err := dataSourceIbmProjectConfigScriptToMap(model.UndeployPostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["undeploy_post_script"] = []map[string]interface{}{undeployPostScriptMap}
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigScriptToMap(model *projectv1.Script) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.ShortDescription != nil {
		modelMap["short_description"] = model.ShortDescription
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigDefinitionResponseToMap(model projectv1.ProjectConfigDefinitionResponseIntf) (map[string]interface{}, error) {
	if _, ok := model.(*projectv1.ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse); ok {
		return dataSourceIbmProjectConfigProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponseToMap(model.(*projectv1.ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse))
	} else if _, ok := model.(*projectv1.ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse); ok {
		return dataSourceIbmProjectConfigProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponseToMap(model.(*projectv1.ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse))
	} else if _, ok := model.(*projectv1.ProjectConfigDefinitionResponse); ok {
		modelMap := make(map[string]interface{})
		model := model.(*projectv1.ProjectConfigDefinitionResponse)
		if model.ComplianceProfile != nil {
			complianceProfileMap, err := dataSourceIbmProjectConfigProjectComplianceProfileToMap(model.ComplianceProfile)
			if err != nil {
				return modelMap, err
			}
			modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
		}
		if model.LocatorID != nil {
			modelMap["locator_id"] = model.LocatorID
		}
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.EnvironmentID != nil {
			modelMap["environment_id"] = model.EnvironmentID
		}
		if model.Authorizations != nil {
			authorizationsMap, err := dataSourceIbmProjectConfigProjectConfigAuthToMap(model.Authorizations)
			if err != nil {
				return modelMap, err
			}
			modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
		}
		if model.Inputs != nil {
			inputs := make(map[string]interface{})
			for k, v := range model.Inputs {
				inputs[k] = fmt.Sprintf("%v", v)
			}
			modelMap["inputs"] = inputs
		}
		if model.Settings != nil {
			settings := make(map[string]interface{})
			for k, v := range model.Settings {
				settings[k] = fmt.Sprintf("%v", v)
			}
			modelMap["settings"] = settings
		}
		if model.ResourceCrns != nil {
			modelMap["resource_crns"] = model.ResourceCrns
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized projectv1.ProjectConfigDefinitionResponseIntf subtype encountered")
	}
}

func dataSourceIbmProjectConfigProjectComplianceProfileToMap(model *projectv1.ProjectComplianceProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = model.InstanceID
	}
	if model.InstanceLocation != nil {
		modelMap["instance_location"] = model.InstanceLocation
	}
	if model.AttachmentID != nil {
		modelMap["attachment_id"] = model.AttachmentID
	}
	if model.ProfileName != nil {
		modelMap["profile_name"] = model.ProfileName
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigAuthToMap(model *projectv1.ProjectConfigAuth) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = model.TrustedProfileID
	}
	if model.Method != nil {
		modelMap["method"] = model.Method
	}
	if model.ApiKey != nil {
		modelMap["api_key"] = model.ApiKey
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponseToMap(model *projectv1.ProjectConfigDefinitionResponseDAConfigDefinitionPropertiesResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ComplianceProfile != nil {
		complianceProfileMap, err := dataSourceIbmProjectConfigProjectComplianceProfileToMap(model.ComplianceProfile)
		if err != nil {
			return modelMap, err
		}
		modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
	}
	if model.LocatorID != nil {
		modelMap["locator_id"] = model.LocatorID
	}
	modelMap["description"] = model.Description
	modelMap["name"] = model.Name
	if model.EnvironmentID != nil {
		modelMap["environment_id"] = model.EnvironmentID
	}
	if model.Authorizations != nil {
		authorizationsMap, err := dataSourceIbmProjectConfigProjectConfigAuthToMap(model.Authorizations)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
	}
	if model.Inputs != nil {
		inputs := make(map[string]interface{})
		for k, v := range model.Inputs {
			inputs[k] = fmt.Sprintf("%v", v)
		}
		modelMap["inputs"] = inputs
	}
	if model.Settings != nil {
		settings := make(map[string]interface{})
		for k, v := range model.Settings {
			settings[k] = fmt.Sprintf("%v", v)
		}
		modelMap["settings"] = settings
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponseToMap(model *projectv1.ProjectConfigDefinitionResponseResourceConfigDefinitionPropertiesResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceCrns != nil {
		modelMap["resource_crns"] = model.ResourceCrns
	}
	modelMap["description"] = model.Description
	modelMap["name"] = model.Name
	if model.EnvironmentID != nil {
		modelMap["environment_id"] = model.EnvironmentID
	}
	if model.Authorizations != nil {
		authorizationsMap, err := dataSourceIbmProjectConfigProjectConfigAuthToMap(model.Authorizations)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
	}
	if model.Inputs != nil {
		inputs := make(map[string]interface{})
		for k, v := range model.Inputs {
			inputs[k] = fmt.Sprintf("%v", v)
		}
		modelMap["inputs"] = inputs
	}
	if model.Settings != nil {
		settings := make(map[string]interface{})
		for k, v := range model.Settings {
			settings[k] = fmt.Sprintf("%v", v)
		}
		modelMap["settings"] = settings
	}
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigVersionSummaryToMap(model *projectv1.ProjectConfigVersionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	definitionMap, err := dataSourceIbmProjectConfigProjectConfigVersionDefinitionSummaryToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["state"] = model.State
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmProjectConfigProjectConfigVersionDefinitionSummaryToMap(model *projectv1.ProjectConfigVersionDefinitionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnvironmentID != nil {
		modelMap["environment_id"] = model.EnvironmentID
	}
	if model.LocatorID != nil {
		modelMap["locator_id"] = model.LocatorID
	}
	return modelMap, nil
}
