import styles from '../styles/App.module.css';

export default function LoginModal({ onLogin, onClose }) {
    return `
        <div class="${styles.modal}">
            <div class="${styles.modalContent}">
                <span class="${styles.close}" onclick="${onClose}">&times;</span>
                <h2>Login</h2>
                <form id="loginForm">
                    <input type="email" name="identifier" placeholder="Email" required />
                    <input type="password" name="password" placeholder="Password" required />
                    <button type="submit">Login</button>
                </form>
                <div id="loginMessage"></div>
            </div>
        </div>
    `;
}
