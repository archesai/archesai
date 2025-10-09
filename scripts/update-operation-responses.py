#!/usr/bin/env python3
"""
Update OpenAPI operation files to use response component refs.
"""

import os
import re
from pathlib import Path

# Mapping of operation patterns to their response components
RESPONSE_MAPPINGS = {
    # Users
    ('users', 'get.yaml'): ('UserListResponse', '../../components/responses/UserListResponse.yaml'),
    ('users/id', 'get.yaml'): ('UserResponse', '../../../components/responses/UserResponse.yaml'),
    ('users/id', 'patch.yaml'): ('UserResponse', '../../../components/responses/UserResponse.yaml'),

    # Me endpoint
    ('me', 'get.yaml'): ('UserResponse', '../../components/responses/UserResponse.yaml'),
    ('me', 'patch.yaml'): ('UserResponse', '../../components/responses/UserResponse.yaml'),

    # Organizations
    ('organizations', 'get.yaml'): ('OrganizationListResponse', '../../components/responses/OrganizationListResponse.yaml'),
    ('organizations/id', 'get.yaml'): ('OrganizationResponse', '../../../components/responses/OrganizationResponse.yaml'),
    ('organizations/id', 'patch.yaml'): ('OrganizationResponse', '../../../components/responses/OrganizationResponse.yaml'),
    ('organizations', 'post.yaml'): ('OrganizationResponse', '../../components/responses/OrganizationResponse.yaml'),

    # Pipelines
    ('pipelines', 'get.yaml'): ('PipelineListResponse', '../../components/responses/PipelineListResponse.yaml'),
    ('pipelines/id', 'get.yaml'): ('PipelineResponse', '../../../components/responses/PipelineResponse.yaml'),
    ('pipelines/id', 'patch.yaml'): ('PipelineResponse', '../../../components/responses/PipelineResponse.yaml'),
    ('pipelines', 'post.yaml'): ('PipelineResponse', '../../components/responses/PipelineResponse.yaml'),

    # Pipeline Steps
    ('pipelines/id/steps', 'get.yaml'): ('PipelineStepListResponse', '../../../../components/responses/PipelineStepListResponse.yaml'),

    # Pipeline Execution Plans
    ('pipelines/id/execution-plans', 'get.yaml'): ('PipelineExecutionPlanResponse', '../../../../components/responses/PipelineExecutionPlanResponse.yaml'),

    # Runs
    ('runs', 'get.yaml'): ('RunListResponse', '../../components/responses/RunListResponse.yaml'),
    ('runs/id', 'get.yaml'): ('RunResponse', '../../../components/responses/RunResponse.yaml'),
    ('runs/id', 'patch.yaml'): ('RunResponse', '../../../components/responses/RunResponse.yaml'),
    ('runs', 'post.yaml'): ('RunResponse', '../../components/responses/RunResponse.yaml'),

    # Tools
    ('tools', 'get.yaml'): ('ToolListResponse', '../../components/responses/ToolListResponse.yaml'),
    ('tools/id', 'get.yaml'): ('ToolResponse', '../../../components/responses/ToolResponse.yaml'),
    ('tools/id', 'patch.yaml'): ('ToolResponse', '../../../components/responses/ToolResponse.yaml'),
    ('tools', 'post.yaml'): ('ToolResponse', '../../components/responses/ToolResponse.yaml'),

    # Artifacts
    ('artifacts', 'get.yaml'): ('ArtifactListResponse', '../../components/responses/ArtifactListResponse.yaml'),
    ('artifacts/id', 'get.yaml'): ('ArtifactResponse', '../../../components/responses/ArtifactResponse.yaml'),
    ('artifacts/id', 'patch.yaml'): ('ArtifactResponse', '../../../components/responses/ArtifactResponse.yaml'),
    ('artifacts', 'post.yaml'): ('ArtifactResponse', '../../components/responses/ArtifactResponse.yaml'),

    # Labels
    ('labels', 'get.yaml'): ('LabelListResponse', '../../components/responses/LabelListResponse.yaml'),
    ('labels/id', 'get.yaml'): ('LabelResponse', '../../../components/responses/LabelResponse.yaml'),
    ('labels/id', 'patch.yaml'): ('LabelResponse', '../../../components/responses/LabelResponse.yaml'),
    ('labels', 'post.yaml'): ('LabelResponse', '../../components/responses/LabelResponse.yaml'),

    # API Keys
    ('api-keys', 'get.yaml'): ('APIKeyListResponse', '../../components/responses/APIKeyListResponse.yaml'),
    ('api-keys/id', 'get.yaml'): ('APIKeyResponse', '../../../components/responses/APIKeyResponse.yaml'),
    ('api-keys/id', 'patch.yaml'): ('APIKeyResponse', '../../../components/responses/APIKeyResponse.yaml'),
    ('api-keys', 'post.yaml'): ('APIKeyResponse', '../../components/responses/APIKeyResponse.yaml'),

    # Members
    ('organizations/organizationID/members', 'get.yaml'): ('MemberListResponse', '../../../../components/responses/MemberListResponse.yaml'),
    ('organizations/organizationID/members/id', 'get.yaml'): ('MemberResponse', '../../../../../components/responses/MemberResponse.yaml'),
    ('organizations/organizationID/members/id', 'patch.yaml'): ('MemberResponse', '../../../../../components/responses/MemberResponse.yaml'),
    ('organizations/organizationID/members', 'post.yaml'): ('MemberResponse', '../../../../components/responses/MemberResponse.yaml'),

    # Invitations
    ('organizations/organizationID/invitations', 'get.yaml'): ('InvitationListResponse', '../../../../components/responses/InvitationListResponse.yaml'),
    ('organizations/organizationID/invitations/id', 'get.yaml'): ('InvitationResponse', '../../../../../components/responses/InvitationResponse.yaml'),
    ('organizations/organizationID/invitations/id', 'patch.yaml'): ('InvitationResponse', '../../../../../components/responses/InvitationResponse.yaml'),
    ('organizations/organizationID/invitations', 'post.yaml'): ('InvitationResponse', '../../../../components/responses/InvitationResponse.yaml'),

    # Sessions
    ('auth/sessions', 'get.yaml'): ('SessionListResponse', '../../../components/responses/SessionListResponse.yaml'),
    ('auth/sessions/id', 'get.yaml'): ('SessionResponse', '../../../../components/responses/SessionResponse.yaml'),
    ('auth/sessions/id', 'patch.yaml'): ('SessionResponse', '../../../../components/responses/SessionResponse.yaml'),

    # Accounts
    ('auth/accounts', 'get.yaml'): ('AccountListResponse', '../../../components/responses/AccountListResponse.yaml'),
    ('auth/accounts/id', 'get.yaml'): ('AccountResponse', '../../../../components/responses/AccountResponse.yaml'),
    ('auth/accounts/id', 'patch.yaml'): ('AccountResponse', '../../../../components/responses/AccountResponse.yaml'),
}


def update_operation_file(file_path: Path):
    """Update an operation file to use response refs."""
    content = file_path.read_text()

    # Find the pattern for this file
    relative_path = str(file_path.relative_to('api/operations'))
    file_name = file_path.name

    # Try to find a matching pattern
    response_ref = None
    for (pattern, fname), (name, ref) in RESPONSE_MAPPINGS.items():
        if pattern in relative_path and fname == file_name:
            response_ref = ref
            break

    if not response_ref:
        return False

    # Find the '200' or '201' response and replace it
    # Pattern to match the full 200/201 response block
    pattern = r"  '2(00|01)':\s*\n(?:    .*\n)*?(?=  '|\Z)"

    def replace_response(match):
        status_code = match.group(1)
        # Special case for 201 Created (login endpoint)
        if status_code == '01' and 'login' in str(file_path):
            return f"  '201':\n    $ref: {response_ref.replace('SessionResponse', 'SessionCreated')}\n"
        return f"  '2{status_code}':\n    $ref: {response_ref}\n"

    new_content = re.sub(pattern, replace_response, content)

    if new_content != content:
        file_path.write_text(new_content)
        return True
    return False


def main():
    """Main function to update all operation files."""
    operations_dir = Path('api/operations')
    updated_count = 0

    for yaml_file in operations_dir.rglob('*.yaml'):
        if update_operation_file(yaml_file):
            print(f"Updated: {yaml_file}")
            updated_count += 1

    print(f"\nTotal files updated: {updated_count}")


if __name__ == '__main__':
    main()
