import _ from 'lodash';
import { isSvg } from './isSvg';

describe('isSvg', () => {
  const validSvgs = [
    `<svg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />Sorry, your browser does not support inline SVG.</svg>`,
    `<svg width="400" height="180"><rect x="50" y="20" width="150" height="150" style="fill:blue;stroke:pink;stroke-width:5;fill-opacity:0.1;stroke-opacity:0.9" />Sorry, your browser does not support inline SVG.</svg>`,
    `<svg height="140" width="500"><ellipse cx="200" cy="80" rx="100" ry="50" style="fill:yellow;stroke:purple;stroke-width:2" />Sorry, your browser does not support inline SVG.</svg>`,
    `<svg height="210" width="500"><polygon points="100,10 40,198 190,78 10,78 160,198" style="fill:lime;stroke:purple;stroke-width:5;fill-rule:evenodd;"/>Sorry, your browser does not support inline SVG.</svg>`,
    `<svg height="90" width="200"><text x="10" y="20" style="fill:red;">Several lines:<tspan x="10" y="45">First line.</tspan><tspan x="10" y="70">Second line.</tspan></text>Sorry, your browser does not support inline SVG.</svg>`,
    `<svg height="60" width="200"><text x="0" y="15" fill="red" transform="rotate(30 20,40)">I love SVG</text>Sorry, your browser does not support inline SVG.</svg>`,
    `<svg width="400" height="200">
      <defs>
        <filter id="MyFilter" filterUnits="userSpaceOnUse" x="0" y="0" width="200" height="120">
          <feGaussianBlur in="SourceAlpha" stdDeviation="4" result="blur" />
          <feOffset in="blur" dx="4" dy="4" result="offsetBlur"/>
          <feSpecularLighting in="blur" surfaceScale="5" specularConstant=".75" specularExponent="20" lighting-color="#bbbbbb" result="specOut">
            <fePointLight x="-5000" y="-10000" z="20000" />
          </feSpecularLighting>
          <feComposite in="specOut" in2="SourceAlpha" operator="in" result="specOut" />
        </filter>
      </defs>
      <rect x="1" y="1" width="198" height="118" fill="#cccccc" />
      <g filter="url(#MyFilter)">
        <path fill="none" stroke="#D90000" stroke-width="10" d="M50,90 C0,90 0,30 50,30 L150,30 C200,30 200,90 150,90 z" />
        <text fill="#FFFFFF" stroke="black" font-size="45" font-family="Verdana" x="52" y="76">SVG</text>
      </g>
      Sorry, your browser does not support inline SVG.
    </svg>`,
  ];

  const invalidSvgs = [
    `<svg></sv>`,
    `<svg><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red"></svg>`,
    `<svg><polygon points="100,10 40,198 190,78 10,78 160,198"style="fill:lime;stroke:purple;stroke-width:5;fill-rule:evenodd;"/></svg>`,
  ];

  it('should return false when svg string is invalid', () => {
    invalidSvgs.forEach(svg => {
      expect(isSvg(svg)).toBe(false);
    });
  });

  it('should return true when svg string is valid', () => {
    validSvgs.forEach(svg => {
      expect(isSvg(svg)).toBe(true);
    });
  });
});
