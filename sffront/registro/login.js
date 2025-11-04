document.addEventListener('DOMContentLoaded', () => {
    //===================================== LIQUID BUBBLE =========================================

    // Efecto liquid glass para la burbuja del menú
    const liquidBubble = document.getElementById('liquidBubble');
    const navItems = document.querySelectorAll('.nav-item');
    
    // Posiciones iniciales de los items del menú
    const itemPositions = {};
    
    // Calcular posiciones de los items
    function calculatePositions() {
      navItems.forEach(item => {
        const rect = item.getBoundingClientRect();
        const headerRect = document.getElementById('mainHeader').getBoundingClientRect();
        
        itemPositions[item.dataset.item] = {
          left: rect.left - headerRect.left,
          width: rect.width
        };
      });
      
      // Posicionar burbuja en el item activo (home por defecto)
      moveBubble('home');
    }
    
    // Mover la burbuja al item seleccionado
    function moveBubble(itemName) {
      const item = itemPositions[itemName];
      if (item) {
        liquidBubble.style.left = `${item.left - 10}px`;
        liquidBubble.style.width = `${item.width + 20}px`;
      }
    }
    
    // Event listeners para los items del menú
    navItems.forEach(item => {
      item.addEventListener('mouseenter', () => {
        moveBubble(item.dataset.item);
      });
      
      item.addEventListener('click', (e) => {
        e.preventDefault();
        moveBubble(item.dataset.item);
        // Scroll suave al hacer clic en los enlaces
        const targetId = item.getAttribute('href');
        if (targetId && targetId !== '#') {
          document.querySelector(targetId).scrollIntoView({
            behavior: 'smooth'
          });
        }
      });
    });
    
    // Calcular posiciones al cargar y al redimensionar
    window.addEventListener('load', calculatePositions);
    window.addEventListener('resize', calculatePositions);

    // Toggle de tema claro/oscuro
    const themeToggle = document.getElementById('themeToggle');
    themeToggle.addEventListener('click', () => {
      document.body.classList.toggle('light-theme');
      
      // Cambiar icono
      const icon = themeToggle.querySelector('i');
      if (document.body.classList.contains('light-theme')) {
        icon.classList.remove('fa-moon');
        icon.classList.add('fa-sun');
      } else {
        icon.classList.remove('fa-sun');
        icon.classList.add('fa-moon');
      }
    });
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
                const response = await fetch("http://192.168.1.66:8000/loginUser", {
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
                        alert(`Error: ${data.error}`);
                        return;
                    }

                    // Si login exitoso, hacer redirect usando JavaScript
                    if (data.redirectTo) {
                        alert(`Redirecting to: ${data.redirectTo}`);
                        // Redirect usando window.location
                        window.location.href = `${data.redirectTo}`;
                        // O si prefieres usar el dominio actual:
                        // window.location.href = data.redirectTo;
                    }
                } else {
                    const errorData = await response.json();
                    alert(`Error: ${errorData.error || 'Error desconocido'}`);
                }
                
            } catch (error) {
                console.error("Error en login:", error);
                alert(`Error al intentar iniciar sesión: ${error.message}`);
            }
        } else {
            alert("Email inválido o password muy corto (mínimo 8 caracteres).");
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
            alert('Por favor, añade un nombre del responsable del sistema');
            return;
        }
        if (!businessName) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>El nombre de la empresa es requerido</strong></p>';
            return;
        }

        const AliasRegex = /^[A-Z\d]{5,20}$/;
        if (!AliasRegex.test(businessAlias)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Tu Alias debe contener solo letras y números, sin espacios o caracteres especiales y un máximos de 20 caracteres.</strong></p>';
            return;
        }

        const TaxNumRegex = /^[A-Z\d]{12,20}$/;
        if (!TaxNumRegex.test(businessTaxNum)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Tu RFC debe contener solo letras y números, sin espacios o caracteres especiales.</strong></p>';
            return;
        }
        
        
        if (!EmailRegex.test(email)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Asegurate que tu email no tiene espacios y sea correcto</strong></p>';
            return;
        }
        const phoneRegex = /^\d{10}$/;
        if (!phoneRegex.test(phone)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Tu telefono a 10 digitos sin espacios ni caracteres especiales</strong></p>';
            return;
        }
        if (!terms) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Lee y acepta el Acuerdo de Términos y Condiciones</strong></p>';
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
            
            try {
                // manda los datos a la api para crear un nuevo cliente que empezara el proceso
                const response = await fetch('http://192.168.1.66:8000/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    },
                    credentials: 'include', // Importante para cookies/CORS
                    body: nuevaEmpresaJson
                });

                if (!response.ok) {
                    throw new Error('Error en la respuesta ' + response.statusText);
                }

                
                // Parsear la respuesta JSON
                const data = await response.json();
                console.log(data);
                // Agregar nueva empresa al arreglo local
                empresas[data.instId] = nuevaEmpresa;
                sessionStorage.setItem('empresas', JSON.stringify(empresas));

                // Redirigir usando el ID del servidor
                window.location.href = '/login';
                
                alert('¡Bienvenido a SignForce! Revisa tu bandeja de entrada.');

            } catch (error) {
                console.error('Error:', error);
                alert(`Error al registrar: ${error.message}`);
            }
            
        } catch (error) {
            console.error('Error:', error);
            alert(`Error al registrar: ${error.message}`);
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