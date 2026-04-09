const _host = typeof window !== 'undefined' ? window.location.hostname : 'localhost';
export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? `http://${_host}:8081`;
export const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL ?? `ws://${_host}:8081/ws`;

export function generateUUID(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID();
  }
  // Fallback for non-secure contexts (HTTP on local network)
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16);
  });
}
