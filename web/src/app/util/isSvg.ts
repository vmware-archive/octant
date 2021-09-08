export const isSvg = (svgString: string): boolean => {
  if (window && window.DOMParser) {
    const parser = new DOMParser();
    const doc = parser.parseFromString(svgString, 'image/svg+xml');

    return doc.getElementsByTagName('parsererror').length === 0;
  }

  // try rendering if window.DOMParser is not available
  return true;
};
