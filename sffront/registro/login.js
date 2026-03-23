document.addEventListener('DOMContentLoaded', () => {
    //===================================== PARTICLES =========================================
    (function createParticles() {
        const container = document.getElementById('particles');
        if (!container) return;
        const count = 20;
        for (let i = 0; i < count; i++) {
            const p = document.createElement('div');
            p.classList.add('particle');
            const size = Math.random() * 3 + 1;
            p.style.width = size + 'px';
            p.style.height = size + 'px';
            p.style.left = Math.random() * 100 + '%';
            p.style.animationDuration = (Math.random() * 20 + 15) + 's';
            p.style.animationDelay = (Math.random() * 15) + 's';
            container.appendChild(p);
        }
    })();

    //===================================== NAV SCROLL =========================================
    const mainHeader = document.getElementById('mainHeader');
    if (mainHeader) {
        window.addEventListener('scroll', () => {
            if (window.scrollY > 50) {
                mainHeader.classList.add('scrolled');
            } else {
                mainHeader.classList.remove('scrolled');
            }
        }, { passive: true });
    }

    //===================================== LIQUID BUBBLE =========================================

    // Efecto liquid glass para la burbuja del menú
    const liquidBubble = document.getElementById('liquidBubble');
    const navItems = document.querySelectorAll('.nav-item');

    // Posiciones iniciales de los items del menú
    const itemPositions = {};

    // Calcular posiciones de los items
    function calculatePositions() {
      const header = document.getElementById('mainHeader');
      if (!header) return;

      navItems.forEach(item => {
        const rect = item.getBoundingClientRect();
        const headerRect = header.getBoundingClientRect();

        const key = item.dataset.item || item.textContent.trim();
        itemPositions[key] = {
          left: rect.left - headerRect.left,
          width: rect.width,
          top: rect.top - headerRect.top
        };
      });

      // Position bubble on first item by default
      const firstKey = Object.keys(itemPositions)[0];
      if (firstKey) moveBubble(firstKey);
    }

    // Mover la burbuja al item seleccionado con smooth animation
    function moveBubble(itemName) {
      const item = itemPositions[itemName];
      if (item && liquidBubble) {
        liquidBubble.style.left = `${item.left - 8}px`;
        liquidBubble.style.width = `${item.width + 16}px`;
        liquidBubble.style.opacity = '1';
      }
    }

    // Hide bubble when mouse leaves nav area
    function hideBubble() {
      if (liquidBubble) {
        liquidBubble.style.opacity = '0';
      }
    }

    // Event listeners para los items del menú
    navItems.forEach(item => {
      item.addEventListener('mouseenter', () => {
        const key = item.dataset.item || item.textContent.trim();
        moveBubble(key);
      });

      item.addEventListener('click', (e) => {
        const key = item.dataset.item || item.textContent.trim();
        moveBubble(key);
        // Allow normal link navigation for login page nav items
        const href = item.getAttribute('href');
        if (href && href.startsWith('#')) {
          e.preventDefault();
          const target = document.querySelector(href);
          if (target) {
            target.scrollIntoView({ behavior: 'smooth' });
          }
        }
      });
    });

    // Hide bubble when mouse leaves the nav item area
    const navItemContainer = liquidBubble ? liquidBubble.parentElement : null;
    if (navItemContainer) {
      navItemContainer.addEventListener('mouseleave', () => {
        hideBubble();
      });
    }

    // Calcular posiciones al cargar y al redimensionar
    window.addEventListener('load', calculatePositions);
    window.addEventListener('resize', calculatePositions);

    // Initial calculation after a brief delay to ensure layout is stable
    setTimeout(calculatePositions, 100);

    //==================================================================================================
    // Validacion representante legal
    const checkbox = document.getElementById('legal-responsible-checkbox');
    const legalContainer = document.getElementById('legal-responsible-container');

    checkbox.addEventListener('change', function() {
        if (this.checked) {
            legalContainer.style.display = 'none';
            // Limpiar los campos cuando se ocultan (opcional)
            document.getElementById('legal-first-name').value = document.getElementById('first-name').value;
            document.getElementById('legal-last-name').value = document.getElementById('last-name').value;
        } else {
            legalContainer.style.display = 'block';
            document.getElementById('legal-first-name').value = '';
            document.getElementById('legal-last-name').value = '';
        }
    });

    // Cambio entre formularios
    const tabs = document.querySelectorAll('.tab');
    const forms = document.querySelectorAll('.form');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const tabName = tab.getAttribute('data-tab');

            // Actualizar tabs
            tabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');

            // Actualizar formularios
            forms.forEach(form => form.classList.remove('active'));
            document.getElementById(`${tabName}Form`).classList.add('active');
        });
    });

    // Cambio a login desde registro
    document.getElementById('switchToLogin').addEventListener('click', (e) => {
        e.preventDefault();
        tabs.forEach(t => t.classList.remove('active'));
        tabs[0].classList.add('active');

        forms.forEach(form => form.classList.remove('active'));
        document.getElementById('loginForm').classList.add('active');
    });

    // Toggle para contraseñas
    const passwordToggles = document.querySelectorAll('.password-toggle');
    passwordToggles.forEach(toggle => {
        toggle.addEventListener('click', () => {
            const input = toggle.previousElementSibling;
            const type = input.getAttribute('type') === 'password' ? 'text' : 'password';
            input.setAttribute('type', type);
            toggle.classList.toggle('fa-eye');
            toggle.classList.toggle('fa-eye-slash');
        });
    });

    // Validación de formularios
    const loginForm = document.getElementById('loginForm');
    const EmailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = document.getElementById('login-email').value.trim();
        const login_password = document.getElementById('login-password').value.trim();
        if (EmailRegex.test(email) && login_password.length >= 8) {
            const loginRequest = {
                account: email,
                password: login_password
            };

            try {
                const response = await fetch("/loginUser", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Accept": "application/json"
                    },
                    credentials: "include", // Importante para las cookies
                    body: JSON.stringify(loginRequest)
                });

                console.log("Status:", response.status);

                if (response.ok) {
                    const data = await response.json();

                    if (data.error) {
                        sfAlert('Error: ${data.error}');
                        return;
                    }

                    // Si login exitoso, hacer redirect usando JavaScript
                    if (data.redirectTo) {
                        sfAlert('Redirecting to: ${data.redirectTo}');
                        // Redirect usando window.location
                        window.location.href = `${data.redirectTo}`;
                        // O si prefieres usar el dominio actual:
                        // window.location.href = data.redirectTo;
                    }
                } else {
                    const errorData = await response.json();
                    sfAlert('Error: ${errorData.error || 'Error desconocido'}');
                }

            } catch (error) {
                console.error("Error en login:", error);
                sfAlert('Error al intentar iniciar sesión: ${error.message}');
            }
        } else {
            sfAlert('Email inválido o password muy corto (mínimo 8 caracteres).');
        }

    });



    // Para Formulario de Registro
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        // Obtener valores del formulario
        const firstName = document.getElementById('first-name').value;
        const lastName = document.getElementById('last-name').value;
        const legalName = document.getElementById('legal-first-name').value;
        const legalLastName = document.getElementById('legal-last-name').value;
        const businessName = document.getElementById('business-name').value;
        const businessAlias = document.getElementById('business-alias').value;
        const businessTaxNum = document.getElementById('business-taxnum').value;
        const email = document.getElementById('register-email').value;
        const phone = document.getElementById('register-phone').value;
        const country = document.getElementById('country').value;
        const city = document.getElementById('city').value;
        const terms = document.getElementById('terms').checked;

        // Validaciones
        if (!legalName || !legalLastName) {
            document.getElementById('notice-box').textContent='Por favor, añade el nombre del responsable legal'; return;;
            return;
        }
        if (!businessName) {
            document.getElementById('notice-box').textContent='El nombre de la empresa es requerido';
            return;
        }

        const AliasRegex = /^[A-Za-z\d\s]{3,50}$/;
        if (!AliasRegex.test(businessAlias)) {
            document.getElementById('notice-box').textContent='El alias solo puede contener letras, números y espacios (3-50 caracteres)';
            return;
        }

        const TaxNumRegex = /^[A-Z\d]{12,20}$/;
        if (!TaxNumRegex.test(businessTaxNum)) {
            document.getElementById('notice-box').textContent='RFC inválido: solo letras y números, 12-20 caracteres';
            return;
        }


        if (!EmailRegex.test(email)) {
            document.getElementById('notice-box').textContent='Email inválido. Verifica que sea correcto';
            return;
        }
        const phoneRegex = /^\d{10}$/;
        if (!phoneRegex.test(phone)) {
            document.getElementById('notice-box').textContent='Teléfono inválido: 10 dígitos sin espacios';
            return;
        }
        if (!terms) {
            document.getElementById('notice-box').textContent='Debes aceptar los términos y condiciones';
            return;
        }

        // SIMULACION DE INSERSCION DE DATOS EN BASE INSTITUTIONS
        try {
            // Simula tabla de empresas o instituciones en DB
            let empresas = JSON.parse(sessionStorage.getItem('empresas')) || {};

            // Datos iniciales para nueva empresa para enviar a DB
            const nuevaEmpresa = {
                firstName,
                lastName,
                legalName,
                legalLastName,
                businessName,
                businessAlias,
                businessTaxNum,
                email,
                phone,
                country,
                city
            };

            const nuevaEmpresaJson = JSON.stringify(nuevaEmpresa);

            // Mostrar loading
            document.getElementById('sf-loading').style.display = 'flex';

            try {
                // manda los datos a la api para crear un nuevo cliente que empezara el proceso
                const response = await fetch('/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    },
                    credentials: 'include', // Importante para cookies/CORS
                    body: nuevaEmpresaJson
                });

                if (!response.ok) {
                    if (response.status === 401) {
                        window.location.href = '/login';
                    }
                    const errorText = await response.text();
                    throw new Error(`Error ${response.status}: ${errorText}`);
                }


                // Parsear la respuesta JSON
                const data = await response.json();
                console.log(data);

                // Verificar si hubo error del servidor
                if (!data.check || data.error) {
                    document.getElementById('sf-loading').style.display = 'none';
                    const nb = document.getElementById('notice-box');
                    nb.textContent = data.error || 'Error al registrar. Intenta de nuevo.';
                    return;
                }

                // Registro exitoso
                empresas[data.instId] = nuevaEmpresa;
                sessionStorage.setItem('empresas', JSON.stringify(empresas));

                // Mostrar pantalla de éxito
                document.getElementById('sf-loading').style.display = 'none';
                document.getElementById('sf-email-sent').textContent = email;
                const successEl = document.getElementById('sf-success');
                successEl.style.display = 'flex';

            } catch (error) {
                console.error('Error:', error);
                document.getElementById('sf-loading').style.display='none'; document.getElementById('notice-box').textContent='Error de conexión. Intenta de nuevo.';
            }

        } catch (error) {
            document.getElementById('sf-loading').style.display = 'none';
            console.error('Error:', error);
            document.getElementById('notice-box').textContent = 'Error al registrar. Intenta de nuevo.';
        }
    });



    // Intersection Observer para animaciones al hacer scroll
    const observerOptions = {
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('fade-in');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    document.querySelectorAll('.fade-in').forEach(element => {
        observer.observe(element);
    });
});