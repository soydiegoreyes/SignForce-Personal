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