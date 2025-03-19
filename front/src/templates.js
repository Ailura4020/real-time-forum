// src/templates.js
import Home from './pages/Home.js';
import About from './pages/About.js';

const templates = {
    Home,
    About,
};

export function renderTemplate(templateName) {
    const template = templates[templateName] || templates.Home;
    return template();
}
