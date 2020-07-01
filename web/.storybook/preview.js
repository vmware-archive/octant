import {addDecorator, moduleMetadata} from '@storybook/angular';
import { SharedModule } from '../src/app/modules/shared/shared.module';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';
import { setConsoleOptions } from '@storybook/addon-console';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {APP_INITIALIZER} from "@angular/core";
import {InitService} from "../src/app/modules/shared/services/init/init.service";
import {RouterModule} from "@angular/router";
import {HttpClientModule} from "@angular/common/http";
import {BrowserModule} from "@angular/platform-browser";
import {AppRoutingModule} from "../src/app/app-routing.module";
import { setCompodocJson } from '@storybook/addon-docs/angular';
import docJson from '../documentation.json';
import { themes } from '@storybook/theming';

setCompodocJson(docJson);

setConsoleOptions({
  panelExclude: [
    /Angular is running in the development mode/,
    /Ignored an update to unaccepted module/,
  ],
});

addDecorator(
  moduleMetadata({
    imports: [
      AppRoutingModule,
      BrowserAnimationsModule,
      BrowserModule,
      HttpClientModule,
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
      RouterModule,
      SharedModule,
    ],
    providers: [
      InitService,
      {
        provide: APP_INITIALIZER,
        useFactory: (initService) => () => initService.init(),
        deps: [InitService],
        multi: true
      },
    ]
  })
);

export const parameters = {
  docs: {
    theme: themes.light,
  },
  options: {
    storySort: (a, b) => { // Show component stories on top
      const first= a[1].kind.substring(0, a[1].kind.indexOf('/'));
      const second= b[1].kind.substring(0, b[1].kind.indexOf('/'));
      return first === second ? a[1].id.localeCompare(b[1].id, { numeric: true }) : a[1].kind.localeCompare(b[1].kind, { numeric: true });
    }
  },
};
