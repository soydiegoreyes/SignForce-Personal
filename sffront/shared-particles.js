(function() {
  function spawnParticles(count) {
    count = count || 35;
    var container = document.createElement('div');
    container.className = 'sf-particles-container';

    var colors = [
      'rgba(181,196,19,.3)',   // verde signforce
      'rgba(181,196,19,.15)',  // verde suave
      'rgba(96,165,250,.2)',   // azul
      'rgba(52,211,153,.2)',   // esmeralda
      'rgba(255,255,255,.08)'  // blanco tenue
    ];

    for (var i = 0; i < count; i++) {
      var dot = document.createElement('div');
      dot.className = 'sf-particle';
      var size = (Math.random() * 4 + 1.5).toFixed(1);
      var left = (Math.random() * 100).toFixed(1);
      var duration = (Math.random() * 16 + 8).toFixed(1);
      var delay = (Math.random() * 15).toFixed(1);
      var drift = (Math.random() * 40 - 20).toFixed(0);
      var driftEnd = (Math.random() * 60 - 30).toFixed(0);
      var peakOpacity = (Math.random() * 0.4 + 0.15).toFixed(2);
      var color = colors[Math.floor(Math.random() * colors.length)];

      dot.style.cssText =
        '--size:' + size + 'px;' +
        '--duration:' + duration + 's;' +
        '--delay:' + delay + 's;' +
        '--drift:' + drift + 'px;' +
        '--drift-end:' + driftEnd + 'px;' +
        '--peak-opacity:' + peakOpacity + ';' +
        '--color:' + color + ';' +
        'left:' + left + '%;';
      container.appendChild(dot);
    }
    document.body.appendChild(container);
  }

  // Auto-init
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() { spawnParticles(); });
  } else {
    spawnParticles();
  }
})();
