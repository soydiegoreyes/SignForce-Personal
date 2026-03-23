/* SignForce — Floating Particles Background v3
   80 particles, pre-spread across viewport, varied sizes */
(function() {
  if (document.querySelector('.sf-particles-container')) return;

  var container = document.createElement('div');
  container.className = 'sf-particles-container';

  var colors = [
    'rgba(181,196,19,.55)',
    'rgba(181,196,19,.40)',
    'rgba(181,196,19,.25)',
    'rgba(96,165,250,.35)',
    'rgba(96,165,250,.20)',
    'rgba(52,211,153,.30)',
    'rgba(255,255,255,.20)',
    'rgba(255,255,255,.12)',
  ];

  var COUNT = 80;

  for (var i = 0; i < COUNT; i++) {
    var dot = document.createElement('div');
    dot.className = 'sf-particle';

    // Size variety: small (1-2.5px), medium (3-5px), large (5-9px)
    var sizeRoll = Math.random();
    var size;
    if (sizeRoll < 0.45) size = Math.random() * 1.5 + 1;       // 45% small
    else if (sizeRoll < 0.78) size = Math.random() * 2 + 3;     // 33% medium
    else size = Math.random() * 4 + 5;                           // 22% large

    var left = (Math.random() * 100).toFixed(1);
    var duration = (Math.random() * 22 + 12).toFixed(1);
    var drift = (Math.random() * 60 - 30).toFixed(0);
    var driftEnd = (Math.random() * 80 - 40).toFixed(0);
    var peakOpacity = size > 5 ? (Math.random() * 0.35 + 0.15).toFixed(2)
                    : size > 3 ? (Math.random() * 0.5 + 0.25).toFixed(2)
                    : (Math.random() * 0.6 + 0.3).toFixed(2);
    var color = colors[Math.floor(Math.random() * colors.length)];
    var twinkle = Math.random() < 0.2;

    // NEGATIVE delay = particle starts mid-animation (already visible on load)
    var negDelay = -(Math.random() * parseFloat(duration)).toFixed(1);

    dot.style.cssText =
      '--size:' + size.toFixed(1) + 'px;' +
      '--duration:' + duration + 's;' +
      '--delay:' + negDelay + 's;' +
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
