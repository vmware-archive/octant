import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { EditorDemoComponent } from './editor.demo';
import { ApiEditorDemoComponent } from './api-editor.demo';
import { AngularEditorDemoComponent } from './angular-editor.demo';

import { MonacoEditorModule } from 'ng-monaco-editor';
import { MonacoEditorConfig, MonacoProviderService } from 'ng-monaco-editor';
import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    MonacoEditorModule.forRoot({
      baseUrl: 'lib',
      defaultOptions: {},
    }),
    FormsModule,
    RouterModule.forChild([{ path: '', component: EditorDemoComponent }]),
  ],
  providers: [MonacoProviderService, MonacoEditorConfig],
  declarations: [
    AngularEditorDemoComponent,
    EditorDemoComponent,
    ApiEditorDemoComponent,
  ],
  exports: [EditorDemoComponent],
})
export class EditorDemoModule {}
