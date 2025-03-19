// src/pages/About.js
import styles from '../styles/App.module.css';

export default function About() {
    return `
        <div class="${styles.about}">
            <h1>About Page</h1>
            <p>This is the about page.</p>
            <a href="/">Go to Home</a>
        </div>
    `;
}
