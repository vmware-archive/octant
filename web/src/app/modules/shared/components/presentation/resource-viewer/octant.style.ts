import { Stylesheet } from 'cytoscape';

const svgHealthy =
  color => `<svg version="1.1" width="36" height="36" viewBox="0 0 36 36" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">
    <title>check-circle-solid</title>
    <path d="M30,18A12,12,0,1,1,18,6,12,12,0,0,1,30,18Zm-4.77-2.16a1.4,1.4,0,0,0-2-2l-6.77,6.77L13,17.16a1.4,1.4,0,0,0-2,2l5.45,5.45Z" fill="${color}"></path>
    <rect x="0" y="0" width="36" height="36" fill-opacity="0"/></svg>`;

const svgWarning =
  color => `<svg version="1.1" width="36" height="36" viewBox="0 0 36 36" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">
    <title>exclamation-triangle-solid</title>
    <path d="M30.33,25.54,20.59,7.6a3,3,0,0,0-5.27,0L5.57,25.54A3,3,0,0,0,8.21,30H27.69a3,3,0,0,0,2.64-4.43ZM16.46,12.74a1.49,1.49,0,0,1,3,0v6.89a1.49,1.49,0,1,1-3,0ZM18,26.25a1.72,1.72,0,1,1,1.72-1.72A1.72,1.72,0,0,1,18,26.25Z" fill="${color}"></path>
    <rect x="0" y="0" width="36" height="36" fill-opacity="0"/></svg>`;

const svgError =
  color => `<svg version="1.1" width="36" height="36" viewBox="0 0 36 36" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">
    <title>exclamation-circle-solid</title>
    <path d="M18,6A12,12,0,1,0,30,18,12,12,0,0,0,18,6Zm-1.49,6a1.49,1.49,0,0,1,3,0v6.89a1.49,1.49,0,1,1-3,0ZM18,25.5a1.72,1.72,0,1,1,1.72-1.72A1.72,1.72,0,0,1,18,25.5Z" fill="${color}"></path>
    <rect x="0" y="0" width="36" height="36" fill-opacity="0"/></svg>`;

const renderImageLight = ele => renderImage(ele, 'black');
const renderImageDark = ele => renderImage(ele, 'white');

const renderImage = (ele, color) => {
  const status = ele.data('status');

  const image =
    status === 'ok'
      ? svgHealthy(color)
      : status === 'error'
      ? svgError(color)
      : svgWarning(color);
  return 'data:image/svg+xml;base64,' + btoa(image);
};

export const ELEMENTS_STYLE: Stylesheet[] = [
  {
    selector: 'node',
    css: {
      shape: 'round-rectangle',
      width: 300,
      height: 100,
      color: 'black',

      'background-color': node => nodeColor(node),
      'background-opacity': 1,
      'overlay-opacity': 0,

      'background-image': renderImageLight,
      'background-width': '36px',
      'background-height': '36px',
      'background-position-x': node => imageOffset(node),
      'background-position-y': '8px',

      'border-color': node => nodeColor(node),
      'border-width': '2px',
      'border-style': 'solid',

      'padding-right': '10px',
      'padding-top': '10px',
      'padding-bottom': '10px',
      'z-index': 1,
    },
  },
  {
    selector: 'node:selected',
    css: {
      'curve-style': 'bezier',
      'border-width': 3,
      'border-color': '#0065ab',
      'border-style': 'solid',
    },
  },
  {
    selector: 'edge',
    css: {
      'curve-style': 'bezier',
      opacity: 0.25,
      width: 3,
      'line-color': 'black',
    },
  },
];

export const ELEMENTS_STYLE_DARK: Stylesheet[] = [
  {
    selector: 'node',
    css: {
      shape: 'round-rectangle',
      width: 300,
      height: 100,
      color: 'white',

      'background-color': node => nodeColor(node),
      'background-opacity': 1,
      'overlay-opacity': 0,

      'background-image': renderImageDark,
      'background-width': '36px',
      'background-height': '36px',
      'background-position-x': node => imageOffset(node),
      'background-position-y': '8px',

      'border-color': node => nodeColor(node),
      'border-width': '2px',
      'border-style': 'solid',

      'padding-right': '10px',
      'padding-top': '10px',
      'padding-bottom': '10px',
      'z-index': 1,
    },
  },
  {
    selector: 'node:selected',
    css: {
      'curve-style': 'bezier',
      'border-width': 3,
      'border-color': '#49afd9',
      'border-style': 'solid',
    },
  },
  {
    selector: 'edge',
    css: {
      'curve-style': 'bezier',
      opacity: 0.25,
      width: 3,
      'line-color': 'white',
    },
  },
];

const imageOffset = node => `${node.width() - 24}px`;

const nodeColor = node => {
  switch (node.data('status')) {
    case 'error':
      return '#EB0E00';
    case 'warning':
      return '#EB7100';
    default:
      return '#62A420';
  }
};

// const polygonPoints = (header: number) => {
//   return `-1, ${header},   -1, -1,   1, -1,  1, 1,  -1, 1,  -1, ${header},   1, ${header}`;
// };
//
