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

    // CARRUSEL CORREGIDO
    const carousel = document.querySelector('.carousel');
    const cards = document.querySelectorAll('.carousel-card');
    const prevBtn = document.querySelector('.carousel-prev');
    const nextBtn = document.querySelector('.carousel-next');
    
    if (!carousel || !cards.length) return;
    
    // Variables para controlar la animación
    let isAnimating = false;
    let currentIndex = 0;
    let isManualControl = false;
    let autoScrollInterval;
    
    // Configuración del carrusel
    const cardWidth = 280;
    const cardGap = 32;
    const cardTotalWidth = cardWidth + cardGap;
    
    // Función para actualizar posiciones de las tarjetas
    function updateCarousel(animate = true) {
      cards.forEach((card, index) => {
        // Remover la animación CSS automática cuando hay control manual
        if (isManualControl) {
          card.style.animation = 'none';
        }
        
        // Calcular nueva posición
        const position = (index - currentIndex) * cardTotalWidth;
        
        if (animate) {
          card.style.transition = 'transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1)';
        } else {
          card.style.transition = 'none';
        }
        
        card.style.transform = `translateX(${position}px)`;
      });
    }
    
    // Función para iniciar auto-scroll
    function startAutoScroll() {
      if (autoScrollInterval) clearInterval(autoScrollInterval);
      
      autoScrollInterval = setInterval(() => {
        if (!isManualControl && !isAnimating) {
          currentIndex = (currentIndex + 1) % cards.length;
          updateCarousel();
        }
      }, 3000); // Cambia cada 3 segundos
    }
    
    // Función para detener auto-scroll
    function stopAutoScroll() {
      if (autoScrollInterval) {
        clearInterval(autoScrollInterval);
        autoScrollInterval = null;
      }
    }
    
    // Pausar animación al hacer hover
    carousel.addEventListener('mouseenter', () => {
      isManualControl = true;
      stopAutoScroll();
      
      // Pausar animaciones CSS si existen
      cards.forEach(card => {
        card.style.animationPlayState = 'paused';
      });
    });
    
    carousel.addEventListener('mouseleave', () => {
      // Reanudar después de un delay
      setTimeout(() => {
        isManualControl = false;
        startAutoScroll();
      }, 1000);
    });
    
    // Navegación con botón anterior
    prevBtn.addEventListener('click', (e) => {
      e.preventDefault();
      if (isAnimating) return;
      
      isAnimating = true;
      isManualControl = true;
      stopAutoScroll();
      
      currentIndex = (currentIndex - 1 + cards.length) % cards.length;
      updateCarousel();
      
      setTimeout(() => {
        isAnimating = false;
      }, 500);
    });
    
    // Navegación con botón siguiente
    nextBtn.addEventListener('click', (e) => {
      e.preventDefault();
      if (isAnimating) return;
      
      isAnimating = true;
      isManualControl = true;
      stopAutoScroll();
      
      currentIndex = (currentIndex + 1) % cards.length;
      updateCarousel();
      
      setTimeout(() => {
        isAnimating = false;
      }, 500);
    });
    
    // Inicializar carrusel
    function initCarousel() {
      // Remover animaciones CSS automáticas inicialmente
      cards.forEach((card, index) => {
        card.style.position = 'absolute';
        card.style.left = '0';
        card.style.animation = 'none';
      });
      
      // Posicionar tarjetas inicialmente
      updateCarousel(false);
      
      // Iniciar auto-scroll después de un delay
      setTimeout(() => {
        if (!isManualControl) {
          startAutoScroll();
        }
      }, 2000);
    }
    
    // Soporte para navegación por teclado
    document.addEventListener('keydown', (e) => {
      if (!carousel.matches(':hover')) return;
      
      if (e.key === 'ArrowLeft') {
        e.preventDefault();
        prevBtn.click();
      } else if (e.key === 'ArrowRight') {
        e.preventDefault();
        nextBtn.click();
      }
    });
    
    // Soporte para touch/swipe en móviles
    let startX = 0;
    let isDragging = false;
    
    carousel.addEventListener('touchstart', (e) => {
      startX = e.touches[0].clientX;
      isDragging = true;
      isManualControl = true;
      stopAutoScroll();
    });
    
    carousel.addEventListener('touchmove', (e) => {
      if (!isDragging) return;
      e.preventDefault();
    });
    
    carousel.addEventListener('touchend', (e) => {
      if (!isDragging) return;
      
      const endX = e.changedTouches[0].clientX;
      const diffX = startX - endX;
      
      if (Math.abs(diffX) > 50) { // Mínimo swipe de 50px
        if (diffX > 0) {
          nextBtn.click();
        } else {
          prevBtn.click();
        }
      }
      
      isDragging = false;
    });
    
    // Inicializar el carrusel
    initCarousel();
    
    // Recalcular en resize
    window.addEventListener('resize', () => {
      updateCarousel(false);
    });
});