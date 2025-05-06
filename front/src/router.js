import { renderHomePage } from './pages/Home.js';
import { renderAboutPage } from './pages/About.js';
import { renderChat } from './pages/Chat.js';
import { renderLoginPage } from './components/Login.js';
import { renderRegistrationPage } from './components/Registration.js';
import { renderPostPage } from './components/Post.js';

const routes = [
    { path: '/', component: renderHomePage },
    { path: '/about', component: renderAboutPage },
    { path: '/chat', component: renderChat },
    { path: '/login', component: renderLoginPage },
    { path: '/register', component: renderRegistrationPage },
    { path: '/post/:id', component: renderPostPage, params: true }
];

export const router = {
    currentRoute: null,
    currentComponent: null,

    init() {
        // Handle initial page load
        this.navigate(window.location.pathname);

        // Handle browser navigation
        window.addEventListener('popstate', () => {
            this.navigate(window.location.pathname, false);
        });

        // Intercept link clicks for SPA navigation
        document.addEventListener('click', (e) => {
            const link = e.target.closest('a');
            const href = link?.getAttribute('href');
            if (link && href && href.startsWith('/')) {
                e.preventDefault();
                // Basic URL sanitization
                const sanitizedUrl = href.replace(/[^\w\s/-]/gi, '');
                this.navigate(sanitizedUrl);
            }
        });
    },

    navigate(path, addToHistory = true) {
        // Don't do anything if we're already on this route
        if (this.currentRoute === path) {
            console.log(`[Router] Already on route: ${path}`);
            return;
        }
        
        // Clean up previous component if needed
        if (this.currentRoute) {
            console.log(`[Router] Cleaning up previous component for route: ${this.currentRoute}`);
            this.cleanupCurrentComponent();
        }
        
        // Update browser history if needed
        if (addToHistory) {
            window.history.pushState({}, '', path);
        }
        
        console.log(`[Router] Navigating to: ${path}`);

        // Find matching route
        let matchedRoute = null;
        let params = {};

        for (const route of routes) {
            if (route.params) {
                // Handle parameterized routes
                const pathParts = path.split('/');
                const routeParts = route.path.split('/');

                if (pathParts.length === routeParts.length) {
                    let isMatch = true;

                    for (let i = 0; i < routeParts.length; i++) {
                        if (routeParts[i].startsWith(':')) {
                            // This is a parameter
                            const paramName = routeParts[i].substring(1);
                            params[paramName] = pathParts[i];
                        } else if (routeParts[i] !== pathParts[i]) {
                            isMatch = false;
                            break;
                        }
                    }

                    if (isMatch) {
                        matchedRoute = route;
                        break;
                    }
                }
            } else if (route.path === path) {
                matchedRoute = route;
                break;
            }
        }

        // Render the matched route or show 404
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = '';

        if (matchedRoute) {
            this.currentRoute = path;
            this.currentComponent = matchedRoute.component;
            matchedRoute.component(mainContent, params);
        } else {
            this.currentComponent = null;
            mainContent.innerHTML = '<div class="error-container"><h2>404 - Page Not Found</h2><p>The page you are looking for does not exist.</p></div>';
        }
            },
            
            cleanupCurrentComponent() {
        // Call cleanup function if the component has one
        if (this.currentComponent) {
            // Handle Chat component special cleanup
            if (this.currentComponent.name === 'renderChat') {
                        // We don't need to import cleanupChat here - it's already available globally
                        // from the existing import in the routes list
                        const cleanupFunction = window.chatCleanupFunction || null;
                        if (cleanupFunction && typeof cleanupFunction === 'function') {
                            cleanupFunction();
                }
            }
            
            // Handle other components with cleanup methods
            if (this.currentComponent.cleanup && typeof this.currentComponent.cleanup === 'function') {
                this.currentComponent.cleanup();
            }
        }
    }
};