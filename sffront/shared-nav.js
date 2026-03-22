/* SignForce — Shared Nav: highlight active link based on URL */
(function() {
  var pathMap = {
    '/users': 'dashboard', '/mydocs': 'mydocs', '/upload': 'upload',
    '/documents': 'mydocs', '/settings': 'settings', '/reports': 'reports',
    '/folders': 'folders', '/templates': 'templates', '/keys': 'keys',
    '/mykeys': 'keys', '/admin': 'dashboard', '/rootpanel': 'dashboard',
    '/approvals': 'dashboard'
  };
  var activeId = pathMap[window.location.pathname] || 'dashboard';
  
  // Remove all active classes, set the right one
  var links = document.querySelectorAll('.sf-nav-link');
  for (var i = 0; i < links.length; i++) {
    links[i].classList.remove('active');
    if (links[i].getAttribute('data-nav') === activeId) {
      links[i].classList.add('active');
    }
  }
})();
