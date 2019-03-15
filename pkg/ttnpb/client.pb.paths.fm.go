// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

var ClientFieldPathsNested = []string{
	"attributes",
	"contact_info",
	"created_at",
	"description",
	"endorsed",
	"grants",
	"ids",
	"ids.client_id",
	"name",
	"redirect_uris",
	"rights",
	"secret",
	"skip_authorization",
	"state",
	"updated_at",
}

var ClientFieldPathsTopLevel = []string{
	"attributes",
	"contact_info",
	"created_at",
	"description",
	"endorsed",
	"grants",
	"ids",
	"name",
	"redirect_uris",
	"rights",
	"secret",
	"skip_authorization",
	"state",
	"updated_at",
}
var ClientsFieldPathsNested = []string{
	"clients",
}

var ClientsFieldPathsTopLevel = []string{
	"clients",
}
var GetClientRequestFieldPathsNested = []string{
	"client_ids",
	"client_ids.client_id",
	"field_mask",
}

var GetClientRequestFieldPathsTopLevel = []string{
	"client_ids",
	"field_mask",
}
var ListClientsRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"field_mask",
	"limit",
	"order",
	"page",
}

var ListClientsRequestFieldPathsTopLevel = []string{
	"collaborator",
	"field_mask",
	"limit",
	"order",
	"page",
}
var CreateClientRequestFieldPathsNested = []string{
	"client",
	"client.attributes",
	"client.contact_info",
	"client.created_at",
	"client.description",
	"client.endorsed",
	"client.grants",
	"client.ids",
	"client.ids.client_id",
	"client.name",
	"client.redirect_uris",
	"client.rights",
	"client.secret",
	"client.skip_authorization",
	"client.state",
	"client.updated_at",
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
}

var CreateClientRequestFieldPathsTopLevel = []string{
	"client",
	"collaborator",
}
var UpdateClientRequestFieldPathsNested = []string{
	"client",
	"client.attributes",
	"client.contact_info",
	"client.created_at",
	"client.description",
	"client.endorsed",
	"client.grants",
	"client.ids",
	"client.ids.client_id",
	"client.name",
	"client.redirect_uris",
	"client.rights",
	"client.secret",
	"client.skip_authorization",
	"client.state",
	"client.updated_at",
	"field_mask",
}

var UpdateClientRequestFieldPathsTopLevel = []string{
	"client",
	"field_mask",
}
var SetClientCollaboratorRequestFieldPathsNested = []string{
	"client_ids",
	"client_ids.client_id",
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.ids",
	"collaborator.ids.ids.organization_ids",
	"collaborator.ids.ids.organization_ids.organization_id",
	"collaborator.ids.ids.user_ids",
	"collaborator.ids.ids.user_ids.email",
	"collaborator.ids.ids.user_ids.user_id",
	"collaborator.rights",
}

var SetClientCollaboratorRequestFieldPathsTopLevel = []string{
	"client_ids",
	"collaborator",
}
