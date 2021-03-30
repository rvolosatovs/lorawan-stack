// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

var ClaimEndDeviceRequestFieldPathsNested = []string{
	"invalidate_authentication_code",
	"source_device",
	"source_device.authenticated_identifiers",
	"source_device.authenticated_identifiers.authentication_code",
	"source_device.authenticated_identifiers.dev_eui",
	"source_device.authenticated_identifiers.join_eui",
	"source_device.qr_code",
	"target_application_ids",
	"target_application_ids.application_id",
	"target_application_server_address",
	"target_application_server_id",
	"target_application_server_kek_label",
	"target_device_id",
	"target_net_id",
	"target_network_server_address",
	"target_network_server_kek_label",
}

var ClaimEndDeviceRequestFieldPathsTopLevel = []string{
	"invalidate_authentication_code",
	"source_device",
	"target_application_ids",
	"target_application_server_address",
	"target_application_server_id",
	"target_application_server_kek_label",
	"target_device_id",
	"target_net_id",
	"target_network_server_address",
	"target_network_server_kek_label",
}
var AuthorizeApplicationRequestFieldPathsNested = []string{
	"api_key",
	"application_ids",
	"application_ids.application_id",
}

var AuthorizeApplicationRequestFieldPathsTopLevel = []string{
	"api_key",
	"application_ids",
}
var CUPSRedirectionFieldPathsNested = []string{
	"current_gateway_key",
	"target_cups_uri",
}

var CUPSRedirectionFieldPathsTopLevel = []string{
	"current_gateway_key",
	"target_cups_uri",
}
var ClaimGatewayRequestFieldPathsNested = []string{
	"collaborator",
	"collaborator.ids",
	"collaborator.ids.organization_ids",
	"collaborator.ids.organization_ids.organization_id",
	"collaborator.ids.user_ids",
	"collaborator.ids.user_ids.email",
	"collaborator.ids.user_ids.user_id",
	"cups_redirection",
	"cups_redirection.current_gateway_key",
	"cups_redirection.target_cups_uri",
	"source_gateway",
	"source_gateway.authenticated_identifiers",
	"source_gateway.authenticated_identifiers.authentication_code",
	"source_gateway.authenticated_identifiers.gateway_eui",
	"source_gateway.qr_code",
	"target_gateway_id",
	"target_gateway_server_address",
}

var ClaimGatewayRequestFieldPathsTopLevel = []string{
	"collaborator",
	"cups_redirection",
	"source_gateway",
	"target_gateway_id",
	"target_gateway_server_address",
}
var AuthorizeGatewayRequestFieldPathsNested = []string{
	"api_key",
	"gateway_ids",
	"gateway_ids.eui",
	"gateway_ids.gateway_id",
}

var AuthorizeGatewayRequestFieldPathsTopLevel = []string{
	"api_key",
	"gateway_ids",
}
var ClaimEndDeviceRequest_AuthenticatedIdentifiersFieldPathsNested = []string{
	"authentication_code",
	"dev_eui",
	"join_eui",
}

var ClaimEndDeviceRequest_AuthenticatedIdentifiersFieldPathsTopLevel = []string{
	"authentication_code",
	"dev_eui",
	"join_eui",
}
var ClaimGatewayRequest_AuthenticatedIdentifiersFieldPathsNested = []string{
	"authentication_code",
	"gateway_eui",
}

var ClaimGatewayRequest_AuthenticatedIdentifiersFieldPathsTopLevel = []string{
	"authentication_code",
	"gateway_eui",
}
