import { Injectable, NgModule, NgZone } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { ContainerComponent } from './components/smart/container/container.component';
import { NamespaceComponent } from '../../components/namespace/namespace.component';
import { PageNotFoundComponent } from '../../components/page-not-found/page-not-found.component';
import { InputFilterComponent } from '../../components/input-filter/input-filter.component';
import { NotifierComponent } from '../../components/notifier/notifier.component';
import { NavigationComponent } from '../../components/navigation/navigation.component';
import { QuickSwitcherComponent } from '../../components/quick-switcher/quick-switcher.component';
import { ThemeSwitchButtonComponent } from '../overview/components/theme-switch/theme-switch-button.component';
import { ClarityModule } from '@clr/angular';
import { HttpClientModule } from '@angular/common/http';
import { Router, RouterModule, Routes } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { NgSelectModule } from '@ng-select/ng-select';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';
import { SharedModule } from '../../shared/shared.module';
import { OverviewComponent } from './components/smart/overview/overview.component';

const routes: Routes = [
  {
    path: '',
    component: ContainerComponent,
    children: [{ path: '**', component: OverviewComponent }],
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
    OverviewComponent,
    QuickSwitcherComponent,
    ThemeSwitchButtonComponent,
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
