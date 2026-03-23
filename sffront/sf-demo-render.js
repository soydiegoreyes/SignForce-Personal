/* SignForce — Demo Renderer v1
   Renders demo data into page elements after DOM loads */
(function() {
  if (!window.SF_DEMO) return;

  function formatDate(d) {
    var dt = new Date(d);
    return dt.toLocaleDateString('es-MX', {day:'2-digit',month:'short',year:'numeric'});
  }

  // Wait for DOM
  document.addEventListener('DOMContentLoaded', function() {
    var path = window.location.pathname;

    // ═══ DASHBOARD ═══
    if (path === '/users' || path === '/dashboard') {
      renderDashboard();
    }
    // ═══ DOCUMENTS ═══ (disabled — uses real API)
    // if (path === '/mydocs' || path === '/documents') {
    //   setTimeout(renderDocuments, 1500);
    // }
    // ═══ FOLDERS ═══
    if (path === '/folders') {
      setTimeout(renderFolders, 800);
    }
    // ═══ TEMPLATES ═══
    if (path === '/templates') {
      setTimeout(renderTemplates, 800);
    }
    // ═══ KEYS ═══
    if (path === '/keys') {
      setTimeout(renderKeys, 800);
    }
    // ═══ REPORTS ═══
    if (path === '/reports') {
      setTimeout(renderReports, 800);
    }
  });

  // ═══ Dashboard Renderer ═══
  function renderDashboard() {
    var s = window.SF_DEMO_STATS;
    // Fill stat cards if they exist
    var statEls = document.querySelectorAll('.stat-number, .stat-value, [data-stat]');
    if (statEls.length === 0) return;

    // Try to fill known IDs
    var fills = {
      'totalDocs': s.totalDocs, 'signedDocs': s.signed, 'pendingDocs': s.pending,
      'rejectedDocs': s.rejected, 'monthDocs': s.thisMonth, 'totalUsers': s.users
    };
    Object.keys(fills).forEach(function(id) {
      var el = document.getElementById(id);
      if (el) el.textContent = fills[id];
    });

    // Render activity feed
    var actContainer = document.getElementById('activityFeed') || document.querySelector('.activity-list, .recent-activity');
    if (actContainer) {
      actContainer.innerHTML = window.SF_DEMO_ACTIVITY.map(function(a) {
        return '<div style="display:flex;align-items:center;gap:.75rem;padding:.75rem;border-radius:12px;background:rgba(255,255,255,.02);border:1px solid rgba(255,255,255,.04);margin-bottom:.5rem;">' +
          '<div style="width:36px;height:36px;border-radius:10px;background:' + a.color + '15;display:flex;align-items:center;justify-content:center;flex-shrink:0;">' +
            '<span class="material-symbols-outlined" style="font-size:18px;color:' + a.color + ';">' + a.icon + '</span></div>' +
          '<div style="flex:1;min-width:0;">' +
            '<div style="font-size:.85rem;font-weight:600;color:rgba(255,255,255,.85);">' + a.action + '</div>' +
            '<div style="font-size:.75rem;color:rgba(255,255,255,.4);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">' + a.doc + ' · ' + a.user + '</div>' +
          '</div>' +
          '<span style="font-size:.7rem;color:rgba(255,255,255,.25);white-space:nowrap;">' + a.time + '</span>' +
        '</div>';
      }).join('');
    }
  }

  // ═══ Documents Renderer ═══
  function renderDocuments() {
    var tbody = document.querySelector('#docsTable tbody, .docs-table tbody, table tbody');
    if (!tbody) return;
    // Only render if empty or has loading
    if (tbody.children.length > 1 && !tbody.querySelector('.loading, .spinner')) return;

    var docs = Object.values(window.SF_DEMO_DOCS);
    var sm = window.SF_STATUS_MAP;

    tbody.innerHTML = docs.map(function(d) {
      var st = sm[d.status] || sm[1];
      return '<tr style="border-bottom:1px solid rgba(255,255,255,.04);">' +
        '<td style="padding:.75rem;"><div style="display:flex;align-items:center;gap:.6rem;">' +
          '<span class="material-symbols-outlined" style="font-size:20px;color:rgba(255,255,255,.3);">description</span>' +
          '<div><div style="font-size:.85rem;font-weight:600;color:rgba(255,255,255,.85);">' + d.nameDoc + '</div>' +
          '<div style="font-size:.72rem;color:rgba(255,255,255,.35);">' + d.pages + ' págs · ' + d.fileSize + '</div></div></div></td>' +
        '<td style="padding:.75rem;"><span style="display:inline-flex;align-items:center;gap:4px;padding:.25rem .6rem;border-radius:8px;background:' + st.bg + ';border:1px solid ' + st.border + ';font-size:.75rem;font-weight:600;color:' + st.color + ';">' +
          '<span style="width:6px;height:6px;border-radius:50%;background:' + st.color + ';"></span>' + st.label + '</span></td>' +
        '<td style="padding:.75rem;font-size:.82rem;color:rgba(255,255,255,.5);">' + d.signers.join(', ') + '</td>' +
        '<td style="padding:.75rem;font-size:.82rem;color:rgba(255,255,255,.4);">' + formatDate(d.createDate) + '</td>' +
        '<td style="padding:.75rem;text-align:center;">' +
          '<button style="background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.06);border-radius:8px;padding:.3rem .6rem;cursor:pointer;color:rgba(255,255,255,.5);font-size:.75rem;font-family:inherit;">Ver</button></td>' +
      '</tr>';
    }).join('');

    // Update counter
    var counter = document.querySelector('.doc-count, #docCount, [data-doc-count]');
    if (counter) counter.textContent = docs.length + ' documentos';

    // Hide loading
    var loading = document.querySelector('.loading-container, .spinner-container, #loadingDocs');
    if (loading) loading.style.display = 'none';
  }

  // ═══ Folders Renderer ═══
  function renderFolders() {
    var grid = document.querySelector('.folders-grid, .folder-list, #foldersContainer, .sf-page-content .grid');
    if (!grid || (grid.children.length > 1 && !grid.querySelector('.loading'))) {
      // Try creating in page content
      grid = document.querySelector('.sf-page-content');
      if (!grid) return;
      var g = document.createElement('div');
      g.style.cssText = 'display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:1rem;';
      g.id = 'demoFolderGrid';
      
      // Remove loading if present
      var loading = grid.querySelector('.loading, .spinner');
      if (loading) loading.remove();
      
      // Check if grid already exists
      if (document.getElementById('demoFolderGrid')) return;
      grid.appendChild(g);
      grid = g;
    }

    grid.innerHTML = window.SF_DEMO_FOLDERS.map(function(f) {
      return '<div class="sf-card folder-card" style="padding:1.25rem;cursor:pointer;transition:all .2s;" onmouseenter="this.style.borderColor=\'rgba(255,255,255,.12)\'" onmouseleave="this.style.borderColor=\'\'">' +
        '<div style="display:flex;align-items:center;gap:.75rem;margin-bottom:1rem;">' +
          '<div style="width:42px;height:42px;border-radius:12px;background:' + f.color + '15;display:flex;align-items:center;justify-content:center;">' +
            '<span class="material-symbols-outlined" style="font-size:22px;color:' + f.color + ';">' + f.icon + '</span></div>' +
          '<div><div style="font-size:.95rem;font-weight:700;color:rgba(255,255,255,.9);">' + f.name + '</div>' +
          '<div style="font-size:.75rem;color:rgba(255,255,255,.35);">Creada ' + formatDate(f.created) + '</div></div></div>' +
        '<div style="display:flex;align-items:center;justify-content:space-between;">' +
          '<span style="font-size:.8rem;color:rgba(255,255,255,.4);">' + f.docs + ' documentos</span>' +
          '<span class="material-symbols-outlined" style="font-size:16px;color:rgba(255,255,255,.2);">chevron_right</span></div>' +
      '</div>';
    }).join('');
  }

  // ═══ Templates Renderer ═══
  function renderTemplates() {
    var container = document.querySelector('.templates-grid, .template-list, #templatesContainer, .sf-page-content');
    if (!container) return;
    
    if (!document.getElementById('demoTemplateGrid')) {
      var loading = container.querySelector('.loading, .spinner');
      if (loading) loading.remove();
      
      var g = document.createElement('div');
      g.id = 'demoTemplateGrid';
      g.style.cssText = 'display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:1rem;';
      container.appendChild(g);
      container = g;
    } else {
      container = document.getElementById('demoTemplateGrid');
    }

    var cats = {Legal:'#a78bfa',Comercial:'#60a5fa',Inmobiliario:'#B5C413',General:'#34d399',Corporativo:'#fbbf24',RH:'#f87171',Financiero:'#fb923c'};

    container.innerHTML = window.SF_DEMO_TEMPLATES.map(function(t) {
      var c = cats[t.category] || '#60a5fa';
      return '<div class="sf-card" style="padding:1.25rem;cursor:pointer;transition:all .2s;" onmouseenter="this.style.transform=\'translateY(-2px)\'" onmouseleave="this.style.transform=\'\'">' +
        '<div style="display:flex;align-items:start;justify-content:space-between;margin-bottom:1rem;">' +
          '<div style="width:42px;height:42px;border-radius:12px;background:rgba(181,196,19,.08);display:flex;align-items:center;justify-content:center;">' +
            '<span class="material-symbols-outlined" style="font-size:22px;color:#B5C413;">draft</span></div>' +
          '<span style="font-size:.68rem;font-weight:600;padding:.2rem .5rem;border-radius:6px;background:' + c + '15;color:' + c + ';border:1px solid ' + c + '25;">' + t.category + '</span></div>' +
        '<div style="font-size:.95rem;font-weight:700;color:rgba(255,255,255,.9);margin-bottom:.4rem;">' + t.name + '</div>' +
        '<div style="font-size:.75rem;color:rgba(255,255,255,.35);margin-bottom:1rem;">' + t.fields + ' campos · Actualizada ' + formatDate(t.updated) + '</div>' +
        '<div style="display:flex;align-items:center;justify-content:space-between;">' +
          '<span style="font-size:.78rem;color:rgba(255,255,255,.4);display:flex;align-items:center;gap:4px;"><span class="material-symbols-outlined" style="font-size:14px;">bar_chart</span> ' + t.uses + ' usos</span>' +
          '<button style="background:rgba(181,196,19,.1);border:1px solid rgba(181,196,19,.2);border-radius:8px;padding:.3rem .7rem;color:#B5C413;font-size:.75rem;font-weight:600;cursor:pointer;font-family:inherit;">Usar</button></div>' +
      '</div>';
    }).join('');
  }

  // ═══ Keys Renderer ═══
  function renderKeys() {
    var container = document.querySelector('.keys-list, #keysContainer, .sf-page-content');
    if (!container) return;

    if (!document.getElementById('demoKeysGrid')) {
      var loading = container.querySelector('.loading, .spinner');
      if (loading) loading.remove();

      var g = document.createElement('div');
      g.id = 'demoKeysGrid';
      g.style.cssText = 'display:flex;flex-direction:column;gap:.75rem;';
      container.appendChild(g);
      container = g;
    } else {
      container = document.getElementById('demoKeysGrid');
    }

    container.innerHTML = window.SF_DEMO_KEYS.map(function(k) {
      var isActive = k.status === 'active';
      var stColor = isActive ? '#34d399' : '#f87171';
      var stLabel = isActive ? 'Activa' : 'Expirada';
      return '<div class="sf-card key-card" style="padding:1.25rem;display:flex;align-items:center;gap:1rem;">' +
        '<div style="width:46px;height:46px;border-radius:12px;background:rgba(181,196,19,.08);display:flex;align-items:center;justify-content:center;flex-shrink:0;">' +
          '<span class="material-symbols-outlined" style="font-size:24px;color:#B5C413;">key</span></div>' +
        '<div style="flex:1;min-width:0;">' +
          '<div style="display:flex;align-items:center;gap:.5rem;margin-bottom:.25rem;">' +
            '<span style="font-size:.95rem;font-weight:700;color:rgba(255,255,255,.9);">' + k.name + '</span>' +
            '<span style="font-size:.68rem;font-weight:600;padding:.15rem .45rem;border-radius:6px;background:' + stColor + '15;color:' + stColor + ';border:1px solid ' + stColor + '25;">' + stLabel + '</span></div>' +
          '<div style="font-size:.78rem;color:rgba(255,255,255,.4);">' + k.type + ' · Expira ' + formatDate(k.expires) + '</div>' +
          '<div style="font-size:.72rem;color:rgba(255,255,255,.25);margin-top:.2rem;">Último uso: ' + formatDate(k.lastUsed) + '</div></div>' +
        '<button style="background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.06);border-radius:8px;padding:.4rem .8rem;color:rgba(255,255,255,.5);font-size:.78rem;cursor:pointer;font-family:inherit;white-space:nowrap;">Administrar</button>' +
      '</div>';
    }).join('');
  }

  // ═══ Reports Renderer ═══
  function renderReports() {
    var container = document.querySelector('.sf-page-content');
    if (!container || document.getElementById('demoReports')) return;

    var r = window.SF_DEMO_REPORTS;
    var div = document.createElement('div');
    div.id = 'demoReports';

    // Monthly chart (bar chart with CSS)
    var maxDocs = Math.max.apply(null, r.monthly.map(function(m){return m.docs;}));
    var barChart = r.monthly.map(function(m) {
      var h = Math.round((m.docs / maxDocs) * 120);
      var h2 = Math.round((m.signed / maxDocs) * 120);
      return '<div style="display:flex;flex-direction:column;align-items:center;gap:.4rem;">' +
        '<div style="display:flex;align-items:flex-end;gap:3px;height:120px;">' +
          '<div style="width:18px;background:rgba(181,196,19,.2);border:1px solid rgba(181,196,19,.3);border-radius:4px 4px 0 0;height:' + h + 'px;transition:height .5s;"></div>' +
          '<div style="width:18px;background:rgba(52,211,153,.2);border:1px solid rgba(52,211,153,.3);border-radius:4px 4px 0 0;height:' + h2 + 'px;transition:height .5s;"></div></div>' +
        '<span style="font-size:.7rem;color:rgba(255,255,255,.4);">' + m.month + '</span></div>';
    }).join('');

    // Top signers
    var signers = r.topSigners.map(function(s,i) {
      var w = Math.round((s.count / r.topSigners[0].count) * 100);
      return '<div style="margin-bottom:.6rem;">' +
        '<div style="display:flex;justify-content:space-between;margin-bottom:.25rem;">' +
          '<span style="font-size:.82rem;color:rgba(255,255,255,.7);">' + (i+1) + '. ' + s.name + '</span>' +
          '<span style="font-size:.78rem;color:rgba(255,255,255,.4);">' + s.count + ' firmas</span></div>' +
        '<div style="height:6px;border-radius:3px;background:rgba(255,255,255,.04);"><div style="height:100%;border-radius:3px;background:linear-gradient(90deg,#B5C413,#34d399);width:' + w + '%;transition:width .5s;"></div></div></div>';
    }).join('');

    // By type
    var types = r.byType.map(function(t) {
      var colors = {Contratos:'#B5C413',NDAs:'#60a5fa',Poderes:'#a78bfa',Actas:'#fbbf24',Convenios:'#34d399',Otros:'#f87171'};
      var c = colors[t.type] || '#60a5fa';
      return '<div style="display:flex;align-items:center;gap:.75rem;padding:.5rem 0;">' +
        '<div style="width:10px;height:10px;border-radius:3px;background:' + c + ';flex-shrink:0;"></div>' +
        '<span style="font-size:.82rem;color:rgba(255,255,255,.7);flex:1;">' + t.type + '</span>' +
        '<span style="font-size:.82rem;color:rgba(255,255,255,.4);">' + t.count + '</span>' +
        '<span style="font-size:.72rem;color:rgba(255,255,255,.25);width:35px;text-align:right;">' + t.pct + '%</span></div>';
    }).join('');

    div.innerHTML =
      '<div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem;margin-bottom:1rem;">' +
        '<div class="sf-card" style="padding:1.25rem;">' +
          '<div style="font-size:.85rem;font-weight:700;color:rgba(255,255,255,.8);margin-bottom:1rem;display:flex;align-items:center;gap:.5rem;"><span class="material-symbols-outlined" style="font-size:18px;color:#B5C413;">bar_chart</span> Documentos por mes</div>' +
          '<div style="display:flex;align-items:flex-end;justify-content:space-around;">' + barChart + '</div>' +
          '<div style="display:flex;gap:1rem;justify-content:center;margin-top:.75rem;">' +
            '<span style="font-size:.7rem;color:rgba(255,255,255,.4);display:flex;align-items:center;gap:4px;"><span style="width:8px;height:8px;border-radius:2px;background:rgba(181,196,19,.3);"></span> Creados</span>' +
            '<span style="font-size:.7rem;color:rgba(255,255,255,.4);display:flex;align-items:center;gap:4px;"><span style="width:8px;height:8px;border-radius:2px;background:rgba(52,211,153,.3);"></span> Firmados</span></div>' +
        '</div>' +
        '<div class="sf-card" style="padding:1.25rem;">' +
          '<div style="font-size:.85rem;font-weight:700;color:rgba(255,255,255,.8);margin-bottom:1rem;display:flex;align-items:center;gap:.5rem;"><span class="material-symbols-outlined" style="font-size:18px;color:#60a5fa;">leaderboard</span> Top firmantes</div>' +
          signers +
        '</div></div>' +
      '<div class="sf-card" style="padding:1.25rem;">' +
        '<div style="font-size:.85rem;font-weight:700;color:rgba(255,255,255,.8);margin-bottom:.75rem;display:flex;align-items:center;gap:.5rem;"><span class="material-symbols-outlined" style="font-size:18px;color:#a78bfa;">donut_large</span> Por tipo de documento</div>' +
        types +
      '</div>';

    container.appendChild(div);
  }

})();
