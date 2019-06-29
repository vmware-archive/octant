import { Injectable } from '@angular/core';
import getAPIBase from '../../../../services/common/getAPIBase';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ActionService {
  constructor(private http: HttpClient) {}

  perform(update: any) {
    const url = [getAPIBase(), 'api/v1/action'].join('/');

    const payload = {
      update: update,
    };

    return this.http.post(url, payload);
  }
}
