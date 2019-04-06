export interface NavigationChild {
  title: string;
  path: string;
  children?: NavigationChild[];
}

export interface Navigation {
  sections: NavigationChild[];
}
