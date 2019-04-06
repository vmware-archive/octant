import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotifierService {
  loading: BehaviorSubject<boolean> = new BehaviorSubject(false);
  error: BehaviorSubject<string> = new BehaviorSubject('');

  constructor() { }
}
