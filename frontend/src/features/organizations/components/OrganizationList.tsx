import { useState } from 'react';
import { Link } from 'react-router';
import { Lozenge } from '../../../shared/design-system';
import type { Organization, OrganizationRole } from '../api/organizationsApi';
import { useOrganizations } from '../hooks/useOrganizations';
import { organizationPath } from '../paths';
import { ListPanel } from './ListPanel';

const roleLabels: Record<OrganizationRole, string> = {
  owner: 'Organization Owner',
  contributor: 'Organization Contributor',
};

function OrganizationRow(organization: Organization) {
  return (
    <li key={organization.id} className="panel__row">
      <Link className="entity-list__link" to={organizationPath(organization.id)}>
        {organization.name}
      </Link>
      <Lozenge>{roleLabels[organization.currentUserRole]}</Lozenge>
    </li>
  );
}

/** The Organizations the signed-in User can see. */
export function OrganizationList() {
  const [page, setPage] = useState(1);
  const query = useOrganizations(page);
  return (
    <ListPanel
      title="Organizations"
      query={query}
      emptyTitle="No Organizations yet"
      emptyDescription="You will see an Organization here once an Organization Owner adds you to it."
      renderItem={OrganizationRow}
      onPageChange={setPage}
    />
  );
}
