// Cambio de tema
const themeToggle = document.getElementById('themeToggle');
const body = document.body;

themeToggle.addEventListener('click', () => {
    const currentTheme = body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    body.setAttribute('data-theme', newTheme);
    sessionStorage.setItem('theme', newTheme);
});

// Cargar tema guardado
const savedTheme = sessionStorage.getItem('theme');
if (savedTheme) {
    body.setAttribute('data-theme', savedTheme);
}

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
    
    //_________________________________________________________________________________________
    // Simula tabla usuarios en DB ______________ borrar cuando se implemente base de datos**
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
            } else if (usuario.tipo === 'maestro_equipo') {
                window.location.href = './../dashboards/dashboard_admin.html';
            } else {
                alert('Tipo de usuario no reconocido');
            }
        } else {
            alert('Credenciales incorrectas');
        }
    }

    // Simular autenticación
    // En un caso real sería un fetch a la API
    /*
    fetch('/api/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Guardar token y redirigir según tipo de usuario
            if (data.user.tipo === 'root') {
                window.location.href = './../dashboards/dashboard-root.html';
            } else if (data.user.tipo === 'admin_equipo') {
                window.location.href = './../dashboards/dashboard-admin.html';
            } else {
                window.location.href = 'dashboard.html';
            }
        } else {
            alert('Credenciales incorrectas');
        }
    });
    */
    //_____________________________________________________________________________________
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
    // Otras validaciones...
    
    
    try {

        // Simula tabla de empresas o instituciones en DB _________________________________________
        // Obtener empresas existentes o iniciar un array vacío
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
        
        //________________________________________________________________
        // Realizar la petición fetch
        /*const response = await fetch('https://tuapi.com/registro', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        
        // Verificar si la respuesta es exitosa
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Error en el registro');
        }
        
        // Procesar respuesta exitosa
        const data = await response.json();*/
        alert('Registro exitoso. ¡Bienvenido a SignForce!');
        
        // Redireccionar o realizar otras acciones después del registro
        // window.location.href = '/dashboard';
        
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