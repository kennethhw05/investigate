package graphql

import (
	"testing"

	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/testutils"
)

func attemptGetAuditByID(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	id := "3a6acaaf-596a-4569-976a-008cb2562c3e"
	audit, err := queryResolver.Audit(testutils.GetTestingContext(role), id)
	if permitted {
		if err != nil {
			t.Fatalf("[TestGetAuditByIDAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if audit.GetID() != id {
				t.Fatalf("[TestGetAuditByIDAs%s] Invalid ID Returned: %s", role.String(), audit.GetID())
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestUpdateAuditAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}
}

func TestGetAuditByIdAsAccessRoleSuperAdmin(t *testing.T) {
	attemptGetAuditByID(models.AccessRoleSuperAdmin, true, t)
}

func TestGetAuditByIdAsAccessRoleAdmin(t *testing.T) {
	attemptGetAuditByID(models.AccessRoleAdmin, true, t)
}

func TestGetAuditByIdAsAccessRoleUser(t *testing.T) {
	attemptGetAuditByID(models.AccessRoleUser, false, t)
}

func TestGetAuditByIdAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptGetAuditByID(models.AccessRoleGuestAPIOnly, false, t)
}

func attemptListAudits(role models.AccessRole, permitted bool, t *testing.T) {
	resolver := setupTestResolver(t)
	defer cleanupTestResolver(resolver)
	queryResolver := resolver.Query()

	filter := models.AuditFilter{}
	filter.EditAction = new(models.EditAction)
	*filter.EditAction = models.EditActionUpdate

	audits, err := queryResolver.AllAudits(testutils.GetTestingContext(role), &filter, nil, nil, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListAuditsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(audits) != 2 {
				t.Fatalf("[TestListAuditsAs%s] Invalid number of Audits Returned: %v", role.String(), len(audits))
			} else {
				audit := audits[0]
				if audit.GetID() != "abb036e3-973f-4428-899d-8a70210cd7c4" {
					t.Fatalf("[TestListAuditsAs%s] Invalid ID Returned: %s", role.String(), audit.GetID())
				}
			}
		}

		filter.EditAction = nil
		page := 0
		perPage := 3

		audits, err = queryResolver.AllAudits(testutils.GetTestingContext(role), &filter, &page, &perPage, nil, nil)
		if err != nil {
			t.Fatalf("[TestListAuditsAs%s] Error from Query Resolver: %s", role.String(), err)
		} else {
			if len(audits) != 3 {
				t.Fatalf("[TestListAuditsAs%s] Invalid number of Audits Returned: %v", role.String(), len(audits))
			}
		}
	} else {
		if err == nil {
			t.Fatalf("[TestListAuditsAs%s] No Error from Query Resolver, Expected Access Denied", role.String())
		}
	}

	// Test Meta Data //
	metadata, err := queryResolver._allAuditsMeta(testutils.GetTestingContext(role), &filter, nil, nil)
	if permitted {
		if err != nil {
			t.Fatalf("[TestListAuditsAs%s] Error from Meta Query Resolver: %s", role.String(), err)
		} else {
			if metadata.Count != 4 {
				t.Fatalf("[TestListAuditsAs%s] Invalid Count Returned: %d, Expected: 4", role.String(), metadata.Count)
			}
		}
	} else if err == nil {
		t.Fatalf("[TestListAuditsAs%s] No Error from Meta Query Resolver, Expected Access Denied", role.String())
	}
}

func TestListAuditsAsAccessRoleSuperAdmin(t *testing.T) {
	attemptListAudits(models.AccessRoleSuperAdmin, true, t)
}

func TestListAuditsAsAccessRoleAdmin(t *testing.T) {
	attemptListAudits(models.AccessRoleAdmin, true, t)
}

func TestListAuditsAsAccessRoleUser(t *testing.T) {
	attemptListAudits(models.AccessRoleUser, false, t)
}

func TestListAuditsAsAccessRoleGuestAPIOnly(t *testing.T) {
	attemptListAudits(models.AccessRoleGuestAPIOnly, false, t)
}
