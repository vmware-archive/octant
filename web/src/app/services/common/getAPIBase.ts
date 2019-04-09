import { environment } from 'src/environments/environment';

// TODO: API_BASE this should be configurable
export default function getAPIBase(): string {
  if (environment.production) {
    return window.location.origin;
  }
  return 'http://localhost:3001';
}
