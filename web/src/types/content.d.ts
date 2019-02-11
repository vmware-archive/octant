interface Metadata {
  type: string
  title?: ContentType[]
}

interface ContentType {
  metadata: Metadata
  config: any
}

type ListContentType = ContentType & {
  config: {
    items: ContentType[]
  }
}

type LinkContentType = ContentType & {
  config: {
    value: string
    ref: string
  }
}

type LabelsContentType = ContentType & {
  config: {
    labels: { [x: string]: string }
  }
}

interface GridPosition {
  x: number
  y: number
  w: number
  h: number
}
