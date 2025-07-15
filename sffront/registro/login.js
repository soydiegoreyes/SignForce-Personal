document.addEventListener('DOMContentLoaded', () => {
    // Efecto liquid glass para la burbuja del menú
    const liquidBubble = document.getElementById('liquidBubble');
    
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
    const registerForm = document.getElementById('registerForm');

    // Para Formulario de Login
    loginForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        
        // Simula tabla usuarios en DB
        let usuarios = JSON.parse(sessionStorage.getItem('usuarios')) || [];
        if (usuarios==[]) {
            alert('No hay usuarios registrados');
            return;
        }
        if (validateEmail(email) && password.length >= 8) {
            const usuario = usuarios.find(u => u.email === email && u.password === password);
                
            if (usuario) {
                // Simulamos que guardamos la sesión
                sessionStorage.setItem('currentUser', JSON.stringify(usuario));
                
                // Redirigir según tipo de usuario
                alert('Inicio de sesión exitoso. Redirigiendo...');

                if (usuario.tipo === 'root') {
                    window.location.href = './../dashboards/dashboard_root.html';
                } else if (usuario.tipo === 'admin_equipo') {
                    window.location.href = './../dashboards/dashboard_admin.html';
                } else if (usuario.tipo === 'miembro_equipo') {
                    window.location.href = './../dashboards/dashboard_user.html';
                } else {
                    alert('Tipo de usuario no reconocido');
                }
            } else {
                alert('Credenciales incorrectas');
            }
        }
    });

    // Para Formulario de Registro
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        // Obtener valores del formulario
        const firstName = document.getElementById('first-name').value;
        const lastName = document.getElementById('last-name').value;
        const businessName = document.getElementById('business-name').value;
        const businessRFC = document.getElementById('business-rfc').value;
        const email = document.getElementById('register-email').value;
        const phone = document.getElementById('register-phone').value;
        const password = document.getElementById('register-password').value;
        const confirm = document.getElementById('register-confirm').value;
        const terms = document.getElementById('terms').checked;
        
        // Validaciones
        if (!firstName || !lastName) {
            alert('Por favor, añade un nombre del responsable del sistema');
            return;
        }
        const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,20}$/;
        if (!passwordRegex.test(password)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>La contraseña debe tener 8-20 caracteres, contener letras y números, sin espacios o caracteres especiales.</strong></p>';
            return;
        }
        if (password !== confirm) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>El password no coincide con la confirmación</strong></p>';
            return;
        }
        if (!businessName) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>El nombre de la empresa es requerido</strong></p>';
            return;
        }
        const RFCRegex = /^[A-Z\d]{13,20}$/;
        if (!RFCRegex.test(businessRFC)) {
            alert('Alguno de tus datos no es correcto.');
            document.getElementById('notice-box').innerHTML = '<p><strong>Tu RFC debe contener solo letras y números, sin espacios o caracteres especiales.</strong></p>';
            return;
        }
        if (!validateEmail(email)) {
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
        
        try {
            // Simula tabla de empresas o instituciones en DB
            let empresas = JSON.parse(sessionStorage.getItem('empresas')) || [];

            // Crear nueva empresa para enviar a DB
            const nuevaEmpresa = {
                id: Date.now().toString(),
                firstName,
                lastName,
                businessName,
                businessRFC,
                email,
                phone,
                password,
                termsAccepted: terms,
                status: 'pendiente_validacion'
            };
            
            // Agregar nueva empresa al arreglo
            empresas.push(nuevaEmpresa);
            console.log(nuevaEmpresa);
            // Guardar el arreglo actualizado en sessionStorage
            sessionStorage.setItem('empresas', JSON.stringify(empresas));

            // Redirigir a la pantalla de validación
            window.location.href = 'validacion.html?empresaId=' + nuevaEmpresa.id;
            
            alert('Registro exitoso. ¡Bienvenido a SignForce!');
            
        } catch (error) {
            console.error('Error:', error);
            alert(`Error al registrar: ${error.message}`);
        }
    });

    // Función de validación de email
    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }

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

document.addEventListener('DOMContentLoaded', () => {
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