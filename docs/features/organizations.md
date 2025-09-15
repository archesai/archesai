# Organizations

## Overview

Organizations in ArchesAI provide multi-tenancy support, enabling teams to collaborate on projects
with proper isolation, member management, and billing integration.

## Core Concepts

### Multi-Tenancy Architecture

- **Complete Data Isolation**: Each organization's data is fully isolated
- **Resource Scoping**: All resources (workflows, content, settings) belong to an organization
- **Cross-Organization Access**: Users can belong to multiple organizations
- **Organization Switching**: Seamless context switching between organizations

### Organization Structure

```typescript
interface Organization {
  id: string;
  name: string;
  slug: string;
  billingEmail: string;
  plan: "free" | "pro" | "enterprise";
  members: Member[];
  settings: OrganizationSettings;
  createdAt: Date;
  updatedAt: Date;
}
```

## Member Management

### Roles and Permissions

#### Owner

- Full organization control
- Billing management
- Delete organization
- Transfer ownership

#### Admin

- Member management
- Settings configuration
- Integration management
- Cannot delete organization

#### Member

- Create and manage own resources
- View organization resources
- Limited settings access

#### Viewer

- Read-only access
- Cannot create resources
- Cannot modify settings

### Invitation System

```typescript
// Invite member
POST /api/v1/organizations/:id/invitations
{
  "email": "new.member@example.com",
  "role": "member",
  "message": "Welcome to our team!"
}

// Accept invitation
POST /api/v1/invitations/:token/accept

// List pending invitations
GET /api/v1/organizations/:id/invitations
```

### Member Operations

```typescript
// List organization members
GET /api/v1/organizations/:id/members

// Update member role
PUT /api/v1/organizations/:id/members/:userID
{
  "role": "admin"
}

// Remove member
DELETE /api/v1/organizations/:id/members/:userID
```

## Billing Integration

### Subscription Plans

#### Free Plan

- Up to 3 members
- 10 GB storage
- 1,000 API calls/month
- Community support

**Pro Plan** ($49/month)

- Up to 20 members
- 100 GB storage
- 50,000 API calls/month
- Email support
- Advanced features

**Enterprise Plan** (Custom)

- Unlimited members
- Custom storage
- Unlimited API calls
- Priority support
- SLA guarantee
- Custom integrations

### Billing Management

```typescript
// Get billing information
GET /api/v1/organizations/:id/billing

// Update payment method
PUT /api/v1/organizations/:id/billing/payment-method
{
  "stripePaymentMethodId": "pm_1234567890"
}

// Get invoices
GET /api/v1/organizations/:id/billing/invoices

// Change plan
POST /api/v1/organizations/:id/billing/upgrade
{
  "plan": "pro"
}
```

### Usage Tracking

```typescript
interface UsageMetrics {
  storage: {
    used: number;
    limit: number;
    unit: "GB";
  };
  apiCalls: {
    used: number;
    limit: number;
    period: "month";
  };
  members: {
    active: number;
    limit: number;
  };
}
```

## Organization Settings

### General Settings

- Organization name and slug
- Logo and branding
- Default timezone
- Language preferences

### Security Settings

- Two-factor authentication requirements
- IP allowlisting
- Session timeout policies
- Audit log retention

### Integration Settings

- Webhook endpoints
- API keys management
- Third-party integrations
- SSO configuration

## API Endpoints

### Organization Management

- `GET /api/v1/organizations` - List user's organizations
- `POST /api/v1/organizations` - Create organization
- `GET /api/v1/organizations/:id` - Get organization details
- `PUT /api/v1/organizations/:id` - Update organization
- `DELETE /api/v1/organizations/:id` - Delete organization

### Member Management

- `GET /api/v1/organizations/:id/members` - List members
- `POST /api/v1/organizations/:id/invitations` - Invite member
- `PUT /api/v1/organizations/:id/members/:userID` - Update member role
- `DELETE /api/v1/organizations/:id/members/:userID` - Remove member

### Settings

- `GET /api/v1/organizations/:id/settings` - Get settings
- `PUT /api/v1/organizations/:id/settings` - Update settings
- `GET /api/v1/organizations/:id/audit-logs` - Get audit logs

## Implementation Details

### Database Schema

```sql
-- Organizations table
CREATE TABLE organizations (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  billing_email VARCHAR(255) NOT NULL,
  plan VARCHAR(50) DEFAULT 'free',
  settings JSONB,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- Organization members
CREATE TABLE organization_members (
  organization_id UUID REFERENCES organizations (id),
  user_id UUID REFERENCES users (id),
  role VARCHAR(50) NOT NULL,
  joined_at TIMESTAMP NOT NULL,
  PRIMARY KEY (organization_id, user_id)
);

-- Invitations
CREATE TABLE organization_invitations (
  id UUID PRIMARY KEY,
  organization_id UUID REFERENCES organizations (id),
  email VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL,
  token VARCHAR(255) UNIQUE NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL
);
```

### Service Layer

```go
type OrganizationService struct {
    repo     OrganizationRepository
    billing  BillingService
    mailer   MailerService
    logger   *slog.Logger
}

func (s *OrganizationService) Create(ctx context.Context, input CreateOrganizationInput) (*Organization, error) {
    // Validate input
    // Create organization
    // Set up default settings
    // Add creator as owner
    // Initialize billing
    // Send welcome email
}

func (s *OrganizationService) InviteMember(ctx context.Context, orgID uuid.UUID, input InviteMemberInput) error {
    // Check permissions
    // Validate member limits
    // Create invitation
    // Send invitation email
}
```

## Workflows and Automation

### Organization Lifecycle

#### Creation Flow

1. User creates organization
2. Default settings applied
3. Creator becomes owner
4. Billing initialized
5. Welcome email sent

#### Deletion Flow

1. Confirm owner identity
2. Export data (optional)
3. Cancel subscriptions
4. Delete all resources
5. Remove all members
6. Send confirmation

### Automated Tasks

- Daily usage calculation
- Monthly billing cycles
- Invitation expiry cleanup
- Audit log rotation
- Storage quota monitoring

## Best Practices

### Scalability

- Use database partitioning by organization ID
- Implement caching for frequently accessed data
- Use connection pooling per organization
- Consider sharding for large deployments

### Security

- Always validate organization context
- Implement row-level security
- Audit all administrative actions
- Encrypt sensitive organization data

### Performance

- Index on organization_id for all tables
- Cache organization settings
- Batch member operations
- Use pagination for member lists

## Monitoring and Analytics

### Key Metrics

- Organizations created/deleted per day
- Member growth rate
- Average organization size
- Plan distribution
- Feature usage by plan

### Alerts

- Unusual deletion activity
- Rapid member additions
- Storage quota exceeded
- Payment failures
- Invitation spam detection

## Testing

### Unit Tests

```bash
go test ./internal/organizations/...
```

### Integration Tests

```bash
go test -tags=integration ./internal/organizations/...
```

### Test Scenarios

- Organization CRUD operations
- Member invitation flow
- Role permission validation
- Billing integration
- Multi-organization access

## Troubleshooting

### Common Issues

#### Cannot Create Organization

- Check user email verification
- Verify organization limit not exceeded
- Ensure unique organization slug

#### Invitation Not Received

- Check spam folder
- Verify email address
- Check invitation expiry
- Resend invitation

#### Billing Issues

- Verify payment method
- Check subscription status
- Review usage limits
- Contact support

## Migration Guide

### Importing Organizations

1. Prepare CSV with organization data
2. Validate data format
3. Run import script
4. Verify member associations
5. Set up billing

### Exporting Organizations

1. Request data export
2. Select export format
3. Include member list
4. Export audit logs
5. Download archive

## Related Documentation

- [Authentication](./auth.md) - User authentication
- [Workflows](./workflows.md) - Workflow automation
- [Billing Guide](../deployment/production.md#billing) - Billing setup
