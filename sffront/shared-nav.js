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
    '/users': 'dashboard', '/mydocs': 'mydocs', '/upload': 'upload',
    '/documents': 'mydocs', '/settings': 'settings', '/reports': 'reports',
    '/folders': 'folders', '/templates': 'templates', '/keys': 'keys',
    '/mykeys': 'keys', '/admin': 'dashboard', '/rootpanel': 'dashboard',
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
        '<div class="sf-sidebar-title">Navegaci\u00f3n</div>' +
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
      '#sf-shared-sidebar{width:250px!important;min-width:250px!important;max-width:250px!important;flex-shrink:0!important;display:flex!important;flex-direction:column!important;gap:.75rem!important}' +
      '.sf-sidebar-card{background:rgba(15,21,36,.4)!important;backdrop-filter:blur(16px)!important;border:1px solid rgba(255,255,255,.06)!important;border-radius:16px!important;padding:1rem!important}' +
      '.sf-sidebar-title{font-size:.65rem!important;text-transform:uppercase!important;letter-spacing:.12em!important;color:rgba(255,255,255,.3)!important;padding:.4rem .6rem!important;font-weight:600!important;margin-bottom:.25rem!important}' +
      '.sf-nav-link{display:flex!important;align-items:center!important;gap:.75rem!important;padding:.6rem .75rem!important;border-radius:10px!important;font-size:.85rem!important;font-weight:500!important;color:rgba(255,255,255,.55)!important;border:1px solid transparent!important;transition:all .25s ease!important;text-decoration:none!important;margin-bottom:2px!important;background:transparent!important}' +
      '.sf-nav-link:hover{color:rgba(255,255,255,.85)!important;background:rgba(255,255,255,.04)!important}' +
      '.sf-nav-link.active{color:#B5C413!important;background:rgba(181,196,19,.08)!important;border-color:rgba(181,196,19,.12)!important}' +
      '.sf-card-icon{width:32px;height:32px;border-radius:8px;display:flex;align-items:center;justify-content:center;flex-shrink:0}' +
      '.sf-layout-normalized{display:flex!important;min-height:100vh!important;padding:80px 1.5rem 2rem!important;gap:1.25rem!important;position:relative!important;z-index:1!important;flex-direction:row!important;align-items:flex-start!important}' +
      '@media(max-width:900px){#sf-shared-sidebar{width:100%!important;min-width:100%!important;max-width:100%!important}}';
    document.head.appendChild(s);
  }

  window.injectSidebar = function(activeId) {
    activeId = activeId || detectActive();
    injectStyles();
    if (document.getElementById('sf-shared-sidebar')) return;

    // Find existing aside and replace
    var existing = document.querySelector('aside');
    if (existing) {
      // Normalize the parent layout container
      var parent = existing.parentElement;
      if (parent) {
        parent.classList.add('sf-layout-normalized');
        // Remove conflicting inline padding/classes
        var cn = parent.className;
        cn = cn.replace(/\b(pt-\d+|pb-\d+|px-\d+|min-h-screen)\b/g, '');
        parent.className = cn;
        parent.style.padding = '';
        parent.style.paddingTop = '';
        parent.style.gap = '';
      }
      existing.outerHTML = buildNavHTML(activeId);
    } else {
      var layout = document.querySelector('.sf-layout, [class*="layout"]');
      if (layout) {
        layout.classList.add('sf-layout-normalized');
        layout.insertAdjacentHTML('afterbegin', buildNavHTML(activeId));
      }
    }
  };
})();
