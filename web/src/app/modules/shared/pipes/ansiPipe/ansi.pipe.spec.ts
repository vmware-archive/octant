import { BrowserModule, DomSanitizer } from '@angular/platform-browser';
import { inject, TestBed } from '@angular/core/testing';
import { AnsiPipe } from './ansi.pipe';
import { SecurityContext } from '@angular/core';

describe('AnsiPipe', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [BrowserModule],
    });
  });

  it('create an instance', inject(
    [DomSanitizer],
    (domSanitizer: DomSanitizer) => {
      const pipe = new AnsiPipe(domSanitizer);
      expect(pipe).toBeTruthy();

      const ansiString = 'Hello \x1B[1;33;41m[o]\x1B[0m';
      const sanitizedValue = pipe.transform(ansiString);

      expect(domSanitizer.sanitize(SecurityContext.HTML, sanitizedValue)).toBe(
        'Hello <span style="font-weight:bold;color:rgb(187,187,0);background-color:rgb(187,0,0)">[o]</span>'
      );
    }
  ));
});
