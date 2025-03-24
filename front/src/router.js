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
            if (link && link.getAttribute('href').startsWith('/')) {
                e.preventDefault();
                this.navigate(link.getAttribute('href'));
            }
        });
    },

    navigate(path, addToHistory = true) {
        // Update browser history if needed
        if (addToHistory) {
            window.history.pushState({}, '', path);
        }

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
            matchedRoute.component(mainContent, params);
        } else {
            mainContent.innerHTML = '<div class="error-container"><h2>404 - Page Not Found</h2><p>The page you are looking for does not exist.</p></div>';
        }
    }
};