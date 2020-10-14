// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { UploaderComponent } from './uploader.component';
import { windowProvider, WindowToken } from '../../../../../window';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { IndicatorComponent } from 'src/app/modules/shared/components/presentation/indicator/indicator.component';
import { EditorComponent } from 'src/app/modules/shared/components/smart/editor/editor.component';

describe('UploaderComponent', () => {
  let component: UploaderComponent;
  let fixture: ComponentFixture<UploaderComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [UploaderComponent, IndicatorComponent, EditorComponent],
        imports: [SharedModule],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(UploaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
