// ====== Estado y utilidades ======
    function $(sel, root = document) { return root.querySelector(sel); }
    function $all(sel, root = document) { return Array.from(root.querySelectorAll(sel)); }

    function getTodayISO() {
      const d = new Date();
      return d.toISOString().slice(0, 10);
    }

    // ====== Petición de endpoint /statusk ======
    async function fetchKeysStatus() {
      // Simula latencia y respuesta desde /statusk
      // Enviar al endpoint de Go
      const response = await fetch('/statusk', {
        method: 'GET',
        credentials: 'include'
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Error: ${response.status} - ${errorText}`);
      }

      const result = await response.json();

      // console.log("Resultado:", result);
      return result;
    }

    // ====== Render de llaves ======
    function renderKeys(keys) {
      const list = document.getElementById('keysList');
      list.innerHTML = '';

      if (keys.length === 0) {
        list.innerHTML = '<div class="glass-card p-block text-center">No tienes llaves registradas todavía.</div>';
        return;
      }

      keys.forEach(key => {
        const keyCard = document.createElement('div');
        keyCard.className = 'key-card glass-panel';
        keyCard.setAttribute('data-idkey', key.idKey);

        const expirationDate = new Date(key.expiration);
        const isExpired = expirationDate < new Date();
        const badgeClass = isExpired
          ? 'border:1px solid rgba(239,68,68,.4); color:#fecaca; background:rgba(239,68,68,.15)'
          : 'border:1px solid rgba(35,197,94,.4); color:#d1fae5; background:rgba(35,197,94,.15)';

        const header = document.createElement('div');
        header.className = 'key-card-header';
        header.id = `key-${key.idKey}`;

        const title = document.createElement('h3');
        title.className = 'key-title';
        title.textContent = key.owner;

        const expBadge = document.createElement('span');
        expBadge.className = 'badge';
        expBadge.style = badgeClass;
        expBadge.textContent = isExpired ? 'Expirada' : 'Vigente';

        header.appendChild(title);
        header.appendChild(expBadge);

        keyCard.appendChild(header);

        const details = document.createElement('div');
        details.className = 'key-details';
        details.innerHTML = `
          <p class="key-detail"><strong>Certificado:</strong> ${key.nameCer}</p>
          <p class="key-detail"><strong>Llave:</strong> ${key.nameKey}</p>
          <p class="key-detail"><strong>ID único:</strong> ${key.subjectUniqueId}</p>
          <p class="key-detail"><strong>Expiración:</strong> ${key.expiration}</p>
          <p class="key-detail"><strong>Subido el:</strong> ${key.uploadedAt}</p>
        `;
        keyCard.appendChild(details);
        const issuerdetails = document.createElement('div');
        issuerdetails.className = 'key-details';
        issuerdetails.innerHTML = formatIssuer(key.issuerData);
        keyCard.appendChild(issuerdetails);

        list.appendChild(keyCard);
      });

      // después de construir todo, marcamos el seleccionado
      const selectedKey = keys.find(k => k.selected);
      if (selectedKey) {
        marcarSeleccionada(selectedKey.idKey);
      } else {
        marcarSeleccionada(null); // ningún seleccionado
      }
    }
    function formatIssuer(issuerRFC4514) {
      if (!issuerRFC4514) {
        return '';
      }

      // 1. Dividir el string por las comas que separan los atributos
      const attributes = issuerRFC4514.split(',');
      
      // 2. Procesar cada atributo
      const budgetsHTML = attributes.map(attribute => {
        // Expresión Regular para encontrar y eliminar el prefijo
        // El patrón es: cualquier cosa que no sea "=" (.[^=]*) seguida de "="
        // y se usa replace para dejar solo el valor (la parte derecha del "=").
        // También se limpia el espacio al inicio del valor si existe (como en ", O=...")
        const value = attribute.replace(/^.[^=]*=/, '').trim();

        // Opcionalmente, puedes obtener el nombre del atributo (el prefijo) para mostrarlo
        const match = attribute.match(/^.[^=]*=/);
        let name = '';
        if (match) {
          // Limpia el '=' final y, en el caso de OIDs, solo deja OID
          name = match[0].replace('=', '').replace(/OID\..*/, 'OID').trim();
        }

        // 3. Crear el HTML para el "budget" (badge)
        // Usamos una clase de Bootstrap para darle estilo
        // Se puede personalizar la clase según tu necesidad (bg-primary, bg-success, etc.)
        // Concatenamos el nombre del atributo y su valor para mayor contexto
        return `<span class="badge bg-secondary me-2 mb-1" title="${name}: ${value}">${value}</span>`;
      }).join(''); // Unir todos los elementos en un solo string

      // 4. Devolver el div contenedor con todos los budgets
      // El 'd-flex flex-wrap' es útil para que se acomoden bien si hay muchos.
      return `<div class="d-flex flex-wrap"><p>Emisor de certificado</p>${budgetsHTML}</div>`;
    }
    function marcarSeleccionada(idSeleccionada) {
      document.querySelectorAll('.key-card').forEach(card => {
        const cardId = card.getAttribute('data-idkey');
        const header = card.querySelector('.key-card-header');

        // quitar botón/badge previos
        const old = card.querySelector('.selected-badge, .select-btn');
        if (old) old.remove();

        if (idSeleccionada && cardId === String(idSeleccionada)) {
          // badge verde
          const selBadge = document.createElement('span');
          selBadge.className = 'badge selected-badge';
          selBadge.style = 'border:1px solid rgba(35,197,94,.4); color:#d1fae5; background:rgba(35,197,94,.15)';
          selBadge.textContent = 'Seleccionada';
          header.appendChild(selBadge);
        } else {
          // botón seleccionar
          const newBtn = document.createElement('button');
          newBtn.type = 'button';
          newBtn.className = 'badge select-btn';
          newBtn.style.cursor = 'pointer';
          newBtn.textContent = 'Seleccionar';
          header.appendChild(newBtn);

          newBtn.addEventListener('click', () => seleccionarKey(cardId, newBtn));
        }
      });
    }

    async function seleccionarKey(idKey, btn) {
      btn.disabled = true;
      const original = btn.textContent;
      btn.textContent = 'Seleccionando...';

      try {
        const r = await fetch('/updatek', {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ idKeyUpdate: idKey })
        });
        if (!r.ok) {
          const t = await r.text();
          throw new Error(`${r.status} - ${t}`);
        }

        // actualiza UI: solo 1 en verde
        marcarSeleccionada(idKey);

      } catch (err) {
        console.error(err);
        alert('Error al seleccionar la llave: ' + err.message);
        btn.disabled = false;
        btn.textContent = original;
      }
    }

    // ====== App init ======
    document.addEventListener('DOMContentLoaded', async () => {
      // Liquid bubble simple: posicionar sobre 'Inicio'
      const bubble = document.getElementById('liquidBubble');
      const mainHeader = document.getElementById('mainHeader');
      const navItems = document.querySelectorAll('.nav-item');
      function moveBubble(target) {
        const rect = target.getBoundingClientRect();
        const headerRect = mainHeader.getBoundingClientRect();
        bubble.style.left = `${rect.left - headerRect.left - 10}px`;
        bubble.style.width = `${rect.width + 20}px`;
      }
      window.addEventListener('load', () => { if (navItems[0]) moveBubble(navItems[0]); });
      navItems.forEach(item => item.addEventListener('mouseenter', () => moveBubble(item)));

      // Tema
      const themeToggle = document.getElementById('themeToggle');
      themeToggle.addEventListener('click', () => {
        document.body.classList.toggle('light-theme');
        themeToggle.textContent = document.body.classList.contains('light-theme') ? '☀️' : '🌙';
      });

      // Cargar y mostrar las llaves
      const keys = await fetchKeysStatus();
      renderKeys(keys);

      // ====== Form credenciales ======
      const credentialsForm = document.getElementById('credentialsForm');
      credentialsForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const keyFile = document.getElementById('keyFile').files[0];
        const certFile = document.getElementById('certFile').files[0];
        const keyPassword = document.getElementById('keyPassword').value;

        if (!keyFile || !certFile || !keyPassword) {
          alert('Por favor, sube la llave, el certificado y escribe la contraseña.');
          return;
        }

        try {
          // Crear FormData para enviar los archivos
          const formData = new FormData();
          formData.append('keyFile', keyFile);
          formData.append('certFile', certFile);
          formData.append('passKey', keyPassword);

          // Enviar al endpoint de Go
          const response = await fetch('/uploadk', {
            method: 'POST',
            body: formData,
            credentials: 'include' // Para incluir las cookies (JWT)
          });

          if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Error: ${response.status} - ${errorText}`);
          }

          const result = await response.json();
          console.log(result);
          if (result.valid) {
            /*
            // Guardar metadata localmente
            const meta = {
              keyName: keyFile.name,
              certName: certFile.name,
              uploadedAt: new Date().toISOString()
            };
            setOnboardingDone(meta);
            */
            
            // Recargar la página para mostrar la nueva llave
            location.reload();
          } else {
            alert('Error: ' + result.message);
          }
        } catch (err) {
          console.error('Error subiendo archivos:', err);
          alert('Error al subir los archivos: ' + err.message);
        }
      });

      // Modal de términos y condiciones
      const termsModal = document.getElementById('termsModal');
      const termsLink = document.getElementById('termsLink');
      const closeTerms = document.getElementById('closeTerms');
      const closeTerms2 = document.getElementById('closeTerms2');
      
      termsLink.addEventListener('click', () => {
        document.getElementById('termsDate').textContent = getTodayISO();
        termsModal.showModal();
      });
      
      closeTerms.addEventListener('click', () => termsModal.close());
      closeTerms2.addEventListener('click', () => termsModal.close());

      // Acciones header
      document.getElementById('goHome').addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));
    });