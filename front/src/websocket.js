import { setupUserClickListener } from "./pages/Chat";

export let socket = null;

// Cache for user nicknames to ensure we always have a valid nickname
const userNicknameCache = new Map();

// Export closeWebSocket function for other modules to use
export function closeWebSocket() {
  if (socket && socket.readyState !== WebSocket.CLOSED) {
    console.log('[WebSocket] Closing connection...');
    socket.close();
    socket = null;
  }
}

// Make closeWebSocket globally accessible for components that need it
window.closeWebSocket = closeWebSocket;

// Set to keep track of displayed users to prevent duplicates
const displayedUsers = new Set();

export function updateConnectedUsers(userContainer, users) {
  console.log('[updateConnectedUsers] Utilisateurs reçus :', users);
  
  // Clear the container and displayed users set
  userContainer.innerHTML = '';
  displayedUsers.clear();
  
  if (!users || !Array.isArray(users) || users.length === 0) {
    userContainer.innerHTML = '<div class="no-users">No users connected</div>';
    return;
  }
  
  // Get current user ID from window scope to ensure it's up-to-date
  const currentUserId = window.currentUser ? window.currentUser.id.toString() : getCurrentUserId();
  console.log('[updateConnectedUsers] Current user ID:', currentUserId);
  
  // Debug current user
  if (window.currentUser) {
    console.log('[updateConnectedUsers] Current user from window:', window.currentUser);
  }
  
  // Filter out current user
  const filteredUsers = users.filter(user => {
    if (!user || !user.id) {
      console.log('Skipping invalid user:', user);
      return false;
    }
    
    const isCurrentUser = user.id.toString() === currentUserId;
    if (isCurrentUser) {
      console.log('Skipping current user:', user);
    }
    
    return !isCurrentUser;
  });
  
  console.log('[updateConnectedUsers] Filtered users:', filteredUsers.length, 'of', users.length);
  
  // If no users after filtering
  if (filteredUsers.length === 0) {
    userContainer.innerHTML = '<div class="no-users">No other users connected</div>';
    return;
  }
  
  // Display filtered users
  filteredUsers.forEach(user => {
    // Add user to the displayed set
    const userId = user.id.toString();
    displayedUsers.add(userId);
    
    console.log('Adding user:', user);
    const userElement = document.createElement('div');
    userElement.id = userId;
    userElement.className = 'user-item';
    // Add CSS class from the module as well
    if (window.chatClasses && window.chatClasses.userItem) {
      userElement.classList.add(window.chatClasses.userItem);
    }
    
    // Add online indicator
    const onlineIndicator = document.createElement('span');
    onlineIndicator.className = 'userOnline';
    userElement.appendChild(onlineIndicator);
    
    // Add user nickname in a span with a class for easy selection
    const nicknameSpan = document.createElement('span');
    nicknameSpan.className = 'user-nickname';
    nicknameSpan.textContent = user.nickname || 'Unknown user';
    userElement.appendChild(nicknameSpan);
    
    // Add user data as attributes for easier access
    userElement.setAttribute('data-nickname', user.nickname || 'Unknown user');
    
    userContainer.appendChild(userElement);
    
    // Update user nickname cache
    userNicknameCache.set(userId, user.nickname);
  });
  
  // Attach click listeners
  setupUserClickListener();
}

// Helper function to get current user ID
function getCurrentUserId() {
  // First try to get it from the window currentUser object
  // which is set by Chat.js and accessible globally
  if (window.currentUser && window.currentUser.id) {
    return window.currentUser.id.toString();
  }
  
  // Try to get from localStorage token
  const token = localStorage.getItem('token');
  if (token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      if (payload && payload.id) {
        return payload.id.toString();
      }
    } catch (e) {
      console.error('[WebSocket] Failed to parse token payload:', e);
    }
  }
  
  // Fallback to parsing from DOM if currentUser is not available
  const userInfoElement = document.getElementById('user-info');
  if (!userInfoElement) return null;
  
  // Parse the user ID from the user info display
  const userIdMatch = userInfoElement.textContent.match(/User ID: (\d+)/);
  if (userIdMatch && userIdMatch[1]) {
    return userIdMatch[1];
  }
  
  return null;
}

export function connectWebSocket(token) {
  if (!token) {
    console.error('[WebSocket] No token provided for connection');
    return;
  }

  // Close existing connection if one exists
  closeWebSocket();
  
  console.log("Token utilisé pour la connexion:", token);
  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  socket.onopen = () => {
    console.log('WebSocket connection established');
      
      // Make sure current user is properly loaded
      if (!window.currentUser) {
        console.warn('[WebSocket] Current user not found in window scope. This might cause issues with sender identification.');
        
        // Try to get it from localStorage as a fallback
        const token = localStorage.getItem('token');
        if (token) {
          try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            if (payload && payload.id) {
              window.currentUser = {
                id: payload.id,
                nickname: payload.nickname || `User ${payload.id}`
              };
              console.log('[WebSocket] Set currentUser from token payload:', window.currentUser);
            }
          } catch (e) {
            console.error('[WebSocket] Failed to parse token payload:', e);
          }
        }
      }
    
    // Clear any existing user list when establishing a new connection
    const userContainer = document.getElementById('connected-users');
    if (userContainer) {
      userContainer.innerHTML = '<div class="loading-users">Loading users...</div>';
    }
    
    // Request user list immediately after connection
    setTimeout(() => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        console.log('[WebSocket] Requesting user list');
        try {
          // Try different message formats since we're not sure which one the server expects
          socket.send(JSON.stringify({ type: 'get_users' }));
          
          // Also send a request for active users after a short delay
          setTimeout(() => {
            if (socket && socket.readyState === WebSocket.OPEN) {
              socket.send(JSON.stringify({ type: 'get_active_users' }));
            }
          }, 300);
        } catch (e) {
          console.error('[WebSocket] Error requesting users:', e);
        }
      }
    }, 500);
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log('[WebSocket] Message received:', data);
      
      // Debug the message content to help identify sender/receiver issues
      console.log('[WebSocket] Message details:', {
        type: data.type,
        from: data.from_id || data.from || 'not specified',
        to: data.to || 'not specified',
        content: data.content ? (data.content.length > 50 ? data.content.substring(0, 50) + '...' : data.content) : 'none'
      });
      throttleMessage(message)
    
      if (data.type === 'user_list' || data.type === 'users') {
        const container = document.getElementById('connected-users');
        if (container) {
          console.log('[WebSocket] Updating user list with:', data.users);
          
          // Ensure users is an array
          const usersList = Array.isArray(data.users) ? data.users : [];
          
          // Check if container is empty
          if (container.children.length === 0 || 
              (container.children.length === 1 && container.children[0].className === 'loading-users')) {
            console.log('[WebSocket] Container was empty, populating from scratch');
          }
          
          updateConnectedUsers(container, usersList);
          
          // Update nickname cache with connected users
          usersList.forEach(user => {
            if (user.id && user.nickname) {
              userNicknameCache.set(user.id.toString(), user.nickname);
            }
          });
        } else {
          console.warn('[WebSocket] Connected users container not found');
        }
      }
      else if (data.type === 'private_message') {
        console.log('[WebSocket] Private message received:', data);
        
        // Get sender ID from different possible properties
        const senderId = data.from_id || data.from || null;
        if (!senderId) {
          console.error('[WebSocket] Message received without sender ID:', data);
          return;
        }
        
        // Extract or look up the sender nickname
        let senderNickname = data.from_nickname || data.sender_nickname || null;
        
        // If nickname not found in message, try to find it in our cache
        if (!senderNickname || senderNickname === 'undefined') {
          console.log('[WebSocket] Nickname not found in message, checking cache for sender:', senderId);
          senderNickname = userNicknameCache.get(senderId.toString());
          
          // If still not found, get it from connected users in DOM
          if (!senderNickname) {
            const userElement = document.getElementById(senderId);
            if (userElement) {
              // Try to get it from the data attribute
              senderNickname = userElement.getAttribute('data-nickname');
              if (!senderNickname) {
                // Or from the nickname span
                const nicknameSpan = userElement.querySelector('.user-nickname');
                senderNickname = nicknameSpan ? nicknameSpan.textContent : `User ${senderId}`;
              }
              console.log('[WebSocket] Found sender nickname in DOM:', senderNickname);
            } else {
              senderNickname = `User ${senderId}`;
              console.log('[WebSocket] Using fallback nickname:', senderNickname);
            }
          } else {
            console.log('[WebSocket] Found sender nickname in cache:', senderNickname);
          }
        }
        
        console.log(`[WebSocket] Processing message from ${senderNickname} (${senderId})`);
        
        // Store the message in conversation history
        import('./pages/Chat.js').then(chatModule => {
          if (chatModule.conversationHistory && typeof chatModule.conversationHistory.addReceivedMessage === 'function') {
            const timestamp = data.timestamp || new Date().toISOString();
            
            // Add to conversation history with verified sender information
            chatModule.conversationHistory.addReceivedMessage(
              senderId,
              senderNickname,
              data.content,
              timestamp
            );
            
            // If we're in a conversation with this user, display the message
            if (window.currentReceiver && 
                window.currentReceiver.id.toString() === senderId.toString()) {
              displayMessage({
                senderId: senderId,
                senderNickname: senderNickname,
                content: data.content,
                timestamp: timestamp
              });
            }
          }
        }).catch(err => {
          console.error('[WebSocket] Error importing Chat module:', err);
          
          // Fallback to simple message display
          displayMessage({
            senderId: senderId,
            senderNickname: senderNickname,
            content: data.content,
            timestamp: data.timestamp || new Date().toISOString()
          });
        });
      }else if (data.type === 'notification') {
        console.log('[WebSocket] 🔔 Notification reçue :', data);
        if (typeof window.showNotificationBadge === 'function') {
          window.showNotificationBadge();
        }
        
        const senderId = data.from;
        const content = data.content;
        const timestamp = data.timestamp || new Date().toISOString();
      
        // Récupère le pseudo si possible
        let senderNickname = userNicknameCache.get(senderId.toString()) || `User ${senderId}`;
      
        // Affiche temporairement une alerte
        alert(`📨 Nouveau message privé de ${senderNickname} : ${content}`);
      
        // TODO : remplacer alert() par un toast/badge plus tard
      }
      
      else {
        console.log('[WebSocket] Unknown message type:', data.type);
      }
    } catch (error) {
      console.error('[WebSocket] Error parsing message:', error, event.data);
    }
  };
  
  socket.onerror = (error) => {
    console.error('[Websocket] Error:', error);
  };

  socket.onclose = (event) => {
    console.log('[Websocket] Connection closed:', event);
  };
}

// Fonction pour retirer un utilisateur de la liste
function removeUserFromList(userID) {
  const userElement = document.getElementById(userID);
  if (userElement) {
    userElement.remove();
  }
}

// Function to get user nickname from ID
async function getUserNickname(userId) {
  // Check cache first
  if (userNicknameCache.has(userId)) {
    return userNicknameCache.get(userId);
  }
  
  // Find in connected users list
  const connectedUsers = document.querySelectorAll('.user-item');
  for (const userEl of connectedUsers) {
    if (userEl.id === userId) {
      const nickname = userEl.textContent;
      userNicknameCache.set(userId, nickname);
      return nickname;
    }
  }
  
  // Fallback to API request if not found in connected users
  try {
    const token = localStorage.getItem('token');
    if (!token) return userId; // Fallback to ID if no token
    
    const response = await fetch(`http://localhost:8080/api/users/${userId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.ok) {
      const userData = await response.json();
      if (userData.data && userData.data.nickname) {
        userNicknameCache.set(userId, userData.data.nickname);
        return userData.data.nickname;
      }
    }
  } catch (error) {
    console.error(`Failed to fetch nickname for user ${userId}:`, error);
  }
  
  // If all else fails, return the user ID
  return userId;
}

export async function displayPrivateMessage(message) {
  // Get nickname for the sender (use provided nickname if available)
  const senderNickname = message.fromNickname || await getUserNickname(message.senderId);
  
  // Check if we're currently in a conversation with this sender
  const isActiveConversation = window.currentReceiver && 
      window.currentReceiver.id.toString() === message.senderId.toString();
  
  // Only display the message if we're in an active conversation with this sender
  if (isActiveConversation) {
    displayMessage({
      senderId: message.senderId,
      senderNickname: senderNickname,
      content: message.content,
      timestamp: message.datetime
    });
  } else {
    // If we're not in a conversation with this sender, we could show a notification
    console.log(`[WebSocket] Received message from ${senderNickname} but not in active conversation`);
    
    // Update conversation history through Chat.js
    import('./pages/Chat.js').then(chatModule => {
      if (chatModule.conversationHistory) {
        chatModule.conversationHistory.addReceivedMessage(
          message.senderId,
          senderNickname,
          message.content,
          message.datetime
        );
      }
    });
  }
}

// Helper function to display a message when Chat module cannot be imported
function displayMessage(message) {
  const messagesContainer = document.getElementById('messages');
  if (!messagesContainer) {
    console.error('[WebSocket] Messages container not found');
    return;
  }
  
  // Format timestamp
  let formattedTime = message.timestamp;
  try {
    const date = new Date(message.timestamp);
    formattedTime = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } catch (e) {
    console.error('[WebSocket] Error formatting timestamp:', e);
  }
  
  // Create and append message element
  const messageElement = document.createElement('div');
  messageElement.className = 'message received';
  
  // Add CSS module classes if available
  if (window.chatClasses && window.chatClasses.message) {
    messageElement.classList.add(window.chatClasses.message);
  }
  if (window.chatClasses && window.chatClasses.received) {
    messageElement.classList.add(window.chatClasses.received);
  }
  
  messageElement.innerHTML = `
    <strong>${message.senderNickname}</strong>: ${message.content}<br>
    <small>${formattedTime}</small>
  `;
  
  messagesContainer.appendChild(messageElement);
  messagesContainer.scrollTop = messagesContainer.scrollHeight;
  
  // Cache this nickname for future use
  if (message.senderId) {
    userNicknameCache.set(message.senderId.toString(), message.senderNickname);
  }
}

export function getSocket() {
  return socket;
}

let lastMessageTime = 0;
const throttleDelay = 1000;  // Délai en ms (1 seconde)

function throttleMessage(message) {
    const now = Date.now();

    // Si le délai minimum est passé, on affiche le message
    if (now - lastMessageTime > throttleDelay) {
        lastMessageTime = now;

        // Ajoute le message au chat
        const messageDiv = document.createElement("div");
        messageDiv.className = message.sender_id === currentUserId ? "message sent" : "message received";
        messageDiv.textContent = message.content;

        // Ajoute le message à l'élément container
        document.getElementById("messagesContainer").appendChild(messageDiv);
    }
}
