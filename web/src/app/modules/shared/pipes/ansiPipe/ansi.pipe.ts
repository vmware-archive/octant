import { Pipe, PipeTransform, SecurityContext } from '@angular/core';
import {
  DomSanitizer,
  SafeHtml,
  SafeResourceUrl,
  SafeScript,
  SafeStyle,
  SafeUrl,
} from '@angular/platform-browser';
import { default as AnsiUp } from 'ansi_up';

@Pipe({
  name: 'ansipipe',
})
export class AnsiPipe implements PipeTransform {
  ansiUp = new AnsiUp();

  constructor(private sanitizer: DomSanitizer) {
    this.ansiUp.escape_for_html = false;
  }

  public transform(
    value: string
  ): SafeHtml | SafeStyle | SafeScript | SafeUrl | SafeResourceUrl {
    
    if (value.includes('\x1B')) {  // ANSI string
      return this.sanitizer.bypassSecurityTrustHtml(
        this.ansiUp.ansi_to_html(value)
      );
    }
    return value;
  }
}
