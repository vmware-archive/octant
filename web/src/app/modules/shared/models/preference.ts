export type PreferenceElement = RadioElement | TextElement;

export enum Operation {
  Equal = 'Equal',
}

export interface Condition {
  lhs: string;
  rhs: any;
  op: Operation;
}

interface Element {
  name: string;
  value: string;
  /**
   * If any of the disable conditions evaluate to true, this
   * element is disabled.
   */
  disableConditions?: Condition[];
}

export interface RadioElement extends Element {
  type: 'radio';
  config: {
    values: { label: string; value: string }[];
  };
}

export interface TextElement extends Element {
  type: 'text';
  config: {
    placeholder: string;
  };
}

export interface PreferenceSection {
  name: string;
  elements: PreferenceElement[];
}

export interface PreferencePanel {
  name: string;
  sections: PreferenceSection[];
}

export interface Preferences {
  panels: PreferencePanel[];
  updateName: string;
}
