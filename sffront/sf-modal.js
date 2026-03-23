/* SignForce — Global Modal System v1
   Replaces confirm() and alert() with glass modals */
(function() {
  // Inject modal CSS
  var style = document.createElement('style');
  style.textContent = `
    .sf-modal-overlay{position:fixed;inset:0;z-index:10000;background:rgba(0,0,0,.6);backdrop-filter:blur(8px);-webkit-backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;opacity:0;transition:opacity .2s ease;pointer-events:none}
    .sf-modal-overlay.sf-modal-show{opacity:1;pointer-events:all}
    .sf-modal-box{background:rgba(10,15,30,.92);backdrop-filter:blur(24px);-webkit-backdrop-filter:blur(24px);border:1px solid rgba(255,255,255,.08);border-radius:20px;padding:2rem;max-width:400px;width:90%;transform:scale(.95) translateY(10px);transition:transform .25s cubic-bezier(.34,1.56,.64,1);box-shadow:0 24px 80px rgba(0,0,0,.5),0 0 1px rgba(255,255,255,.1),inset 0 1px 0 rgba(255,255,255,.05)}
    .sf-modal-show .sf-modal-box{transform:scale(1) translateY(0)}
    .sf-modal-icon{width:48px;height:48px;border-radius:14px;display:flex;align-items:center;justify-content:center;margin-bottom:1.25rem}
    .sf-modal-icon.sf-modal-warn{background:rgba(251,191,36,.1)}
    .sf-modal-icon.sf-modal-danger{background:rgba(248,113,113,.1)}
    .sf-modal-icon.sf-modal-info{background:rgba(96,165,250,.1)}
    .sf-modal-icon.sf-modal-success{background:rgba(52,211,153,.1)}
    .sf-modal-title{font-size:1.05rem;font-weight:700;color:#fff;margin-bottom:.5rem}
    .sf-modal-msg{font-size:.88rem;color:rgba(255,255,255,.55);line-height:1.5;margin-bottom:1.5rem}
    .sf-modal-actions{display:flex;gap:.75rem;justify-content:flex-end}
    .sf-modal-btn{padding:.55rem 1.15rem;border-radius:12px;font-size:.85rem;font-weight:600;cursor:pointer;font-family:inherit;border:none;transition:all .2s}
    .sf-modal-btn:hover{transform:translateY(-1px)}
    .sf-modal-btn-cancel{background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.08)!important;color:rgba(255,255,255,.6)}
    .sf-modal-btn-cancel:hover{background:rgba(255,255,255,.1)}
    .sf-modal-btn-confirm{background:rgba(181,196,19,.15);border:1px solid rgba(181,196,19,.25)!important;color:#B5C413}
    .sf-modal-btn-confirm:hover{background:rgba(181,196,19,.25);box-shadow:0 0 20px rgba(181,196,19,.15)}
    .sf-modal-btn-danger{background:rgba(248,113,113,.15);border:1px solid rgba(248,113,113,.25)!important;color:#f87171}
    .sf-modal-btn-danger:hover{background:rgba(248,113,113,.25);box-shadow:0 0 20px rgba(248,113,113,.15)}
    .sf-modal-btn-ok{background:rgba(96,165,250,.15);border:1px solid rgba(96,165,250,.25)!important;color:#60a5fa}
    .sf-modal-btn-ok:hover{background:rgba(96,165,250,.25);box-shadow:0 0 20px rgba(96,165,250,.15)}
  `;
  document.head.appendChild(style);

  // Create overlay element (reused)
  var overlay = document.createElement('div');
  overlay.className = 'sf-modal-overlay';
  document.body.appendChild(overlay);

  function showModal(opts) {
    return new Promise(function(resolve) {
      var iconColors = {warn:'#fbbf24',danger:'#f87171',info:'#60a5fa',success:'#34d399'};
      var iconNames = {warn:'warning',danger:'error',info:'info',success:'check_circle'};
      var type = opts.type || 'warn';

      overlay.innerHTML = '<div class="sf-modal-box">' +
        '<div class="sf-modal-icon sf-modal-' + type + '">' +
          '<span class="material-symbols-outlined" style="font-size:24px;color:' + iconColors[type] + '">' + iconNames[type] + '</span>' +
        '</div>' +
        '<div class="sf-modal-title">' + (opts.title || '') + '</div>' +
        '<div class="sf-modal-msg">' + (opts.message || '') + '</div>' +
        '<div class="sf-modal-actions">' +
          (opts.showCancel !== false ? '<button class="sf-modal-btn sf-modal-btn-cancel" id="sf-modal-cancel">' + (opts.cancelText || 'Cancelar') + '</button>' : '') +
          '<button class="sf-modal-btn ' + (opts.confirmClass || 'sf-modal-btn-confirm') + '" id="sf-modal-ok">' + (opts.confirmText || 'Aceptar') + '</button>' +
        '</div>' +
      '</div>';

      requestAnimationFrame(function() {
        overlay.classList.add('sf-modal-show');
      });

      var cancelBtn = document.getElementById('sf-modal-cancel');
      var okBtn = document.getElementById('sf-modal-ok');

      function close(val) {
        overlay.classList.remove('sf-modal-show');
        setTimeout(function() { overlay.innerHTML = ''; }, 250);
        resolve(val);
      }

      if (cancelBtn) cancelBtn.onclick = function() { close(false); };
      okBtn.onclick = function() { close(true); };

      // ESC to cancel
      function onKey(e) {
        if (e.key === 'Escape') { close(false); document.removeEventListener('keydown', onKey); }
      }
      document.addEventListener('keydown', onKey);
    });
  }

  // Override native confirm
  window._nativeConfirm = window.confirm;
  window.confirm = function(msg) {
    // Can't make synchronous confirm async, so we return true
    // and handle via sfConfirm for new code
    return window._nativeConfirm(msg);
  };

  // Global async confirm — use this in new code
  window.sfConfirm = function(opts) {
    if (typeof opts === 'string') opts = {message: opts, title: 'Confirmar'};
    return showModal(opts);
  };

  // Global async alert
  window.sfAlert = function(opts) {
    if (typeof opts === 'string') opts = {message: opts, title: 'Aviso'};
    opts.showCancel = false;
    opts.confirmClass = opts.confirmClass || 'sf-modal-btn-ok';
    opts.confirmText = opts.confirmText || 'Entendido';
    return showModal(opts);
  };
})();
