import { renderHook, waitFor } from '@testing-library/react';
import { useHealth } from './useHealth';

describe('useHealth', () => {
  it('starts loading and becomes online when the API answers', async () => {
    const fetchHealth = vi.fn().mockResolvedValue({ status: 'ok' });
    const { result } = renderHook(() => useHealth(fetchHealth));

    expect(result.current).toEqual({ kind: 'loading' });
    await waitFor(() => {
      expect(result.current).toEqual({ kind: 'online' });
    });
  });

  it('becomes offline with the error message when the request fails', async () => {
    const fetchHealth = vi.fn().mockRejectedValue(new Error('Network down'));
    const { result } = renderHook(() => useHealth(fetchHealth));

    await waitFor(() => {
      expect(result.current).toEqual({ kind: 'offline', message: 'Network down' });
    });
  });

  it('uses a generic message for non-Error rejections', async () => {
    const fetchHealth = vi.fn().mockRejectedValue('nope');
    const { result } = renderHook(() => useHealth(fetchHealth));

    await waitFor(() => {
      expect(result.current).toEqual({ kind: 'offline', message: 'Unknown error' });
    });
  });

  it('aborts the request on unmount and ignores the late failure', async () => {
    let rejectRequest: (reason: unknown) => void = () => undefined;
    let receivedSignal: AbortSignal | undefined;
    const fetchHealth = vi.fn((signal: AbortSignal) => {
      receivedSignal = signal;
      return new Promise<{ status: string }>((_resolve, reject) => {
        rejectRequest = reject;
      });
    });
    const { result, unmount } = renderHook(() => useHealth(fetchHealth));

    unmount();
    rejectRequest(new Error('aborted'));
    await Promise.resolve();

    expect(receivedSignal?.aborted).toBe(true);
    expect(result.current).toEqual({ kind: 'loading' });
  });
});
