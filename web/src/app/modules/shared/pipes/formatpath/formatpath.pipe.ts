import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'formatpath',
})
export class FormatPathPipe implements PipeTransform {
  transform(path: string): string {
    if (!path.startsWith('/')) {
      return '/' + path;
    }

    return path;
  }
}
