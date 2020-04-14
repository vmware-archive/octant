import { NgModule } from '@angular/core';
import { DocumentationRoutingModule } from './documentation-routing.module';
import { DocumentationComponent } from './documentation.component';
import { ComponentOverviewComponent } from './component-overview/component-overview.component';

@NgModule({
  imports: [DocumentationRoutingModule],
  declarations: [DocumentationComponent, ComponentOverviewComponent],
})
export class DocumentationModule {}
