import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ClarityModule } from '@clr/angular';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import {
  MonacoEditorModule,
  MonacoEditorConfig,
  MonacoProviderService,
} from 'ng-monaco-editor';
import { DocumentationModule } from './documentation/documentation.module';

@NgModule({
  declarations: [AppComponent, HomeComponent],
  imports: [
    ClarityModule,
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    DocumentationModule,
    MonacoEditorModule.forRoot({
      // Angular CLI currently does not handle assets with hashes. We manage it by manually adding
      // version numbers to force library updates:
      baseUrl: 'lib',
      defaultOptions: {},
    }),
    AppRoutingModule,
  ],
  providers: [MonacoProviderService, MonacoEditorConfig],
  bootstrap: [AppComponent],
})
export class AppModule {}
