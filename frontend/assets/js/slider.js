// Slider/Carousel JavaScript
class Slider {
    constructor(container, options = {}) {
        this.container = container;
        this.slides = container.querySelectorAll('.slide');
        this.currentSlide = 0;
        this.isAutoPlaying = options.autoPlay !== false;
        this.autoPlayInterval = options.autoPlayInterval || 5000;
        this.transitionDuration = options.transitionDuration || 500;
        this.autoPlayTimer = null;
        
        this.init();
    }
    
    init() {
        if (this.slides.length === 0) return;
        
        // Create navigation elements
        this.createNavigation();
        
        // Show first slide
        this.showSlide(0);
        
        // Start autoplay if enabled
        if (this.isAutoPlaying) {
            this.startAutoPlay();
        }
        
        // Add event listeners
        this.addEventListeners();
        
        // Add touch/swipe support
        this.addTouchSupport();
    }
    
    createNavigation() {
        // Create dots navigation
        const dotsContainer = document.createElement('div');
        dotsContainer.className = 'slider-dots';
        
        this.slides.forEach((_, index) => {
            const dot = document.createElement('button');
            dot.className = 'slider-dot';
            dot.setAttribute('data-slide', index);
            dot.addEventListener('click', () => this.goToSlide(index));
            dotsContainer.appendChild(dot);
        });
        
        this.container.appendChild(dotsContainer);
        this.dots = dotsContainer.querySelectorAll('.slider-dot');
        
        // Create arrow navigation
        const prevBtn = document.createElement('button');
        prevBtn.className = 'slider-arrow slider-prev';
        prevBtn.innerHTML = '<i class="fas fa-chevron-left"></i>';
        prevBtn.addEventListener('click', () => this.prevSlide());
        
        const nextBtn = document.createElement('button');
        nextBtn.className = 'slider-arrow slider-next';
        nextBtn.innerHTML = '<i class="fas fa-chevron-right"></i>';
        nextBtn.addEventListener('click', () => this.nextSlide());
        
        this.container.appendChild(prevBtn);
        this.container.appendChild(nextBtn);
        
        // Add CSS for navigation
        this.addNavigationStyles();
    }
    
    addNavigationStyles() {
        const style = document.createElement('style');
        style.textContent = `
            .slider {
                position: relative;
            }
            
            .slider-dots {
                position: absolute;
                bottom: 20px;
                left: 50%;
                transform: translateX(-50%);
                display: flex;
                gap: 10px;
                z-index: 10;
            }
            
            .slider-dot {
                width: 12px;
                height: 12px;
                border-radius: 50%;
                border: 2px solid rgba(255, 255, 255, 0.5);
                background: transparent;
                cursor: pointer;
                transition: all 0.3s ease;
            }
            
            .slider-dot.active {
                background: white;
                border-color: white;
            }
            
            .slider-dot:hover {
                background: rgba(255, 255, 255, 0.3);
            }
            
            .slider-arrow {
                position: absolute;
                top: 50%;
                transform: translateY(-50%);
                background: rgba(0, 0, 0, 0.3);
                border: none;
                color: white;
                width: 50px;
                height: 50px;
                border-radius: 50%;
                cursor: pointer;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 18px;
                transition: all 0.3s ease;
                z-index: 10;
            }
            
            .slider-arrow:hover {
                background: rgba(0, 0, 0, 0.6);
                transform: translateY(-50%) scale(1.1);
            }
            
            .slider-prev {
                left: 20px;
            }
            
            .slider-next {
                right: 20px;
            }
            
            .slide {
                opacity: 0;
                transition: opacity ${this.transitionDuration}ms ease-in-out;
            }
            
            .slide.active {
                opacity: 1;
            }
            
            @media (max-width: 768px) {
                .slider-arrow {
                    width: 40px;
                    height: 40px;
                    font-size: 14px;
                }
                
                .slider-prev {
                    left: 10px;
                }
                
                .slider-next {
                    right: 10px;
                }
                
                .slider-dots {
                    bottom: 15px;
                }
                
                .slider-dot {
                    width: 10px;
                    height: 10px;
                }
            }
        `;
        
        document.head.appendChild(style);
    }
    
    showSlide(index) {
        // Hide all slides
        this.slides.forEach(slide => {
            slide.classList.remove('active');
        });
        
        // Remove active class from all dots
        this.dots.forEach(dot => {
            dot.classList.remove('active');
        });
        
        // Show current slide
        this.slides[index].classList.add('active');
        this.dots[index].classList.add('active');
        
        this.currentSlide = index;
    }
    
    nextSlide() {
        const nextIndex = (this.currentSlide + 1) % this.slides.length;
        this.showSlide(nextIndex);
        this.restartAutoPlay();
    }
    
    prevSlide() {
        const prevIndex = this.currentSlide === 0 ? this.slides.length - 1 : this.currentSlide - 1;
        this.showSlide(prevIndex);
        this.restartAutoPlay();
    }
    
    goToSlide(index) {
        if (index >= 0 && index < this.slides.length) {
            this.showSlide(index);
            this.restartAutoPlay();
        }
    }
    
    startAutoPlay() {
        this.stopAutoPlay();
        this.autoPlayTimer = setInterval(() => {
            this.nextSlide();
        }, this.autoPlayInterval);
    }
    
    stopAutoPlay() {
        if (this.autoPlayTimer) {
            clearInterval(this.autoPlayTimer);
            this.autoPlayTimer = null;
        }
    }
    
    restartAutoPlay() {
        if (this.isAutoPlaying) {
            this.startAutoPlay();
        }
    }
    
    addEventListeners() {
        // Pause autoplay on hover
        this.container.addEventListener('mouseenter', () => {
            this.stopAutoPlay();
        });
        
        this.container.addEventListener('mouseleave', () => {
            if (this.isAutoPlaying) {
                this.startAutoPlay();
            }
        });
        
        // Keyboard navigation
        document.addEventListener('keydown', (e) => {
            if (this.container.contains(document.activeElement) || this.container.matches(':hover')) {
                if (e.key === 'ArrowLeft') {
                    e.preventDefault();
                    this.prevSlide();
                } else if (e.key === 'ArrowRight') {
                    e.preventDefault();
                    this.nextSlide();
                }
            }
        });
        
        // Visibility change (pause when tab is not visible)
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.stopAutoPlay();
            } else if (this.isAutoPlaying) {
                this.startAutoPlay();
            }
        });
    }
    
    addTouchSupport() {
        let startX = 0;
        let endX = 0;
        let isDragging = false;
        
        this.container.addEventListener('touchstart', (e) => {
            startX = e.touches[0].clientX;
            isDragging = true;
            this.stopAutoPlay();
        }, { passive: true });
        
        this.container.addEventListener('touchmove', (e) => {
            if (!isDragging) return;
            endX = e.touches[0].clientX;
        }, { passive: true });
        
        this.container.addEventListener('touchend', (e) => {
            if (!isDragging) return;
            
            const diffX = startX - endX;
            const threshold = 50;
            
            if (Math.abs(diffX) > threshold) {
                if (diffX > 0) {
                    this.nextSlide();
                } else {
                    this.prevSlide();
                }
            }
            
            isDragging = false;
            if (this.isAutoPlaying) {
                this.startAutoPlay();
            }
        }, { passive: true });
    }
    
    // Public methods
    play() {
        this.isAutoPlaying = true;
        this.startAutoPlay();
    }
    
    pause() {
        this.isAutoPlaying = false;
        this.stopAutoPlay();
    }
    
    destroy() {
        this.stopAutoPlay();
        // Remove event listeners and elements
        this.container.innerHTML = '';
    }
}

// Initialize all sliders on page load
document.addEventListener('DOMContentLoaded', function() {
    const sliders = document.querySelectorAll('.slider');
    
    sliders.forEach(slider => {
        new Slider(slider, {
            autoPlay: true,
            autoPlayInterval: 5000,
            transitionDuration: 500
        });
    });
});

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = Slider;
} 