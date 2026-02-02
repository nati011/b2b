package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"marketplace/internal/config"
	resourcedomain "marketplace/internal/infra/authz/resource/domain"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		resourcesPath = flag.String("resources", "config/resources.yaml", "Path to resources.yaml file")
		rolesPath     = flag.String("roles", "config/roles.yaml", "Path to roles.yaml file")
		outputPath    = flag.String("output", "", "Path to output updated roles.yaml (default: overwrites input)")
		dryRun        = flag.Bool("dry-run", false, "Show what would be changed without modifying files")
		autoSuperadmin = flag.Bool("auto-superadmin", true, "Automatically assign all new resources to superadmin role")
	)
	flag.Parse()

	// Load resources
	resources, err := loadResources(*resourcesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading resources: %v\n", err)
		os.Exit(1)
	}

	// Load roles
	roles, err := loadRoles(*rolesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading roles: %v\n", err)
		os.Exit(1)
	}

	// Analyze and distribute
	changes := analyzeAndDistribute(resources, roles, *autoSuperadmin)

	if len(changes) == 0 {
		fmt.Println("✓ No new resources or actions found. All roles are up to date.")
		return
	}

	// Print changes
	fmt.Println("\n📋 Changes to be applied:")
	fmt.Println(strings.Repeat("=", 80))
	for _, change := range changes {
		fmt.Printf("\n%s\n", change)
	}
	fmt.Println(strings.Repeat("=", 80))

	if *dryRun {
		fmt.Println("\n🔍 Dry-run mode: No files were modified.")
		return
	}

	// Write updated roles
	output := *outputPath
	if output == "" {
		output = *rolesPath
	}

	if err := writeRoles(output, roles); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing roles file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Successfully updated roles file: %s\n", output)
}

type ResourceAction struct {
	Resource string
	Action   string
}

type Change struct {
	Role     string
	Resource string
	Actions  []string
}

func (c Change) String() string {
	return fmt.Sprintf("Role: %s\n  Resource: %s\n  New Actions: %s", 
		c.Role, c.Resource, strings.Join(c.Actions, ", "))
}

func loadResources(path string) ([]resourcedomain.ManifestResource, error) {
	return resourcedomain.LoadResources(path)
}

func loadRoles(path string) ([]config.RoleDefinition, error) {
	return config.LoadRolesFromFile(path)
}

func analyzeAndDistribute(resources []resourcedomain.ManifestResource, roles []config.RoleDefinition, autoSuperadmin bool) []Change {
	var changes []Change

	// Build map of existing permissions per role
	rolePermissions := make(map[string]map[ResourceAction]bool)
	for _, role := range roles {
		rolePermissions[role.ID] = make(map[ResourceAction]bool)
		for _, perm := range role.Permissions {
			for _, action := range perm.Actions {
				rolePermissions[role.ID][ResourceAction{
					Resource: perm.Resource,
					Action:   action,
				}] = true
			}
		}
	}

	// Find superadmin role
	var superadminRole *config.RoleDefinition
	for i := range roles {
		if roles[i].ID == "superadmin" {
			superadminRole = &roles[i]
			break
		}
	}

	// Process each resource
	for _, resource := range resources {
		// Get all unique actions for this resource (deduplicate)
		resourceActions := make(map[string]bool)
		for _, action := range resource.Actions {
			resourceActions[action.Name] = true
		}

		// Check each role
		for i := range roles {
			role := &roles[i]
			existingPerms := rolePermissions[role.ID]

			// Check if this resource is already in the role
			var existingPermPattern *config.PermissionPattern
			for j := range role.Permissions {
				if role.Permissions[j].Resource == resource.Code {
					existingPermPattern = &role.Permissions[j]
					break
				}
			}

			// Find missing actions that should be assigned to this role
			var missingActions []string
			for actionName := range resourceActions {
				ra := ResourceAction{
					Resource: resource.Code,
					Action:   actionName,
				}
				// Check if action is missing AND should be assigned to this role
				if !existingPerms[ra] && shouldAssignAction(role.ID, resource.Code, actionName) {
					missingActions = append(missingActions, actionName)
				}
			}

			if len(missingActions) == 0 {
				continue
			}

			// Auto-assign to superadmin if enabled
			if autoSuperadmin && role.ID == "superadmin" && superadminRole != nil {
				if existingPermPattern == nil {
					// Add new resource entry
					superadminRole.Permissions = append(superadminRole.Permissions, config.PermissionPattern{
						Resource: resource.Code,
						Actions:  missingActions,
					})
				} else {
					// Add missing actions to existing entry
					existingActions := make(map[string]bool)
					for _, a := range existingPermPattern.Actions {
						existingActions[a] = true
					}
					for _, newAction := range missingActions {
						if !existingActions[newAction] {
							existingPermPattern.Actions = append(existingPermPattern.Actions, newAction)
						}
					}
				}
				sort.Strings(missingActions)
				changes = append(changes, Change{
					Role:     role.Name,
					Resource: resource.Code,
					Actions:  missingActions,
				})
				continue
			}

			// For other roles, check if they should get this resource based on patterns
			shouldAssign := shouldAssignToRole(role, resource)

			if shouldAssign {
				if existingPermPattern == nil {
					// Add new resource entry
					role.Permissions = append(role.Permissions, config.PermissionPattern{
						Resource: resource.Code,
						Actions:  missingActions,
					})
				} else {
					// Add missing actions to existing entry
					existingActions := make(map[string]bool)
					for _, a := range existingPermPattern.Actions {
						existingActions[a] = true
					}
					for _, newAction := range missingActions {
						if !existingActions[newAction] {
							existingPermPattern.Actions = append(existingPermPattern.Actions, newAction)
						}
					}
				}
				sort.Strings(missingActions)
				changes = append(changes, Change{
					Role:     role.Name,
					Resource: resource.Code,
					Actions:  missingActions,
				})
			}
		}
	}

	return changes
}

// shouldAssignToRole determines if a resource should be assigned to a role based on patterns
// This is conservative - only assigns resources that make sense for each role
func shouldAssignToRole(role *config.RoleDefinition, resource resourcedomain.ManifestResource) bool {
	roleID := strings.ToLower(role.ID)
	resourceCode := strings.ToLower(resource.Code)

	// Supplier role: only gets supplier resources (view/update) and products (CRUD)
	if roleID == "supplier" {
		if resourceCode == "supplier" {
			return true // Supplier can view/update their own profile
		}
		if resourceCode == "products" {
			return true // Supplier can manage products
		}
		if resourceCode == "order" {
			return true // Supplier can view orders (already assigned)
		}
		if resourceCode == "bank-account" {
			return true // Supplier can manage their own bank accounts
		}
		return false
	}

	// Customer role: only gets customer resources (view/update), products (view), and orders (view/create)
	if roleID == "customer" {
		if resourceCode == "customer" {
			return true // Customer can view/update their own profile
		}
		if resourceCode == "products" {
			return true // Customer can view products (already assigned)
		}
		if resourceCode == "order" {
			return true // Customer can view/create orders (already assigned)
		}
		return false
	}

	// By default, don't auto-assign to non-superadmin roles
	return false
}

// shouldAssignAction determines if a specific action should be assigned to a role
func shouldAssignAction(roleID, resourceCode, actionName string) bool {
	roleID = strings.ToLower(roleID)
	resourceCode = strings.ToLower(resourceCode)
	actionName = strings.ToLower(actionName)

	// Supplier role restrictions
	if roleID == "supplier" {
		if resourceCode == "supplier" {
			// Supplier can only view and update their own profile, not create/delete
			return actionName == "view" || actionName == "update"
		}
		if resourceCode == "products" {
			// Supplier can do full CRUD on products
			return true
		}
		if resourceCode == "order" {
			// Supplier can only view orders, not create/update
			return actionName == "view"
		}
		if resourceCode == "bank-account" {
			// Supplier can do full CRUD on their own bank accounts
			return true
		}
	}

	// Customer role restrictions
	if roleID == "customer" {
		if resourceCode == "customer" {
			// Customer can only view and update their own profile, not create/delete
			return actionName == "view" || actionName == "update"
		}
		if resourceCode == "products" {
			// Customer can only view products, not create/update/delete
			return actionName == "view"
		}
		if resourceCode == "order" {
			// Customer can view and create orders, but not update (status updates are done by supplier/admin)
			return actionName == "view" || actionName == "create"
		}
	}

	// Superadmin gets everything (handled separately)
	if roleID == "superadmin" {
		return true
	}

	return false
}

func writeRoles(path string, roles []config.RoleDefinition) error {
	rolesFile := config.RolesFile{
		DefaultRoles: roles,
	}

	data, err := yaml.Marshal(&rolesFile)
	if err != nil {
		return fmt.Errorf("failed to marshal roles: %w", err)
	}

	// Add header comment
	header := `# Default Roles Configuration
# This file contains default roles with their permission patterns.
# These roles are bootstrapped during application initialization.
# 
# NOTE: This file can be updated by the distributor tool.
# Run 'go run cmd/distributor/main.go' to update roles with new resources.

`
	data = append([]byte(header), data...)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write roles file: %w", err)
	}

	return nil
}

