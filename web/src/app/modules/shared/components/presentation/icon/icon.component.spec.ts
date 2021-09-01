import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IconView } from '../../../models/content';

import { IconComponent } from './icon.component';

describe('IconComponent', () => {
  let component: IconComponent;
  let fixture: ComponentFixture<IconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [IconComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('uses custom icon when customIcon svg is valid', () => {
    component.view = {
      metadata: {
        type: 'icon',
      },
      config: {
        shape: 'custom',
        solid: true,
        customSvg: '<svg></svg>',
      },
    } as IconView;
    fixture.detectChanges();

    let view = component.view as IconView;
    expect(view.config.shape).toEqual('custom');
  });

  it('uses placeholder icon when customIcon svg is invalid', () => {
    component.view = {
      metadata: {
        type: 'icon',
      },
      config: {
        shape: 'custom',
        solid: true,
        customSvg: '<svg><circle><foo>',
      },
    } as IconView;
    fixture.detectChanges();

    let view = component.view as IconView;
    expect(view.config.shape).toEqual('times');
  });
});
