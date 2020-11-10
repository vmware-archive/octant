import { BrowserModule, DomSanitizer } from '@angular/platform-browser';
import { inject, TestBed } from '@angular/core/testing';
import { StringEscapePipe } from './string.escape.pipe';

describe('StringEscapePipe', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [BrowserModule],
    });
  });

  const testCases = [
    {
      input: `NOOP test`,
      expected: 'NOOP test',
    },
    {
      input: `<a>First test</a>`,
      expected: '&lt;a&gt;First test&lt;/a&gt;',
    },
    {
      input: '<div>><a>Multiple levels test</a></div>',
      expected:
        '&lt;div&gt;&gt;&lt;a&gt;Multiple levels test&lt;/a&gt;&lt;/div&gt;',
    },
    {
      input: 'Message\n More \n...',
      expected: 'Message\\n More \\n...',
    },
  ];

  it('Input strings are escaped properly', inject(
    [DomSanitizer],
    (domSanitizer: DomSanitizer) => {
      const pipe = new StringEscapePipe(domSanitizer);
      expect(pipe).toBeTruthy();

      testCases.forEach(testCase => {
        const sanitizedValue = pipe.transform(testCase.input);
        expect(sanitizedValue).toBe(testCase.expected);
      });
    }
  ));
});
