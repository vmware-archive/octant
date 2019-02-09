interface NamespaceOption {
  label: string
  value: string
}

interface NavigationSectionType {
  title: string
  path: string
  children: NavigationSectionType[]
}

interface Navigation {
  sections: NavigationSectionType[]
}
