// Toggle de tema
const themeToggle = document.getElementById('themeToggle');
const body = document.body;

themeToggle.addEventListener('click', () => {
    const currentTheme = body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    body.setAttribute('data-theme', newTheme);
    
    // Cambiar ícono
    const icon = themeToggle.querySelector('i');
    icon.classList.toggle('fa-moon');
    icon.classList.toggle('fa-sun');
    
    // Guardar preferencia
    sessionStorage.setItem('theme', newTheme);
});

// Cargar tema guardado
const savedTheme = sessionStorage.getItem('theme') || 'dark';
body.setAttribute('data-theme', savedTheme);

// Configurar ícono inicial
const icon = themeToggle.querySelector('i');
if (savedTheme === 'light') {
    icon.classList.remove('fa-moon');
    icon.classList.add('fa-sun');
}

// Animación al hacer scroll
const fadeElements = document.querySelectorAll('.fade-in');

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.animationPlayState = 'running';
            observer.unobserve(entry.target);
        }
    });
}, { threshold: 0.1 });

fadeElements.forEach(element => {
    observer.observe(element);
});