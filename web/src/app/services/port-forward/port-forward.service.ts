import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Port } from 'src/app/models/content';
import getAPIBase from '../common/getAPIBase';

const API_BASE = getAPIBase();

@Injectable({
  providedIn: 'root'
})
export class PortForwardService {
  constructor(private http: HttpClient) { }

  public create(port: Port) {
    return this.http.post(`${API_BASE}/api/v1/content/overview/port-forwards`, {
      apiVersion: port.config.apiVersion,
      kind: port.config.kind,
      name: port.config.name,
      namespace: port.config.namespace,
      port: port.config.port,
    });
  }

  public remove(port: Port) {
    return this.http.delete(`${API_BASE}/api/v1/content/overview/port-forwards/${port.config.state.id}`);
  }
}
