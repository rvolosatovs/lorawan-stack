// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"

	types "github.com/gogo/protobuf/types"
)

func (dst *Organization) SetFields(src *Organization, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "created_at":
			if len(subs) > 0 {
				return fmt.Errorf("'created_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CreatedAt = src.CreatedAt
			} else {
				var zero time.Time
				dst.CreatedAt = zero
			}
		case "updated_at":
			if len(subs) > 0 {
				return fmt.Errorf("'updated_at' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UpdatedAt = src.UpdatedAt
			} else {
				var zero time.Time
				dst.UpdatedAt = zero
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "description":
			if len(subs) > 0 {
				return fmt.Errorf("'description' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Description = src.Description
			} else {
				var zero string
				dst.Description = zero
			}
		case "attributes":
			if len(subs) > 0 {
				return fmt.Errorf("'attributes' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Attributes = src.Attributes
			} else {
				dst.Attributes = nil
			}
		case "contact_info":
			if len(subs) > 0 {
				return fmt.Errorf("'contact_info' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ContactInfo = src.ContactInfo
			} else {
				dst.ContactInfo = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *Organizations) SetFields(src *Organizations, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organizations":
			if len(subs) > 0 {
				return fmt.Errorf("'organizations' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Organizations = src.Organizations
			} else {
				dst.Organizations = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *GetOrganizationRequest) SetFields(src *GetOrganizationRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListOrganizationsRequest) SetFields(src *ListOrganizationsRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "collaborator":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationOrUserIdentifiers
				if (src == nil || src.Collaborator == nil) && dst.Collaborator == nil {
					continue
				}
				if src != nil {
					newSrc = src.Collaborator
				}
				if dst.Collaborator != nil {
					newDst = dst.Collaborator
				} else {
					newDst = &OrganizationOrUserIdentifiers{}
					dst.Collaborator = newDst
				}
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					dst.Collaborator = nil
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}
		case "order":
			if len(subs) > 0 {
				return fmt.Errorf("'order' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Order = src.Order
			} else {
				var zero string
				dst.Order = zero
			}
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *CreateOrganizationRequest) SetFields(src *CreateOrganizationRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization":
			if len(subs) > 0 {
				var newDst, newSrc *Organization
				if src != nil {
					newSrc = &src.Organization
				}
				newDst = &dst.Organization
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Organization = src.Organization
				} else {
					var zero Organization
					dst.Organization = zero
				}
			}
		case "collaborator":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationOrUserIdentifiers
				if src != nil {
					newSrc = &src.Collaborator
				}
				newDst = &dst.Collaborator
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					var zero OrganizationOrUserIdentifiers
					dst.Collaborator = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *UpdateOrganizationRequest) SetFields(src *UpdateOrganizationRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization":
			if len(subs) > 0 {
				var newDst, newSrc *Organization
				if src != nil {
					newSrc = &src.Organization
				}
				newDst = &dst.Organization
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Organization = src.Organization
				} else {
					var zero Organization
					dst.Organization = zero
				}
			}
		case "field_mask":
			if len(subs) > 0 {
				return fmt.Errorf("'field_mask' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.FieldMask = src.FieldMask
			} else {
				var zero types.FieldMask
				dst.FieldMask = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListOrganizationAPIKeysRequest) SetFields(src *ListOrganizationAPIKeysRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *GetOrganizationAPIKeyRequest) SetFields(src *GetOrganizationAPIKeyRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "key_id":
			if len(subs) > 0 {
				return fmt.Errorf("'key_id' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.KeyID = src.KeyID
			} else {
				var zero string
				dst.KeyID = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *CreateOrganizationAPIKeyRequest) SetFields(src *CreateOrganizationAPIKeyRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "rights":
			if len(subs) > 0 {
				return fmt.Errorf("'rights' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Rights = src.Rights
			} else {
				dst.Rights = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *UpdateOrganizationAPIKeyRequest) SetFields(src *UpdateOrganizationAPIKeyRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "api_key":
			if len(subs) > 0 {
				var newDst, newSrc *APIKey
				if src != nil {
					newSrc = &src.APIKey
				}
				newDst = &dst.APIKey
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.APIKey = src.APIKey
				} else {
					var zero APIKey
					dst.APIKey = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *ListOrganizationCollaboratorsRequest) SetFields(src *ListOrganizationCollaboratorsRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "limit":
			if len(subs) > 0 {
				return fmt.Errorf("'limit' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Limit = src.Limit
			} else {
				var zero uint32
				dst.Limit = zero
			}
		case "page":
			if len(subs) > 0 {
				return fmt.Errorf("'page' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Page = src.Page
			} else {
				var zero uint32
				dst.Page = zero
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *GetOrganizationCollaboratorRequest) SetFields(src *GetOrganizationCollaboratorRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "collaborator":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationOrUserIdentifiers
				if src != nil {
					newSrc = &src.OrganizationOrUserIdentifiers
				}
				newDst = &dst.OrganizationOrUserIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationOrUserIdentifiers = src.OrganizationOrUserIdentifiers
				} else {
					var zero OrganizationOrUserIdentifiers
					dst.OrganizationOrUserIdentifiers = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *SetOrganizationCollaboratorRequest) SetFields(src *SetOrganizationCollaboratorRequest, paths ...string) error {
	for name, subs := range _processPaths(paths) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				var newDst, newSrc *OrganizationIdentifiers
				if src != nil {
					newSrc = &src.OrganizationIdentifiers
				}
				newDst = &dst.OrganizationIdentifiers
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.OrganizationIdentifiers = src.OrganizationIdentifiers
				} else {
					var zero OrganizationIdentifiers
					dst.OrganizationIdentifiers = zero
				}
			}
		case "collaborator":
			if len(subs) > 0 {
				var newDst, newSrc *Collaborator
				if src != nil {
					newSrc = &src.Collaborator
				}
				newDst = &dst.Collaborator
				if err := newDst.SetFields(newSrc, subs...); err != nil {
					return err
				}
			} else {
				if src != nil {
					dst.Collaborator = src.Collaborator
				} else {
					var zero Collaborator
					dst.Collaborator = zero
				}
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
