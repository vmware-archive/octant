import { async, TestBed } from '@angular/core/testing';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';
import { ApplyYAMLComponent } from '../../components/smart/apply-yaml/apply-yaml.component';
import { FilterTextPipe } from './filtertext.pipe';

describe('FilterTextPipe', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ApplyYAMLComponent, OverlayScrollbarsComponent],
    });
  }));
  it('create an instance', () => {
    const pipe = new FilterTextPipe();
    expect(pipe).toBeTruthy();
  });
});
