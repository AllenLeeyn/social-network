import { currentState, MESSENGER_DISPLAY, showMessage, showTab, renderDisplay } from "./main.js";
import { templateUserList, templateChat, templateChatHistory, templateChatMessage } from "./template.js";

const USER_LIST = document.getElementById("userList");

let socket;

export function openWebSocket() {
    socket = new WebSocket("/ws");
    socket.onopen = function() {
        console.log("Websocket opened.");
    };
    socket.onmessage = onMessageHandler;
    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
    socket.onclose = function() {
        console.log("Websocket closed.");
    };
}

function onMessageHandler(event) {
    const data = JSON.parse(event.data)                                                                                                                                                                                                                                    
    
    if (data.action === "start") {
        addUserListItems(data)

    } else if (data.action === "online") {
        updateUserListItems(data.id, "online")

    } else if (data.action === "offline") {
        updateUserListItems(data.id, "offline")

    } else if (data.action === "messageSendOK") {
        const messageInput = document.getElementById('message-input');
        messageInput.value = "";

    } else if (data.action === "messageHistory") {
        processMessageHistory(data);

    } else if (data.action === "message") {
        processMessage(data);

    } else if (data.action === "typing") {
        processTyping(data);
    }
}

let typingTimer;
function processTyping (data) {
    if (data.senderID !== currentState.chatID) return;
    const typingIndicator = document.getElementById('typing-indicator');

    if (typingIndicator.textContent !== "") {
        clearTimeout(typingTimer);
    } else {
        typingIndicator.textContent = `${currentState.chat} is typing`;
        typingIndicator.innerHTML += `<svg width="30" height="10" viewBox="0 0 50 10" xmlns="http://www.w3.org/2000/svg">
            <circle cx="10" cy="8" r="3" fill="hsl(0, 0%, 80%)">
              <animate attributeName="r" values="3;5;3" dur="1.2s" repeatCount="indefinite" begin="0s"/>
            </circle>
            <circle cx="25" cy="8" r="3" fill="hsl(0, 0%, 80%)">
              <animate attributeName="r" values="3;5;3" dur="1.2s" repeatCount="indefinite" begin="0.1s"/>
            </circle>
            <circle cx="40" cy="8" r="3" fill="hsl(0, 0%, 80%)">
              <animate attributeName="r" values="3;5;3" dur="1.2s" repeatCount="indefinite" begin="0.2s"/>
            </circle>
          </svg>`;
    }
    
    typingTimer = setTimeout(() => {
        typingIndicator.textContent = ""; 
    }, 800);
}

function processMessage(data) {
    const MESSAGE_CONTAINER = document.getElementById("message-container");

    if (data.receiverID === currentState.id && data.senderID === currentState.chatID) {
        
        MESSAGE_CONTAINER.innerHTML += templateChatMessage(data);
        MESSAGE_CONTAINER.scrollTop = MESSAGE_CONTAINER.scrollHeight;
        sendMessage("messageAck", currentState.chatID, "");

    } else if (data.senderID === currentState.id) {
        MESSAGE_CONTAINER.innerHTML += templateChatMessage(data);
        MESSAGE_CONTAINER.scrollTop = MESSAGE_CONTAINER.scrollHeight;

    } else if (data.receiverID === currentState.id && data.senderID !== currentState.chatID) {
        const clientElement = document.getElementById(`user-${data.senderID}`);
        clientElement.classList.add("unread");
    }

    let clientElement = document.getElementById(`user-${data.receiverID}`).parentElement;
    if (data.receiverID === currentState.id) {
        clientElement = document.getElementById(`user-${data.senderID}`).parentElement;
    }
    USER_LIST.prepend(clientElement);

    const typingIndicator = document.getElementById('typing-indicator');
    if (typingIndicator) typingIndicator.textContent = "";
}

function processMessageHistory(data) {
    const MESSAGE_CONTAINER = document.getElementById("message-container");
    const oldScrollPosition = MESSAGE_CONTAINER.scrollTop - MESSAGE_CONTAINER.scrollHeight;

    MESSAGE_CONTAINER.innerHTML = templateChatHistory(data.content) + MESSAGE_CONTAINER.innerHTML;
    sendMessage("messageAck", currentState.chatID, "");
    document.getElementById(`user-${currentState.chatID}`).classList.remove("unread");
    
    MESSAGE_CONTAINER.scrollTop = oldScrollPosition + MESSAGE_CONTAINER.scrollHeight;
}

function updateUserListItems(tgtID, action) {
    const clientElement = document.getElementById(`user-${tgtID}`);
    if (action === "online") {
        clientElement.classList.add(action);
    } else {
        clientElement.classList.remove("online");
    }
}

function addUserListItems(data) {
    const clientList = data.allClients.map((name, index) => ({
        name: name, id: data.clientIDs[index]
    }));

    USER_LIST.innerHTML = templateUserList(clientList);
    data.onlineClients.forEach(client=>{
        const clientElement = document.getElementById(`user-${client}`);
        clientElement.classList.add("online");
    });

    if (Array.isArray(data.unreadMsgClients))
    data.unreadMsgClients.forEach(client=>{
        const clientElement = document.getElementById(`user-${client}`);
        clientElement.classList.add("unread");
    });

    const listItems = document.querySelectorAll(".user-item");
    listItems.forEach(item => addUserListItemListener(item));
}

function addUserListItemListener(item) {
    const userName = item.textContent;
    const userID = item.dataset.id;

    item.onclick = (event) => {
        event.preventDefault();

        currentState.display = MESSENGER_DISPLAY;
        currentState.chat = userName;
        currentState.chatID = parseInt(userID);
        showTab("chat", userName);
        renderDisplay();

        isThrottled = false;
        MESSENGER_DISPLAY.innerHTML = templateChat(userID);
        sendMessage("messageReq", userID, "-1");

        document.getElementById('message-input').onkeydown = handleKeyDown;
        document.getElementById('submit-message').onclick = submitMessage;
        document.getElementById("message-container").onscroll = (event) => {messageScollHandler(event, userID)};
    }
}

let isThrottled = false;
function messageScollHandler(event, tgtID){
    if (event.target.scrollTop > 20 || isThrottled) return;
    isThrottled = true;

    const firstChild = event.target.children[0];
    if (firstChild.textContent === " End of history ") return;

    setTimeout(() => {
        isThrottled = false;
    }, 300);

    const msgId = firstChild.dataset.id;
    sendMessage("messageReq", tgtID, msgId)
}

function submitMessage() {
    const receiverID = document.getElementById('submit-message').dataset.id;
    const messageInput = document.getElementById('message-input');
    const messageText = messageInput.value.trim();

    if (messageText === "") {
            showMessage('Message cannot be empty!');
            return;
    }
    sendMessage("message", receiverID, messageText);
}

function sendMessage(action, receiverID, content){
    const message = {
        action: action,
        receiverID: parseInt(receiverID),
        content: content
    }
    socket.send(JSON.stringify(message))
}

let isTyping = false;
function handleKeyDown(event) {
    if (event.key === 'Enter') {
        submitMessage();
    } else {
        if (!isTyping){
            isTyping = true;
            sendMessage("typing", currentState.chatID, "");
        }
        setTimeout(() => {
            isTyping = false;
        }, 600); 
    }
}
