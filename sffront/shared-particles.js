/* SignForce — Floating Particles Background v2
   63 particles (80% more), varied sizes, glow effects, multi-layer */
(function() {
  if (document.querySelector('.sf-particles-container')) return;

  var container = document.createElement('div');
  container.className = 'sf-particles-container';

  var colors = [
    'rgba(181,196,19,.55)',   // verde signforce bright
    'rgba(181,196,19,.40)',   // verde signforce medium
    'rgba(181,196,19,.25)',   // verde signforce soft
    'rgba(96,165,250,.35)',   // azul
    'rgba(96,165,250,.20)',   // azul soft
    'rgba(52,211,153,.30)',   // esmeralda
    'rgba(255,255,255,.20)', // blanco
    'rgba(255,255,255,.12)', // blanco tenue
  ];

  var COUNT = 80;

  for (var i = 0; i < COUNT; i++) {
    var dot = document.createElement('div');
    dot.className = 'sf-particle';

    // More size variety: small dots (1-2px), medium (2-4px), large accents (4-7px)
    var sizeRoll = Math.random();
    var size;
    if (sizeRoll < 0.5) size = Math.random() * 1.5 + 1;       // 50% small: 1-2.5px
    else if (sizeRoll < 0.85) size = Math.random() * 2 + 2.5;  // 35% medium: 2.5-4.5px
    else size = Math.random() * 3 + 4;                          // 15% large: 4-7px

    var left = (Math.random() * 100).toFixed(1);
    var duration = (Math.random() * 20 + 10).toFixed(1);  // 10-30s (slower, more relaxed)
    var delay = (Math.random() * 20).toFixed(1);           // stagger up to 20s
    var drift = (Math.random() * 60 - 30).toFixed(0);      // wider horizontal drift
    var driftEnd = (Math.random() * 80 - 40).toFixed(0);
    var peakOpacity = size > 4 ? (Math.random() * 0.4 + 0.2).toFixed(2)   // large = softer
                    : size > 2.5 ? (Math.random() * 0.5 + 0.3).toFixed(2) // medium
                    : (Math.random() * 0.6 + 0.3).toFixed(2);             // small = brighter
    var color = colors[Math.floor(Math.random() * colors.length)];

    // Some particles pulse/twinkle
    var twinkle = Math.random() < 0.2;

    dot.style.cssText =
      '--size:' + size.toFixed(1) + 'px;' +
      '--duration:' + duration + 's;' +
      '--delay:' + delay + 's;' +
      '--drift:' + drift + 'px;' +
      '--drift-end:' + driftEnd + 'px;' +
      '--peak-opacity:' + peakOpacity + ';' +
      '--color:' + color + ';' +
      'left:' + left + '%;';

    if (twinkle) dot.classList.add('sf-particle-twinkle');

    container.appendChild(dot);
  }

  document.body.appendChild(container);
})();
