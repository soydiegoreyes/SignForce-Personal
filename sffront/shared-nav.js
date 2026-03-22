(function() {
  var navItems = [
    { id: 'dashboard',  label: 'Dashboard',        icon: 'dashboard',    href: '/users' },
    { id: 'mydocs',     label: 'Mis Documentos',   icon: 'description',  href: '/mydocs' },
    { id: 'settings',   label: 'Configuración',    icon: 'settings',     href: '/settings' },
    { id: 'reports',    label: 'Reportes',          icon: 'assessment',   href: '/reports' },
    { id: 'upload',     label: 'Subir Documento',   icon: 'upload_file',  href: '/upload' },
    { id: 'folders',    label: 'Mis Folders',       icon: 'folder',       href: '/folders' },
    { id: 'templates',  label: 'Plantillas',        icon: 'article',      href: '/templates' },
    { id: 'keys',       label: 'Mis Llaves',        icon: 'key',          href: '/keys' }
  ];

  var pathMap = {
    '/users': 'dashboard',
    '/mydocs': 'mydocs',
    '/upload': 'mydocs',
    '/documents': 'mydocs',
    '/settings': 'settings',
    '/reports': 'reports',
    '/folders': 'folders',
    '/templates': 'templates',
    '/keys': 'keys',
    '/mykeys': 'keys',
    '/admin': 'dashboard',
    '/rootpanel': 'dashboard',
    '/approvals': 'dashboard'
  };

  function detectActive() {
    return pathMap[window.location.pathname] || 'dashboard';
  }

  function buildNavHTML(activeId) {
    var links = '';
    for (var i = 0; i < navItems.length; i++) {
      var item = navItems[i];
      var cls = item.id === activeId ? 'sf-nav-link active' : 'sf-nav-link';
      links += '<a class="' + cls + '" href="' + item.href + '"><span class="material-symbols-outlined">' + item.icon + '</span> ' + item.label + '</a>';
    }

    return '<aside class="sf-sidebar" id="sf-shared-sidebar">' +
      '<div class="sf-sidebar-card" style="padding:1.1rem 1.25rem;">' +
        '<div style="display:flex;align-items:center;gap:.6rem;">' +
          '<div class="sf-card-icon" style="background:rgba(96,165,250,.1);"><span class="material-symbols-outlined" style="font-size:1.1rem;color:#60a5fa;">person</span></div>' +
          '<div><div style="font-size:.95rem;font-weight:700;">Mi Panel</div><div style="font-size:.68rem;color:rgba(255,255,255,.35);">SignForce</div></div>' +
        '</div>' +
      '</div>' +
      '<nav class="sf-sidebar-card" style="flex:1;">' +
        '<div class="sf-sidebar-title">Navegación</div>' +
        links +
        '<div style="flex:1;"></div>' +
        '<div style="border-top:1px solid rgba(255,255,255,.06);padding-top:.75rem;margin-top:.75rem;">' +
          '<div style="display:flex;align-items:center;gap:5px;"><span style="width:6px;height:6px;border-radius:50%;background:#34d399;box-shadow:0 0 8px rgba(52,211,153,.4);"></span><span style="font-size:.7rem;color:rgba(255,255,255,.3);">v2.1.0</span></div>' +
        '</div>' +
      '</nav>' +
    '</aside>';
  }

  function injectStyles() {
    if (document.getElementById('sf-nav-styles')) return;
    var s = document.createElement('style');
    s.id = 'sf-nav-styles';
    s.textContent =
      '.sf-sidebar{width:250px;flex-shrink:0;display:flex;flex-direction:column;gap:.75rem}' +
      '.sf-sidebar-card{background:rgba(15,21,36,.4);backdrop-filter:blur(16px);border:1px solid rgba(255,255,255,.06);border-radius:16px;padding:1rem}' +
      '.sf-sidebar-title{font-size:.65rem;text-transform:uppercase;letter-spacing:.12em;color:rgba(255,255,255,.3);padding:.4rem .6rem;font-weight:600;margin-bottom:.25rem}' +
      '.sf-nav-link{display:flex;align-items:center;gap:.75rem;padding:.6rem .75rem;border-radius:10px;font-size:.85rem;font-weight:500;color:rgba(255,255,255,.55);border:1px solid transparent;transition:all .25s ease;text-decoration:none;margin-bottom:2px}' +
      '.sf-nav-link:hover{color:rgba(255,255,255,.85);background:rgba(255,255,255,.04)}' +
      '.sf-nav-link.active{color:#B5C413;background:rgba(181,196,19,.08);border-color:rgba(181,196,19,.12)}' +
      '.sf-card-icon{width:32px;height:32px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0}' +
      '@media(max-width:900px){.sf-sidebar{width:100%}}';
    document.head.appendChild(s);
  }

  window.injectSidebar = function(activeId) {
    activeId = activeId || detectActive();
    injectStyles();

    // If already injected, skip
    if (document.getElementById('sf-shared-sidebar')) return;

    // Find any existing aside/sidebar and replace
    var existing = document.querySelector('aside.sf-sidebar, aside');
    if (existing) {
      existing.outerHTML = buildNavHTML(activeId);
      return;
    }

    // Fallback: find layout container and prepend
    var layout = document.querySelector('.sf-layout, .bg.flex, [class*="layout"]');
    if (layout) {
      layout.insertAdjacentHTML('afterbegin', buildNavHTML(activeId));
    }
  };
})();
