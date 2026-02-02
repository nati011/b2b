# Resource Distributor Tool

The distributor tool automatically distributes new resources and actions from `resources.yaml` to appropriate roles in `roles.yaml`.

## Usage

### Basic Usage

```bash
# Dry-run mode (preview changes without modifying files)
go run cmd/distributor/main.go --dry-run

# Apply changes to roles.yaml
go run cmd/distributor/main.go

# Specify custom paths
go run cmd/distributor/main.go \
  --resources config/resources.yaml \
  --roles config/roles.yaml \
  --output config/roles.yaml
```

### Command Line Options

- `--resources` (default: `config/resources.yaml`): Path to resources manifest file
- `--roles` (default: `config/roles.yaml`): Path to roles configuration file
- `--output` (default: same as `--roles`): Path to output updated roles file
- `--dry-run`: Preview changes without modifying files
- `--auto-superadmin` (default: `true`): Automatically assign all new resources to superadmin role

## How It Works

1. **Loads Resources**: Reads all resources and actions from `resources.yaml`
2. **Loads Roles**: Reads current role definitions from `roles.yaml`
3. **Analyzes Gaps**: Identifies new resources/actions not assigned to any role
4. **Distributes**: Assigns resources to roles based on:
   - **Superadmin**: Gets all new resources and actions automatically
   - **Supplier**: Gets supplier resources (view/update), products (full CRUD), orders (view)
   - **Customer**: Gets customer resources (view/update), products (view), orders (view/create)
5. **Updates File**: Writes updated roles back to `roles.yaml`

## Distribution Rules

### Super Administrator
- Gets **all** resources and actions automatically
- Full system access

### Supplier Role
- `supplier` resource: `view`, `update` (can view/update own profile)
- `products` resource: `view`, `create`, `update`, `delete` (full CRUD)
- `order` resource: `view` (can view orders for their products)

### Customer Role
- `customer` resource: `view`, `update` (can view/update own profile)
- `products` resource: `view` (can browse products)
- `order` resource: `view`, `create` (can view and create orders)

## Example Output

```
📋 Changes to be applied:
================================================================================

Role: Super Administrator
  Resource: auth_credentials
  New Actions: login
================================================================================

✓ Successfully updated roles file: config/roles.yaml
```

## When to Use

Run the distributor tool when:
- New resources are added to `resources.yaml` (via resourcegen)
- New actions are added to existing resources
- You want to ensure all roles have appropriate permissions
- After updating HTTP route annotations

## Integration with Bootstrap

The updated `roles.yaml` file is automatically used by the application bootstrap process (`internal/app/bootstrap.go`) when the application starts. The bootstrap process:
1. Loads resources from `resources.yaml`
2. Loads roles from `roles.yaml`
3. Creates/updates permissions and role assignments in the database

## Notes

- The tool is **idempotent** - safe to run multiple times
- Existing permissions are preserved
- Only missing permissions are added
- The tool respects role-specific restrictions (e.g., customers can't delete products)




