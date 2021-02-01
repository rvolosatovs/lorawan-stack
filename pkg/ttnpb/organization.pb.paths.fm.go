// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

var OrganizationFieldPathsNested = []string{
	"attributes",
	"contact_info",
	"created_at",
	"deleted_at",
	"description",
	"ids",
	"ids.organization_id",
	"name",
	"updated_at",
}

var OrganizationFieldPathsTopLevel = []string{
	"attributes",
	"contact_info",
	"created_at",
	"deleted_at",
	"description",
	"ids",
	"name",
	"updated_at",
}
var OrganizationsFieldPathsNested = []string{
	"organizations",
}

var OrganizationsFieldPathsTopLevel = []string{
	"organizations",
}
var GetOrganizationRequestFieldPathsNested = []string{
	"field_mask",
	"organization_ids",
	"organization_ids.organization_id",
}

var GetOrganizationRequestFieldPathsTopLevel = []string{
	"field_mask",
	"organization_ids",
}
var ListOrganizationsRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"field_mask",
	"include_deleted",
	"limit",
	"order",
	"page",
}

var ListOrganizationsRequestFieldPathsTopLevel = []string{
	"collaborator",
	"field_mask",
	"include_deleted",
	"limit",
	"order",
	"page",
}
var CreateOrganizationRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"organization",
	"organization.attributes",
	"organization.contact_info",
	"organization.created_at",
	"organization.deleted_at",
	"organization.description",
	"organization.ids",
	"organization.ids.organization_id",
	"organization.name",
	"organization.updated_at",
}

var CreateOrganizationRequestFieldPathsTopLevel = []string{
	"collaborator",
	"organization",
}
var UpdateOrganizationRequestFieldPathsNested = []string{
	"field_mask",
	"organization",
	"organization.attributes",
	"organization.contact_info",
	"organization.created_at",
	"organization.deleted_at",
	"organization.description",
	"organization.ids",
	"organization.ids.organization_id",
	"organization.name",
	"organization.updated_at",
}

var UpdateOrganizationRequestFieldPathsTopLevel = []string{
	"field_mask",
	"organization",
}
var ListOrganizationAPIKeysRequestFieldPathsNested = []string{
	"limit",
	"organization_ids",
	"organization_ids.organization_id",
	"page",
}

var ListOrganizationAPIKeysRequestFieldPathsTopLevel = []string{
	"limit",
	"organization_ids",
	"page",
}
var GetOrganizationAPIKeyRequestFieldPathsNested = []string{
	"key_id",
	"organization_ids",
	"organization_ids.organization_id",
}

var GetOrganizationAPIKeyRequestFieldPathsTopLevel = []string{
	"key_id",
	"organization_ids",
}
var CreateOrganizationAPIKeyRequestFieldPathsNested = []string{
	"name",
	"organization_ids",
	"organization_ids.organization_id",
	"rights",
}

var CreateOrganizationAPIKeyRequestFieldPathsTopLevel = []string{
	"name",
	"organization_ids",
	"rights",
}
var UpdateOrganizationAPIKeyRequestFieldPathsNested = []string{
	"api_key",
	"api_key.id",
	"api_key.key",
	"api_key.name",
	"api_key.rights",
	"organization_ids",
	"organization_ids.organization_id",
}

var UpdateOrganizationAPIKeyRequestFieldPathsTopLevel = []string{
	"api_key",
	"organization_ids",
}
var ListOrganizationCollaboratorsRequestFieldPathsNested = []string{
	"limit",
	"organization_ids",
	"organization_ids.organization_id",
	"page",
}

var ListOrganizationCollaboratorsRequestFieldPathsTopLevel = []string{
	"limit",
	"organization_ids",
	"page",
}
var GetOrganizationCollaboratorRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"organization_ids",
	"organization_ids.organization_id",
}

var GetOrganizationCollaboratorRequestFieldPathsTopLevel = []string{
	"collaborator",
	"organization_ids",
}
var SetOrganizationCollaboratorRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.ids",
	"collaborator.ids.ids.organization_ids",
	"collaborator.ids.ids.organization_ids.organization_id",
	"collaborator.ids.ids.user_ids",
	"collaborator.ids.ids.user_ids.email",
	"collaborator.ids.ids.user_ids.user_id",
	"collaborator.rights",
	"organization_ids",
	"organization_ids.organization_id",
}

var SetOrganizationCollaboratorRequestFieldPathsTopLevel = []string{
	"collaborator",
	"organization_ids",
}
