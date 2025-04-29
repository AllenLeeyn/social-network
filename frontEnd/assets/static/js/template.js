import { currentState } from "./main.js";

export const templateNoFound = (value) => `
  <tbody><tr><td>
    <div class="post-card">
      <h3>No ${value} found</h3>
    </div>
  </tbody></tr></td>`;

export const templateNewPost = (categories) => `
  <form id="newPostForm">
    <div class="input-group">
      <input type="text" name="title" placeholder="Title" required />
    </div>
    <div class="input-group">
      <textarea name="content" placeholder="Write your content here..." rows="10" required ></textarea>
    </div>
    <div class="input-group">
      <h4>Click to select categories:</h4>
      <div class="checkbox-group">
        ${categories.map((cat, index) => `
          <div class="checkbox-item">
            <input type="checkbox" id="category${index}" name="categories" value=${index}>
            <label for="category${index}">${cat}</label>
          </div>
          `).join('')}
      </div>
    </div>
    <div class="input-group">
      <button class="new-post" id="newPostSubmit" type="submit">Create Post</button>
    </div>
  </form>`;

export const templatePostCard = (post) => `
  <div>
    <div class="post-header">
      <div class="post-meta">
        <a href="/profile?id=${post.userID}" class="post-author">${post.userName}</a>
        <div class="post-time">${relativeTime(post.createdAt)}</div>
      </div>
    </div>
    <div class="post-content">
      <h3>
        <a href="/post?id=${post.ID}">${post.title}</a>
      </h3>
      <pre>${post.content}</pre>
    </div>
    <div class="post-actions" data-id=${post.ID} data-state="${post.rating}" data-for="post">
      <button class="icon-button like-button" data-id=${post.ID} data-for="post">
        <i class="fas fa-thumbs-up"></i> <span>${post.likeCount}</span>
      </button>
      <button class="icon-button dislike-button" data-id=${post.ID} data-for="post">
        <i class="fas fa-thumbs-down"></i> <span>${post.dislikeCount}</span>
      </button>
      <button class="icon-button">
        <i class="fas fa-comment" href="/post?id=${post.ID}"></i> <span>${post.commentCount}</span>
      </button>
      <p class="icon-button">
        <span>${post.catNames}</span>
      </p>
    </div>
  </div>`;

export const templateCommentCard = (comment) => `
  <div>
    <div class="post-header">
      <div class="post-meta">
        <a href="/profile?id=${comment.userID}" class="post-author">${comment.userName}</a>
        <div class="post-time">${relativeTime(comment.createdAt)}</div>
      </div>
    </div>
    <div class="post-content">
      <pre>${comment.content}</pre>
    </div>
    <div class="post-actions" data-id=${comment.ID} data-state="${comment.rating}" data-for="comment">
      <button class="icon-button like-button" data-id=${comment.ID} data-for="comment">
      <i class="fas fa-thumbs-up"></i> <span>${comment.likeCount}</span>
      </button>
      <button class="icon-button dislike-button" data-id=${comment.ID} data-for="comment">
      <i class="fas fa-thumbs-down"></i> <span>${comment.dislikeCount}</span>
      </button>
    </div>
  </div>`;

export const templateCommentForm = (post) => `
  <div id="comment-form">
    <textarea id="comment-input" placeholder="Write your comment here..." rows="3"></textarea>
    <button id="submit-comment" data-id=${post.ID}>Post Comment</button>
  </div>`;

export const templateProfileCard = (data) => {
  let result = `
  <div class="post-header">
    <div class="post-meta">
      <h3>Name: ${data.name}</h3>
    </div>
  </div>
  <div>
    <h4>Posts:</h4>
  </div>`;

  if (!Array.isArray(data.posts) || data.posts.length === 0){
    result+=templateNoFound("post");
  } else {
    data.posts.forEach(post => {
      result += templatePostCard(post);
    });
  }
  return result;
};

export const templateCategoriesList = (categories) => {
  let result = `
    <li><a href="/posts" class="category-item active">All</a></li>
    <li><a href="/posts?filterBy=createdBy" class="category-item">My posts</a></li>
    <li><a href="/posts?filterBy=likedBy" class="category-item">Liked posts</a></li>`;

  categories.forEach((category, index) => {
    result += `
    <li><a href="/posts?filterBy=category&id=${index}" class="category-item">${category}</a></li>`;
  })
  return result;
};

export const templateUserList = (clientList) => {
  let result = ``;

  clientList.forEach((client) => {
  result += `
    <li><a class="user-item" data-id=${client.id} id="user-${client.id}">${client.name}</a></li>`;
  })
  return result;
};

export const templateChat = (userId) => `
  <tr><tc><td>
    <div id="message-container"></div>
    
    <div><span id="typing-indicator" class="typing-indicator"></div>

    <div id="message-form">
      <textarea id="message-input" rows="3"></textarea>
      <button id="submit-message" data-id=${userId}>Send</button>
    </div>
  </td></tc></tr>
`;

export const templateChatHistory = (messages) => {
  let result = ``;
  if (!Array.isArray(messages)) return `<div class="nuetral"> End of history </div>`
  
  messages.forEach(msg => {
    result = templateChatMessage(msg) + result
  });
  return result;
};

export const templateChatMessage = (message) => `
    <div class="${(message.senderID === currentState.id) ? "sender": "receiver"}" data-id=${message.ID}>
      <pre>${message.content}</pre>
      <p>[${(message.senderID === currentState.id) ? currentState.user: currentState.chat}] ${formatTime(message.createdAt)}</p></div>`
;

function relativeTime(timestamp) {
  const timeObj = new Date(timestamp);
  const now = Date.now()
  const diff = now - timeObj;

  if (diff < 60000) return `moments ago`;
  if (diff < 3600000) return `${Math.floor(diff / 60000)} minutes ago`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} hours ago`;
  return `${Math.floor(diff / 86400000)} days ago`;
};

function formatTime(timestamp) {
  const timeObj = new Date(timestamp);

  const year = timeObj.getFullYear();
  const month = String(timeObj.getMonth() + 1).padStart(2, '0');
  const day = String(timeObj.getDate()).padStart(2, '0');
  const hours = String(timeObj.getHours()).padStart(2, '0');
  const minutes = String(timeObj.getMinutes()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}`;
}