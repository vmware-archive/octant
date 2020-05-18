import { withKnobs } from '@storybook/addon-knobs';
import { addDecorator, moduleMetadata } from '@storybook/angular';
import { SharedModule } from '../src/app/modules/shared/shared.module';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';

// enable knobs
addDecorator(withKnobs);

addDecorator(
  moduleMetadata({
    imports: [
      SharedModule,
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
    ],
  })
);
