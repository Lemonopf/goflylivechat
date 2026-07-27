const CHAT_WIDGET = {
    API_URL: "",
    AGENT_ID: "",
    AUTO_OPEN: true,
    DISPLAY_MODE: 1,
    TOKEN: "",
    THEME: "",
    USER_NAME: "",
    USER_AVATAR: "",
    isChatOpen: false,
    originalPageTitle: document.title,
    chatWindowTitle: "在线客服",
    isOffline: false,
    iframeId: "chat-widget-iframe",
    containerId: "chat-widget-container"
};

CHAT_WIDGET.initialize = function(config) {
    // Apply configuration
    for (let key in config) {
        if (this.hasOwnProperty(key)) {
            this[key] = config[key];
        }
    }

    // Normalize URL by removing trailing slash
    if (this.API_URL) {
        this.API_URL = this.API_URL.replace(/\/$/, "");
    }
    this.THEME = this.normalizeTheme(this.THEME || this.detectTheme());

    // Add required CSS styles
    this.injectStyles();

    // Display the chat button
    this.createChatButton();

    // Set up event handlers
    this.setupEventHandlers();
};

CHAT_WIDGET.injectStyles = function() {
    const style = document.createElement('style');
    style.textContent = `
        #chat-widget-button,
        #chat-widget-container {
            --cw-primary-50: #f0fdfa;
            --cw-primary-100: #ccfbf1;
            --cw-primary-500: #14b8a6;
            --cw-primary-600: #0d9488;
            --cw-bg: #ffffff;
            --cw-border: #e2e8f0;
            --cw-text: #0f172a;
            --cw-muted: #64748b;
            --cw-shadow: 0 20px 50px rgba(15, 23, 42, 0.16);
            font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "Helvetica Neue", Arial, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
        }
        #chat-widget-button[data-chat-widget-theme="dark"],
        #chat-widget-container[data-chat-widget-theme="dark"] {
            --cw-bg: #0f172a;
            --cw-border: #334155;
            --cw-text: #f8fafc;
            --cw-muted: #94a3b8;
            --cw-shadow: 0 24px 60px rgba(0, 0, 0, 0.45);
        }
        #chat-widget-button {
            position: fixed;
            bottom: 20px;
            right: 20px;
            width: 60px;
            height: 60px;
            background: linear-gradient(135deg, #14b8a6 0%, #0d9488 100%);
            color: white;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            box-shadow: 0 0 20px rgba(20, 184, 166, 0.25), 0 2px 10px rgba(15, 23, 42, 0.18);
            z-index: 9999;
            transition: transform 0.18s ease, box-shadow 0.18s ease;
        }

        #chat-widget-button:hover {
            transform: translateY(-1px);
            box-shadow: 0 0 30px rgba(20, 184, 166, 0.34), 0 10px 24px rgba(15, 23, 42, 0.22);
        }
        
        #chat-widget-button .notification-badge {
            position: absolute;
            top: -5px;
            right: -5px;
            background-color: #FF5722;
            color: white;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
        }
        
        #chat-widget-container {
            position: fixed;
            bottom: 90px;
            right: 20px;
            width: 350px;
            height: 500px;
            background: var(--cw-bg);
            border: 1px solid var(--cw-border);
            border-radius: 12px;
            box-shadow: var(--cw-shadow);
            display: none;
            flex-direction: column;
            z-index: 9998;
            overflow: hidden;
            color: var(--cw-text);
        }
        
        #chat-widget-header {
            min-height: 48px;
            padding: 0 14px;
            background: linear-gradient(135deg, #14b8a6 0%, #0d9488 100%);
            color: white;
            display: flex;
            align-items: center;
            font-size: 14px;
            font-weight: 600;
        }
        
        #chat-widget-iframe {
            flex: 1;
            border: none;
            background: var(--cw-bg);
        }
        
        .close-button {
            margin-left: auto;
            cursor: pointer;
            width: 30px;
            height: 30px;
            line-height: 28px;
            border-radius: 8px;
            font-size: 20px;
            text-align: center;
            transition: background-color 0.18s ease;
        }

        .close-button:hover {
            background: rgba(255, 255, 255, 0.16);
        }
        @media (max-width: 800px) {
          #chat-widget-container {
            width: 100% !important;
            right: 0 !important;
            bottom: 80px !important;
          }
        }
    `;
    document.head.appendChild(style);
};

CHAT_WIDGET.normalizeTheme = function(theme) {
    return theme === "dark" ? "dark" : "light";
};

CHAT_WIDGET.detectTheme = function() {
    if (document.documentElement && document.documentElement.classList.contains("dark")) {
        return "dark";
    }
    try {
        const savedTheme = window.localStorage && window.localStorage.getItem("theme");
        if (savedTheme === "dark" || savedTheme === "light") {
            return savedTheme;
        }
    } catch (e) {}
    if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
        return "dark";
    }
    return "light";
};

CHAT_WIDGET.createChatButton = function() {
    const button = document.createElement('div');
    button.id = 'chat-widget-button';
    button.setAttribute("data-chat-widget-theme", this.THEME);
    button.innerHTML = `
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M21 15C21 15.5304 20.7893 16.0391 20.4142 16.4142C20.0391 16.7893 19.5304 17 19 17H7L3 21V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H19C19.5304 3 20.0391 3.21071 20.4142 3.58579C20.7893 3.96086 21 4.46957 21 5V15Z" fill="white"/>
        </svg>
        <div class="notification-badge" style="display: none;">0</div>
    `;
    document.body.appendChild(button);

    button.addEventListener('click', () => {
        this.openChatWindow();
    });

    // Open automatically if configured
    if (this.AUTO_OPEN) {
        setTimeout(() => {
            this.openChatWindow();
        }, 3000);
    }
};

CHAT_WIDGET.openChatWindow = function() {
    if (this.isChatOpen) return;

    const badge = document.querySelector('#chat-widget-button .notification-badge');
    badge.style.display = 'none';
    badge.textContent = '0';

    // Create container if it doesn't exist
    if (!document.getElementById(this.containerId)) {
        const container = document.createElement('div');
        container.id = this.containerId;
        container.setAttribute("data-chat-widget-theme", this.THEME);
        container.innerHTML = `
            <div id="chat-widget-header">
                <span>${this.chatWindowTitle}</span>
                <span class="close-button">×</span>
            </div>
            <iframe id="${this.iframeId}" src="${this.buildChatUrl()}"></iframe>
        `;
        document.body.appendChild(container);

        // Add close button handler
        document.querySelector(`#${this.containerId} .close-button`).addEventListener('click', () => {
            this.closeChatWindow();
        });
    }

    // Show the chat window
    document.getElementById(this.containerId).style.display = 'flex';
    this.isChatOpen = true;

    // Hide the floating button
    document.getElementById('chat-widget-button').style.display = 'none';
};

CHAT_WIDGET.closeChatWindow = function() {
    document.getElementById(this.containerId).style.display = 'none';
    this.isChatOpen = false;
    document.getElementById('chat-widget-button').style.display = 'flex';
};

CHAT_WIDGET.buildChatUrl = function() {
    let url = `${this.API_URL}/livechat?customer_id=${this.AGENT_ID}`;
    url += `&theme=${encodeURIComponent(this.THEME)}`;
    url += `&ui_mode=embedded`;

    // TOKEN：sub2api 用户登录 token，传入即为登录态（后端校验真伪）
    if (this.TOKEN) {
        url += `&token=${encodeURIComponent(this.TOKEN)}`;
    }
    // USER_NAME / USER_AVATAR 通过 extra 参数传递（base64 编码的 JSON）
    const extra = {};
    if (this.USER_NAME) {
        extra.visitorName = this.USER_NAME;
    }
    if (this.USER_AVATAR) {
        extra.visitorAvatar = this.USER_AVATAR;
    }
    if (Object.keys(extra).length > 0) {
        url += `&extra=${encodeURIComponent(btoa(unescape(encodeURIComponent(JSON.stringify(extra)))))}`;
    }

    return url;
};

CHAT_WIDGET.setupEventHandlers = function() {
    // Handle messages from the chat iframe
    window.addEventListener('message', (e) => {
        if (!e.data || !e.data.type) return;

        switch (e.data.type) {
            case 'new_message':
            case 'message':
                this.handleIncomingMessage(e.data);
                break;
            case 'close_chat':
                this.closeChatWindow();
                break;
        }
    });
};

CHAT_WIDGET.handleIncomingMessage = function(data) {
    // Update notification badge
    const badge = document.querySelector('#chat-widget-button .notification-badge');
    let count = parseInt(badge.textContent || '0');

    if (!this.isChatOpen) {
        count++;
        badge.textContent = count;
        badge.style.display = 'flex';

        // Flash title if window is not focused
        if (!document.hasFocus()) {
            this.notifyWithTitleFlash();
        }
    }
};

CHAT_WIDGET.notifyWithTitleFlash = function() {
    let isFlashing = true;
    const flashInterval = setInterval(() => {
        document.title = isFlashing ? "新消息！" : this.originalPageTitle;
        isFlashing = !isFlashing;
    }, 1000);

    // Stop flashing when window regains focus
    window.addEventListener('focus', () => {
        clearInterval(flashInterval);
        document.title = this.originalPageTitle;
    }, { once: true });
};
