interface ContentType {
  type: string;
  label: string;
  data: {
    value: string;
  }
}

type ListContentType = ContentType & {
  data: {
    items: ContentType[];
  };
}

type LinkContentType = ContentType & {
  data: {
    value: string;
    ref: string;
  };
}

interface ContentSection {
  title: string;
  items: ContentType[];
}

interface BaseContent {
  type: string;
  title: string;
}

type ContentSummary = BaseContent & {
  sections: ContentSection[];
}

type ContentTable = BaseContent

type Content = ContentSummary | ContentTable
