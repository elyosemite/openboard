import { Link, useParams } from 'react-router';
import { errorMessage } from '../../../shared/api/httpClient';
import { Alert, Card, Spinner } from '../../../shared/design-system';
import { useWorkspace } from '../hooks/useOrganizations';
import { organizationPath } from '../paths';
import './Organizations.css';
import { visibilityLabels } from './visibility';
import { VisibilityLozenge } from './VisibilityLozenge';

const dateFormat = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' });

export function WorkspacePage() {
  const { workspaceId = '' } = useParams();
  const workspace = useWorkspace(workspaceId);

  if (workspace.isPending) {
    return <Spinner />;
  }
  if (workspace.isError) {
    return <Alert tone="danger">{errorMessage(workspace.error)}</Alert>;
  }
  const { name, visibility, createdAt, organizationId } = workspace.data;
  return (
    <div className="page-stack">
      <header className="page-header">
        <h1 className="page-header__title">{name}</h1>
        <VisibilityLozenge visibility={visibility} />
      </header>
      <Card title="Details">
        <dl className="details">
          <dt>Visibility</dt>
          <dd>{visibilityLabels[visibility]}</dd>
          <dt>Created</dt>
          <dd>{dateFormat.format(new Date(createdAt))}</dd>
        </dl>
      </Card>
      <Link className="btn btn--link" to={organizationPath(organizationId)}>
        Back to Organization
      </Link>
    </div>
  );
}
