import { useId, useState, type SubmitEvent } from 'react';
import { ApiError, errorMessage } from '../../../shared/api/httpClient';
import { Alert, Button, TextField } from '../../../shared/design-system';
import type { WorkspaceVisibility } from '../api/organizationsApi';
import { useCreateWorkspace } from '../hooks/useOrganizations';
import { visibilityLabels } from './visibility';

const visibilityHelp: Record<WorkspaceVisibility, string> = {
  restricted: 'Only people added to this Workspace and Organization Owners can enter.',
  open: 'Everyone in the Organization can enter.',
};

const visibilities: WorkspaceVisibility[] = ['restricted', 'open'];

/** The message for the name field, when the failure concerns the name. */
function nameErrorOf(error: unknown): string | undefined {
  if (!(error instanceof ApiError)) {
    return undefined;
  }
  if (error.status === 409) {
    return 'Another Workspace in this Organization already uses this name.';
  }
  return error.problem?.errors?.find(({ field }) => field === 'name')?.message;
}

interface CreateWorkspaceFormProps {
  organizationId: string;
  onClose: () => void;
}

export function CreateWorkspaceForm({ organizationId, onClose }: CreateWorkspaceFormProps) {
  const create = useCreateWorkspace(organizationId);
  const visibilityLabelId = useId();
  const [name, setName] = useState('');
  const [visibility, setVisibility] = useState<WorkspaceVisibility>('restricted');

  const nameError = create.isError ? nameErrorOf(create.error) : undefined;

  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    create.mutate({ organizationId, name, visibility }, { onSuccess: onClose });
  }

  return (
    <form className="list-form" aria-label="New Workspace" onSubmit={handleSubmit}>
      {create.isError && !nameError && <Alert tone="danger">{errorMessage(create.error)}</Alert>}
      <TextField
        label="Name"
        name="name"
        required
        maxLength={255}
        value={name}
        error={nameError}
        onChange={(e) => {
          setName(e.target.value);
        }}
      />
      <div className="field">
        <span id={visibilityLabelId} className="field__label">
          Visibility
        </span>
        <div className="segmented" role="group" aria-labelledby={visibilityLabelId}>
          {visibilities.map((option) => (
            <button
              key={option}
              type="button"
              className="btn"
              aria-pressed={visibility === option}
              onClick={() => {
                setVisibility(option);
              }}
            >
              {visibilityLabels[option]}
            </button>
          ))}
        </div>
        <span className="field__help">{visibilityHelp[visibility]}</span>
      </div>
      <div className="list-form__actions">
        <Button type="submit" busy={create.isPending}>
          Create Workspace
        </Button>
        <Button variant="secondary" onClick={onClose}>
          Cancel
        </Button>
      </div>
    </form>
  );
}
