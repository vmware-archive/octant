import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'truncate',
})
export class TruncatePipe implements PipeTransform {
  transform(name: string): string {
    if (name.length > 28) {
      return (
        name.substr(0, 10) + '...' + name.substr(name.length - 14, name.length)
      );
    }
    return name;
  }
}
