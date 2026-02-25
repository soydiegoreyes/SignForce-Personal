
document.addEventListener('DOMContentLoaded', () => {
    // Efecto liquid glass para la burbuja del menú
    const liquidBubble = document.getElementById('liquidBubble');
    const navItems = document.querySelectorAll('.nav-item');
    
    // Carrusel functionality
    let currentSlideIndex = 0;
    const slides = document.querySelectorAll('.carousel-slide');
    const indicators = document.querySelectorAll('.indicator');
    const totalSlides = slides.length;

    function showSlide(n) {
        slides.forEach(slide => slide.classList.remove('active'));
        indicators.forEach(indicator => indicator.classList.remove('active'));

        if (n >= totalSlides) currentSlideIndex = 0;
        if (n < 0) currentSlideIndex = totalSlides - 1;

        slides[currentSlideIndex].classList.add('active');
        if (indicators[currentSlideIndex]) {
            indicators[currentSlideIndex].classList.add('active');
        }
    }

    function changeSlide(n) {
        currentSlideIndex += n;
        showSlide(currentSlideIndex);
        pauseAutoPlay();
    }

    function currentSlide(n) {
        currentSlideIndex = n - 1;
        showSlide(currentSlideIndex);
        pauseAutoPlay();
    }

    function autoPlay() {
        currentSlideIndex++;
        showSlide(currentSlideIndex);
    }

    let autoPlayInterval = setInterval(autoPlay, 5000);

    function pauseAutoPlay() {
        clearInterval(autoPlayInterval);
        setTimeout(() => {
            autoPlayInterval = setInterval(autoPlay, 5000);
        }, 10000);
    }

    // Event listeners para el carrusel
    document.querySelector('.carousel-prev').addEventListener('click', () => changeSlide(-1));
    document.querySelector('.carousel-next').addEventListener('click', () => changeSlide(1));
    
    indicators.forEach((indicator, index) => {
        indicator.addEventListener('click', () => {
            currentSlide(index + 1);
            pauseAutoPlay();
        });
    });

    // Navegación por teclado
    document.addEventListener('keydown', function(e) {
        if (e.key === 'ArrowLeft') {
            changeSlide(-1);
        } else if (e.key === 'ArrowRight') {
            changeSlide(1);
        }
    });

    // Soporte para gestos táctiles
    let startX = 0;
    let endX = 0;

    document.querySelector('.carousel-container').addEventListener('touchstart', function(e) {
        startX = e.touches[0].clientX;
    });

    document.querySelector('.carousel-container').addEventListener('touchend', function(e) {
        endX = e.changedTouches[0].clientX;
        
        if (startX - endX > 50) {
            changeSlide(1);
        } else if (endX - startX > 50) {
            changeSlide(-1);
        }
    });

    // Efecto liquid glass para el menú
    const itemPositions = {};
    
    function calculatePositions() {
        navItems.forEach(item => {
            const rect = item.getBoundingClientRect();
            const headerRect = document.getElementById('mainHeader').getBoundingClientRect();
            
            itemPositions[item.dataset.item] = {
                left: rect.left - headerRect.left,
                width: rect.width
            };
        });
        
        moveBubble('home');
    }
    
    function moveBubble(itemName) {
        const item = itemPositions[itemName];
        if (item && liquidBubble) {
            liquidBubble.style.left = `${item.left - 10}px`;
            liquidBubble.style.width = `${item.width + 20}px`;
        }
    }
    
    navItems.forEach(item => {
        item.addEventListener('mouseenter', () => {
            moveBubble(item.dataset.item);
        });
        
        item.addEventListener('click', (e) => {
            e.preventDefault();
            moveBubble(item.dataset.item);
            
            const targetId = item.getAttribute('href');
            if (targetId && targetId !== '#') {
                document.querySelector(targetId).scrollIntoView({
                    behavior: 'smooth'
                });
            }
        });
    });

    window.addEventListener('load', calculatePositions);
    window.addEventListener('resize', calculatePositions);

    // Intersection Observer para animaciones
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

    // Contador animado para las estadísticas
    const statNumbers = document.querySelectorAll('.stat-number');
    const observerStats = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateNumbers();
                observerStats.unobserve(entry.target);
            }
        });
    }, { threshold: 0.5 });

    const statsSection = document.querySelector('.hero-stats');
    if (statsSection) {
        observerStats.observe(statsSection);
    }

    function animateNumbers() {
        statNumbers.forEach(stat => {
            const target = parseInt(stat.getAttribute('data-target'));
            const duration = 2000;
            const step = target / (duration / 16);
            let current = 0;
            
            const timer = setInterval(() => {
                current += step;
                if (current >= target) {
                    stat.textContent = target + '+';
                    clearInterval(timer);
                } else {
                    stat.textContent = Math.floor(current) + '+';
                }
            }, 16);
        });
    }
});
 // Funcionalidad del modal de video
const videoModal = document.getElementById('videoModal');
const demoVideo = document.getElementById('demoVideo');
const openVideoModalBtn = document.getElementById('openVideoModal');
const closeModalBtn = document.getElementById('closeModal');
const playPauseBtn = document.getElementById('playPauseBtn');
const muteBtn = document.getElementById('muteBtn');
const fullscreenBtn = document.getElementById('fullscreenBtn');
const progressBar = document.getElementById('progressBar');
const progressFill = document.getElementById('progressFill');
const videoTime = document.getElementById('videoTime');

// Abrir modal
if (openVideoModalBtn) {
    openVideoModalBtn.addEventListener('click', () => {
        videoModal.classList.add('active');
        document.body.style.overflow = 'hidden'; // Previene scroll
        demoVideo.play();
    });
}

// Cerrar modal
if (closeModalBtn) {
    closeModalBtn.addEventListener('click', closeVideoModal);
}

// Cerrar modal al hacer clic fuera del contenido
videoModal.addEventListener('click', (e) => {
    if (e.target === videoModal) {
        closeVideoModal();
    }
});

// Cerrar con tecla Escape
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && videoModal.classList.contains('active')) {
        closeVideoModal();
    }
});

function closeVideoModal() {
    videoModal.classList.remove('active');
    document.body.style.overflow = 'auto';
    demoVideo.pause();
    demoVideo.currentTime = 0;
}

// Control de reproducción
if (playPauseBtn && demoVideo) {
    playPauseBtn.addEventListener('click', () => {
        if (demoVideo.paused) {
            demoVideo.play();
            playPauseBtn.innerHTML = '<i class="fas fa-pause"></i> Pausar';
        } else {
            demoVideo.pause();
            playPauseBtn.innerHTML = '<i class="fas fa-play"></i> Reproducir';
        }
    });
}

// Control de volumen
if (muteBtn && demoVideo) {
    muteBtn.addEventListener('click', () => {
        demoVideo.muted = !demoVideo.muted;
        muteBtn.innerHTML = demoVideo.muted 
            ? '<i class="fas fa-volume-mute"></i> Sin Sonido'
            : '<i class="fas fa-volume-up"></i> Sonido';
    });
}

// Pantalla completa
if (fullscreenBtn) {
    fullscreenBtn.addEventListener('click', () => {
        if (!document.fullscreenElement) {
            videoModal.requestFullscreen().catch(err => {
                console.log(`Error al intentar pantalla completa: ${err.message}`);
            });
            fullscreenBtn.innerHTML = '<i class="fas fa-compress"></i> Salir';
        } else {
            document.exitFullscreen();
            fullscreenBtn.innerHTML = '<i class="fas fa-expand"></i> Pantalla Completa';
        }
    });
}

// Actualizar barra de progreso
if (demoVideo && progressBar && progressFill && videoTime) {
    demoVideo.addEventListener('timeupdate', updateProgressBar);
    
    progressBar.addEventListener('click', (e) => {
        const rect = progressBar.getBoundingClientRect();
        const pos = (e.clientX - rect.left) / rect.width;
        demoVideo.currentTime = pos * demoVideo.duration;
    });

    demoVideo.addEventListener('loadedmetadata', () => {
        updateVideoTime();
    });
}

function updateProgressBar() {
    if (demoVideo.duration) {
        const percent = (demoVideo.currentTime / demoVideo.duration) * 100;
        progressFill.style.width = `${percent}%`;
        updateVideoTime();
    }
}

function updateVideoTime() {
    const currentMinutes = Math.floor(demoVideo.currentTime / 60);
    const currentSeconds = Math.floor(demoVideo.currentTime % 60);
    const totalMinutes = Math.floor(demoVideo.duration / 60);
    const totalSeconds = Math.floor(demoVideo.duration % 60);
    
    videoTime.textContent = 
        `${currentMinutes.toString().padStart(2, '0')}:${currentSeconds.toString().padStart(2, '0')} / ` +
        `${totalMinutes.toString().padStart(2, '0')}:${totalSeconds.toString().padStart(2, '0')}`;
}

// Actualizar botón de play/pause cuando el video termina
if (demoVideo) {
    demoVideo.addEventListener('ended', () => {
        playPauseBtn.innerHTML = '<i class="fas fa-redo"></i> Repetir';
    });

    demoVideo.addEventListener('play', () => {
        playPauseBtn.innerHTML = '<i class="fas fa-pause"></i> Pausar';
    });

    demoVideo.addEventListener('pause', () => {
        playPauseBtn.innerHTML = '<i class="fas fa-play"></i> Reproducir';
    });
}