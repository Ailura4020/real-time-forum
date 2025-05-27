import { api } from '../main.js';
import classes from '../styles/Chat.module.css';
import { chatTemplate } from "../templates.js"
import { connectWebSocket } from '../websocket.js';
import { getSocket } from '../websocket.js';

// Store current user data
export let currentUser = null;
export let currentReceiver = null;

// Gestion de la pagination dans le chat
let currentPage = 0;           // Page actuelle
const PAGE_SIZE = 10;          // Nombre de messages par page
let isLoadingMessages = false; // Évite les appels en double
let hasMoreMessages = true;    // Pour arrêter le chargement quand on a tout récupéré

export async function renderChat(container) {
    // console.log("CLASSES", classes) // CSS modules
    // Make classes available globally for other modules
    window.chatClasses = classes;

    // Export critical functions to window for access from other modules
    window.setCurrentReceiver = setCurrentReceiver;
    window.setupUserClickListener = setupUserClickListener;
    window.renderRecentConversations = renderRecentConversations;
    window.Chat = {
        setCurrentReceiver,
        setupUserClickListener,
        renderRecentConversations,
        conversationHistory
    };

    // Clear the container before rendering
    container.innerHTML = '';

    // Set up a cleanup function for when the component is unmounted
    window.addEventListener('beforeunload', cleanupChat);

    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = chatTemplate(classes); // Use the template with the classes
    container.appendChild(tempDiv);

    // Add WebSocket status indicator
    let wsStatus = document.createElement('div');
    wsStatus.id = 'ws-status';
    wsStatus.className = `${classes['ws-status']} ws-status connecting`;
    wsStatus.textContent = 'Connecting...';
    container.prepend(wsStatus);
    window.updateWsStatus = function(status, text) {
        wsStatus.className = `${classes['ws-status']} ws-status ${status}`;
        wsStatus.textContent = text;
    };

    const token = localStorage.getItem('token');
    if (token) {
        // Fetch user data before connecting to WebSocket
        // to ensure current user is available for filtering
        const userData = await fetchUserData();

        // Store current user data
        if (userData && userData.data) {
            currentUser = {
                id: userData.data.id,
                nickname: userData.data.nickname
            };

            // Also make it available globally for other modules
            window.currentUser = currentUser;

            // Load conversation history from local storage and from server
            await loadConversationsFromDatabase();
            renderRecentConversations();
        }

        // Ensure WebSocket is connected after status indicator is set
        connectWebSocket(token);
    }

    // Display user ID and nickname if available
    const userInfoElement = container.querySelector('#user-info');
    if (currentUser) {
        userInfoElement.innerHTML = `
            <p>User ID: ${currentUser.id}</p>
            <p>Nickname: ${currentUser.nickname}</p>
        `;
    } else {
        userInfoElement.innerHTML = '<p>User not logged in or data could not be retrieved.</p>';
    }

    const sendButton = document.getElementById('send-button');
    sendButton.addEventListener('click', async () => {
        console.log('[Chat] Attempting to send message to', currentReceiver);

        // Check current user is set
        if (!currentUser || !currentUser.id) {
            console.error('[Chat] Current user not set, cannot send message');
            return;
        }

        const messageContent = document.getElementById('messageInput').value;
        if (currentReceiver && messageContent) {
            const message = {
                type: 'private_message',
                to: currentReceiver.id,
                from_id: currentUser.id, // Explicitly include the sender ID
                content: messageContent
            };

            // Send the message via WebSocket
            console.log('[Chat] Sending message:', message);

            const socket = getSocket();
            if (socket && socket.readyState === WebSocket.OPEN) {
                // Send the message via WebSocket
                socket.send(JSON.stringify(message));

                // Create current timestamp
                const now = new Date();
                const formattedDate = now.toISOString().replace('T', ' ').substring(0, 19);

                // Ensure the receiver has a proper nickname
                if (!currentReceiver.nickname || currentReceiver.nickname === 'undefined') {
                    const userElement = document.getElementById(currentReceiver.id);
                    if (userElement) {
                        const nicknameSpan = userElement.querySelector('.user-nickname');
                        if (nicknameSpan) {
                            currentReceiver.nickname = nicknameSpan.textContent;
                        } else if (userElement.getAttribute('data-nickname')) {
                            currentReceiver.nickname = userElement.getAttribute('data-nickname');
                        } else {
                            currentReceiver.nickname = `User ${currentReceiver.id}`;
                        }
                    } else {
                        currentReceiver.nickname = `User ${currentReceiver.id}`;
                    }
                    console.log('[Chat] Updated receiver nickname to:', currentReceiver.nickname);
                }

                // Save message to conversation history
                conversationHistory.addMessage(
                    currentReceiver.id,
                    currentReceiver.nickname,
                    messageContent,
                    now.toISOString()
                );

                // Also display the message locally for the sender
                const messagesContainer = document.getElementById('messages');
                if (messagesContainer) {
                    // Use currentUser data for the sender info
                    let senderNickname = currentUser ? currentUser.nickname : "Me";

                    // Add the message to the chat
                    const messageElement = document.createElement('div');
                    messageElement.className = 'message sent';
                    if (classes.message) {
                        messageElement.classList.add(classes.message);
                    }
                    if (classes.sent) {
                        messageElement.classList.add(classes.sent);
                    }

                    messageElement.innerHTML = `
                        <strong>${senderNickname} (You)</strong>: ${messageContent} <br>
                        <small>${formattedDate}</small>
                    `;
                    messagesContainer.appendChild(messageElement);

                    // Scroll to the bottom
                    messagesContainer.scrollTop = messagesContainer.scrollHeight;
                }
            } else {
                console.log('[Chat] WebSocket is not open.');
            }

            // Reset the input field
            document.getElementById('messageInput').value = '';
        }
    });
}


// Store conversation history
export const conversationHistory = {
    conversations: [],

    // Add a message to history
    addMessage(userId, nickname, message, timestamp) {
        let conversation = this.getConversation(userId);
        if (!conversation) {
            conversation = {
                userId,
                nickname,
                messages: [],
                lastUpdated: timestamp,
                unread: false
            };
            this.conversations.push(conversation);
        } else {
            conversation.lastUpdated = timestamp;
        }
        conversation.unread = false;
        conversation.messages.push({
            content: message,
            timestamp,
            fromCurrentUser: true
        });
        this.sortConversations();
        renderRecentConversations();
    },

    // Add a received message
    addReceivedMessage(userId, nickname, message, timestamp) {
        let conversation = this.getConversation(userId);
        if (!conversation) {
            conversation = {
                userId,
                nickname,
                messages: [],
                lastUpdated: timestamp,
                unread: true
            };
            this.conversations.push(conversation);
        } else {
            conversation.lastUpdated = timestamp;
        }
        conversation.unread = true;
        conversation.messages.push({
            content: message,
            timestamp,
            fromCurrentUser: false
        });
        this.sortConversations();
        renderRecentConversations();
    },

    // Get a conversation by user ID
    getConversation(userId) {
        if (userId === null || userId === undefined) {
            console.warn('[Chat] getConversation called with invalid userId', userId);
            return null;
        }
        const userIdStr = userId.toString();
        return this.conversations.find(c => c.userId && c.userId.toString() === userIdStr) || null;
    },

    // Sort conversations by last updated time
    sortConversations() {
        this.conversations.sort((a, b) => new Date(b.lastUpdated) - new Date(a.lastUpdated));
    }
};

// Function to clean up chat resources
export function cleanupChat() {
    // Reset chat state
    currentUser = null;
    currentReceiver = null;

    // Close WebSocket connection - we'll use the global instance
    // that's already imported at the top of the file
    if (typeof closeWebSocket === 'function') {
        closeWebSocket();
    }

    console.log('[Chat] Chat component cleaned up');
}

// Make the cleanup function globally available for the router
window.chatCleanupFunction = cleanupChat;


// Function to fetch user data
async function fetchUserData() {
    const token = localStorage.getItem('token');
    if (!token) return null;

    try {
        // Adjust the endpoint as necessary
        return await api.get('/user');
    } catch (error) {
        console.error('Error fetching user data:', error);
        return null;
    }
}

// Function to load conversation history from database
async function loadConversationsFromDatabase() {
    if (!currentUser) return;

    console.log('[Chat] Loading conversation history from database for current user:', currentUser.id, '/', currentUser.nickname);

    try {
        // Clear existing conversations to avoid duplicates
        conversationHistory.conversations = [];

        // Fetch previous conversations from server
        const token = localStorage.getItem('token');
        if (!token) throw new Error('No auth token found');
        const response = await api.get('/messages/all', {
            headers: { Authorization: `Bearer ${token}` }
        });
        console.log('[Chat] Conversation response:', response);

        // Use response.conversations directly (not response.data.conversations)
        if (response && Array.isArray(response.conversations)) {
            const conversations = response.conversations;

            conversations.forEach(conversation => {
                if (!conversation.user_id || !conversation.messages || !Array.isArray(conversation.messages) || conversation.messages.length === 0) {
                    console.log('[Chat] Skipping invalid conversation:', conversation);
                    return;
                }
                const userId = conversation.user_id.toString();
                let nickname = conversation.nickname || `User ${userId}`;
                let lastUpdated = conversation.last_updated || new Date().toISOString();

                // Map messages to expected format
                const messages = conversation.messages.map(msg => ({
                    content: msg.content,
                    timestamp: msg.timestamp,
                    fromCurrentUser: !!msg.is_from_current_user
                }));

                // Sort messages by timestamp (oldest first)
                messages.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

                conversationHistory.conversations.push({
                    userId,
                    nickname,
                    messages,
                    lastUpdated
                });
            });

            // Sort conversations by last updated time
            conversationHistory.sortConversations();

            console.log('[Chat] Loaded', conversationHistory.conversations.length, 'conversations from database');
        } else {
            throw new Error('Invalid response format from server');
        }
    } catch (error) {
        console.error('[Chat] Error loading conversations from database:', error);
        // Show error to user
        const container = document.getElementById('recent-conversations');
        if (container) {
            container.innerHTML = '<div class="error">Failed to load conversations. Please try again later.</div>';
        }
    }
}

// Function to render recent conversations list
function renderRecentConversations() {
    const container = document.getElementById('recent-conversations');
    if (!container) return;

    if (!conversationHistory.conversations || conversationHistory.conversations.length === 0) {
        container.innerHTML = '<div class="no-conversations">No recent conversations</div>';
        return;
    }

    container.innerHTML = '';

    console.log('[Chat] Rendering conversations:', conversationHistory.conversations.length);

    // Display conversations (most recent first)
    conversationHistory.conversations.forEach(conv => {
        // Skip conversations with no messages
        if (!conv.messages || conv.messages.length === 0) {
            console.log('[Chat] Skipping empty conversation for user:', conv.userId);
            return;
        }

        // Ensure we have a valid nickname (not undefined)
        if (!conv.nickname || conv.nickname === 'undefined') {
            console.log('[Chat] Fixing missing nickname for conversation:', conv.userId);

            // Try to find the nickname in the users list
            const userElement = document.querySelector(`.user-item[id="${conv.userId}"]`);
            if (userElement && userElement.getAttribute('data-nickname')) {
                conv.nickname = userElement.getAttribute('data-nickname');
            } else {
                // Fallback nickname
                conv.nickname = `User ${conv.userId}`;
            }
        }

        const convElement = document.createElement('div');
        convElement.className = 'conversationItem';

        // Also add class from CSS module if available
        if (classes.conversationItem) {
            convElement.classList.add(classes.conversationItem);
        }

        // Add active class if this is the current conversation
        if (currentReceiver && currentReceiver.id &&
            currentReceiver.id.toString() === conv.userId.toString()) {
            convElement.classList.add('active');
            if (classes.active) {
                convElement.classList.add(classes.active);
            }
        }

        // Affichage du pseudo en gras si unread
        const nicknameClass = conv.unread ? 'unread' : '';
        // Get last message preview
        const lastMessage = conv.messages[conv.messages.length - 1];
        const lastMessageTime = new Date(lastMessage.timestamp);
        const formattedTime = formatTime(lastMessageTime);

        // Determine if last message is from current user (for styling)
        const isFromCurrentUser = lastMessage.fromCurrentUser;

        // Create a preview that indicates who sent the last message
        const previewPrefix = isFromCurrentUser ? 'You: ' : '';

        convElement.innerHTML = `
            <div><span class="${nicknameClass}">${conv.nickname}</span></div>
            <div class="${classes.lastMessagePreview}">${previewPrefix}${lastMessage.content}</div>
            <div class="${classes.conversationTime}">${formattedTime}</div>
        `;

        // Store user data as attributes for easier access
        convElement.setAttribute('data-user-id', conv.userId);
        convElement.setAttribute('data-nickname', conv.nickname);

        // Add click event to open conversation
        convElement.addEventListener('click', () => {
            setCurrentReceiver({
                id: conv.userId,
                nickname: conv.nickname
            });

            // Display conversation history
            displayConversationHistory(conv);
        });

        container.appendChild(convElement);
    });

    console.log('[Chat] Rendered conversations:', container.children.length);
}

// Helper function to format time for display
function formatTime(date) {
    const now = new Date();
    const yesterday = new Date(now);
    yesterday.setDate(yesterday.getDate() - 1);

    // Today: show only time
    if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    // Yesterday: show "Yesterday"
    else if (date.toDateString() === yesterday.toDateString()) {
        return 'Yesterday';
    }
    // Other days: show date
    else {
        return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
    }
}

// Function to display conversation history in the messages area
function displayConversationHistory(conversation) {
    const messagesContainer = document.getElementById('messages');
    if (!messagesContainer) return;

    // Clear current messages
    messagesContainer.innerHTML = '';

    // If no messages, show empty state
    if (!conversation.messages || conversation.messages.length === 0) {
        messagesContainer.innerHTML = '<div class="empty-conversation">No messages yet. Start typing to begin conversation.</div>';
        return;
    }

    // Group messages by date for better organization
    const messagesByDate = groupMessagesByDate(conversation.messages);

    // Display grouped messages with date separators
    Object.keys(messagesByDate).forEach(date => {
        // Add date separator
        const dateHeader = document.createElement('div');
        dateHeader.className = classes.conversationSeparator;
        dateHeader.innerHTML = `<span>${formatDateHeader(date)}</span>`;
        messagesContainer.appendChild(dateHeader);

        // Display messages for this date
        messagesByDate[date].forEach(msg => {
            const messageElement = document.createElement('div');
            messageElement.className = msg.fromCurrentUser ? 'message sent' : 'message received';

            // Apply CSS module classes if available
            if (classes.message) {
                messageElement.classList.add(classes.message);
            }
            if (classes[msg.fromCurrentUser ? 'sent' : 'received']) {
                messageElement.classList.add(classes[msg.fromCurrentUser ? 'sent' : 'received']);
            }

            // Format the timestamp
            let formattedTime = '';
            try {
                const timestamp = new Date(msg.timestamp);
                formattedTime = timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
            } catch (e) {
                console.error('[Chat] Error formatting message timestamp:', e);
                formattedTime = 'Unknown time';
            }

            // Get sender nickname - ensure it's never undefined
            let senderNickname;

            if (msg.fromCurrentUser) {
                // Current user's nickname
                senderNickname = currentUser && currentUser.nickname
                    ? currentUser.nickname
                    : 'You';
                senderNickname += ' (You)';
            } else {
                // Conversation partner's nickname
                senderNickname = conversation.nickname
                    ? conversation.nickname
                    : msg.senderNickname
                        ? msg.senderNickname
                        : `User ${conversation.userId}`;
            }

            messageElement.innerHTML = `
                <strong>${senderNickname}</strong>: ${msg.content}
                <small>${formattedTime}</small>
            `;
            messagesContainer.appendChild(messageElement);
        });
    });

    // Scroll to bottom
    messagesContainer.scrollTop = messagesContainer.scrollHeight;

}

function throttle(fn, delay) {
    let lastCall = 0;
    return function (...args) {
        const now = Date.now();
        if (now - lastCall >= delay) {
            lastCall = now;
            fn.apply(this, args);
        }
    };
}

// Appel à l'API paginée
async function fetchOlderMessages(userId, page) {
    const offset = page * PAGE_SIZE;

    try {
        const res = await api.get(`/messages/${userId}?offset=${offset}&limit=${PAGE_SIZE}`);
        return res.data || [];
    } catch (err) {
        console.error("[Chat] Erreur lors du chargement des anciens messages :", err);
        return [];
    }
}

// Ajoute les anciens messages en haut du conteneur
function prependMessagesToConversation(conversation, messages) {
    const container = document.getElementById('messages');
    const oldScrollHeight = container.scrollHeight;

    messages.reverse().forEach(msg => {
        const fromCurrentUser = msg.from_id?.toString() === currentUser.id.toString();
        conversation.messages.unshift({
            content: msg.content,
            timestamp: msg.timestamp,
            fromCurrentUser: fromCurrentUser,
            senderId: msg.from_id
        });
    });

    // Re-affiche les messages
    displayConversationHistory(conversation);

    // Ajuste le scroll pour rester au bon endroit
    const newScrollHeight = container.scrollHeight;
    container.scrollTop = newScrollHeight - oldScrollHeight;
}


// Helper function to group messages by date
function groupMessagesByDate(messages) {
    const groups = {};

    messages.forEach(msg => {
        const date = new Date(msg.timestamp);
        const dateKey = date.toISOString().split('T')[0]; // YYYY-MM-DD format

        if (!groups[dateKey]) {
            groups[dateKey] = [];
        }

        groups[dateKey].push(msg);
    });

    return groups;
}

// Helper function to format date headers
function formatDateHeader(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const yesterday = new Date(now);
    yesterday.setDate(yesterday.getDate() - 1);

    // Format as "Today", "Yesterday", or date
    if (date.toDateString() === now.toDateString()) {
        return 'Today';
    } else if (date.toDateString() === yesterday.toDateString()) {
        return 'Yesterday';
    } else {
        // Format as: "Mon, Jan 1, 2023"
        return date.toLocaleDateString(undefined, {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });
    }
}

export function setCurrentReceiver(user) {
    if (!user || !user.id) {
        console.warn('[Chat] Invalid user object for receiver');
        return;
    }

    // Ensure user.id is a string
    user.id = user.id.toString();

    // Fix undefined nickname
    if (!user.nickname || user.nickname === 'undefined') {
        // Try to find the nickname in the users list
        const userElement = document.querySelector(`.user-item[id="${user.id}"]`);
        if (userElement) {
            const nicknameElement = userElement.querySelector('.user-nickname');
            if (nicknameElement) {
                user.nickname = nicknameElement.textContent;
            } else if (userElement.getAttribute('data-nickname')) {
                user.nickname = userElement.getAttribute('data-nickname');
            } else {
                user.nickname = `User ${user.id}`;
            }
        } else {
            user.nickname = `User ${user.id}`;
        }
    }

    console.log("[Chat] Setting current receiver:", user);

    // Supprime la pastille rouge si elle existe (notification de nouveau message)
    const userElement = document.getElementById(user.id);
    if (userElement && userElement.classList.contains('user-notification')) {
        userElement.classList.remove('user-notification');
    }


    if (typeof window.hideNotificationBadge === 'function') {
        window.hideNotificationBadge();
    }


    currentReceiver = user;
    // Also set in window scope for access from other modules
    window.currentReceiver = user;

    if (!currentReceiver) {
        console.warn('[Chat] No recipient selected.');

        const recipientElement = document.getElementById('chat-recipient');
        if (recipientElement) {
            recipientElement.innerHTML = `<span class="recipient-label">Select a user to start a chat</span>`;
        }
        return;
    }

    const recipientElement = document.getElementById('chat-recipient');
    if (recipientElement) {
        // Check if user is online
        const userElements = document.querySelectorAll('.user-item');
        let isOnline = false;
        userElements.forEach(el => {
            if (el.id === user.id.toString()) {
                isOnline = true;
            }
        });

        const onlineStatus = isOnline ?
            '<span style="color: #4CAF50; margin-left: 10px;">● Online</span>' :
            '<span style="color: #888; margin-left: 10px;">● Offline</span>';

        recipientElement.innerHTML = `<strong>${user.nickname}</strong> ${onlineStatus}`;
    }
    console.log("[Chat] Conversation opened with", user);

    // Mark conversation as active in UI using data attributes for more reliable matching
    const conversationElements = document.querySelectorAll('.conversationItem');
    conversationElements.forEach(el => {
        el.classList.remove('active');
        if (classes.active) {
            el.classList.remove(classes.active);
        }

        const elUserId = el.getAttribute('data-user-id');
        // Match by user ID which is more reliable than nickname
        if (elUserId === user.id.toString()) {
            el.classList.add('active');
            if (classes.active) {
                el.classList.add(classes.active);
            }
        }
    });

    // Focus message input for immediate typing
    setTimeout(() => {
        const messageInput = document.getElementById('messageInput');
        if (messageInput) {
            messageInput.focus();
        }
    }, 100);

    // Load conversation history if it exists
    let conversation = conversationHistory.getConversation(user.id);

    // If conversation exists but has no nickname, update it
    if (conversation && (!conversation.nickname || conversation.nickname === 'undefined')) {
        conversation.nickname = user.nickname;
    }

    if (conversation) {
        displayConversationHistory(conversation);
        // Ajoute ici le scroll listener une seule fois
        const messagesContainer = document.getElementById('messages');
        messagesContainer.addEventListener('scroll', throttle(async () => {
            if (messagesContainer.scrollTop <= 10 && !isLoadingMessages && hasMoreMessages) {
                isLoadingMessages = true;
                console.log("[Chat] Chargement de messages supplémentaires...");

                const olderMessages = await fetchOlderMessages(conversation.userId, currentPage + 1);
                if (olderMessages && olderMessages.length > 0) {
                    prependMessagesToConversation(conversation, olderMessages);
                    currentPage++;
                } else {
                    hasMoreMessages = false;
                    console.log("[Chat] Plus de messages à charger.");
                }

                isLoadingMessages = false;
            }
        }, 400));
    } else {
        // Create a new empty conversation
        conversation = {
            userId: user.id,
            nickname: user.nickname,
            messages: [],
            lastUpdated: new Date().toISOString()
        };

        // Add to history
        conversationHistory.conversations.push(conversation);
    }

    // Marquer la conversation comme lue
    if (conversation) {
        conversation.unread = false;
        renderRecentConversations();
    }
}

// notification conversation 
// --- Gestion de la pastille rouge ---
function addNotificationToConversation(conversationId) {
    const conversationElement = document.getElementById(`conversation-${conversationId}`);
    if (conversationElement && !conversationElement.classList.contains('user-notification')) {
        conversationElement.classList.add('user-notification');
    }
}

function openConversation(conversationId) {
    const conversationElement = document.getElementById(`conversation-${conversationId}`);
    if (conversationElement) {
        conversationElement.classList.remove('user-notification');
    }

    // Ici, ajoute le code pour charger et afficher les messages
    loadConversationMessages(conversationId); // à adapter selon ton code
}

export function getCurrentReceiver() {
    return currentReceiver;
}


export function setupUserClickListener() {

    const userElements = document.querySelectorAll('.user-item');
    if (userElements.length === 0) {
        console.log("No users found.");
        return;  // If no users, stop here
    }

    userElements.forEach(el => {
        el.addEventListener('click', () => {
            const userId = el.id;
            // Extract nickname without the online indicator
            const nicknameSpan = el.querySelector('span:not(.userOnline)');
            const nickname = nicknameSpan ? nicknameSpan.textContent : el.textContent.trim();

            console.log('[Chat] User selected:', userId, nickname);
            setCurrentReceiver({ id: userId, nickname });
        });
    });
}




