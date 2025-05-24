var isLoadingMessages = false;
let ws = null;
let nusers;
let currentChatUser = null;
let messageOffset = 0;
let allMessagesLoaded = false;
let previousScrollHeight = 0;
let hasMoreMessages = true;
let tptimeout;

function connectWebSocket() {
    ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        ws.send(JSON.stringify({ type: "get_users" }));
    };
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        switch (data.type) {
            case "users_update":
                convupdate(data.users);
                nusers = data.users;
                break;
            case "new_post":
                // console.log(window.location.pathname);
                if (window.location.pathname == "/"){
                    fetchPosts(data.islogged);
                }
                break;
            case "chat":
                if (data.sender === currentChatUser ||
                    (data.receiver === localStorage.getItem('username') && !currentChatUser)) {
                    displayChatMessage(data);
                }
                convupdate(nusers);
                break;
            case "chat_history":
                displayChatHistory(data);
                break;
            case "update_conversations":
                convupdate(nusers);
                break;
            case "typing":
                if (data.sender === currentChatUser) {
                    istyping(data.sender, data.isTyping);
                }
                break;
        }
    };
    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        setTimeout(connectWebSocket, 3000);
    };

    ws.onclose = () => {
        setTimeout(connectWebSocket, 3000);
    };
}


function convupdate(users = null) {
    const currentUser = localStorage.getItem('username');
    const container = document.getElementById('conversations-list');
    if (!container) return;
    
    fetch(`/api/conversations`)
        .then(response => response.json())
        .then(conversations => {
            // console.log(conversations);
            if (users) {
                const onlineUsers = users.filter(u => u !== currentUser);
                // console.log(onlineUsers);
                let offlineusers = conversations.users.filter(usr => !onlineUsers.includes(usr));
                offlineusers = offlineusers.filter(u => u != currentUser);
                // console.log(offlineusers, "offffffine");
                onlineUsers.forEach(user => {
                    // conversations.users.filter(us => us !== user);
                    if (!conversations.conv.some(c => c.user === user)) {
                        conversations.conv.push({
                            user: user,
                            last_message: null,
                            unread: 0,
                            enligne: true
                        });
                        // console.log(conversations.users)
                    }
                    // console.log("lloodsfsdofosdof");
                    conversations.conv.forEach(c => {
                        if (c.user === user) {
                            c.enligne = true;
                            // console.log(c.enligne, "eeeeeeeeeeeeeeeeeeeeeeeeeee");
                        }
                    })
                });
                conversations.conv.forEach(c => {
                    offlineusers = offlineusers.filter( u => u !== c.user);
                });
                for (i=0;i<offlineusers.length;i++){
                    conversations.conv.push({
                        user : offlineusers[i],
                        last_message:null,
                        unread: 0,
                        enligne:false
                    })
                }
                // console.log(offlineusers,"offfffffffff without conv");
                // conversations.users.forEach{us => {

                // }}
            }

            conversations.conv.sort((a, b) => {
                if (a.last_message && b.last_message) {
                    return new Date(b.last_message.created_at) - new Date(a.last_message.created_at);
                } else if (a.last_message) {
                    return -1;
                } else if (b.last_message) {
                    return 1;
                } else {
                    return a.user.localeCompare(b.user);
                }
            });
            // console.log(conversations);
            container.innerHTML = conversations.conv.map(conv =>
                `<div class="conversation ${conv.user === currentChatUser ? 'active' : ''} ${conv.unread > 0 ? 'unread' : ''}" 
             onclick="startChat('${conv.user}')">
            <div class="user-name">${conv.user}</div>
            ${conv.last_message ? `
                <div class="last-message">${conv.last_message}</div>
                <div class="message-time">${formatTime(conv.msg_time)}</div>
            ` : ''}
            ${conv.unread > 0 ? `<div class="unread-count">${conv.unread}</div>` : ''}
            <div class="${conv.enligne == true ?  'online' : 'offline'}">${conv.enligne == true ?  'online' : 'offline'}</div>

        </div>`
            ).join('');
        })
        .catch(error => console.error("Error fetching conversations:", error));
}


function displayChatMessage(message) {
    const chatMessages = document.getElementById('chat-messages');
    if (!chatMessages) return;

    const messageElement = document.createElement('div');
    messageElement.className = `message ${message.sender === localStorage.getItem('username') ? 'sent' : 'received'}`;
    messageElement.innerHTML = `
        <div class="message-sender"><a href="/profile?user=${message.sender}" onclick="navevent(event)">${message.sender}</a></div>
        <div class="message-content">${message.content}</a></div>
        <div class="message-time">${formatTime(message.created_at)}</div>
    `;

    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}




function startChat(username) {
    navigate('/chat');
    currentChatUser = username;
    document.getElementById('chat-with-user').textContent = `Chat with ${username}`;

    const messageInput = document.getElementById('chatMessage');
    const sendButton = document.getElementById('sendMessage');
    messageInput.disabled = false;
    sendButton.disabled = false;

    messageInput.addEventListener('input', function() {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: "typing",
                sender: localStorage.getItem('username'),
                receiver: username,
                isTyping: true
            }));

            clearTimeout(tptimeout);
            
            tptimeout = setTimeout(() => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({
                        type: "typing",
                        sender: localStorage.getItem('username'),
                        receiver: username,
                        isTyping: false
                    }));
                }
            }, 2000);
        }
    });


    sendButton.onclick = function() {
        const message = messageInput.value.trim();
        if (message && ws) {
            const chatMessages = document.getElementById('chat-messages');
            const messageElement = document.createElement('div');
            messageElement.className = 'message sent';
            messageElement.innerHTML = `
                <div class="message-sender">${localStorage.getItem('username')}</div>
                <div class="message-content">${message}</div>
                <div class="message-time">${formatTime(new Date().toISOString())}</div>
            `;
            chatMessages.appendChild(messageElement);
            chatMessages.scrollTop = chatMessages.scrollHeight;

            ws.send(JSON.stringify({
                type: "chat",
                sender: localStorage.getItem('username'),
                receiver: username,
                content: message
            }));

            messageInput.value = '';
        }
    };

    messageInput.onkeypress = function(e) {
        if (e.key === 'Enter') {
            sendButton.click();
        }
    };
    loadConversation(username);
    // convupdate(nusers);
    setupscrolllst();
}


function istyping(username, isTyping) {
    const chatHeader = document.getElementById('chat-with-user');
    if (!chatHeader) return;

    if (isTyping) {
        chatHeader.textContent = `${username} is typing...`;
    } else {
        chatHeader.textContent = `Chat with ${username}`;
    }
}
function setupscrolllst() {
    const chatMessages = document.getElementById('chat-messages');
    if (chatMessages) {
        console.log("Setting up scroll listener!");
        chatMessages.addEventListener('scroll', throttle(() => {
            if (chatMessages.scrollTop < 10 && !isLoadingMessages && !allMessagesLoaded) {
                console.log("  Conditions met, calling loadMoreMessages()");
                loadMoreMessages();
            }
        }, 150));
    }
}


function createMessageElement(msg) {
    const div = document.createElement('div');
    div.className = `message ${msg.sender === localStorage.getItem('username') ? 'sent' : 'received'}`;
    div.innerHTML = `
        <div class="message-sender">${msg.sender}</div>
        <div class="message-content">${msg.content}</div>
        <div class="message-time">${formatTime(msg.created_at)}</div>
    `;
    return div;
}

function formatTime(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit'
    });
}

function isScrolledToBottom() {
    const chatMessages = document.getElementById('chat-messages');
    return chatMessages.scrollTop + chatMessages.clientHeight >= chatMessages.scrollHeight - 20;
}

function loadConversation(user) {
    messageOffset = 0;
    allMessagesLoaded = false;
    isLoadingMessages = false;
    currentChatUser = user;

    const chatMessages = document.getElementById('chat-messages');
    chatMessages.innerHTML = '';
    

    document.getElementById('chat-with-user').textContent = `Chat with ${user}`;
    document.getElementById('chatMessage').disabled = false;
    document.getElementById('sendMessage').disabled = false;

    fetch(`/api/mark-read?user=${user}`, {
            method: 'POST'
        })
        .then(() => convupdate(nusers))
        .catch(err => console.error("Error marking messages as read:", err));

    loadMoreMessages();
}

function loadMoreMessages() {
    if (isLoadingMessages || allMessagesLoaded || !currentChatUser) {
        console.log("loadMoreMessages aborted:", {
            isLoadingMessages,
            allMessagesLoaded,
            currentChatUser
        });
        return;
    }

    isLoadingMessages = true;
    console.log("loadMoreMessages called", {
        messageOffset
    });

    ws.send(JSON.stringify({
        type: "get_history",
        user: currentChatUser,
        offset: messageOffset,
        limit: 10
    }));

    const chatMessages = document.getElementById('chat-messages');
    const scrollHeightBefore = chatMessages.scrollHeight;
    const scrollTopBefore = chatMessages.scrollTop;

    window.restoreScrollPosition = function() {
        const scrollHeightAfter = chatMessages.scrollHeight;
        chatMessages.scrollTop = scrollTopBefore + (scrollHeightAfter - scrollHeightBefore);
        delete window.restoreScrollPosition;
    };
}


function displayChatHistory(data) {
    if (data.user !== currentChatUser) return;

    const chatMessages = document.getElementById('chat-messages');
    console.log("displayChatHistory called", data); // Log the data

    if (!data.messages || data.messages.length === 0) {
        allMessagesLoaded = true;
        console.log("  No more messages from server");
        if (messageOffset === 0) {
            chatMessages.innerHTML = '<div class="no-messages">No messages yet</div>';
        }
    } else {
        const fragment = document.createDocumentFragment();
        data.messages.forEach(msg => {
            fragment.appendChild(createMessageElement(msg));
        });

        if (messageOffset > 0) {
            chatMessages.insertBefore(fragment, chatMessages.firstChild);
            if (window.restoreScrollPosition) {
                window.restoreScrollPosition();
            }

        } else {
            chatMessages.innerHTML = '';
            chatMessages.appendChild(fragment);
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        messageOffset += data.messages.length;
        console.log("  Received", data.messages.length, "messages. new offset:", messageOffset);
        if (data.messages.length < 10) {
            allMessagesLoaded = true;
            console.log("  Fewer than 10 messages, assuming all loaded");
        }
    }

    isLoadingMessages = false;
    console.log("  displayChatHistory finished", {
        isLoadingMessages,
        allMessagesLoaded,
        messageOffset
    });
}

function throttle(func, limit) {
    let lastFunc;
    let lastRan;
    return function() {
        const context = this;
        const args = arguments;
        if (!lastRan) {
            func.apply(context, args);
            lastRan = Date.now();
        } else {
            clearTimeout(lastFunc);
            lastFunc = setTimeout(function() {
                if ((Date.now() - lastRan) >= limit) {
                    func.apply(context, args);
                    lastRan = Date.now();
                }
            }, limit - (Date.now() - lastRan));
        }
    };
}