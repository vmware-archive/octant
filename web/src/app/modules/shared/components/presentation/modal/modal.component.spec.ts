// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ModalComponent } from './modal.component';
import { SharedModule } from '../../../shared.module';
import { ModalService } from '../../../services/modal/modal.service';
import { windowProvider, WindowToken } from '../../../../../window';
import { OctantTooltipComponent } from '../octant-tooltip/octant-tooltip';

describe('ModalComponent', () => {
  let component: ModalComponent;
  let fixture: ComponentFixture<ModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [SharedModule],
      declarations: [ModalComponent, OctantTooltipComponent],
      providers: [
        { provide: ModalService },
        { provide: WindowToken, useFactory: windowProvider },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
