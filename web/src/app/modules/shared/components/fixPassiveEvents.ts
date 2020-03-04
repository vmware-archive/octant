/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  * SPDX-License-Identifier: Apache-2.0
 *
 */

export default function fixPassiveEvents() {
  const relevantEvents = [
    'scroll',
    'wheel',
    'touchstart',
    'touchmove',
    'touchenter',
    'touchend',
    'touchleave',
    'mouseout',
    'mouseleave',
    'mouseup',
    'mousedown',
    'mousemove',
    'mouseenter',
    'mousewheel',
    'mouseover',
  ];

  const overwriteAddEvent = superMethod => {
    EventTarget.prototype.addEventListener = function(
      type,
      listener,
      options: any
    ) {
      if (relevantEvents.includes(type)) {
        let newCapture = false;
        if (options) {
          if (typeof options === 'object') {
            newCapture = options.capture || false;
          } else {
            newCapture = options || false;
          }
        }
        options = { passive: false, capture: newCapture };
      }
      superMethod.call(this, type, listener, options);
    };
  };

  overwriteAddEvent(EventTarget.prototype.addEventListener);
}
