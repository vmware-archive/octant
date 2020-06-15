import { Injectable, NgModule } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { ContainerComponent } from './components/smart/container/container.component';
import { NamespaceComponent } from './components/smart/namespace/namespace.component';
import { PageNotFoundComponent } from './components/smart/page-not-found/page-not-found.component';
import { InputFilterComponent } from './components/smart/input-filter/input-filter.component';
import { NotifierComponent } from './components/smart/notifier/notifier.component';
import { NavigationComponent } from './components/smart/navigation/navigation.component';
import { QuickSwitcherComponent } from './components/smart/quick-switcher/quick-switcher.component';
import { ThemeSwitchButtonComponent } from './components/smart/theme-switch/theme-switch-button.component';
import { UploaderComponent } from './components/smart/uploader/uploader.component';
import { ClarityModule } from '@clr/angular';
import { HttpClientModule } from '@angular/common/http';
import { RouterModule, Routes } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { NgSelectModule } from '@ng-select/ng-select';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';
import { SharedModule } from '../shared/shared.module';
import { ContentComponent } from './components/smart/content/content.component';
import { FilterTextPipe } from './pipes/filtertext/filtertext.pipe';

const routes: Routes = [
  {
    path: '',
    component: ContainerComponent,
    children: [{ path: '**', component: ContentComponent }],
  },
];

@Injectable()
export class UnstripTrailingSlashLocation extends Location {
  public static stripTrailingSlash(url: string): string {
    return url;
  }
}

@NgModule({
  declarations: [
    ContainerComponent,
    NamespaceComponent,
    PageNotFoundComponent,
    InputFilterComponent,
    NotifierComponent,
    NavigationComponent,
    ContentComponent,
    QuickSwitcherComponent,
    ThemeSwitchButtonComponent,
    UploaderComponent,
    FilterTextPipe,
  ],
  imports: [
    CommonModule,
    ClarityModule,
    HttpClientModule,
    FormsModule,
    NgSelectModule,
    MarkdownModule.forRoot({
      markedOptions: {
        provide: MarkedOptions,
        useValue: {
          gfm: true,
          tables: true,
          breaks: true,
          pedantic: false,
          sanitize: false,
          smartLists: true,
          smartypants: false,
        },
      },
    }),
    SharedModule,

    // routing must come last
    RouterModule.forChild(routes),
  ],
  providers: [
    {
      provide: Location,
      useClass: UnstripTrailingSlashLocation,
    },
  ],
})
export class SugarloafModule {}
