import { ResourceObject } from './schema'

export default class ResourceNode {
  constructor(
    private readonly object: ResourceObject,
    private readonly isSelected: boolean,
  ) {}

  toDescriptor(): any {
    let nodeClass = `node-${this.object.status}`
    if (this.isSelected) {
      nodeClass += ` selected`
    }

    return {
      label: `${this.title()}${this.subTitle()}`,
      labelType: 'html',
      class: `${nodeClass}`,
    }
  }

  title(): string {
    return `<div class="resource-name">${this.object.name}</div>`
  }

  subTitle(): string {
    return `<div class="resource-type">${this.object.apiVersion} ${
      this.object.kind
    }</div>`
  }
}
