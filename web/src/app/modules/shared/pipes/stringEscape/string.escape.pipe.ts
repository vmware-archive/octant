import { Pipe, PipeTransform, SecurityContext } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Pipe({
  name: 'escapepipe',
})
export class StringEscapePipe implements PipeTransform {
  constructor(private sanitizer: DomSanitizer) {}

  public transform(value: string): SafeHtml {
    return this.sanitizer.sanitize(
      SecurityContext.HTML,
      this.escapePipe(value)
    );
  }

  escapePipe(str: string): string {
    return str
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/\n/g, '\\n');
  }
}
