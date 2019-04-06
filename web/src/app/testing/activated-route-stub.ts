import { convertToParamMap, ParamMap, Params, UrlSegment } from '@angular/router';
import { ReplaySubject } from 'rxjs';

/**
 * An ActivateRoute test double with a `paramMap` observable.
 * Use the `setParamMap()` method to add the next `paramMap` value.
 */
export class ActivatedRouteStub {

  constructor(initialParams?: Params) {
    this.setParamMap(initialParams);
  }
  // Use a ReplaySubject to share previous values with subscribers
  // and pump new values into the `paramMap` observable
  private subject = new ReplaySubject<ParamMap>();

  /** The mock paramMap observable */
  readonly paramMap = this.subject.asObservable();

  private urlSubject = new ReplaySubject<UrlSegment[]>();
  readonly url = this.urlSubject.asObservable();

  /** Set the paramMap observable's next value */
  setParamMap(params?: Params) {
    this.subject.next(convertToParamMap(params));
  }
}
